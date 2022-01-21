// Copyright 2020 The PipeCD Authors.
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

	"github.com/stretchr/testify/require"
)

const manifest = `
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
      name: helloworld-v050-0b13751
      annotations:
        autoscaling.knative.dev/maxScale: '1'
    spec:
      containerConcurrency: 80
      timeoutSeconds: 300
      containers:
      - image: gcr.io/pipecd/helloworld:v0.5.0
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
  - revisionName: helloworld-v050-0b13751
    percent: 100
`

func TestServiceManifest(t *testing.T) {
	sm, err := ParseServiceManifest([]byte(manifest))
	require.NoError(t, err)
	require.NotEmpty(t, sm)

	// SetRevision
	err = sm.SetRevision("helloworld-v050-0b13751")
	require.NoError(t, err)

	// UpdateTraffic
	traffics := []RevisionTraffic{
		{
			RevisionName: "helloworld-v050-0b13751",
			Percent:      50,
		},
		{
			RevisionName: "helloworld-v050-cb01dce",
			Percent:      50,
		},
	}
	err = sm.UpdateTraffic(traffics)
	require.NoError(t, err)

	// YamlBytes
	data, err := sm.YamlBytes()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// AddLabels
	labels := map[string]string{
		LabelPiped:       "hoge",
		LabelApplication: "foo",
	}
	sm.AddLabels(labels)

	// Labels
	require.Len(t, sm.Labels(), 4)
}

func TestLoadServiceManifest(t *testing.T) {
	// Success
	sm, err := loadServiceManifest("testdata/new_manifest.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, sm)

	// Failed
	sm, err = loadServiceManifest("testdata/not_found")
	require.Error(t, err)
}

func TestParseServiceManifest(t *testing.T) {
	// Success
	data := []byte(manifest)
	sm, err := ParseServiceManifest(data)
	require.NoError(t, err)
	require.Equal(t, "helloworld", sm.Name)

	// Failed
	data = []byte("error")
	_, err = ParseServiceManifest(data)
	require.Error(t, err)
}

func TestDecideRevisionName(t *testing.T) {
	data := []byte(manifest)
	sm, err := ParseServiceManifest(data)
	require.NoError(t, err)

	name, err := DecideRevisionName(sm, "bbdc2ed674ce4fd987")
	require.NoError(t, err)
	require.Equal(t, "helloworld-v050-bbdc2ed", name)
}

func TestFindImageTag(t *testing.T) {
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
      name: helloworld-v050-0b13751
      annotations:
        autoscaling.knative.dev/maxScale: '1'
    spec:
      containerConcurrency: 80
      timeoutSeconds: 300
      containers:
      - image: gcr.io/pipecd/helloworld:v0.5.0
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
  - revisionName: helloworld-v050-0b13751
    percent: 100
`,
			want:    "v0.5.0",
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
      name: helloworld-v050-0b13751
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
      name: helloworld-v050-0b13751
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
  - revisionName: helloworld-v050-0b13751
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
