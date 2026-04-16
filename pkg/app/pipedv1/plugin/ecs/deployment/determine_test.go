// Copyright 2026 The PipeCD Authors.
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
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func TestParseContainerImage(t *testing.T) {
	tests := []struct {
		name  string
		image string
		want  containerImage
	}{
		// No registry, just name and optional tag/digest
		{
			name:  "name and tag only",
			image: "nginx:1.21",
			want:  containerImage{name: "nginx", tag: "1.21"},
		},
		{
			name:  "name only, no tag",
			image: "nginx",
			want:  containerImage{name: "nginx"},
		},
		// Registry with domain
		{
			name:  "registry with domain and tag",
			image: "gcr.io/myproject/myapp:v1.0",
			want:  containerImage{name: "myapp", tag: "v1.0"},
		},
		{
			name:  "ECR registry with tag",
			image: "123456789.dkr.ecr.us-east-1.amazonaws.com/myapp:latest",
			want:  containerImage{name: "myapp", tag: "latest"},
		},
		// Registry with port: the colon in "host:port" must not be parsed as a tag separator
		{
			name:  "registry with port and tag",
			image: "my-registry:5000/app:latest",
			want:  containerImage{name: "app", tag: "latest"},
		},
		{
			name:  "registry with port, no tag",
			image: "my-registry:5000/app",
			want:  containerImage{name: "app"},
		},
		// Digest
		{
			name:  "digest only, no tag",
			image: "nginx@sha256:abcdef1234567890",
			want:  containerImage{name: "nginx", digest: "sha256:abcdef1234567890"},
		},
		{
			name:  "registry with digest",
			image: "gcr.io/myproject/myapp@sha256:abcdef1234567890",
			want:  containerImage{name: "myapp", digest: "sha256:abcdef1234567890"},
		},
		// Multi-level path
		{
			name:  "multi-level path with tag",
			image: "gcr.io/project-id/subpath/app:1.0",
			want:  containerImage{name: "app", tag: "1.0"},
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
		name    string
		taskDef types.TaskDefinition
		want    []sdk.ArtifactVersion
	}{
		{
			name: "single container with tag",
			taskDef: types.TaskDefinition{
				ContainerDefinitions: []types.ContainerDefinition{
					{Image: aws.String("nginx:1.21")},
				},
			},
			want: []sdk.ArtifactVersion{
				{Name: "nginx", Version: "1.21", URL: "nginx:1.21"},
			},
		},
		{
			name: "two containers with different images",
			taskDef: types.TaskDefinition{
				ContainerDefinitions: []types.ContainerDefinition{
					{Image: aws.String("nginx:1.21")},
					{Image: aws.String("redis:7.0")},
				},
			},
			want: []sdk.ArtifactVersion{
				{Name: "nginx", Version: "1.21", URL: "nginx:1.21"},
				{Name: "redis", Version: "7.0", URL: "redis:7.0"},
			},
		},
		{
			name: "two containers with same image are deduplicated",
			taskDef: types.TaskDefinition{
				ContainerDefinitions: []types.ContainerDefinition{
					{Image: aws.String("nginx:1.21")},
					{Image: aws.String("nginx:1.21")},
				},
			},
			want: []sdk.ArtifactVersion{
				{Name: "nginx", Version: "1.21", URL: "nginx:1.21"},
			},
		},
		{
			name: "container with nil image is skipped",
			taskDef: types.TaskDefinition{
				ContainerDefinitions: []types.ContainerDefinition{
					{Image: nil},
					{Image: aws.String("nginx:1.21")},
				},
			},
			want: []sdk.ArtifactVersion{
				{Name: "nginx", Version: "1.21", URL: "nginx:1.21"},
			},
		},
		{
			name: "container with empty image is skipped",
			taskDef: types.TaskDefinition{
				ContainerDefinitions: []types.ContainerDefinition{
					{Image: aws.String("")},
					{Image: aws.String("nginx:1.21")},
				},
			},
			want: []sdk.ArtifactVersion{
				{Name: "nginx", Version: "1.21", URL: "nginx:1.21"},
			},
		},
		{
			name: "registry with port",
			taskDef: types.TaskDefinition{
				ContainerDefinitions: []types.ContainerDefinition{
					{Image: aws.String("my-registry:5000/app:v2.0")},
				},
			},
			want: []sdk.ArtifactVersion{
				{Name: "app", Version: "v2.0", URL: "my-registry:5000/app:v2.0"},
			},
		},
		{
			name: "digest image uses digest as version",
			taskDef: types.TaskDefinition{
				ContainerDefinitions: []types.ContainerDefinition{
					{Image: aws.String("nginx@sha256:abcdef1234567890")},
				},
			},
			want: []sdk.ArtifactVersion{
				{Name: "nginx", Version: "sha256:abcdef1234567890", URL: "nginx@sha256:abcdef1234567890"},
			},
		},
		{
			name:    "empty container definitions",
			taskDef: types.TaskDefinition{},
			want:    []sdk.ArtifactVersion{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineVersions(tt.taskDef)
			// Sort both slices by URL for deterministic comparison since map iteration is unordered
			sort.Slice(got, func(i, j int) bool { return got[i].URL < got[j].URL })
			sort.Slice(tt.want, func(i, j int) bool { return tt.want[i].URL < tt.want[j].URL })
			assert.Equal(t, tt.want, got)
		})
	}
}
