// Copyright 2023 The PipeCD Authors.
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
	t.Parallel()

	testcases := []struct {
		name         string
		appDir       string
		includes     []string
		excludes     []string
		changedFiles []string
		expected     bool
	}{
		{
			name:   "not touched",
			appDir: "app/demo",
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
			name:   "touched in app dir",
			appDir: "app/demo",
			changedFiles: []string{
				"app/hello.txt",
				"app/demo/deployment.yaml",
			},
			expected: true,
		},
		{
			name:   "touched in app dir while not in includes",
			appDir: "app/demo",
			includes: []string{
				"charts/demo",
				"charts/bar/*",
			},
			changedFiles: []string{
				"app/hello.txt",
				"app/demo/deployment.yaml",
			},
			expected: true,
		},
		{
			name:   "touched in app dir but listed in excluded",
			appDir: "app/demo",
			excludes: []string{
				"*/hello.txt",
			},
			changedFiles: []string{
				"app/hello.txt",
				"app/demo/deployment.yaml",
			},
			expected: false,
		},
		{
			name:   "touched in the includes",
			appDir: "app/demo",
			includes: []string{
				"charts/demo",
				"charts/bar/*",
			},
			changedFiles: []string{
				"app/hello.txt",
				"charts/bar/deployment.yaml",
			},
			expected: true,
		},
		{
			name:   "touched in the includes, not excludes",
			appDir: "app/demo",
			excludes: []string{
				"foo/*",
			},
			includes: []string{
				"charts/demo",
				"charts/bar/*",
			},
			changedFiles: []string{
				"app/hello.txt",
				"charts/bar/deployment.yaml",
			},
			expected: true,
		},
		{
			name:   "touched in the excludes",
			appDir: "app/demo",
			excludes: []string{
				"charts/demo",
				"charts/bar/*",
			},
			changedFiles: []string{
				"app/hello.txt",
				"charts/bar/deployment.yaml",
			},
			expected: false,
		},
		{
			name:   "touched in both the excludes and includes",
			appDir: "app/demo",
			excludes: []string{
				"charts/demo",
				"charts/bar/*",
			},
			includes: []string{
				"charts/bar/deployment.yaml",
			},
			changedFiles: []string{
				"app/hello.txt",
				"charts/bar/deployment.yaml",
			},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := isTouchedByChangedFiles(tc.appDir, tc.includes, tc.excludes, tc.changedFiles)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}
