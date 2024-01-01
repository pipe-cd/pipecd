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

package cloudrun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

const serviceManifest = `
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
  observedGeneration: 5
  conditions:
  - type: Ready
    status: 'False'
    reason: RevisionFailed
    message: Revision helloworld-v010-1234567 is not ready.
    lastTransitionTime: '2022-01-31T06:18:57.242172Z'
  - type: ConfigurationsReady
    status: 'False'
    reason: ContainerMissing
    message: Image 'gcr.io/pipecd/helloworld:v0.1.0' not found.
    lastTransitionTime: '2022-01-31T06:18:57.177493Z'
  - type: RoutesReady
    status: 'False'
    reason: RevisionFailed
    message: Revision helloworld-v010-1234567 is not ready.
    lastTransitionTime: '2022-01-31T06:18:57.242172Z'
  latestReadyRevisionName: helloworld-v010-1234567
  latestCreatedRevisionName: helloworld-v010-1234567
  traffic:
  - revisionName: helloworld-v010-1234567
    percent: 100
`

func TestServiceManifest(t *testing.T) {
	t.Parallel()

	sm, err := ParseServiceManifest([]byte(serviceManifest))
	require.NoError(t, err)
	require.NotEmpty(t, sm)

	// SetRevision
	err = sm.SetRevision("helloworld-v010-1234567")
	require.NoError(t, err)

	// UpdateTraffic
	traffics := []RevisionTraffic{
		{
			RevisionName: "helloworld-v010-1234567",
			Percent:      50,
		},
		{
			RevisionName: "helloworld-v011-2345678",
			Percent:      50,
		},
	}
	err = sm.UpdateTraffic(traffics)
	require.NoError(t, err)

	// YamlBytes
	data, err := sm.YamlBytes()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// AddLabels
	labels := map[string]string{
		LabelPiped:       "hoge",
		LabelApplication: "foo",
	}
	sm.AddLabels(labels)

	// Labels
	assert.Len(t, sm.Labels(), 4)

	// AppID
	id, ok := sm.AppID()
	assert.True(t, ok)
	assert.Equal(t, "foo", id)

	// RunService
	got, err := sm.RunService()
	require.NoError(t, err)
	assert.NotEmpty(t, got)

	// AddRevisionLabels
	err = sm.AddRevisionLabels(labels)
	require.NoError(t, err)

	labels[LabelRevisionName] = "revision"
	err = sm.AddRevisionLabels(labels)
	require.NoError(t, err)

	// RevisionLabels
	v := sm.RevisionLabels()
	assert.Equal(t, labels, v)
}

func TestParseServiceManifest(t *testing.T) {
	t.Parallel()

	// Success
	data := []byte(serviceManifest)
	sm, err := ParseServiceManifest(data)
	require.NoError(t, err)
	require.Equal(t, "helloworld", sm.Name)

	// Failure
	data = []byte("error")
	_, err = ParseServiceManifest(data)
	require.Error(t, err)
}

func TestDecideRevisionName(t *testing.T) {
	t.Parallel()

	data := []byte(serviceManifest)
	sm, err := ParseServiceManifest(data)
	require.NoError(t, err)

	name, err := DecideRevisionName(sm, "12345678912345678")
	require.NoError(t, err)
	require.Equal(t, "helloworld-v010-1234567", name)
}

func TestFindImageTag(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		manifest string
		want     string
		wantErr  bool
	}{
		{
			name: "ok",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
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
			want:    "v0.1.0",
			wantErr: false,
		},
		{
			name: "err: containers missing",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
spec:
  template:
    metadata:
      name: helloworld-v010-1234567
      annotations:
        autoscaling.knative.dev/maxScale: '1'
    spec:
      containerConcurrency: 80
      timeoutSeconds: 300
`,
			want:    "",
			wantErr: true,
		},
		{
			name: "err: image missing",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
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
      - args:
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
			want:    "",
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.manifest)
			sm, err := ParseServiceManifest(data)
			require.NoError(t, err)

			got, err := FindImageTag(sm)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestFindArtifactVersions(t *testing.T) {
	testcases := []struct {
		name     string
		manifest string
		want     []*model.ArtifactVersion
		wantErr  bool
	}{
		{
			name: "ok",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
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
			want: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
					Version: "v0.1.0",
					Name:    "helloworld",
					Url:     "gcr.io/pipecd/helloworld:v0.1.0",
				},
			},
			wantErr: false,
		},
		{
			name: "err: containers missing",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
spec:
  template:
    metadata:
      name: helloworld-v010-1234567
      annotations:
        autoscaling.knative.dev/maxScale: '1'
    spec:
      containerConcurrency: 80
      timeoutSeconds: 300
`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "err: image missing",
			manifest: `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
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
      - args:
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
			want:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.manifest)
			sm, err := ParseServiceManifest(data)
			require.NoError(t, err)

			got, err := FindArtifactVersions(sm)
			require.Equal(t, tc.wantErr, err != nil)
			require.Equal(t, tc.want, got)
		})
	}
}
