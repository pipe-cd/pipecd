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
)

const revisionManifest = `
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
`

func TestRevisionManifest(t *testing.T) {
	t.Parallel()

	rm, err := ParseRevisionManifest([]byte(revisionManifest))
	require.NoError(t, err)
	require.NotEmpty(t, rm)

	// YamlBytes
	data, err := rm.YamlBytes()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// RunRevision
	got, err := rm.RunRevision()
	require.NoError(t, err)
	assert.NotEmpty(t, got)
}
