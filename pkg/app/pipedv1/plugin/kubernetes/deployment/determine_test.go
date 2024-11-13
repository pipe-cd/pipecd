// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deployment

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func mustUnmarshalYAML[T any](t *testing.T, data []byte) T {
	t.Helper()

	// Convert YAML to JSON.
	// we define structs without defining UnmarshalYAML method, so we can't use yaml.Unmarshal directly.
	// Instead, we convert YAML to JSON and then unmarshal JSON to the struct.
	j, err := yaml.YAMLToJSON(data)
	require.NoError(t, err)

	// then, unmarshal JSON to the struct.
	var m T
	require.NoError(t, json.Unmarshal(j, &m))

	return m
}

func mustParseManifests(t *testing.T, data string) []provider.Manifest {
	t.Helper()

	manifests, err := provider.ParseManifests(data)
	require.NoError(t, err)

	return manifests
}

func TestParseContainerImage(t *testing.T) {
	tests := []struct {
		name  string
		image string
		want  containerImage
	}{
		{
			name:  "image with tag",
			image: "nginx:1.19.3",
			want:  containerImage{name: "nginx", tag: "1.19.3"},
		},
		{
			name:  "image without tag",
			image: "nginx",
			want:  containerImage{name: "nginx", tag: ""},
		},
		{
			name:  "image with tag and registry",
			image: "docker.io/nginx:1.19.3",
			want:  containerImage{name: "nginx", tag: "1.19.3"},
		},
		{
			name:  "image with tag and repository",
			image: "myrepo/nginx:1.19.3",
			want:  containerImage{name: "nginx", tag: "1.19.3"},
		},
		{
			name:  "image with tag, registry and repository",
			image: "docker.io/myrepo/nginx:1.19.3",
			want:  containerImage{name: "nginx", tag: "1.19.3"},
		},
		{
			name:  "image without tag, with registry and repository",
			image: "docker.io/myrepo/nginx",
			want:  containerImage{name: "nginx", tag: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseContainerImage(tt.image)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestDetermineVersions(t *testing.T) {
	tests := []struct {
		name      string
		manifests []string
		want      []*model.ArtifactVersion
		wantErr   bool
	}{
		{
			name: "single manifest with one container",
			manifests: []string{
				`
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
`,
			},
			want: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
					Version: "1.19.3",
					Name:    "nginx",
					Url:     "nginx:1.19.3",
				},
			},
		},
		{
			name: "multiple manifests with multiple containers",
			manifests: []string{
				`
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
`,
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`,
			},
			want: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
					Version: "1.19.3",
					Name:    "nginx",
					Url:     "nginx:1.19.3",
				},
				{
					Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
					Version: "6.0.9",
					Name:    "redis",
					Url:     "redis:6.0.9",
				},
			},
		},
		{
			name: "manifest with duplicate images",
			manifests: []string{
				`
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
        - name: nginx
          image: nginx:1.19.3
`,
			},
			want: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
					Version: "1.19.3",
					Name:    "nginx",
					Url:     "nginx:1.19.3",
				},
			},
		},
		{
			name: "manifest with no containers",
			manifests: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: empty-deployment
spec:
  template:
    spec:
      containers: []
`,
			},
			want: []*model.ArtifactVersion{},
		},
		{
			name: "manifest with missing image field",
			manifests: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: missing-image-deployment
spec:
  template:
    spec:
      containers:
        - name: nginx
`,
			},
			want: []*model.ArtifactVersion{},
		},
		{
			name: "manifest with non-string image field",
			manifests: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: non-string-image-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: 12345
`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "manifest with no containers field",
			manifests: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: no-containers-deployment
spec:
  template:
    spec: {}
`,
			},
			want:    []*model.ArtifactVersion{},
			wantErr: false,
		},
		{
			name: "manifest with invalid containers field -- returns error",
			manifests: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: no-containers-deployment
spec:
  template:
    spec:
      containers: "invalid-containers-field"
`,
			},
			wantErr: true,
		},
		{
			name: "manifest with invalid containers field -- skipped",
			manifests: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: no-containers-deployment
spec:
  template:
    spec:
      containers:
        - "invalid-containers-field"
`,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var manifests []provider.Manifest
			for _, data := range tt.manifests {
				manifests = append(manifests, mustUnmarshalYAML[provider.Manifest](t, []byte(strings.TrimSpace(data))))
			}
			got, err := determineVersions(manifests)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestFindManifests(t *testing.T) {
	tests := []struct {
		name      string
		kind      string
		nameField string
		manifests []string
		want      []provider.Manifest
	}{
		{
			name: "find by kind",
			kind: "Deployment",
			manifests: []string{
				`
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
`,
				`
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
`,
			},
			want: []provider.Manifest{
				mustUnmarshalYAML[provider.Manifest](t, []byte(strings.TrimSpace(`
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
`))),
			},
		},
		{
			name:      "find by kind and name",
			kind:      "Deployment",
			nameField: "nginx-deployment",
			manifests: []string{
				`
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
`,
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`,
			},
			want: []provider.Manifest{
				mustUnmarshalYAML[provider.Manifest](t, []byte(strings.TrimSpace(`
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
`))),
			},
		},
		{
			name: "no match",
			kind: "StatefulSet",
			manifests: []string{
				`
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
`,
			},
			want: []provider.Manifest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var manifests []provider.Manifest
			for _, data := range tt.manifests {
				manifests = append(manifests, mustUnmarshalYAML[provider.Manifest](t, []byte(strings.TrimSpace(data))))
			}
			got := findManifests(tt.kind, tt.nameField, manifests)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestFindWorkloadManifests(t *testing.T) {
	tests := []struct {
		name      string
		manifests []string
		refs      []config.K8sResourceReference
		want      []provider.Manifest
	}{
		{
			name: "default to Deployment kind",
			manifests: []string{
				`
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
`,
				`
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
`,
			},
			refs: nil,
			want: []provider.Manifest{
				mustUnmarshalYAML[provider.Manifest](t, []byte(strings.TrimSpace(`
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
`))),
			},
		},
		{
			name: "specified kind and name",
			manifests: []string{
				`
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
`,
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`,
			},
			refs: []config.K8sResourceReference{
				{
					Kind: "Deployment",
					Name: "nginx-deployment",
				},
			},
			want: []provider.Manifest{
				mustUnmarshalYAML[provider.Manifest](t, []byte(strings.TrimSpace(`
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
`))),
			},
		},
		{
			name: "specified kind only",
			manifests: []string{
				`
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
`,
				`
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis-statefulset
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`,
			},
			refs: []config.K8sResourceReference{
				{
					Kind: "StatefulSet",
				},
			},
			want: []provider.Manifest{
				mustUnmarshalYAML[provider.Manifest](t, []byte(strings.TrimSpace(`
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis-statefulset
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`))),
			},
		},
		{
			name: "no match",
			manifests: []string{
				`
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
`,
			},
			refs: []config.K8sResourceReference{
				{
					Kind: "StatefulSet",
					Name: "redis-statefulset",
				},
			},
			want: []provider.Manifest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var manifests []provider.Manifest
			for _, data := range tt.manifests {
				manifests = append(manifests, mustUnmarshalYAML[provider.Manifest](t, []byte(strings.TrimSpace(data))))
			}
			got := findWorkloadManifests(manifests, tt.refs)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestFindUpdatedWorkloads(t *testing.T) {
	tests := []struct {
		name string
		olds []string
		news []string
		want []workloadPair
	}{
		{
			name: "single updated workload",
			olds: []string{
				`
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
`,
			},
			news: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.4
`,
			},
			want: []workloadPair{
				{
					old: mustParseManifests(t, `
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
`)[0],
					new: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.4
`)[0],
				},
			},
		},
		{
			name: "multiple updated workloads",
			olds: []string{
				`
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
`,
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`,
			},
			news: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.4
`,
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.10
`,
			},
			want: []workloadPair{
				{
					old: mustParseManifests(t, `
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
`)[0],
					new: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.4
`)[0],
				},
				{
					old: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`)[0],
					new: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.10
`)[0],
				},
			},
		},
		{
			name: "no updated workloads",
			olds: []string{
				`
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
`,
			},
			news: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:7.0.0
`,
			},
			want: []workloadPair{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldManifests := mustParseManifests(t, strings.Join(tt.olds, "\n---\n"))
			newManifests := mustParseManifests(t, strings.Join(tt.news, "\n---\n"))
			got := findUpdatedWorkloads(oldManifests, newManifests)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestFindConfigs(t *testing.T) {
	tests := []struct {
		name      string
		manifests []string
		want      map[provider.ResourceKey]provider.Manifest
	}{
		{
			name: "find ConfigMap and Secret",
			manifests: []string{
				`
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`,
				`
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`,
				`
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
`,
			},
			want: map[provider.ResourceKey]provider.Manifest{
				{
					APIVersion: "v1",
					Kind:       "ConfigMap",
					Name:       "my-config",
					Namespace:  "default",
				}: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`)[0],
				{
					APIVersion: "v1",
					Kind:       "Secret",
					Name:       "my-secret",
					Namespace:  "default",
				}: mustParseManifests(t, `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`)[0],
			},
		},
		{
			name: "no ConfigMap or Secret",
			manifests: []string{
				`
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
`,
			},
			want: map[provider.ResourceKey]provider.Manifest{},
		},
		{
			name: "only ConfigMap",
			manifests: []string{
				`
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`,
			},
			want: map[provider.ResourceKey]provider.Manifest{
				{
					APIVersion: "v1",
					Kind:       "ConfigMap",
					Name:       "my-config",
					Namespace:  "default",
				}: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`)[0],
			},
		},
		{
			name: "only Secret",
			manifests: []string{
				`
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`,
			},
			want: map[provider.ResourceKey]provider.Manifest{
				{
					APIVersion: "v1",
					Kind:       "Secret",
					Name:       "my-secret",
					Namespace:  "default",
				}: mustParseManifests(t, `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`)[0],
			},
		},
		{
			name: "non-default namespace",
			manifests: []string{
				`
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: custom-namespace
data:
  key: value
`,
				`
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: custom-namespace
data:
  key: dmFsdWU=
`,
			},
			want: map[provider.ResourceKey]provider.Manifest{
				{
					APIVersion: "v1",
					Kind:       "ConfigMap",
					Name:       "my-config",
					Namespace:  "custom-namespace",
				}: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: custom-namespace
data:
  key: value
`)[0],
				{
					APIVersion: "v1",
					Kind:       "Secret",
					Name:       "my-secret",
					Namespace:  "custom-namespace",
				}: mustParseManifests(t, `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: custom-namespace
data:
  key: dmFsdWU=
`)[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var manifests []provider.Manifest
			for _, data := range tt.manifests {
				manifests = append(manifests, mustParseManifests(t, data)...)
			}
			got := findConfigsAndSecrets(manifests)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckImageChange(t *testing.T) {
	tests := []struct {
		name   string
		old    string
		new    string
		want   string
		wantOk bool
	}{
		{
			name: "image updated",
			old: `
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
`,
			new: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.4
`,
			want:   "Sync progressively because of updating image nginx from 1.19.3 to 1.19.4",
			wantOk: true,
		},
		{
			name: "image name changed",
			old: `
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
`,
			new: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`,
			want:   "Sync progressively because of updating image nginx:1.19.3 to redis:6.0.9",
			wantOk: true,
		},
		{
			name: "no image change",
			old: `
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
`,
			new: `
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
`,
			want:   "",
			wantOk: false,
		},
		{
			name: "multiple image updates",
			old: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
      - name: redis
        image: redis:6.0.9
`,
			new: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.4
      - name: redis
        image: redis:6.0.10
`,
			want:   "Sync progressively because of updating image nginx from 1.19.3 to 1.19.4, image redis from 6.0.9 to 6.0.10",
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldManifests := mustParseManifests(t, tt.old)
			newManifests := mustParseManifests(t, tt.new)
			logger := zap.NewNop() // or use a real logger if available
			diffs, err := provider.Diff(oldManifests[0], newManifests[0], logger)
			require.NoError(t, err)

			got, ok := checkImageChange(diffs.Nodes())
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}
