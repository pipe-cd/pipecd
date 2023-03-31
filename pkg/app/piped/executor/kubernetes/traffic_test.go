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

package kubernetes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/config"
)

func TestGenerateVirtualServiceManifest(t *testing.T) {
	t.Parallel()

	exec := &deployExecutor{
		appCfg: &config.KubernetesApplicationSpec{
			VariantLabel: config.KubernetesVariantLabel{
				Key:           "pipecd.dev/variant",
				PrimaryValue:  "primary",
				BaselineValue: "baseline",
				CanaryValue:   "canary",
			},
		},
	}
	testcases := []struct {
		name           string
		manifestFile   string
		editableRoutes []string
		expectedFile   string
	}{
		{
			name:         "apply all routes",
			manifestFile: "testdata/virtual-service.yaml",
			expectedFile: "testdata/generated-virtual-service.yaml",
		},
		{
			name:           "apply only speficied routes",
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

			generatedManifest, err := exec.generateVirtualServiceManifest(manifests[0], "helloworld", tc.editableRoutes, 30, 20)
			assert.NoError(t, err)

			expectedManifests, err := provider.LoadManifestsFromYAMLFile(tc.expectedFile)
			require.NoError(t, err)
			require.Equal(t, 1, len(expectedManifests))

			expected, err := expectedManifests[0].YamlBytes()
			require.NoError(t, err)
			got, err := generatedManifest.YamlBytes()
			require.NoError(t, err)

			assert.EqualValues(t, string(expected), string(got))
		})
	}
}

func TestCheckVariantSelectorInService(t *testing.T) {
	t.Parallel()

	const (
		variantLabel   = "pipecd.dev/variant"
		primaryVariant = "primary"
	)
	testcases := []struct {
		name     string
		manifest string
		expected error
	}{
		{
			name: "missing variant selector",
			manifest: `
apiVersion: v1
kind: Service
metadata:
    name: simple
spec:
    selector:
        app: simple
`,
			expected: fmt.Errorf("missing pipecd.dev/variant key in spec.selector"),
		},
		{
			name: "wrong variant",
			manifest: `
apiVersion: v1
kind: Service
metadata:
    name: simple
spec:
    selector:
        app: simple
        pipecd.dev/variant: canary
`,
			expected: fmt.Errorf("require primary but got canary for pipecd.dev/variant key in spec.selector"),
		},
		{
			name: "ok",
			manifest: `
apiVersion: v1
kind: Service
metadata:
    name: simple
spec:
    selector:
        app: simple
        pipecd.dev/variant: primary
`,
			expected: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.ParseManifestsFromGit(tc.manifest)
			require.NoError(t, err)
			require.Equal(t, 1, len(manifests))

			err = checkVariantSelectorInService(manifests[0], variantLabel, primaryVariant)
			assert.Equal(t, tc.expected, err)
		})
	}
}
