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

package ecs

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestParseTaskDefinition(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       string
		expected    types.TaskDefinition
		expectedErr bool
	}{
		{
			name: "yaml format input",
			input: `
family: nginx-canary-fam-1
compatibilities:
  - FARGATE
networkMode: awsvpc
memory: 512
cpu: 256
`,
			expected: types.TaskDefinition{
				Family:          aws.String("nginx-canary-fam-1"),
				Compatibilities: []types.Compatibility{types.CompatibilityFargate},
				NetworkMode:     types.NetworkModeAwsvpc,
				Memory:          aws.String("512"),
				Cpu:             aws.String("256"),
			},
		},
		{
			name: "json format input",
			input: `
{
  "family": "nginx-canary-fam-1",
  "compatibilities": [
    "FARGATE"
  ],
  "networkMode": "awsvpc",
  "memory": 512,
  "cpu": 256
}
`,
			expected: types.TaskDefinition{
				Family:          aws.String("nginx-canary-fam-1"),
				Compatibilities: []types.Compatibility{types.CompatibilityFargate},
				NetworkMode:     types.NetworkModeAwsvpc,
				Memory:          aws.String("512"),
				Cpu:             aws.String("256"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseTaskDefinition([]byte(tc.input))
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestFindArtifactVersions(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       []byte
		expected    []*model.ArtifactVersion
		expectedErr bool
	}{
		{
			name: "ok",
			input: []byte(`
{
	"family": "nginx-canary-fam-1",
	"compatibilities": [
		"FARGATE"
	],
	"networkMode": "awsvpc",
	"memory": 512,
	"cpu": 256,
	"containerDefinitions" : [
		{
			"image": "gcr.io/pipecd/helloworld:v1.0.0",
			"name": "helloworld",
			"portMappings": [ 
				{ 
				"containerPort": 80,
				"hostPort": 9085,
				"protocol": "tcp"
				}
			]
		}
	]
}
`),
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
					Version: "v1.0.0",
					Name:    "helloworld",
					Url:     "gcr.io/pipecd/helloworld:v1.0.0",
				},
			},
			expectedErr: false,
		},
		{
			name: "missing containerDefinitions",
			input: []byte(`
{
	"family": "nginx-canary-fam-1",
	"compatibilities": [
		"FARGATE"
	],
	"networkMode": "awsvpc",
	"memory": 512,
	"cpu": 256,
}
`),
			expected:    nil,
			expectedErr: true,
		},
		{
			name: "missing image name",
			input: []byte(`
{
	"family": "nginx-canary-fam-1",
	"compatibilities": [
		"FARGATE"
	],
	"networkMode": "awsvpc",
	"memory": 512,
	"cpu": 256,
	"containerDefinitions" : [
		{
			"image": "gcr.io/pipecd/:v1.0.0",
			"name": "helloworld",
			"portMappings": [ 
				{ 
				"containerPort": 80,
				"hostPort": 9085,
				"protocol": "tcp"
				}
			]
		}
	]
}
`),
			expected:    nil,
			expectedErr: true,
		},
		{
			name: "multiple containers",
			input: []byte(`
{
	"family": "nginx-canary-fam-1",
	"compatibilities": [
		"FARGATE"
	],
	"networkMode": "awsvpc",
	"memory": 512,
	"cpu": 256,
	"containerDefinitions" : [
		{
			"image": "gcr.io/pipecd/helloworld:v1.0.0",
			"name": "helloworld",
			"portMappings": [ 
				{ 
				"containerPort": 80,
				"hostPort": 9085,
				"protocol": "tcp"
				}
			]
		},
		{
			"image": "gcr.io/pipecd/my-service:v1.0.0",
			"name": "my-service",
			"portMappings": [ 
				{ 
				"containerPort": 80,
				"hostPort": 9090,
				"protocol": "tcp"
				}
			]
		}
	]
}
`),
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
					Version: "v1.0.0",
					Name:    "helloworld",
					Url:     "gcr.io/pipecd/helloworld:v1.0.0",
				},
				{
					Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
					Version: "v1.0.0",
					Name:    "my-service",
					Url:     "gcr.io/pipecd/my-service:v1.0.0",
				},
			},
			expectedErr: false,
		},
		{
			name: "multiple containers with the same image",
			input: []byte(`
{
	"family": "nginx-canary-fam-1",
	"compatibilities": [
		"FARGATE"
	],
	"networkMode": "awsvpc",
	"memory": 512,
	"cpu": 256,
	"containerDefinitions" : [
		{
			"image": "gcr.io/pipecd/helloworld:v1.0.0",
			"name": "helloworld",
			"portMappings": [ 
				{ 
				"containerPort": 80,
				"hostPort": 9085,
				"protocol": "tcp"
				}
			]
		},
		{
			"image": "gcr.io/pipecd/helloworld:v1.0.0",
			"name": "helloworld-02",
			"portMappings": [ 
				{ 
				"containerPort": 80,
				"hostPort": 9091,
				"protocol": "tcp"
				}
			]
		}
	]
}
`),
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
					Version: "v1.0.0",
					Name:    "helloworld",
					Url:     "gcr.io/pipecd/helloworld:v1.0.0",
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			td, _ := parseTaskDefinition(tc.input)
			versions, err := FindArtifactVersions(td)
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.ElementsMatch(t, tc.expected, versions)
		})
	}
}
