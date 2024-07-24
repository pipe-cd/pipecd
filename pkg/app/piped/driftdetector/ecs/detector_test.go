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

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
	"github.com/pipe-cd/pipecd/pkg/diff"
	"github.com/stretchr/testify/assert"
)

func TestIgnoreDesiredCountDiff(t *testing.T) {
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

			ignore := ignoreDesiredCountDiff(diff)
			assert.Equal(t, tc.ignoreDiff, ignore)
		})
	}

}
