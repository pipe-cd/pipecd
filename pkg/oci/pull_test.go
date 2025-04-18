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

func TestPullFileFromRegistry(t *testing.T) {
	t.Parallel()

	// TODO: Use dockertest or something to start a local registry for testing.
	const ociURL = "oci://localhost:5001/test"

	testcases := pushTestFiles(t, t.TempDir(), ociURL)

	for platform, content := range testcases {
		t.Run(fmt.Sprintf("platform=%s", platform), func(t *testing.T) {
			t.Parallel()

			workDir := t.TempDir()
			dst, err := os.CreateTemp(workDir, "test.txt")
			if err != nil {
				t.Fatalf("could not create temporary file: %s", err)
			}
			defer os.Remove(dst.Name())

			if err := PullFileFromRegistry(t.Context(), workDir, dst, ociURL, true, platform.OS, platform.Arch, "text/plain"); err != nil {
				t.Fatalf("could not pull file from OCI: %s", err)
			}

			got, err := os.ReadFile(dst.Name())
			if err != nil {
				t.Fatalf("could not read file: %s", err)
			}

			if string(got) != content {
				t.Fatalf("file content is not expected: %s", string(got))
			}
		})
	}
}
