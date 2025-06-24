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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/backoff"
	"github.com/pipe-cd/pipecd/pkg/oci"
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
		// TODO: add context.Context as a argument for DownloadBinary.
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := oci.PullFileFromRegistry(
			ctx,
			destDir,
			tmpFile,
			sourceURL,
			oci.WithTargetOS(runtime.GOOS),
			oci.WithTargetArch(runtime.GOARCH),
			oci.WithMediaType(oci.MediaTypePipedPlugin),
		); err != nil {
			return "", fmt.Errorf("could not pull file from OCI (%w)", err)
		}
	case "http", "https":
		req, err := http.NewRequest("GET", sourceURL, nil)
		if err != nil {
			return "", fmt.Errorf("could not create request (%w)", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("HTTP GET %s failed (%w)", sourceURL, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("HTTP GET %s failed with error %d", sourceURL, resp.StatusCode)
		}

		if _, err = io.Copy(tmpFile, resp.Body); err != nil {
			return "", fmt.Errorf("could not copy from %s to %s (%w)", sourceURL, tmpName, err)
		}

	case "file":
		data, err := os.ReadFile(u.Path)
		if err != nil {
			return "", fmt.Errorf("could not read file %s (%w)", u.Path, err)
		}

		if _, err = tmpFile.Write(data); err != nil {
			return "", fmt.Errorf("could not write to %s (%w)", tmpName, err)
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
