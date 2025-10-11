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

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
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
			msg:    "Waiting for remaining",
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

func int32Ptr(i int32) *int32 { return &i }

func Test_statefulSetHealthStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		obj     *appsv1.StatefulSet
		want    sdk.ResourceHealthStatus
		wantMsg string
	}{
		{
			name: "ObservedGeneration is zero",
			obj: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 2},
				Status:     appsv1.StatefulSetStatus{ObservedGeneration: 0},
			},
			want:    sdk.ResourceHealthStateUnhealthy,
			wantMsg: "Waiting for statefulset spec update to be observed",
		},
		{
			name: "Generation > ObservedGeneration",
			obj: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 2},
				Status:     appsv1.StatefulSetStatus{ObservedGeneration: 1},
			},
			want:    sdk.ResourceHealthStateUnhealthy,
			wantMsg: "Waiting for statefulset spec update to be observed",
		},
		{
			name: "Replicas is nil",
			obj: &appsv1.StatefulSet{
				Status: appsv1.StatefulSetStatus{ObservedGeneration: 1},
			},
			want:    sdk.ResourceHealthStateUnhealthy,
			wantMsg: "The number of desired replicas is unspecified",
		},
		{
			name: "ReadyReplicas != Spec.Replicas",
			obj: &appsv1.StatefulSet{
				Spec:   appsv1.StatefulSetSpec{Replicas: int32Ptr(3)},
				Status: appsv1.StatefulSetStatus{ObservedGeneration: 1, ReadyReplicas: 2},
			},
			want:    sdk.ResourceHealthStateUnhealthy,
			wantMsg: "The number of ready replicas (2) is different from the desired number (3)",
		},
		{
			name: "Partitioned rollout in progress",
			obj: &appsv1.StatefulSet{
				Spec: appsv1.StatefulSetSpec{
					Replicas: int32Ptr(5),
					UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
						Type: appsv1.RollingUpdateStatefulSetStrategyType,
						RollingUpdate: &appsv1.RollingUpdateStatefulSetStrategy{
							Partition: int32Ptr(2),
						},
					},
				},
				Status: appsv1.StatefulSetStatus{
					ObservedGeneration: 1,
					ReadyReplicas:      5,
					UpdatedReplicas:    2,
				},
			},
			want:    sdk.ResourceHealthStateUnhealthy,
			wantMsg: "Waiting for partitioned roll out to finish because 2 out of 3 new pods have been updated",
		},
		{
			name: "UpdateRevision != CurrentRevision",
			obj: &appsv1.StatefulSet{
				Spec: appsv1.StatefulSetSpec{Replicas: int32Ptr(2)},
				Status: appsv1.StatefulSetStatus{
					ObservedGeneration: 1,
					ReadyReplicas:      2,
					UpdateRevision:     "rev2",
					CurrentRevision:    "rev1",
					UpdatedReplicas:    2,
				},
			},
			want:    sdk.ResourceHealthStateUnhealthy,
			wantMsg: "Waiting for statefulset rolling update to complete 2 pods at revision rev2",
		},
		{
			name: "Healthy statefulset",
			obj: &appsv1.StatefulSet{
				Spec: appsv1.StatefulSetSpec{Replicas: int32Ptr(2)},
				Status: appsv1.StatefulSetStatus{
					ObservedGeneration: 1,
					ReadyReplicas:      2,
					UpdateRevision:     "rev1",
					CurrentRevision:    "rev1",
					UpdatedReplicas:    2,
				},
			},
			want:    sdk.ResourceHealthStateHealthy,
			wantMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotMsg := statefulSetHealthStatus(tt.obj)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantMsg, gotMsg)
		})
	}
}
