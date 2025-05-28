package deployment

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

func Test_findConfigMapManifests(t *testing.T) {
	tests := []struct {
		name      string
		manifests []provider.Manifest
		want      []provider.Manifest
	}{
		{
			name: "found ConfigMap",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
			want: mustParseManifests(t, strings.TrimSpace(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
		},
		{
			name: "no match",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
`)),
			want: []provider.Manifest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := findConfigMapManifests(tt.manifests)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_findSecretManifests(t *testing.T) {
	tests := []struct {
		name      string
		manifests []provider.Manifest
		want      []provider.Manifest
	}{
		{
			name: "found Secret",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: v1
kind: Secret
metadata:
  name: nginx-secret
data:
  password: dGVzdA==
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
			want: mustParseManifests(t, strings.TrimSpace(`
apiVersion: v1
kind: Secret
metadata:
  name: nginx-secret
data:
  password: dGVzdA==
`)),
		},
		{
			name: "no match",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
			want: []provider.Manifest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := findSecretManifests(tt.manifests)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_patchManifest(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		manifests     string
		patch         config.K8sResourcePatch
		expectedError error
	}{
		{
			name:      "one op",
			manifests: "testdata/patch_manifest/patch_configmap.yaml",
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
			manifests: "testdata/patch_manifest/patch_configmap_multi_ops.yaml",
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
			manifests: "testdata/patch_manifest/patch_configmap_field.yaml",
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
			manifests: "testdata/patch_manifest/patch_configmap_field_multi_ops.yaml",
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
			t.Parallel()

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

func Test_patchManifests(t *testing.T) {
	testcases := []struct {
		name          string
		manifests     []provider.Manifest
		patches       []config.K8sResourcePatch
		expected      []provider.Manifest
		expectedError error
	}{
		{
			name: "no patches",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
        env:
        - name: VALUE
          value: none
`)),
			patches: []config.K8sResourcePatch{},
			expected: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
        env:
        - name: VALUE
          value: none
`)),
		},
		{
			name: "no manifest for the given patch",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
        env:
        - name: VALUE
          value: none
`)),
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
			expectedError: errors.New("no manifest matches the given patch: kind=Deployment, name=deployment-2"),
		},
		{
			name: "multiple patches",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
        env:
        - name: VALUE
          value: none
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
data:
  key1: value-1
  key2: value-2
`)),
			patches: []config.K8sResourcePatch{
				{
					Target: config.K8sResourcePatchTarget{
						K8sResourceReference: config.K8sResourceReference{
							Kind: "Deployment",
							Name: "deployment-1",
						},
					},
					Ops: []config.K8sResourcePatchOp{
						{
							Op:    config.K8sResourcePatchOpYAMLReplace,
							Path:  "$.spec.template.spec.containers[0].env[0].value",
							Value: "patched",
						},
					},
				},
				{
					Target: config.K8sResourcePatchTarget{
						K8sResourceReference: config.K8sResourceReference{
							Kind: "ConfigMap",
							Name: "config",
						},
					},
					Ops: []config.K8sResourcePatchOp{
						{
							Op:    config.K8sResourcePatchOpYAMLReplace,
							Path:  "$.data.key1",
							Value: "patched",
						},
						{
							Op:    config.K8sResourcePatchOpYAMLReplace,
							Path:  "$.data.key2",
							Value: "patched",
						},
					},
				},
			},
			expected: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
        env:
        - name: VALUE
          value: patched
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
data:
  key1: patched
  key2: patched
`)),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := patchManifests(tc.manifests, tc.patches, patchManifest)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}
