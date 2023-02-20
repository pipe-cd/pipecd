// Copyright 2023 The PipeCD Authors.
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

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/cloudrun"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const serviceManifest = `
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld
  uid: service-uid
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
`

func TestAddBuiltinLabels(t *testing.T) {
	t.Parallel()

	var (
		hash         = "commit-hash"
		pipedID      = "piped-id"
		appID        = "app-id"
		revisionName = "revision-name"
	)
	sm, err := provider.ParseServiceManifest([]byte(serviceManifest))
	require.NoError(t, err)

	ok := addBuiltinLabels(sm, hash, pipedID, appID, revisionName, nil)
	require.True(t, ok)

	want := map[string]string{
		provider.LabelManagedBy:   provider.ManagedByPiped,
		provider.LabelPiped:       pipedID,
		provider.LabelApplication: appID,
		provider.LabelCommitHash:  hash,
	}
	got := sm.Labels()
	assert.Equal(t, want, got)

	want = map[string]string{
		provider.LabelManagedBy:    provider.ManagedByPiped,
		provider.LabelPiped:        pipedID,
		provider.LabelApplication:  appID,
		provider.LabelCommitHash:   hash,
		provider.LabelRevisionName: revisionName,
	}
	got = sm.RevisionLabels()
	assert.Equal(t, want, got)
}
