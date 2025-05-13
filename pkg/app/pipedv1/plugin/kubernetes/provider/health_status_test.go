// Copyright 2025 The PipeCD Authors.
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func TestDeploymentHealthStatus(t *testing.T) {
	int32Ptr := func(i int32) *int32 { return &i }

	tests := []struct {
		name   string
		obj    *appsv1.Deployment
		health sdk.ResourceHealthStatus
		msg    string
	}{
		{
			name: "paused",
			obj: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{Paused: true, Replicas: int32Ptr(3)},
			},
			health: sdk.ResourceHealthStateUnknown,
			msg:    "Deployment is paused",
		},
		{
			name: "generation mismatch",
			obj: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Generation: 2},
				Spec:       appsv1.DeploymentSpec{Replicas: int32Ptr(3)},
				Status:     appsv1.DeploymentStatus{ObservedGeneration: 1},
			},
			health: sdk.ResourceHealthStateUnknown,
			msg:    "Waiting for rollout to finish because observed deployment generation less than desired generation",
		},
		{
			name: "progress deadline exceeded",
			obj: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{Replicas: int32Ptr(3)},
				Status: appsv1.DeploymentStatus{
					Conditions: []appsv1.DeploymentCondition{{
						Type:   appsv1.DeploymentProgressing,
						Reason: "ProgressDeadlineExceeded",
					}},
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "exceeded its progress deadline",
		},
		{
			name: "unspecified replicas",
			obj: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{},
			},
			health: sdk.ResourceHealthStateUnknown,
			msg:    "The number of desired replicas is unspecified",
		},
		{
			name: "not enough replicas",
			obj: &appsv1.Deployment{
				Spec:   appsv1.DeploymentSpec{Replicas: int32Ptr(5)},
				Status: appsv1.DeploymentStatus{Replicas: 3, AvailableReplicas: 2},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for remaining",
		},
		{
			name: "old replicas pending termination",
			obj: &appsv1.Deployment{
				Spec:   appsv1.DeploymentSpec{Replicas: int32Ptr(3)},
				Status: appsv1.DeploymentStatus{Replicas: 3, UpdatedReplicas: 2, AvailableReplicas: 3},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "old replicas are pending termination",
		},
		{
			name: "not enough available replicas",
			obj: &appsv1.Deployment{
				Spec:   appsv1.DeploymentSpec{Replicas: int32Ptr(4)},
				Status: appsv1.DeploymentStatus{Replicas: 4, UpdatedReplicas: 4, AvailableReplicas: 2},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for remaining",
		},
		{
			name: "healthy",
			obj: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "healthy-deploy"},
				Spec:       appsv1.DeploymentSpec{Replicas: int32Ptr(2)},
				Status:     appsv1.DeploymentStatus{Replicas: 2, UpdatedReplicas: 2, AvailableReplicas: 2},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h, msg := deploymentHealthStatus(tt.obj)
			assert.Equal(t, tt.health, h)
			if tt.msg != "" {
				require.NotEmpty(t, msg)
				assert.Contains(t, msg, tt.msg)
			}
		})
	}
}
