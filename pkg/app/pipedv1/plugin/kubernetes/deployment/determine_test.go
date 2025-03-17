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
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
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
func Test_determineVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		manifests []string
		want      []sdk.ArtifactVersion
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
			want: []sdk.ArtifactVersion{
				{
					Kind:    sdk.ArtifactKindContainerImage,
					Version: "1.19.3",
					Name:    "nginx",
					URL:     "nginx:1.19.3",
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
			want: []sdk.ArtifactVersion{
				{
					Kind:    sdk.ArtifactKindContainerImage,
					Version: "1.19.3",
					Name:    "nginx",
					URL:     "nginx:1.19.3",
				},
				{
					Kind:    sdk.ArtifactKindContainerImage,
					Version: "6.0.9",
					Name:    "redis",
					URL:     "redis:6.0.9",
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
			want: []sdk.ArtifactVersion{
				{
					Kind:    sdk.ArtifactKindContainerImage,
					Version: "1.19.3",
					Name:    "nginx",
					URL:     "nginx:1.19.3",
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
			want: []sdk.ArtifactVersion{},
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
			want: []sdk.ArtifactVersion{},
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
			want: []sdk.ArtifactVersion{},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var manifests []provider.Manifest
			for _, data := range tt.manifests {
				manifests = append(manifests, mustUnmarshalYAML[provider.Manifest](t, []byte(strings.TrimSpace(data))))
			}
			got := determineVersions(manifests)
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
			want: mustParseManifests(t, strings.TrimSpace(`
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
`)),
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
			want: mustParseManifests(t, strings.TrimSpace(`
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
`)),
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
				manifests = append(manifests, mustParseManifests(t, strings.TrimSpace(data))...)
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
			want: mustParseManifests(t, strings.TrimSpace(`
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
`)),
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
			want: mustParseManifests(t, strings.TrimSpace(`
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
`)),
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
			want: mustParseManifests(t, strings.TrimSpace(`
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
`)),
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
				manifests = append(manifests, mustParseManifests(t, data)...)
			}
			got := findWorkloadManifests(manifests, tt.refs)
			assert.ElementsMatch(t, tt.want, got)
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
		{
			name: "change the order cause multi-image update",
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
      - name: redis
        image: redis:6.0.9
      - name: nginx
        image: nginx:1.19.3
`,
			want:   "Sync progressively because of updating image nginx:1.19.3 to redis:6.0.9, image redis:6.0.9 to nginx:1.19.3",
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

func TestCheckReplicasChange(t *testing.T) {
	tests := []struct {
		name        string
		old         string
		new         string
		wantBefore  string
		wantAfter   string
		wantChanged bool
	}{
		{
			name: "replicas updated",
			old: `
apiVersion: apps/v1
kind: Deployment
metadata:
    ame: nginx-deployment
spec:
  replicas: 3
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
  replicas: 5
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`,
			wantBefore:  "3",
			wantAfter:   "5",
			wantChanged: true,
		},
		{
			name: "replicas not changed",
			old: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
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
  replicas: 3
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`,
			wantBefore:  "",
			wantAfter:   "",
			wantChanged: false,
		},
		{
			name: "replicas field removed",
			old: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
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
			wantBefore:  "3",
			wantAfter:   "<nil>",
			wantChanged: true,
		},
		{
			name: "replicas field added",
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
  replicas: 3
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`,
			wantBefore:  "<nil>",
			wantAfter:   "3",
			wantChanged: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldManifests := mustParseManifests(t, tt.old)
			newManifests := mustParseManifests(t, tt.new)
			logger := zap.NewNop() // or use a real logger if available
			diffs, err := provider.Diff(oldManifests[0], newManifests[0], logger)
			require.NoError(t, err)

			before, after, changed := checkReplicasChange(diffs.Nodes())
			assert.Equal(t, tt.wantChanged, changed, "changed")
			assert.Equal(t, tt.wantBefore, before, "before")
			assert.Equal(t, tt.wantAfter, after, "after")
		})
	}
}

func Test_determineStrategy(t *testing.T) {
	tests := []struct {
		name         string
		olds         []string
		news         []string
		workloadRefs []config.K8sResourceReference
		wantStrategy sdk.SyncStrategy
		wantSummary  string
	}{
		{
			name: "no old workloads",
			olds: []string{},
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
        image: nginx:1.19.3
`,
			},
			wantStrategy: sdk.SyncStrategyQuickSync,
			wantSummary:  "Quick sync by applying all manifests because it was unable to find the currently running workloads",
		},
		{
			name: "no new workloads",
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
			news:         []string{},
			wantStrategy: sdk.SyncStrategyQuickSync,
			wantSummary:  "Quick sync by applying all manifests because it was unable to find workloads in the new manifests",
		},
		{
			name: "image updated",
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
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because of updating image nginx from 1.19.3 to 1.19.4",
		},
		{
			name: "configmap updated",
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
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
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
        image: nginx:1.19.3
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: new-value
`,
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because ConfigMap my-config was updated",
		},
		{
			name: "scaling",
			olds: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
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
  replicas: 5
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`,
			},
			wantStrategy: sdk.SyncStrategyQuickSync,
			wantSummary:  "Quick sync to scale Deployment/nginx-deployment from 3 to 5",
		},
		{
			name: "scaling nil to 1",
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
  replicas: 1
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`,
			},
			wantStrategy: sdk.SyncStrategyQuickSync,
			wantSummary:  "Quick sync to scale Deployment/nginx-deployment from <nil> to 1",
		},
		{
			name: "scaling 1 to nil",
			olds: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 1
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
        image: nginx:1.19.3
`,
			},
			wantStrategy: sdk.SyncStrategyQuickSync,
			wantSummary:  "Quick sync to scale Deployment/nginx-deployment from 1 to <nil>",
		},
		{
			name: "multiple updated scaling",
			olds: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
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
  replicas: 2
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
  replicas: 5
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
  replicas: 4
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`,
			},
			wantStrategy: sdk.SyncStrategyQuickSync,
			wantSummary:  "Quick sync to scale Deployment/nginx-deployment from 3 to 5, Deployment/redis-deployment from 2 to 4",
		},
		{
			name: "multiple updated image tags",
			olds: []string{
				`
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
			},
			news: []string{
				`
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
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because of updating image nginx from 1.19.3 to 1.19.4, image redis from 6.0.9 to 6.0.10",
		},
		{
			name: "change order of containers",
			olds: []string{
				`
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
			},
			news: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
      - name: nginx
        image: nginx:1.19.3
`,
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because of updating image nginx:1.19.3 to redis:6.0.9, image redis:6.0.9 to nginx:1.19.3",
		},
		{
			name: "workloadRef specified",
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
			workloadRefs: []config.K8sResourceReference{
				{
					Kind: "Deployment",
					Name: "nginx-deployment",
				},
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because of updating image nginx from 1.19.3 to 1.19.4",
		},
		{
			name: "workload spec.template.spec.containers.name changed",
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
      - name: nginx-new
        image: nginx:1.19.3
`,
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because pod template of workload nginx-deployment was changed",
		},
		{
			name: "workload spec.template.spec.containers.resources changed",
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
        resources:
          limits:
          cpu: "500m"
          memory: "512Mi"
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
        image: nginx:1.19.3
        resources:
          limits:
          cpu: "1"
          memory: "1Gi"
`,
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because pod template of workload nginx-deployment was changed",
		},
		{
			name: "configmap deleted",
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
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
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
        image: nginx:1.19.3
`,
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because 1 configmap/secret deleted",
		},
		{
			name: "secret deleted",
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
`, `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
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
        image: nginx:1.19.3
`,
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because 1 configmap/secret deleted",
		},
		{
			name: "configmap added",
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
        image: nginx:1.19.3
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`,
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because new 1 configmap/secret added",
		},
		{
			name: "secret added",
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
        image: nginx:1.19.3
`, `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`,
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because new 1 configmap/secret added",
		},
		{
			name: "configmap deleted and added",
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
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: old-config
  namespace: default
data:
  key: old-value
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
        image: nginx:1.19.3
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: new-config
  namespace: default
data:
  key: new-value
`,
			},
			wantStrategy: sdk.SyncStrategyPipelineSync,
			wantSummary:  "Sync progressively because ConfigMap old-config was deleted",
		},
		{
			name: "namespace default specified and unspecified",
			olds: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: default
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
        image: nginx:1.19.3
`,
			},
			wantStrategy: sdk.SyncStrategyQuickSync,
			wantSummary:  "Quick sync by applying all manifests",
		},
		{
			name: "no changes",
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
        image: nginx:1.19.3
`,
			},
			wantStrategy: sdk.SyncStrategyQuickSync,
			wantSummary:  "Quick sync by applying all manifests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var oldManifests, newManifests []provider.Manifest
			for _, data := range tt.olds {
				oldManifests = append(oldManifests, mustParseManifests(t, strings.TrimSpace(data))...)
			}
			for _, data := range tt.news {
				newManifests = append(newManifests, mustParseManifests(t, strings.TrimSpace(data))...)
			}
			logger := zap.NewNop()
			gotStrategy, gotSummary := determineStrategy(oldManifests, newManifests, tt.workloadRefs, logger)
			assert.Equal(t, tt.wantStrategy, gotStrategy)
			assert.Equal(t, tt.wantSummary, gotSummary)
		})
	}
}
