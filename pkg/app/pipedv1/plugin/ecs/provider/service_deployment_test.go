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

package provider

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateServiceDeployment(t *testing.T) {
	t.Parallel()

	d := func(arn string, s types.ServiceDeploymentStatus) *types.ServiceDeploymentBrief {
		return &types.ServiceDeploymentBrief{ServiceDeploymentArn: aws.String(arn), Status: s}
	}
	tests := []struct {
		name     string
		latest   *types.ServiceDeploymentBrief
		prev     *string
		wantDone bool
		wantErr  bool
	}{
		{name: "no deployment yet", latest: nil},
		{name: "still the previous one", latest: d("old", types.ServiceDeploymentStatusSuccessful), prev: aws.String("old")},
		{name: "in progress", latest: d("new", types.ServiceDeploymentStatusInProgress), prev: aws.String("old")},
		{name: "successful", latest: d("new", types.ServiceDeploymentStatusSuccessful), prev: aws.String("old"), wantDone: true},
		{name: "first deployment", latest: d("new", types.ServiceDeploymentStatusSuccessful), wantDone: true},
		{name: "stopped", latest: d("new", types.ServiceDeploymentStatusStopped), wantDone: true, wantErr: true},
		{name: "rollback successful is a failure", latest: d("new", types.ServiceDeploymentStatusRollbackSuccessful), wantDone: true, wantErr: true},
		{name: "rollback failed", latest: d("new", types.ServiceDeploymentStatusRollbackFailed), wantDone: true, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			done, err := evaluateServiceDeployment(tc.latest, tc.prev)
			assert.Equal(t, tc.wantDone, done)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
