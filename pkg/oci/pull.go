// Copyright 2025 The PipeCD Authors.
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

package oci

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

// PullFileFromRegistry pulls a file from an OCI registry and writes it to the provided destination writer.
// It supports options for insecure connections and platform/media type targeting.
func PullFileFromRegistry(ctx context.Context, workdir string, dst io.Writer, sourceURL string, opts ...PullOption) error {
	options := &PullOptions{
		insecure: false,
	}
	for _, opt := range opts {
		opt.applyPullOption(options)
	}

	r, ref, err := parseOCIURL(sourceURL)
	if err != nil {
		return fmt.Errorf("could not parse OCI URL %s (%w)", sourceURL, err)
	}

	repo, err := remote.NewRepository(r)
	if err != nil {
		return fmt.Errorf("could not create repository (%w)", err)
	}

	repo.PlainHTTP = options.insecure

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

	return copyOCIArtifact(ctx, dst, desc, store, options.targetOS, options.targetArch, options.mediaType, options.artifactType)
}

// parseOCIURL parses an OCI URL and returns the repository and reference parts.
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
		// some OCI URLs are like oci://example.com/test:v1@sha256:1234567890
		repo, tag, _ := strings.Cut(repo, ":")
		return u.Host + repo, tag + "@" + ref, nil
	}

	repo, ref, ok = strings.Cut(u.Path, ":")
	if ok {
		return u.Host + repo, ref, nil
	}

	return u.Host + u.Path, "latest", nil
}

// copyOCIArtifact resolves the given OCI artifact and copies it to the destination writer.
func copyOCIArtifact(ctx context.Context, dst io.Writer, desc ocispec.Descriptor, fetcher content.Fetcher, targetOS, targetArch, mediaType, artifactType string) error {
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
			if targetOS != "" && targetOS != m.Platform.OS {
				continue
			}
			if targetArch != "" && targetArch != m.Platform.Architecture {
				continue
			}
			return copyOCIArtifact(ctx, dst, m, fetcher, targetOS, targetArch, mediaType, artifactType)
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

		if artifactType != "" && artifactType != manifest.ArtifactType {
			return fmt.Errorf("artifact type mismatch: %s != %s", manifest.ArtifactType, artifactType)
		}

		for _, layer := range manifest.Layers {
			if mediaType != "" && mediaType != layer.MediaType {
				continue
			}

			r, err = fetcher.Fetch(ctx, layer)
			if err != nil {
				return fmt.Errorf("could not fetch OCI layer (%w)", err)
			}
			defer r.Close()

			if _, err := io.Copy(dst, r); err != nil {
				return fmt.Errorf("could not copy OCI layer (%w)", err)
			}
		}

		return nil

	default:
		return fmt.Errorf("unsupported media type %s", desc.MediaType)
	}
}
