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

package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	testcases := []struct {
		name     string
		yamlFile string
	}{
		{
			name:     "nodiff",
			yamlFile: "testdata/diff_nodiff.yaml",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := LoadManifestsFromYAMLFile(tc.yamlFile)
			require.NoError(t, err)
			require.Equal(t, 2, len(manifests))

			r := Diff(manifests[0], manifests[1])
			assert.Equal(t, r.Path, "hello")
		})
	}
}
