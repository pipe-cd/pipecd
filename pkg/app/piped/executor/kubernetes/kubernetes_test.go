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

func TestGenerateServiceManifests(t *testing.T) {
	testcases := []struct {
		name          string
		manifestsFile string
	}{
		{
			name:          "Update selector and change type to ClusterIP",
			manifestsFile: "testdata/services.yaml",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.LoadManifestsFromYAMLFile(tc.manifestsFile)
			require.NoError(t, err)
			require.Equal(t, 2, len(manifests))

			generatedManifests, err := generateVariantServiceManifests(manifests[:1], "canary-variant", "canary")
			require.NoError(t, err)
			require.Equal(t, 1, len(generatedManifests))

			assert.Equal(t, manifests[1], generatedManifests[0])
		})
	}
}

func TestGenerateWorkloadManifests(t *testing.T) {
	testcases := []struct {
		name           string
		manifestsFile  string
		configmapsFile string
		secretsFile    string
	}{
		{
			name:          "No configmap and secret",
			manifestsFile: "testdata/no-config-deployments.yaml",
		},
		{
			name:           "Has configmap and secret",
			manifestsFile:  "testdata/deployments.yaml",
			configmapsFile: "testdata/configmaps.yaml",
			secretsFile:    "testdata/secrets.yaml",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.LoadManifestsFromYAMLFile(tc.manifestsFile)
			require.NoError(t, err)
			require.Equal(t, 2, len(manifests))

			var configmaps, secrets []provider.Manifest
			if tc.configmapsFile != "" {
				configmaps, err = provider.LoadManifestsFromYAMLFile(tc.configmapsFile)
				require.NoError(t, err)
			}
			if tc.secretsFile != "" {
				secrets, err = provider.LoadManifestsFromYAMLFile(tc.secretsFile)
				require.NoError(t, err)
			}

			generatedManifests, err := generateVariantWorkloadManifests(manifests[:1], configmaps, secrets, "canary-variant", "canary", func(r *int32) int32 {
				return *r - 1
			})
			require.NoError(t, err)
			require.Equal(t, 1, len(generatedManifests))

			assert.Equal(t, manifests[1], generatedManifests[0])
		})
	}
}
