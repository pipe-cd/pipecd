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
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes/providertest"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
)

type fakeLogPersister struct{}

func (l *fakeLogPersister) Write(_ []byte) (int, error) {
	return 0, nil
}
func (l *fakeLogPersister) Info(_ string)                       {}
func (l *fakeLogPersister) Infof(_ string, _ ...interface{})    {}
func (l *fakeLogPersister) Success(_ string)                    {}
func (l *fakeLogPersister) Successf(_ string, _ ...interface{}) {}
func (l *fakeLogPersister) Error(_ string)                      {}
func (l *fakeLogPersister) Errorf(_ string, _ ...interface{})   {}

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

func TestCheckVariantSelectorInWorkload(t *testing.T) {
	testcases := []struct {
		name     string
		manifest string
		expected error
	}{
		{
			name: "missing variant in selector",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			expected: fmt.Errorf("missing pipecd.dev/variant key in spec.selector.matchLabels"),
		},
		{
			name: "missing variant in template labels",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
`,
			expected: fmt.Errorf("missing pipecd.dev/variant key in spec.template.metadata.labels"),
		},
		{
			name: "wrong variant in selector",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: canary
  template:
    metadata:
      labels:
        app: simple
`,
			expected: fmt.Errorf("require primary but got canary for pipecd.dev/variant key in spec.selector.matchLabels"),
		},
		{
			name: "wrong variant in temlate labels",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: canary
`,
			expected: fmt.Errorf("require primary but got canary for pipecd.dev/variant key in spec.template.metadata.labels"),
		},
		{
			name: "ok",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: primary
`,
			expected: nil,
		},
	}

	expected := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: primary
`
	generatedManifests, err := provider.ParseManifests(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(generatedManifests))

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.ParseManifests(tc.manifest)
			require.NoError(t, err)
			require.Equal(t, 1, len(manifests))

			err = checkVariantSelectorInWorkload(manifests[0], primaryVariant)
			assert.Equal(t, tc.expected, err)

			err = ensureVariantSelectorInWorkload(manifests[0], primaryVariant)
			assert.NoError(t, err)
			assert.Equal(t, generatedManifests[0], manifests[0])
		})
	}

}

func TestDeleteResources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name      string
		executor  *Executor
		resources []provider.ResourceKey
		wantErr   bool
	}{
		{
			name:      "no resource to delete",
			wantErr:   false,
			resources: []provider.ResourceKey{},
			executor: &Executor{
				Input: executor.Input{
					LogPersister: &fakeLogPersister{},
					Logger:       zap.NewNop(),
				},
			},
		},
		{
			name:    "not found resource to delete",
			wantErr: false,
			resources: []provider.ResourceKey{
				{
					Name: "foo",
				},
			},
			executor: &Executor{
				Input: executor.Input{
					LogPersister: &fakeLogPersister{},
					Logger:       zap.NewNop(),
				},
				provider: func() provider.Provider {
					p := providertest.NewMockProvider(ctrl)
					p.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(provider.ErrNotFound)
					return p
				}(),
			},
		},
		{
			name:    "unable to delete",
			wantErr: true,
			resources: []provider.ResourceKey{
				{
					Name: "foo",
				},
			},
			executor: &Executor{
				Input: executor.Input{
					LogPersister: &fakeLogPersister{},
					Logger:       zap.NewNop(),
				},
				provider: func() provider.Provider {
					p := providertest.NewMockProvider(ctrl)
					p.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("unexpected error"))
					return p
				}(),
			},
		},
		{
			name:    "successfully deletion",
			wantErr: false,
			resources: []provider.ResourceKey{
				{
					Name: "foo",
				},
			},
			executor: &Executor{
				Input: executor.Input{
					LogPersister: &fakeLogPersister{},
					Logger:       zap.NewNop(),
				},
				provider: func() provider.Provider {
					p := providertest.NewMockProvider(ctrl)
					p.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
					return p
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := tt.executor.deleteResources(ctx, tt.resources)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
