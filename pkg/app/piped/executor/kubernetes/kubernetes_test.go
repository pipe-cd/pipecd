// Copyright 2024 The PipeCD Authors.
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
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/pipe-cd/pipecd/pkg/app/piped/metadatastore"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes/kubernetestest"
	"github.com/pipe-cd/pipecd/pkg/config"
)

type fakeLogPersister struct{}

func (l *fakeLogPersister) Write(_ []byte) (int, error)         { return 0, nil }
func (l *fakeLogPersister) Info(_ string)                       {}
func (l *fakeLogPersister) Infof(_ string, _ ...interface{})    {}
func (l *fakeLogPersister) Success(_ string)                    {}
func (l *fakeLogPersister) Successf(_ string, _ ...interface{}) {}
func (l *fakeLogPersister) Error(_ string)                      {}
func (l *fakeLogPersister) Errorf(_ string, _ ...interface{})   {}

type fakeMetadataStore struct{}

func (m *fakeMetadataStore) Shared() metadatastore.Store {
	return &fakeMetadataSharedStore{}
}

func (m *fakeMetadataStore) Stage(stageID string) metadatastore.Store {
	return &fakeMetadataStageStore{}
}

type fakeMetadataSharedStore struct{}

func (m *fakeMetadataSharedStore) Get(_ string) (string, bool)                           { return "", false }
func (m *fakeMetadataSharedStore) Put(_ context.Context, _, _ string) error              { return nil }
func (m *fakeMetadataSharedStore) PutMulti(_ context.Context, _ map[string]string) error { return nil }

type fakeMetadataStageStore struct{}

func (m *fakeMetadataStageStore) Get(_ string) (string, bool)                           { return "", false }
func (m *fakeMetadataStageStore) Put(_ context.Context, _, _ string) error              { return nil }
func (m *fakeMetadataStageStore) PutMulti(_ context.Context, _ map[string]string) error { return nil }

func TestGenerateServiceManifests(t *testing.T) {
	t.Parallel()

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

			generatedManifests, err := generateVariantServiceManifests(manifests[:1], "pipecd.dev/variant", "canary-variant", "canary")
			require.NoError(t, err)
			require.Equal(t, 1, len(generatedManifests))

			assert.Equal(t, manifests[1], generatedManifests[0])
		})
	}
}

func TestGenerateVariantWorkloadManifests(t *testing.T) {
	t.Parallel()

	const (
		variantLabel  = "pipecd.dev/variant"
		canaryVariant = "canary-variant"
	)
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

			calculator := func(r *int32) int32 {
				return *r - 1
			}
			generatedManifests, err := generateVariantWorkloadManifests(
				manifests[:1],
				configmaps,
				secrets,
				variantLabel,
				canaryVariant,
				"canary",
				calculator,
			)
			require.NoError(t, err)
			require.Equal(t, 1, len(generatedManifests))

			assert.Equal(t, manifests[1], generatedManifests[0])
		})
	}
}

func TestCheckVariantSelectorInWorkload(t *testing.T) {
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

			err = checkVariantSelectorInWorkload(manifests[0], variantLabel, primaryVariant)
			assert.Equal(t, tc.expected, err)

			err = ensureVariantSelectorInWorkload(manifests[0], variantLabel, primaryVariant)
			assert.NoError(t, err)
			assert.Equal(t, generatedManifests[0], manifests[0])
		})
	}

}

func TestApplyManifests(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	testcases := []struct {
		name      string
		applier   provider.Applier
		manifest  string
		namespace string
		wantErr   bool
	}{

		{
			name: "unable to apply manifest",
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().ApplyManifest(gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))
				return p
			}(),
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
			namespace: "",
			wantErr:   true,
		},
		{
			name: "unable to replace manifest",
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().ReplaceManifest(gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))
				return p
			}(),
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  annotations:
    pipecd.dev/sync-by-replace: "enabled"
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			namespace: "",
			wantErr:   true,
		},
		{
			name: "unable to force replace manifest",
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().ForceReplaceManifest(gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))
				return p
			}(),
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  annotations:
    pipecd.dev/force-sync-by-replace: "enabled"
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			namespace: "",
			wantErr:   true,
		},
		{
			name: "unable to create manifest",
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().ReplaceManifest(gomock.Any(), gomock.Any()).Return(provider.ErrNotFound)
				p.EXPECT().CreateManifest(gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))
				return p
			}(),
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  annotations:
    pipecd.dev/sync-by-replace: "enabled"
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			namespace: "",
			wantErr:   true,
		},
		{
			name: "unable to create manifest",
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().ForceReplaceManifest(gomock.Any(), gomock.Any()).Return(provider.ErrNotFound)
				p.EXPECT().CreateManifest(gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))
				return p
			}(),
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  annotations:
    pipecd.dev/force-sync-by-replace: "enabled"
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			namespace: "",
			wantErr:   true,
		},
		{
			name: "successfully apply manifest",
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().ApplyManifest(gomock.Any(), gomock.Any()).Return(nil)
				return p
			}(),
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
			namespace: "",
			wantErr:   false,
		},
		{
			name: "successfully replace manifest",
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().ReplaceManifest(gomock.Any(), gomock.Any()).Return(nil)
				return p
			}(),
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  annotations:
    pipecd.dev/sync-by-replace: "enabled"
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			namespace: "",
			wantErr:   false,
		},
		{
			name: "successfully force replace manifest",
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().ForceReplaceManifest(gomock.Any(), gomock.Any()).Return(nil)
				return p
			}(),
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  annotations:
    pipecd.dev/force-sync-by-replace: "enabled"
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			namespace: "",
			wantErr:   false,
		},
		{
			name: "successfully create manifest",
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().ReplaceManifest(gomock.Any(), gomock.Any()).Return(provider.ErrNotFound)
				p.EXPECT().CreateManifest(gomock.Any(), gomock.Any()).Return(nil)
				return p
			}(),
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  annotations:
    pipecd.dev/sync-by-replace: "enabled"
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			namespace: "",
			wantErr:   false,
		},
		{
			name: "successfully force create manifest",
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().ForceReplaceManifest(gomock.Any(), gomock.Any()).Return(provider.ErrNotFound)
				p.EXPECT().CreateManifest(gomock.Any(), gomock.Any()).Return(nil)
				return p
			}(),
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  annotations:
    pipecd.dev/force-sync-by-replace: "enabled"
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			namespace: "",
			wantErr:   false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			manifests, err := provider.ParseManifests(tc.manifest)
			require.NoError(t, err)
			ag := &applierGroup{defaultApplier: tc.applier}
			err = applyManifests(ctx, ag, manifests, tc.namespace, &fakeLogPersister{})
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestDeleteResources(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	testcases := []struct {
		name      string
		applier   provider.Applier
		resources []provider.ResourceKey
		wantErr   bool
	}{
		{
			name:      "no resource to delete",
			wantErr:   false,
			resources: []provider.ResourceKey{},
		},
		{
			name:    "not found resource to delete",
			wantErr: false,
			resources: []provider.ResourceKey{
				{
					Name: "foo",
				},
			},
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(provider.ErrNotFound)
				return p
			}(),
		},
		{
			name:    "unable to delete",
			wantErr: true,
			resources: []provider.ResourceKey{
				{
					Name: "foo",
				},
			},
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("unexpected error"))
				return p
			}(),
		},
		{
			name:    "successfully deletion",
			wantErr: false,
			resources: []provider.ResourceKey{
				{
					Name: "foo",
				},
			},
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				return p
			}(),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ag := &applierGroup{defaultApplier: tc.applier}
			err := deleteResources(ctx, ag, tc.resources, &fakeLogPersister{})
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestAnnotateConfigHash(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		manifests     string
		expected      string
		expectedError error
	}{
		{
			name: "empty list",
		},
		{
			name: "one config",
			manifests: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-by-config-change
  labels:
    app: canary-by-config-change
spec:
  replicas: 2
  selector:
    matchLabels:
      app: canary-by-config-change
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: canary-by-config-change
        pipecd.dev/variant: primary
    spec:
      containers:
        - name: helloworld
          image: gcr.io/pipecd/helloworld:v0.5.0
          args:
            - server
          ports:
            - containerPort: 9085
          volumeMounts:
            - name: config
              mountPath: /etc/pipecd-config
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: canary-by-config-change
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: canary-by-config-change
data:
  two: "2"
`,
			expected: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-by-config-change
  labels:
    app: canary-by-config-change
spec:
  replicas: 2
  selector:
    matchLabels:
      app: canary-by-config-change
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: canary-by-config-change
        pipecd.dev/variant: primary
      annotations:
        pipecd.dev/config-hash: 75c9m2btb6
    spec:
      containers:
        - name: helloworld
          image: gcr.io/pipecd/helloworld:v0.5.0
          args:
            - server
          ports:
            - containerPort: 9085
          volumeMounts:
            - name: config
              mountPath: /etc/pipecd-config
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: canary-by-config-change
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: canary-by-config-change
data:
  two: "2"
`,
		},
		{
			name: "multiple configs",
			manifests: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-by-config-change
  labels:
    app: canary-by-config-change
spec:
  replicas: 2
  selector:
    matchLabels:
      app: canary-by-config-change
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: canary-by-config-change
        pipecd.dev/variant: primary
    spec:
      containers:
        - name: helloworld
          image: gcr.io/pipecd/helloworld:v0.5.0
          args:
            - server
          ports:
            - containerPort: 9085
          volumeMounts:
            - name: config
              mountPath: /etc/pipecd-config
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: canary-by-config-change
        - name: secret
          secret:
            secretName: secret-1
        - name: unmanaged-config
          configMap:
            name: unmanaged-config
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: canary-by-config-change
data:
  two: "2"
---
apiVersion: v1
kind: Secret
metadata:
  name: secret-1
type: my-type
data:
  "one": "Mg=="
`,
			expected: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-by-config-change
  labels:
    app: canary-by-config-change
spec:
  replicas: 2
  selector:
    matchLabels:
      app: canary-by-config-change
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: canary-by-config-change
        pipecd.dev/variant: primary
      annotations:
        pipecd.dev/config-hash: t7dtkdm455
    spec:
      containers:
        - name: helloworld
          image: gcr.io/pipecd/helloworld:v0.5.0
          args:
            - server
          ports:
            - containerPort: 9085
          volumeMounts:
            - name: config
              mountPath: /etc/pipecd-config
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: canary-by-config-change
        - name: secret
          secret:
            secretName: secret-1
        - name: unmanaged-config
          configMap:
            name: unmanaged-config
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: canary-by-config-change
data:
  two: "2"
---
apiVersion: v1
kind: Secret
metadata:
  name: secret-1
type: my-type
data:
  "one": "Mg=="
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.ParseManifests(tc.manifests)
			require.NoError(t, err)

			expected, err := provider.ParseManifests(tc.expected)
			require.NoError(t, err)

			err = annotateConfigHash(manifests)
			assert.Equal(t, expected, manifests)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestPatchManifest(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		manifests     string
		patch         config.K8sResourcePatch
		expectedError error
	}{
		{
			name:      "one op",
			manifests: "testdata/patch_configmap.yaml",
			patch: config.K8sResourcePatch{
				Ops: []config.K8sResourcePatchOp{
					{
						Op:    config.K8sResourcePatchOpYAMLReplace,
						Path:  "$.data.key1",
						Value: "value-1",
					},
				},
			},
		},
		{
			name:      "multi ops",
			manifests: "testdata/patch_configmap_multi_ops.yaml",
			patch: config.K8sResourcePatch{
				Ops: []config.K8sResourcePatchOp{
					{
						Op:    config.K8sResourcePatchOpYAMLReplace,
						Path:  "$.data.key1",
						Value: "value-1",
					},
					{
						Op:    config.K8sResourcePatchOpYAMLReplace,
						Path:  "$.data.key2",
						Value: "value-2",
					},
				},
			},
		},
		{
			name:      "one op with a given field",
			manifests: "testdata/patch_configmap_field.yaml",
			patch: config.K8sResourcePatch{
				Target: config.K8sResourcePatchTarget{
					DocumentRoot: "$.data.envoy-config",
				},
				Ops: []config.K8sResourcePatchOp{
					{
						Op:    config.K8sResourcePatchOpYAMLReplace,
						Path:  "$.admin.address.socket_address.port_value",
						Value: "9096",
					},
				},
			},
		},
		{
			name:      "multi ops with a given field",
			manifests: "testdata/patch_configmap_field_multi_ops.yaml",
			patch: config.K8sResourcePatch{
				Target: config.K8sResourcePatchTarget{
					DocumentRoot: "$.data.envoy-config",
				},
				Ops: []config.K8sResourcePatchOp{
					{
						Op:    config.K8sResourcePatchOpYAMLReplace,
						Path:  "$.admin.address.socket_address.port_value",
						Value: "19095",
					},
					{
						Op:    config.K8sResourcePatchOpYAMLReplace,
						Path:  "$.static_resources.clusters[1].load_assignment.endpoints[0].lb_endpoints[0].endpoint.address.socket_address.port_value",
						Value: "19081",
					},
					{
						Op:    config.K8sResourcePatchOpYAMLReplace,
						Path:  "$.static_resources.clusters[1].type",
						Value: "DNS",
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.LoadManifestsFromYAMLFile(tc.manifests)
			require.NoError(t, err)

			if tc.expectedError == nil {
				require.Equal(t, 2, len(manifests))
			} else {
				require.Equal(t, 1, len(manifests))
			}

			got, err := patchManifest(manifests[0], tc.patch)
			require.Equal(t, tc.expectedError, err)

			expectedBytes, err := manifests[1].YamlBytes()
			require.NoError(t, err)

			gotBytes, err := got.YamlBytes()
			require.NoError(t, err)

			if tc.expectedError == nil {
				assert.Equal(t, string(expectedBytes), string(gotBytes))
			}
		})
	}
}

func TestPatchManifests(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       []provider.Manifest
		patches     []config.K8sResourcePatch
		expected    []provider.Manifest
		expectedErr error
	}{
		{
			name: "no patches",
			input: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						Kind: "Deployment",
						Name: "deployment-1",
					},
				},
			},
			expected: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						Kind: "Deployment",
						Name: "deployment-1",
					},
				},
			},
		},
		{
			name: "no manifest for the given patch",
			input: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						Kind: "Deployment",
						Name: "deployment-1",
					},
				},
			},
			patches: []config.K8sResourcePatch{
				{
					Target: config.K8sResourcePatchTarget{
						K8sResourceReference: config.K8sResourceReference{
							Kind: "Deployment",
							Name: "deployment-2",
						},
					},
				},
			},
			expectedErr: errors.New("no manifest matches the given patch: kind=Deployment, name=deployment-2"),
		},
		{
			name: "multiple patches",
			input: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						Kind: "Deployment",
						Name: "deployment-1",
					},
				},
				{
					Key: provider.ResourceKey{
						Kind: "Deployment",
						Name: "deployment-2",
					},
				},
				{
					Key: provider.ResourceKey{
						Kind: "ConfigMap",
						Name: "configmap-1",
					},
				},
			},
			patches: []config.K8sResourcePatch{
				{
					Target: config.K8sResourcePatchTarget{
						K8sResourceReference: config.K8sResourceReference{
							Kind: "ConfigMap",
							Name: "configmap-1",
						},
					},
				},
				{
					Target: config.K8sResourcePatchTarget{
						K8sResourceReference: config.K8sResourceReference{
							Kind: "Deployment",
							Name: "deployment-1",
						},
					},
				},
				{
					Target: config.K8sResourcePatchTarget{
						K8sResourceReference: config.K8sResourceReference{
							Kind: "ConfigMap",
							Name: "configmap-1",
						},
					},
				},
			},
			expected: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						Kind:      "Deployment",
						Name:      "deployment-1",
						Namespace: "+",
					},
				},
				{
					Key: provider.ResourceKey{
						Kind: "Deployment",
						Name: "deployment-2",
					},
				},
				{
					Key: provider.ResourceKey{
						Kind:      "ConfigMap",
						Name:      "configmap-1",
						Namespace: "++",
					},
				},
			},
		},
	}

	patcher := func(m provider.Manifest, cfg config.K8sResourcePatch) (*provider.Manifest, error) {
		out := m
		out.Key.Namespace = fmt.Sprintf("%s+", out.Key.Namespace)
		return &out, nil
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := patchManifests(tc.input, tc.patches, patcher)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}
