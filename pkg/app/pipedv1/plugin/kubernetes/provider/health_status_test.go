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
	corev1 "k8s.io/api/core/v1"
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

func TestReplicaSetHealthStatus(t *testing.T) {
	int32Ptr := func(i int32) *int32 { return &i }

	tests := []struct {
		name   string
		obj    *appsv1.ReplicaSet
		health sdk.ResourceHealthStatus
		msg    string
	}{
		{
			name: "healthy replicaset",
			obj: &appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Spec:       appsv1.ReplicaSetSpec{Replicas: int32Ptr(3)},
				Status: appsv1.ReplicaSetStatus{
					ObservedGeneration: 1,
					Replicas:           3,
					AvailableReplicas:  3,
					ReadyReplicas:      3,
				},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "",
		},
		{
			name: "observed generation is 0",
			obj: &appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Spec:       appsv1.ReplicaSetSpec{Replicas: int32Ptr(3)},
				Status: appsv1.ReplicaSetStatus{
					ObservedGeneration: 0,
					Replicas:           3,
					AvailableReplicas:  3,
					ReadyReplicas:      3,
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for rollout to finish because observed replica set generation less than desired generation",
		},
		{
			name: "generation mismatch",
			obj: &appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 2},
				Spec:       appsv1.ReplicaSetSpec{Replicas: int32Ptr(3)},
				Status: appsv1.ReplicaSetStatus{
					ObservedGeneration: 1,
					Replicas:           3,
					AvailableReplicas:  3,
					ReadyReplicas:      3,
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for rollout to finish because observed replica set generation less than desired generation",
		},
		{
			name: "replica failure condition true",
			obj: &appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Spec:       appsv1.ReplicaSetSpec{Replicas: int32Ptr(3)},
				Status: appsv1.ReplicaSetStatus{
					ObservedGeneration: 1,
					Replicas:           3,
					AvailableReplicas:  3,
					ReadyReplicas:      3,
					Conditions: []appsv1.ReplicaSetCondition{
						{
							Type:    appsv1.ReplicaSetReplicaFailure,
							Status:  corev1.ConditionTrue,
							Message: "Failed to create pod: insufficient resources",
						},
					},
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Failed to create pod: insufficient resources",
		},
		{
			name: "nil replicas spec",
			obj: &appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Spec:       appsv1.ReplicaSetSpec{Replicas: nil},
				Status: appsv1.ReplicaSetStatus{
					ObservedGeneration: 1,
					Replicas:           0,
					AvailableReplicas:  0,
					ReadyReplicas:      0,
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "The number of desired replicas is unspecified",
		},
		{
			name: "insufficient available replicas",
			obj: &appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Spec:       appsv1.ReplicaSetSpec{Replicas: int32Ptr(3)},
				Status: appsv1.ReplicaSetStatus{
					ObservedGeneration: 1,
					Replicas:           3,
					AvailableReplicas:  2, // Less than desired
					ReadyReplicas:      2,
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for remaining 1/3 replicas to be available",
		},
		{
			name: "ready replicas mismatch",
			obj: &appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Spec:       appsv1.ReplicaSetSpec{Replicas: int32Ptr(3)},
				Status: appsv1.ReplicaSetStatus{
					ObservedGeneration: 1,
					Replicas:           3,
					AvailableReplicas:  3,
					ReadyReplicas:      2, // Less than desired
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "The number of ready replicas (2) is different from the desired number (3)",
		},
		{
			name: "zero replicas - healthy",
			obj: &appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Spec:       appsv1.ReplicaSetSpec{Replicas: int32Ptr(0)},
				Status: appsv1.ReplicaSetStatus{
					ObservedGeneration: 1,
					Replicas:           0,
					AvailableReplicas:  0,
					ReadyReplicas:      0,
				},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotMsg := replicaSetHealthStatus(tt.obj)
			assert.Equal(t, tt.health, got)
			assert.Equal(t, tt.msg, gotMsg)
		})
	}
}

func TestDaemonSetHealthStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		obj    *appsv1.DaemonSet
		health sdk.ResourceHealthStatus
		msg    string
	}{
		{
			name: "observed generation is 0",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status:     appsv1.DaemonSetStatus{ObservedGeneration: 0},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for rollout to finish because observed daemon set generation less than desired generation",
		},
		{
			name: "generation mismatch",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 2, Name: "test-daemonset"},
				Status:     appsv1.DaemonSetStatus{ObservedGeneration: 1},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for rollout to finish because observed daemon set generation less than desired generation",
		},
		{
			name: "updated number scheduled less than desired",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 5,
					UpdatedNumberScheduled: 3,
					NumberAvailable:        5,
					NumberMisscheduled:     0,
					NumberUnavailable:      0,
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for daemon set \"test-daemonset\" rollout to finish because 3 out of 5 new pods have been updated",
		},
		{
			name: "number available less than desired",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 5,
					UpdatedNumberScheduled: 5,
					NumberAvailable:        3,
					NumberMisscheduled:     0,
					NumberUnavailable:      0,
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for daemon set \"test-daemonset\" rollout to finish because 3 of 5 updated pods are available",
		},
		{
			name: "number misscheduled greater than 0",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 5,
					UpdatedNumberScheduled: 5,
					NumberAvailable:        5,
					NumberMisscheduled:     2,
					NumberUnavailable:      0,
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "2 nodes that are running the daemon pod, but are not supposed to run the daemon pod",
		},
		{
			name: "number unavailable greater than 0",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 5,
					UpdatedNumberScheduled: 5,
					NumberAvailable:        5,
					NumberMisscheduled:     0,
					NumberUnavailable:      1,
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "1 nodes that should be running the daemon pod and have none of the daemon pod running and available",
		},
		{
			name: "healthy daemonset",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 5,
					UpdatedNumberScheduled: 5,
					NumberAvailable:        5,
					NumberMisscheduled:     0,
					NumberUnavailable:      0,
				},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "",
		},
		{
			name: "healthy daemonset with zero desired pods",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 0,
					UpdatedNumberScheduled: 0,
					NumberAvailable:        0,
					NumberMisscheduled:     0,
					NumberUnavailable:      0,
				},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "",
		},
		{
			name: "multiple issues - should report first one (updated number scheduled)",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 5,
					UpdatedNumberScheduled: 3, // This should be reported first
					NumberAvailable:        2, // This would be reported second
					NumberMisscheduled:     1, // This would be reported third
					NumberUnavailable:      1, // This would be reported fourth
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for daemon set \"test-daemonset\" rollout to finish because 3 out of 5 new pods have been updated",
		},
		{
			name: "multiple issues - should report second one (number available)",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 5,
					UpdatedNumberScheduled: 5, // This is OK
					NumberAvailable:        2, // This should be reported
					NumberMisscheduled:     1, // This would be reported second
					NumberUnavailable:      1, // This would be reported third
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Waiting for daemon set \"test-daemonset\" rollout to finish because 2 of 5 updated pods are available",
		},
		{
			name: "multiple issues - should report third one (number misscheduled)",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 5,
					UpdatedNumberScheduled: 5, // This is OK
					NumberAvailable:        5, // This is OK
					NumberMisscheduled:     2, // This should be reported
					NumberUnavailable:      1, // This would be reported second
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "2 nodes that are running the daemon pod, but are not supposed to run the daemon pod",
		},
		{
			name: "multiple issues - should report fourth one (number unavailable)",
			obj: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1, Name: "test-daemonset"},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 5,
					UpdatedNumberScheduled: 5, // This is OK
					NumberAvailable:        5, // This is OK
					NumberMisscheduled:     0, // This is OK
					NumberUnavailable:      2, // This should be reported
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "2 nodes that should be running the daemon pod and have none of the daemon pod running and available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotMsg := daemonSetHealthStatus(tt.obj)
			assert.Equal(t, tt.health, got)
			assert.Equal(t, tt.msg, gotMsg)
		})
	}
}

func TestPodHealthStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		obj    *corev1.Pod
		health sdk.ResourceHealthStatus
		msg    string
	}{
		{
			name: "healthy pod with RestartPolicyAlways and Running phase",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodRunning,
					Message: "Pod is running",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Running: &corev1.ContainerStateRunning{},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "Pod is running",
		},
		{
			name: "healthy pod with RestartPolicyAlways and Succeeded phase",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodSucceeded,
					Message: "Pod completed successfully",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Terminated: &corev1.ContainerStateTerminated{
									ExitCode: 0,
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "Pod completed successfully",
		},
		{
			name: "unhealthy pod with RestartPolicyAlways and container error",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodRunning,
					Message: "Pod is running",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "ErrImagePull",
									Message: "Failed to pull image: image not found",
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Failed to pull image: image not found",
		},
		{
			name: "unhealthy pod with RestartPolicyAlways and container backoff",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodRunning,
					Message: "Pod is running",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "CrashLoopBackOff",
									Message: "Container is crashing repeatedly",
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Container is crashing repeatedly",
		},
		{
			name: "unhealthy pod with RestartPolicyAlways and multiple container errors",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodRunning,
					Message: "Pod is running",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "ErrImagePull",
									Message: "Failed to pull image: image not found",
								},
							},
						},
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "CrashLoopBackOff",
									Message: "Container is crashing repeatedly",
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Failed to pull image: image not found, Container is crashing repeatedly",
		},
		{
			name: "pod with RestartPolicyAlways but no error conditions",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodRunning,
					Message: "Pod is running",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "ContainerCreating",
									Message: "Container is being created",
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "Pod is running",
		},
		{
			name: "pod with RestartPolicyOnFailure and Running phase",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodRunning,
					Message: "Pod is running",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "ErrImagePull",
									Message: "Failed to pull image: image not found",
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "Pod is running",
		},
		{
			name: "pod with RestartPolicyNever and Running phase",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodRunning,
					Message: "Pod is running",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "ErrImagePull",
									Message: "Failed to pull image: image not found",
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "Pod is running",
		},
		{
			name: "pod with RestartPolicyAlways and Pending phase",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodPending,
					Message: "Pod is pending",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "ErrImagePull",
									Message: "Failed to pull image: image not found",
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Failed to pull image: image not found",
		},
		{
			name: "pod with RestartPolicyAlways and Failed phase",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodFailed,
					Message: "Pod failed",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "ErrImagePull",
									Message: "Failed to pull image: image not found",
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Failed to pull image: image not found",
		},
		{
			name: "pod with RestartPolicyAlways and no container statuses",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:             corev1.PodRunning,
					Message:           "Pod is running",
					ContainerStatuses: []corev1.ContainerStatus{},
				},
			},
			health: sdk.ResourceHealthStateHealthy,
			msg:    "Pod is running",
		},
		{
			name: "pod with RestartPolicyAlways and container with error suffix",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodRunning,
					Message: "Pod is running",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "ImagePullError",
									Message: "Failed to pull image",
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Failed to pull image",
		},
		{
			name: "pod with RestartPolicyAlways and container with backoff suffix",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
				Status: corev1.PodStatus{
					Phase:   corev1.PodRunning,
					Message: "Pod is running",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason:  "ImagePullBackOff",
									Message: "Backing off from pulling image",
								},
							},
						},
					},
				},
			},
			health: sdk.ResourceHealthStateUnhealthy,
			msg:    "Backing off from pulling image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotMsg := podHealthStatus(tt.obj)
			assert.Equal(t, tt.health, got)
			assert.Equal(t, tt.msg, gotMsg)
		})
	}
}
