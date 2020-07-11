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

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
)

func TestGenerateVirtualServiceManifest(t *testing.T) {
	testcases := []struct {
		name           string
		manifestFile   string
		editableRoutes []string
		expectedFile   string
	}{
		{
			name:         "generated correct manifest",
			manifestFile: "testdata/virtual-service.yaml",
			expectedFile: "testdata/generated-virtual-service.yaml",
		},
		{
			name:           "generated correct manifest",
			manifestFile:   "testdata/virtual-service.yaml",
			editableRoutes: []string{"only-primary-destination"},
			expectedFile:   "testdata/generated-virtual-service-for-editable-routes.yaml",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.LoadManifestsFromYAMLFile(tc.manifestFile)
			require.NoError(t, err)
			require.Equal(t, 1, len(manifests))

			err = generateVirtualServiceManifest(manifests[0], "helloworld", tc.editableRoutes, 30, 20)
			assert.NoError(t, err)

			expectedManifests, err := provider.LoadManifestsFromYAMLFile(tc.expectedFile)
			require.NoError(t, err)
			require.Equal(t, 1, len(expectedManifests))

			expected, err := expectedManifests[0].YamlBytes()
			require.NoError(t, err)
			got, err := manifests[0].YamlBytes()
			require.NoError(t, err)

			assert.EqualValues(t, expected, got)
		})
	}
}
