// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lifecycle

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/pipe-cd/pipecd/pkg/backoff"
)

const runBinaryRetryCount = 3

type Command struct {
	cmd       *exec.Cmd
	stoppedCh chan struct{}
	result    atomic.Pointer[error]
}

func (c *Command) IsRunning() bool {
	select {
	case _, notClosed := <-c.stoppedCh:
		return notClosed
	default:
		return true
	}
}

func (c *Command) GracefulStop(period time.Duration) error {
	// For graceful shutdown, we send SIGTERM signal to old Piped process
	// and wait grace-period of time before force killing it.
	c.cmd.Process.Signal(syscall.SIGTERM)
	timer := time.NewTimer(period)

	select {
	case <-timer.C:
		c.cmd.Process.Kill()
		<-c.stoppedCh
		if perr := c.result.Load(); perr != nil {
			return *perr
		}
		return nil
	case <-c.stoppedCh:
		if perr := c.result.Load(); perr != nil {
			return *perr
		}
		return nil
	}
}

func RunBinary(ctx context.Context, execPath string, args []string) (*Command, error) {
	cmd, err := backoff.NewRetry(runBinaryRetryCount, backoff.NewConstant(5*time.Second)).Do(ctx, func() (interface{}, error) {
		cmd := exec.Command(execPath, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return nil, err
		}

		c := &Command{
			cmd:       cmd,
			stoppedCh: make(chan struct{}),
			result:    atomic.Pointer[error]{},
		}
		go func() {
			err := cmd.Wait()
			c.result.Store(&err)
			close(c.stoppedCh)
		}()

		return c, nil
	})

	if err != nil {
		return nil, err
	}

	return cmd.(*Command), nil // The return type is always *Command.
}

// DownloadBinary downloads a file from the given URL into the specified path
// this also marks it executable and returns its full path.
func DownloadBinary(sourceURL, destDir, destFile string, logger *zap.Logger) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("could not create directory %s (%w)", destDir, err)
	}
	destPath := filepath.Join(destDir, destFile)

	// If the destination is already existing, just return its path.
	if _, err := os.Stat(destPath); err == nil {
		return destPath, nil
	}

	// Make a temporary file to save downloaded data.
	tmpFile, err := os.CreateTemp(destDir, "download")
	if err != nil {
		return "", fmt.Errorf("could not create temporary file (%w)", err)
	}
	var (
		tmpName = tmpFile.Name()
		done    = false
	)

	defer func() {
		tmpFile.Close()
		if !done {
			os.Remove(tmpName)
		}
	}()

	logger.Info("downloading binary", zap.String("url", sourceURL))

	u, err := url.Parse(sourceURL)
	if err != nil {
		return "", fmt.Errorf("could not parse URL %s (%w)", sourceURL, err)
	}

	switch u.Scheme {
	case "oci":
		if err := downloadOCI(context.TODO(), destDir, tmpFile, sourceURL, false, runtime.GOOS, runtime.GOARCH); err != nil {
			return "", fmt.Errorf("could not download from %s to %s (%w)", sourceURL, tmpName, err)
		}

	case "http", "https":
		if err := downloadHTTP(tmpFile, sourceURL); err != nil {
			return "", fmt.Errorf("could not download from %s to %s (%w)", sourceURL, tmpName, err)
		}

	case "file":
		if err := downloadFile(tmpFile, u.Path); err != nil {
			return "", fmt.Errorf("could not download from %s to %s (%w)", sourceURL, tmpName, err)
		}

	default:
		return "", fmt.Errorf("unsupported file scheme %s", u.Scheme)
	}

	if err := os.Chmod(tmpName, 0755); err != nil {
		return "", fmt.Errorf("could not chmod file %s (%w)", tmpName, err)
	}

	if err := os.Rename(tmpName, destPath); err != nil {
		return "", fmt.Errorf("could not move %s to %s (%w)", tmpName, destPath, err)
	}

	done = true
	return destPath, nil
}

func downloadHTTP(dst io.Writer, sourceURL string) error {
	req, err := http.NewRequest("GET", sourceURL, nil)
	if err != nil {
		return fmt.Errorf("could not create request (%w)", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP GET %s failed (%w)", sourceURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP GET %s failed with error %d", sourceURL, resp.StatusCode)
	}

	if _, err = io.Copy(dst, resp.Body); err != nil {
		return fmt.Errorf("could not copy from %s (%w)", sourceURL, err)
	}

	return nil
}

func downloadFile(dst io.Writer, source string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("could not read file %s (%w)", source, err)
	}

	if _, err = dst.Write(data); err != nil {
		return fmt.Errorf("could not write to %s (%w)", source, err)
	}

	return nil
}

func parseOCIURL(sourceURL string) (repo string, ref string, _ error) {
	u, err := url.Parse(sourceURL)
	if err != nil {
		return "", "", fmt.Errorf("could not parse URL %s (%w)", sourceURL, err)
	}

	if u.Scheme != "oci" {
		return "", "", fmt.Errorf("unsupported scheme %s", u.Scheme)
	}

	if u.Host == "" {
		return "", "", fmt.Errorf("host is required")
	}

	if u.Path == "" {
		return "", "", fmt.Errorf("path is required")
	}

	if !strings.HasPrefix(u.Path, "/") {
		return "", "", fmt.Errorf("path must start with a slash")
	}

	repo, ref, ok := strings.Cut(u.Path, "@")
	if ok {
		return u.Host + repo, ref, nil
	}

	repo, ref, ok = strings.Cut(u.Path, ":")
	if ok {
		return u.Host + repo, ref, nil
	}

	return u.Host + u.Path, "latest", nil
}

func downloadOCI(ctx context.Context, workdir string, dst io.Writer, sourceURL string, insecure bool, targetOS, targetArch string) error {
	r, ref, err := parseOCIURL(sourceURL)
	if err != nil {
		return fmt.Errorf("could not parse OCI URL %s (%w)", sourceURL, err)
	}

	repo, err := remote.NewRepository(r)
	if err != nil {
		return fmt.Errorf("could not create repository (%w)", err)
	}

	repo.PlainHTTP = insecure

	d, err := os.MkdirTemp(workdir, "oci-pull")
	if err != nil {
		return fmt.Errorf("could not create temporary directory (%w)", err)
	}
	defer os.RemoveAll(d)

	store, err := file.New(d)
	if err != nil {
		return fmt.Errorf("could not create file system (%w)", err)
	}
	defer store.Close()

	store.AllowPathTraversalOnWrite = false
	store.DisableOverwrite = true

	desc, err := oras.Copy(ctx, repo, ref, store, "", oras.DefaultCopyOptions)
	if err != nil {
		return fmt.Errorf("could not copy OCI image (%w)", err)
	}

	return copyOCIArtifact(ctx, dst, desc, store, targetOS, targetArch)
}

func copyOCIArtifact(ctx context.Context, dst io.Writer, desc ocispec.Descriptor, fetcher content.Fetcher, targetOS, targetArch string) error {
	switch desc.MediaType {
	case ocispec.MediaTypeImageIndex:
		r, err := fetcher.Fetch(ctx, desc)
		if err != nil {
			return fmt.Errorf("could not fetch OCI image index (%w)", err)
		}
		defer r.Close()

		var idx ocispec.Index
		if err := json.NewDecoder(r).Decode(&idx); err != nil {
			return fmt.Errorf("could not decode OCI image index (%w)", err)
		}

		for _, m := range idx.Manifests {
			if targetOS == m.Platform.OS && targetArch == m.Platform.Architecture {
				return copyOCIArtifact(ctx, dst, m, fetcher, targetOS, targetArch)
			}
		}

		return fmt.Errorf("no matching manifest found")

	case ocispec.MediaTypeImageManifest:
		r, err := fetcher.Fetch(ctx, desc)
		if err != nil {
			return fmt.Errorf("could not fetch OCI image manifest (%w)", err)
		}
		defer r.Close()

		var manifest ocispec.Manifest
		if err := json.NewDecoder(r).Decode(&manifest); err != nil {
			return fmt.Errorf("could not decode OCI image manifest (%w)", err)
		}

		if len(manifest.Layers) != 1 {
			return fmt.Errorf("expected exactly one layer, got %d", len(manifest.Layers))
		}

		layer := manifest.Layers[0]
		r, err = fetcher.Fetch(ctx, layer)
		if err != nil {
			return fmt.Errorf("could not fetch OCI layer (%w)", err)
		}
		defer r.Close()

		if _, err := io.Copy(dst, r); err != nil {
			return fmt.Errorf("could not copy OCI layer (%w)", err)
		}

		return nil

	default:
		return fmt.Errorf("unsupported media type %s", desc.MediaType)
	}
}
