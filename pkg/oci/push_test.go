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

func TestPushFilesToRegistry(t *testing.T) {
	workDir := t.TempDir()

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

		if _, err := f.WriteString(fmt.Sprintf("test %s %s", platform.OS, platform.Arch)); err != nil {
			t.Fatalf("could not write to temporary file: %s", err)
		}

		artifactFiles[platform] = f.Name()
		if err := f.Close(); err != nil {
			t.Fatalf("could not close temporary file: %s", err)
		}
	}

	const ociURL = "oci://localhost:5001/test"

	artifact := &Artifact{
		MediaType: "text/plain",
		FilePaths: artifactFiles,
	}

	if err := PushFilesToRegistry(t.Context(), workDir, artifact, ociURL, true); err != nil {
		t.Fatalf("could not push files to OCI: %s", err)
	}
}

func TestPullFileFromRegistry(t *testing.T) {
	t.Parallel()

	TestPushFilesToRegistry(t)

	const ociURL = "oci://localhost:5001/test"

	testcases := []struct {
		platform Platform
		content  string
	}{
		{Platform{OS: "linux", Arch: "amd64"}, "test linux amd64"},
		{Platform{OS: "linux", Arch: "arm64"}, "test linux arm64"},
		{Platform{OS: "darwin", Arch: "amd64"}, "test darwin amd64"},
		{Platform{OS: "darwin", Arch: "arm64"}, "test darwin arm64"},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("platform=%s", tc.platform), func(t *testing.T) {
			t.Parallel()

			workDir := t.TempDir()
			dst, err := os.CreateTemp(workDir, "test.txt")
			if err != nil {
				t.Fatalf("could not create temporary file: %s", err)
			}
			defer os.Remove(dst.Name())

			if err := PullFileFromRegistry(t.Context(), workDir, dst, ociURL, true, tc.platform.OS, tc.platform.Arch, "text/plain"); err != nil {
				t.Fatalf("could not pull file from OCI: %s", err)
			}

			content, err := os.ReadFile(dst.Name())
			if err != nil {
				t.Fatalf("could not read file: %s", err)
			}

			if string(content) != tc.content {
				t.Fatalf("file content is not expected: %s", string(content))
			}
		})
	}
}
