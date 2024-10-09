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
	"sigs.k8s.io/yaml"

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
