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
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
	"github.com/pipe-cd/pipecd/pkg/diff"
	"github.com/stretchr/testify/assert"
)

func TestIgnoreParameters(t *testing.T) {
	t.Parallel()

	livestate := provider.ECSManifests{
		ServiceDefinition: &types.Service{
			CreatedAt: aws.Time(time.Now()),
			CreatedBy: aws.String("test-createdby"),
			Events: []types.ServiceEvent{
				{
					Id: aws.String("test-event"),
				},
			},
			LoadBalancers: []types.LoadBalancer{
				{
					LoadBalancerName: aws.String("test-lb"),
				},
			},
			NetworkConfiguration: &types.NetworkConfiguration{
				AwsvpcConfiguration: &types.AwsVpcConfiguration{
					AssignPublicIp: types.AssignPublicIpDisabled,
					Subnets:        []string{"0_test-subnet", "1_test-subnet"}, // sorted
					SecurityGroups: []string{"1_test-sg", "0_test-sg"},
				},
			},
			PendingCount:    3,
			PlatformFamily:  aws.String("LINUX"),
			PlatformVersion: aws.String("1.4"),
			RunningCount:    10,
			RoleArn:         aws.String("test-role-arn"),
			ServiceArn:      aws.String("test-service-arn"),
			Status:          aws.String("ACTIVE"),
			Tags: []types.Tag{
				{
					Key:   aws.String("a_test-tag"),
					Value: aws.String("test-value-a"),
				},
				{
					Key:   aws.String("pipecd-dev-managed-by"),
					Value: aws.String("piped"),
				},
				{
					Key:   aws.String("pipecd-dev-piped"),
					Value: aws.String("test-piped"),
				},
				{
					Key:   aws.String("pipecd-dev-application"),
					Value: aws.String("test-application"),
				},
				{
					Key:   aws.String("pipecd-dev-commit-hash"),
					Value: aws.String("test-commit-hash"),
				},
				{
					Key:   aws.String("test-tag_b"),
					Value: aws.String("test-value-b"),
				},
			},
			TaskDefinition: aws.String("test-taskdef"),
			TaskSets: []types.TaskSet{
				{
					Id: aws.String("test-taskset"),
				},
			},
		},
		TaskDefinition: &types.TaskDefinition{
			Compatibilities: []types.Compatibility{types.CompatibilityEc2, types.CompatibilityFargate},
			ContainerDefinitions: []types.ContainerDefinition{
				{
					Essential: aws.Bool(false),
					PortMappings: []types.PortMapping{
						{
							HostPort: aws.Int32(80),
							Protocol: types.TransportProtocolTcp,
						},
						{
							HostPort: aws.Int32(443),
							Protocol: types.TransportProtocolTcp,
						},
					},
				},
				{
					Essential: aws.Bool(true),
					PortMappings: []types.PortMapping{
						{
							HostPort: aws.Int32(80),
							Protocol: types.TransportProtocolTcp,
						},
					},
				},
			},
			RegisteredAt:       aws.Time(time.Now()),
			RegisteredBy:       aws.String("test-registeredby"),
			RequiresAttributes: []types.Attribute{},
			Revision:           10,
			Status:             types.TaskDefinitionStatusActive,
			TaskDefinitionArn:  aws.String("test-taskdef-arn"),
		},
	}

	headManifest := provider.ECSManifests{
		ServiceDefinition: &types.Service{
			NetworkConfiguration: &types.NetworkConfiguration{
				AwsvpcConfiguration: &types.AwsVpcConfiguration{
					Subnets:        []string{"1_test-subnet", "0_test-subnet"}, // not sorted
					SecurityGroups: []string{"1_test-sg", "0_test-sg"},
				},
			},
			Tags: []types.Tag{
				// Currently, tags are ignored.
				{
					Key:   aws.String("c_test-tag"),
					Value: aws.String("test-value-c"),
				},
			},
		},
		TaskDefinition: &types.TaskDefinition{
			ContainerDefinitions: []types.ContainerDefinition{
				{
					Essential: aws.Bool(false),
					PortMappings: []types.PortMapping{
						// HostPort will be ignored
						// Protocol will be automatically tcp
						{}, {},
					},
				},
				{
					// Use default value for 'Essential'
					PortMappings: []types.PortMapping{
						// HostPort will be ignored
						// Protocol will be automatically tcp
						{},
					},
				},
			},
		},
	}

	ignoreParameters(livestate, headManifest)
	result, err := provider.Diff(
		livestate,
		headManifest,
		diff.WithEquateEmpty(),
		diff.WithIgnoreAddingMapKeys(),
		diff.WithCompareNumberAndNumericString(),
	)

	assert.NoError(t, err)
	assert.Equal(t, false, result.Diff.HasDiff())
}

func TestIgnoreAutoScalingDiff(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title                   string
		desiredCountInManifest  int
		desiredCountInLiveState int
		hasOtherDiff            bool
		ignoreDiff              bool
	}{
		{
			title:                   "n:n not ignore diff of another field",
			desiredCountInLiveState: 5,
			desiredCountInManifest:  5,
			hasOtherDiff:            true,
			ignoreDiff:              false,
		},
		{
			title:                   "0:n not ignore",
			desiredCountInLiveState: 0,
			desiredCountInManifest:  5,
			hasOtherDiff:            false,
			ignoreDiff:              false,
		},
		{
			title:                   "n:0 ignore (autoscaling is enabled)",
			desiredCountInLiveState: 5,
			desiredCountInManifest:  0,
			hasOtherDiff:            false,
			ignoreDiff:              true,
		},
		{
			title:                   "m:n not ignore",
			desiredCountInLiveState: 5,
			desiredCountInManifest:  10,
			hasOtherDiff:            false,
			ignoreDiff:              false,
		},
		{
			title:                   "no diff: not ignore (should be handled in advance)",
			desiredCountInLiveState: 5,
			desiredCountInManifest:  5,
			hasOtherDiff:            false,
			ignoreDiff:              false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			livestate := provider.ECSManifests{
				ServiceDefinition: &types.Service{
					DesiredCount: int32(tc.desiredCountInLiveState),
				},
			}
			headManifest := provider.ECSManifests{
				ServiceDefinition: &types.Service{
					DesiredCount: int32(tc.desiredCountInManifest),
				},
			}
			if tc.hasOtherDiff {
				// Add a differed field other than DesiredCount.
				headManifest.ServiceDefinition.EnableExecuteCommand = true
			}

			diff, err := provider.Diff(
				livestate,
				headManifest,
				diff.WithEquateEmpty(),
				diff.WithIgnoreAddingMapKeys(),
				diff.WithCompareNumberAndNumericString(),
			)
			assert.NoError(t, err)

			ignore := ignoreAutoScalingDiff(diff)
			assert.Equal(t, tc.ignoreDiff, ignore)
		})
	}

}
