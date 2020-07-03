// Copyright 2020 The PipeCD Authors.
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

package trigger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTouchedByChangedFiles(t *testing.T) {
	testcases := []struct {
		name           string
		appDir         string
		dependencyDirs []string
		changedFiles   []string
		expected       bool
	}{
		{
			name:           "not touched",
			appDir:         "app/demo",
			dependencyDirs: nil,
			changedFiles: []string{
				"app/hello.txt",
				"app/foo/deployment.yaml",
			},
			expected: false,
		},
		{
			name:   "not touched in dir whose name does not match exactly",
			appDir: "app/demo",
			changedFiles: []string{
				"app/demo-2",
			},
			expected: false,
		},
		{
			name:           "touched in app dir",
			appDir:         "app/demo",
			dependencyDirs: nil,
			changedFiles: []string{
				"app/hello.txt",
				"app/demo/deployment.yaml",
			},
			expected: true,
		},
		{
			name:   "touched in dependency dir",
			appDir: "app/demo",
			dependencyDirs: []string{
				"charts/demo",
				"charts/bar",
			},
			changedFiles: []string{
				"app/hello.txt",
				"charts/bar/deployment.yaml",
			},
			expected: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := isTouchedByChangedFiles(tc.appDir, tc.dependencyDirs, tc.changedFiles)
			assert.Equal(t, tc.expected, got)
		})
	}
}
