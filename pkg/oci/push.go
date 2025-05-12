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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

// Platform represents an OS/Arch platform for an OCI artifact.
type Platform struct {
	OS   string
	Arch string
}

// Artifact describes an artifact to be pushed to an OCI registry.
type Artifact struct {
	// MediaType is the media type of the artifact.
	// e.g. "application/vnd.pipecd.piped.plugin"
	MediaType string
	// ArtifactType is the type of the artifact.
	// e.g. "application/vnd.pipecd.piped.plugin+type"
	ArtifactType string
	// FilePaths maps platforms to file paths.
	FilePaths map[Platform]string
}

// PushFilesToRegistry pushes files described by the artifact to the target OCI registry URL.
// It supports options for insecure connections.
func PushFilesToRegistry(ctx context.Context, workDir string, artifact *Artifact, targetURL string, opts ...PushOption) error {
	options := &PushOptions{
		insecure: false,
	}
	for _, opt := range opts {
		opt.applyPushOption(options)
	}

	repo, ref, err := parseOCIURL(targetURL)
	if err != nil {
		return fmt.Errorf("could not parse OCI URL %s: %w", targetURL, err)
	}

	r, err := remote.NewRepository(repo)
	if err != nil {
		return fmt.Errorf("could not create repository %s: %w", repo, err)
	}

	r.PlainHTTP = options.insecure

	if options.username != "" || options.password != "" {
		r.Client = &auth.Client{
			Client: retry.DefaultClient,
			Header: http.Header{
				"User-Agent": {"oras-go"},
			},
			Credential: func(_ context.Context, _ string) (auth.Credential, error) {
				return auth.Credential{
					Username: options.username,
					Password: options.password,
				}, nil
			},
		}
	}

	descriptors := make([]ocispec.Descriptor, 0, len(artifact.FilePaths))
	for platform, path := range artifact.FilePaths {
		d, err := pushFile(ctx, workDir, r, path, artifact.MediaType, artifact.ArtifactType, ref)
		if err != nil {
			return fmt.Errorf("could not push file %s: %w", path, err)
		}
		d.Platform = &ocispec.Platform{
			OS:           platform.OS,
			Architecture: platform.Arch,
		}
		descriptors = append(descriptors, d)
	}

	index := ocispec.Index{
		Versioned: specs.Versioned{
			SchemaVersion: 2,
		},
		MediaType: ocispec.MediaTypeImageIndex,
		Manifests: descriptors,
	}

	b, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("could not marshal index: %w", err)
	}

	desc := content.NewDescriptorFromBytes(ocispec.MediaTypeImageIndex, b)

	if err := r.Push(ctx, desc, bytes.NewReader(b)); err != nil {
		return fmt.Errorf("could not push index: %w", err)
	}
	if err := r.Tag(ctx, desc, ref); err != nil {
		return fmt.Errorf("could not tag index: %w", err)
	}

	return nil
}

// pushFile pushes a single file to the repository and returns its descriptor.
func pushFile(ctx context.Context, workDir string, repo *remote.Repository, path, mediaType, artifactType, ref string) (ocispec.Descriptor, error) {
	dir, err := os.MkdirTemp(workDir, "")
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not create temporary directory: %w", err)
	}
	defer os.RemoveAll(dir)

	fs, err := file.New(dir)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not create file system: %w", err)
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not get absolute path: %w", err)
	}

	desc, err := fs.Add(ctx, filepath.Base(path), mediaType, path)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not add file to file system: %w", err)
	}

	manifest, err := oras.PackManifest(ctx, fs, oras.PackManifestVersion1_1, artifactType, oras.PackManifestOptions{
		Layers: []ocispec.Descriptor{desc},
	})
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not pack manifest: %w", err)
	}

	if err = fs.Tag(ctx, manifest, ref); err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not tag manifest: %w", err)
	}

	d, err := oras.Copy(ctx, fs, ref, repo, ref, oras.DefaultCopyOptions)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not copy manifest to repository: %w", err)
	}

	return d, nil
}
