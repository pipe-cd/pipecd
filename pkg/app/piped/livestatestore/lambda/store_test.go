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

package lambda

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/lambda"
	"github.com/stretchr/testify/assert"
)

func TestConvertToManifest(t *testing.T) {
	t.Parallel()

	f := &lambda.GetFunctionOutput{
		Code: &types.FunctionCodeLocation{
			ImageUri: aws.String("test-image-uri"),
		},
		Configuration: &types.FunctionConfiguration{
			FunctionName:     aws.String("test-function"),
			Architectures:    []types.Architecture{types.ArchitectureArm64},
			Handler:          aws.String("test-handler"),
			Runtime:          types.RuntimeGo1x,
			EphemeralStorage: &types.EphemeralStorage{Size: aws.Int32(1024)},
			Timeout:          aws.Int32(60),
			MemorySize:       aws.Int32(128),
			Environment: &types.EnvironmentResponse{
				Variables: map[string]string{
					"env1": "value1",
					"env2": "value2",
				},
			},
			VpcConfig: &types.VpcConfigResponse{
				SubnetIds:        []string{"subnet-1", "subnet-2"},
				SecurityGroupIds: []string{"sg-1", "sg-2"},
			},
			Role: aws.String("test-role"),
			Layers: []types.Layer{
				{Arn: aws.String("layer-1")},
				{Arn: aws.String("layer-2")},
			},
		},
		Tags: map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		},
	}

	expected := &provider.FunctionManifest{
		Kind:       "LambdaFunction",
		APIVersion: "pipecd.dev/v1beta1",
		Spec: provider.FunctionManifestSpec{
			Name:     "test-function",
			ImageURI: "test-image-uri",
			Role:     "test-role",
			Handler:  "test-handler",
			Architectures: []provider.Architecture{
				{Name: "arm64"},
			},
			EphemeralStorage: &provider.EphemeralStorage{
				Size: 1024,
			},
			Runtime: "go1.x",
			Memory:  128,
			Timeout: 60,
			Tags: map[string]string{
				"tag1": "value1",
				"tag2": "value2",
			},
			Environments: map[string]string{
				"env1": "value1",
				"env2": "value2",
			},
			VPCConfig: &provider.VPCConfig{
				SecurityGroupIDs: []string{"sg-1", "sg-2"},
				SubnetIDs:        []string{"subnet-1", "subnet-2"},
			},
			Layers: []string{
				"layer-1",
				"layer-2",
			},
		},
	}

	t.Run("convert successfully", func(t *testing.T) {
		fm := convertToManifest(f)
		assert.Equal(t, expected, fm)
	})
}
