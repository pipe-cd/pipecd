// Copyright 2022 The PipeCD Authors.
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

package cloudrun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestLoadServiceManifest(t *testing.T) {
	t.Parallel()

	const (
		appDir      = "testdata"
		serviceFile = "new_manifest.yaml"
	)
	// Success
	got, err := LoadServiceManifest(appDir, serviceFile)
	require.NoError(t, err)
	assert.NotEmpty(t, got)

	// Failure
	_, err = LoadServiceManifest(appDir, "")
	assert.Error(t, err)
}

func TestMakeManagedByPipedSelector(t *testing.T) {
	t.Parallel()

	want := "pipecd-dev-managed-by=piped"
	got := MakeManagedByPipedSelector()
	assert.Equal(t, want, got)
}

func TestMakeRevisionNamesSelector(t *testing.T) {
	t.Parallel()

	names := []string{"test-1", "test-2", "test-3"}
	got := MakeRevisionNamesSelector(names)
	want := "pipecd-dev-revision-name in (test-1,test-2,test-3)"
	assert.Equal(t, want, got)
}

func TestService(t *testing.T) {
	t.Parallel()

	sm, err := ParseServiceManifest([]byte(serviceManifest))
	require.NoError(t, err)

	svc, err := sm.RunService()
	require.NoError(t, err)

	// ServiceManifest
	s := (*Service)(svc)
	got, err := s.ServiceManifest()
	require.NoError(t, err)
	assert.Equal(t, sm, got)

	// UID
	id, ok := s.UID()
	assert.True(t, ok)
	assert.Equal(t, "service-uid", id)

	// ActiveRevisionNames
	names := s.ActiveRevisionNames()
	assert.Len(t, names, 1)
}

func TestService_HealthStatus(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		manifest string
		want     model.CloudRunResourceState_HealthStatus
	}{
		{
			name: "healthy",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
  uid: service-uid
  labels:
    cloud.googleapis.com/location: asia-northeast1
    pipecd-dev-managed-by: piped
  annotations:
    run.googleapis.com/ingress: all
    run.googleapis.com/ingress-status: all
spec:
  template:
    metadata:
      name: helloworld-v010-1234567
      annotations:
        autoscaling.knative.dev/maxScale: '1'
    spec:
      containerConcurrency: 80
      timeoutSeconds: 300
      containers:
      - image: gcr.io/pipecd/helloworld:v0.1.0
        args:
        - server
        ports:
        - name: http1
          containerPort: 9085
        resources:
          limits:
            cpu: 1000m
            memory: 128Mi
  traffic:
  - revisionName: helloworld-v010-1234567
    percent: 100
status:
  observedGeneration: 65
  conditions:
  - type: Ready
    status: 'True'
    lastTransitionTime: '2021-09-15T06:56:22.222303Z'
  - type: ConfigurationsReady
    status: 'True'
    lastTransitionTime: '2021-09-15T06:55:41.885793Z'
  - type: RoutesReady
    status: 'True'
    lastTransitionTime: '2021-09-15T06:56:22.338031Z'
  latestReadyRevisionName: helloworld-v010-1234567
  latestCreatedRevisionName: helloworld-v010-1234567
  traffic:
  - revisionName: helloworld-v010-1234567
    percent: 100
`,
			want: model.CloudRunResourceState_HEALTHY,
		},
		{
			name: "unknown: unable to find status",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
  uid: service-uid
  labels:
    cloud.googleapis.com/location: asia-northeast1
    pipecd-dev-managed-by: piped
  annotations:
    run.googleapis.com/ingress: all
    run.googleapis.com/ingress-status: all
spec:
  template:
    metadata:
      name: helloworld-v010-1234567
      annotations:
        autoscaling.knative.dev/maxScale: '1'
    spec:
      containerConcurrency: 80
      timeoutSeconds: 300
      containers:
      - image: gcr.io/pipecd/helloworld:v0.1.0
        args:
        - server
        ports:
        - name: http1
          containerPort: 9085
        resources:
          limits:
            cpu: 1000m
            memory: 128Mi
  traffic:
  - revisionName: helloworld-v010-1234567
    percent: 100
`,
			want: model.CloudRunResourceState_UNKNOWN,
		},
		{
			name: "unknown: unable to parse status",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
  uid: service-uid
  labels:
    cloud.googleapis.com/location: asia-northeast1
    pipecd-dev-managed-by: piped
  annotations:
    run.googleapis.com/ingress: all
    run.googleapis.com/ingress-status: all
spec:
  template:
    metadata:
      name: helloworld-v010-1234567
      annotations:
        autoscaling.knative.dev/maxScale: '1'
    spec:
      containerConcurrency: 80
      timeoutSeconds: 300
      containers:
      - image: gcr.io/pipecd/helloworld:v0.1.0
        args:
        - server
        ports:
        - name: http1
          containerPort: 9085
        resources:
          limits:
            cpu: 1000m
            memory: 128Mi
  traffic:
  - revisionName: helloworld-v010-1234567
    percent: 100
status:
  observedGeneration: 65
  conditions:
  - type: Ready
    status: 'Unknown'
    lastTransitionTime: '2021-09-15T06:56:22.222303Z'
  - type: ConfigurationsReady
    status: 'False'
    lastTransitionTime: '2021-09-15T06:55:41.885793Z'
  - type: RoutesReady
    status: 'True'
    lastTransitionTime: '2021-09-15T06:56:22.338031Z'
  latestReadyRevisionName: helloworld-v010-1234567
  latestCreatedRevisionName: helloworld-v010-1234567
  traffic:
  - revisionName: helloworld-v010-1234567
    percent: 100
`,
			want: model.CloudRunResourceState_OTHER,
		},
		{
			name: "unhealthy",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
  uid: service-uid
  labels:
    cloud.googleapis.com/location: asia-northeast1
    pipecd-dev-managed-by: piped
  annotations:
    run.googleapis.com/ingress: all
    run.googleapis.com/ingress-status: all
spec:
  template:
    metadata:
      name: helloworld-v010-1234567
      annotations:
        autoscaling.knative.dev/maxScale: '1'
    spec:
      containerConcurrency: 80
      timeoutSeconds: 300
      containers:
      - image: gcr.io/pipecd/helloworld:v0.1.0
        args:
        - server
        ports:
        - name: http1
          containerPort: 9085
        resources:
          limits:
            cpu: 1000m
            memory: 128Mi
  traffic:
  - revisionName: helloworld-v010-1234567
    percent: 100
status:
  observedGeneration: 65
  conditions:
  - type: Ready
    status: 'False'
    lastTransitionTime: '2021-09-15T06:56:22.222303Z'
  - type: ConfigurationsReady
    status: 'Unknown'
    lastTransitionTime: '2021-09-15T06:55:41.885793Z'
  - type: RoutesReady
    status: 'Unknown'
    lastTransitionTime: '2021-09-15T06:56:22.338031Z'
  latestReadyRevisionName: helloworld-v010-1234567
  latestCreatedRevisionName: helloworld-v010-1234567
  traffic:
  - revisionName: helloworld-v010-1234567
    percent: 100
`,
			want: model.CloudRunResourceState_OTHER,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.manifest)
			sm, err := ParseServiceManifest(data)
			require.NoError(t, err)

			svc, err := sm.RunService()
			require.NoError(t, err)

			s := (*Service)(svc)
			got, _ := s.StatusConditions().HealthStatus()
			require.Equal(t, tc.want, got)
		})
	}
}

func TestRevision(t *testing.T) {
	t.Parallel()

	rm, err := ParseRevisionManifest([]byte(revisionManifest))
	require.NoError(t, err)
	require.NotEmpty(t, rm)

	rev, err := rm.RunRevision()
	require.NoError(t, err)

	r := (*Revision)(rev)
	got, err := r.RevisionManifest()
	require.NoError(t, err)
	assert.Equal(t, rm, got)
}

func TestRevision_HealthStatus(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		manifest string
		want     model.CloudRunResourceState_HealthStatus
	}{
		{
			name: "healthy",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Revision
metadata:
  name: helloworld-v010-1234567
  namespace: '0123456789'
  selfLink: /apis/serving.knative.dev/v1/namespaces/0123456789/revisions/helloworld-v010-1234567
  uid: 0123-456-789-101112-13141516
  resourceVersion: AAAAAAA
  generation: 1
  creationTimestamp: '2022-01-28T07:46:53.981805Z'
  labels:
    serving.knative.dev/route: helloworld
    serving.knative.dev/configuration: helloworld
    serving.knative.dev/configurationGeneration: '3'
    serving.knative.dev/service: helloworld
    serving.knative.dev/serviceUid: 0123-456-789-101112-13141516
    cloud.googleapis.com/location: asia-northeast1
  annotations:
    serving.knative.dev/creator: example@foo.iam.gserviceaccount.com
    autoscaling.knative.dev/maxScale: '1'
    run.googleapis.com/cpu-throttling: 'true'
  ownerReferences:
  - kind: Configuration
    name: helloworld
    uid: 0123-456-789-101112-13141516
    apiVersion: serving.knative.dev/v1
    controller: true
    blockOwnerDeletion: true
spec:
  containerConcurrency: 80
  timeoutSeconds: 300
  serviceAccountName: example@foo.iam.gserviceaccount.com
  containers:
  - image: gcr.io/pipecd/helloworld:v0.1.0
    args:
    - server
    ports:
    - name: http1
      containerPort: 9085
    resources:
      limits:
        cpu: 1000m
        memory: 128Mi
status:
  observedGeneration: 1
  conditions:
  - type: Ready
    status: 'True'
    lastTransitionTime: '2022-01-28T07:46:58.929438Z'
  - type: Active
    status: 'True'
    lastTransitionTime: '2022-01-28T07:47:04.722527Z'
    severity: Info
  - type: ContainerHealthy
    status: 'True'
    lastTransitionTime: '2022-01-28T07:46:58.929438Z'
  - type: ResourcesAvailable
    status: 'True'
    lastTransitionTime: '2022-01-28T07:46:58.150114Z'
  logUrl: https://console.cloud.google.com/logs
  imageDigest: gcr.io/pipecd/helloworld@sha256:abcdefg
`,
			want: model.CloudRunResourceState_HEALTHY,
		},
		{
			name: "unknown: unable to find status",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Revision
metadata:
  name: helloworld-v010-1234567
  namespace: '0123456789'
  selfLink: /apis/serving.knative.dev/v1/namespaces/0123456789/revisions/helloworld-v010-1234567
  uid: 0123-456-789-101112-13141516
  resourceVersion: AAAAAAA
  generation: 1
  creationTimestamp: '2022-01-28T07:46:53.981805Z'
  labels:
    serving.knative.dev/route: helloworld
    serving.knative.dev/configuration: helloworld
    serving.knative.dev/configurationGeneration: '3'
    serving.knative.dev/service: helloworld
    serving.knative.dev/serviceUid: 0123-456-789-101112-13141516
    cloud.googleapis.com/location: asia-northeast1
  annotations:
    serving.knative.dev/creator: example@foo.iam.gserviceaccount.com
    autoscaling.knative.dev/maxScale: '1'
    run.googleapis.com/cpu-throttling: 'true'
  ownerReferences:
  - kind: Configuration
    name: helloworld
    uid: 0123-456-789-101112-13141516
    apiVersion: serving.knative.dev/v1
    controller: true
    blockOwnerDeletion: true
spec:
  containerConcurrency: 80
  timeoutSeconds: 300
  serviceAccountName: example@foo.iam.gserviceaccount.com
  containers:
  - image: gcr.io/pipecd/helloworld:v0.1.0
    args:
    - server
    ports:
    - name: http1
      containerPort: 9085
    resources:
      limits:
        cpu: 1000m
        memory: 128Mi
`,
			want: model.CloudRunResourceState_UNKNOWN,
		},
		{
			name: "unknown: unable to parse status",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Revision
metadata:
  name: helloworld-v010-1234567
  namespace: '0123456789'
  selfLink: /apis/serving.knative.dev/v1/namespaces/0123456789/revisions/helloworld-v010-1234567
  uid: 0123-456-789-101112-13141516
  resourceVersion: AAAAAAA
  generation: 1
  creationTimestamp: '2022-01-28T07:46:53.981805Z'
  labels:
    serving.knative.dev/route: helloworld
    serving.knative.dev/configuration: helloworld
    serving.knative.dev/configurationGeneration: '3'
    serving.knative.dev/service: helloworld
    serving.knative.dev/serviceUid: 0123-456-789-101112-13141516
    cloud.googleapis.com/location: asia-northeast1
  annotations:
    serving.knative.dev/creator: example@foo.iam.gserviceaccount.com
    autoscaling.knative.dev/maxScale: '1'
    run.googleapis.com/cpu-throttling: 'true'
  ownerReferences:
  - kind: Configuration
    name: helloworld
    uid: 0123-456-789-101112-13141516
    apiVersion: serving.knative.dev/v1
    controller: true
    blockOwnerDeletion: true
spec:
  containerConcurrency: 80
  timeoutSeconds: 300
  serviceAccountName: example@foo.iam.gserviceaccount.com
  containers:
  - image: gcr.io/pipecd/helloworld:v0.1.0
    args:
    - server
    ports:
    - name: http1
      containerPort: 9085
    resources:
      limits:
        cpu: 1000m
        memory: 128Mi
status:
  observedGeneration: 1
  conditions:
  - type: Ready
    status: 'Unknown'
    lastTransitionTime: '2022-01-28T07:46:58.929438Z'
  - type: Active
    status: 'True'
    lastTransitionTime: '2022-01-28T07:47:04.722527Z'
    severity: Info
  - type: ContainerHealthy
    status: 'True'
    lastTransitionTime: '2022-01-28T07:46:58.929438Z'
  - type: ResourcesAvailable
    status: 'True'
    lastTransitionTime: '2022-01-28T07:46:58.150114Z'
  logUrl: https://console.cloud.google.com/logs
  imageDigest: gcr.io/pipecd/helloworld@sha256:abcdefg
`,
			want: model.CloudRunResourceState_UNKNOWN,
		},
		{
			name: "unhealthy",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Revision
metadata:
  name: helloworld-v010-1234567
  namespace: '0123456789'
  selfLink: /apis/serving.knative.dev/v1/namespaces/0123456789/revisions/helloworld-v010-1234567
  uid: 0123-456-789-101112-13141516
  resourceVersion: AAAAAAA
  generation: 1
  creationTimestamp: '2022-01-28T07:46:53.981805Z'
  labels:
    serving.knative.dev/route: helloworld
    serving.knative.dev/configuration: helloworld
    serving.knative.dev/configurationGeneration: '3'
    serving.knative.dev/service: helloworld
    serving.knative.dev/serviceUid: 0123-456-789-101112-13141516
    cloud.googleapis.com/location: asia-northeast1
  annotations:
    serving.knative.dev/creator: example@foo.iam.gserviceaccount.com
    autoscaling.knative.dev/maxScale: '1'
    run.googleapis.com/cpu-throttling: 'true'
  ownerReferences:
  - kind: Configuration
    name: helloworld
    uid: 0123-456-789-101112-13141516
    apiVersion: serving.knative.dev/v1
    controller: true
    blockOwnerDeletion: true
spec:
  containerConcurrency: 80
  timeoutSeconds: 300
  serviceAccountName: example@foo.iam.gserviceaccount.com
  containers:
  - image: gcr.io/pipecd/helloworld:v0.1.0
    args:
    - server
    ports:
    - name: http1
      containerPort: 9085
    resources:
      limits:
        cpu: 1000m
        memory: 128Mi
status:
  observedGeneration: 1
  conditions:
  - type: Ready
    status: 'False'
    lastTransitionTime: '2022-01-28T07:46:58.929438Z'
  - type: Active
    status: 'True'
    lastTransitionTime: '2022-01-28T07:47:04.722527Z'
    severity: Info
  - type: ContainerHealthy
    status: 'True'
    lastTransitionTime: '2022-01-28T07:46:58.929438Z'
  - type: ResourcesAvailable
    status: 'True'
    lastTransitionTime: '2022-01-28T07:46:58.150114Z'
  logUrl: https://console.cloud.google.com/logs
  imageDigest: gcr.io/pipecd/helloworld@sha256:abcdefg
`,
			want: model.CloudRunResourceState_OTHER,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.manifest)
			rm, err := ParseRevisionManifest(data)
			require.NoError(t, err)

			rev, err := rm.RunRevision()
			require.NoError(t, err)

			r := (*Revision)(rev)
			got, _ := r.StatusConditions().HealthStatus()
			require.Equal(t, tc.want, got)
		})
	}
}
