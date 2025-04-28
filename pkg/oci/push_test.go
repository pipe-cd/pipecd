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
	"fmt"
	"os"
	"testing"
)

func pushTestFiles(t *testing.T, workDir, ociURL string) map[Platform]string {
	artifactsDir, err := os.MkdirTemp(workDir, "artifacts")
	if err != nil {
		t.Fatalf("could not create temporary directory: %s", err)
	}
	defer os.RemoveAll(artifactsDir)

	artifactFiles := make(map[Platform]string)

	for _, platform := range []Platform{
		{OS: "linux", Arch: "amd64"},
		{OS: "linux", Arch: "arm64"},
		{OS: "darwin", Arch: "amd64"},
		{OS: "darwin", Arch: "arm64"},
	} {
		f, err := os.CreateTemp(artifactsDir, fmt.Sprintf("%s-%s.txt", platform.OS, platform.Arch))
		if err != nil {
			t.Fatalf("could not create temporary file: %s", err)
		}

		if _, err := fmt.Fprintf(f, "test %s %s", platform.OS, platform.Arch); err != nil {
			t.Fatalf("could not write to temporary file: %s", err)
		}

		artifactFiles[platform] = f.Name()
		if err := f.Close(); err != nil {
			t.Fatalf("could not close temporary file: %s", err)
		}
	}

	artifact := &Artifact{
		ArtifactType: "application/vnd.pipecd.test+type",
		MediaType:    "text/plain",
		FilePaths:    artifactFiles,
	}

	if err := PushFilesToRegistry(t.Context(), workDir, artifact, ociURL, WithInsecure()); err != nil {
		t.Fatalf("could not push files to OCI: %s", err)
	}

	results := make(map[Platform]string)
	for platform := range artifactFiles {
		results[platform] = fmt.Sprintf("test %s %s", platform.OS, platform.Arch)
	}
	return results
}

func TestPushFilesToRegistry(t *testing.T) {
	t.Parallel()

	// OCI_REGISTRY_HOST is set by TestMain in main_test.go
	ociURL := fmt.Sprintf("oci://%s/test-push", os.Getenv("OCI_REGISTRY_HOST"))

	workDir := t.TempDir()

	pushTestFiles(t, workDir, ociURL)
}
