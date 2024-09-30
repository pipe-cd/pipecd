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

package provider

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindReferencingConfigMapsInDeployment(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		manifest string
		expected []string
	}{
		{
			name: "no configmap",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  labels:
    app: simple
spec:
  replicas: 2
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: primary
      annotations:
        sidecar.istio.io/inject: "false"
    spec:
      containers:
      - name: helloworld
        image: gcr.io/pipecd/helloworld:v0.5.0
        args:
          - server
        ports:
        - containerPort: 9085
`,
			expected: nil,
		},
		{
			name: "one configmap",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-by-config-change
  labels:
    app: canary-by-config-change
spec:
  replicas: 2
  selector:
    matchLabels:
      app: canary-by-config-change
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: canary-by-config-change
        pipecd.dev/variant: primary
    spec:
      containers:
        - name: helloworld
          image: gcr.io/pipecd/helloworld:v0.5.0
          args:
            - server
          ports:
            - containerPort: 9085
          volumeMounts:
            - name: config
              mountPath: /etc/pipecd-config
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: canary-by-config-change
`,
			expected: []string{
				"canary-by-config-change",
			},
		},
		{
			name: "multiple configmaps",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-by-config-change
  labels:
    app: canary-by-config-change
spec:
  replicas: 2
  selector:
    matchLabels:
      app: canary-by-config-change
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: canary-by-config-change
        pipecd.dev/variant: primary
    spec:
      initContainers:
        - name: init
          image: gcr.io/pipecd/helloworld:v0.5.0
          env:
            - name: env1
              valueFrom:
                configMapKeyRef:
                  name: init-configmap-1
                  key: key1
      containers:
        - name: helloworld
          image: gcr.io/pipecd/helloworld:v0.5.0
          args:
            - server
          ports:
            - containerPort: 9085
          env:
            - name: env1
              valueFrom:
                configMapKeyRef:
                  name: configmap-1
                  key: key1
            - name: env2
              valueFrom:
                configMapKeyRef:
                  name: configmap-2
                  key: key2
          volumeMounts:
            - name: config
              mountPath: /etc/pipecd-config
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: canary-by-config-change
        - name: config2
          configMap:
            name: configmap-2
`,
			expected: []string{
				"canary-by-config-change",
				"configmap-1",
				"configmap-2",
				"init-configmap-1",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := ParseManifests(tc.manifest)
			require.NoError(t, err)
			require.Equal(t, 1, len(manifests))

			d := &appsv1.Deployment{}
			err = manifests[0].ConvertToStructuredObject(d)
			require.NoError(t, err)

			out := FindReferencingConfigMapsInDeployment(d)
			assert.Equal(t, tc.expected, out)
		})
	}
}

func TestFindReferencingSecretsInDeployment(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		manifest string
		expected []string
	}{
		{
			name: "no secret",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  labels:
    app: simple
spec:
  replicas: 2
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: primary
      annotations:
        sidecar.istio.io/inject: "false"
    spec:
      containers:
      - name: helloworld
        image: gcr.io/pipecd/helloworld:v0.5.0
        args:
          - server
        ports:
        - containerPort: 9085
`,
			expected: nil,
		},
		{
			name: "one secret",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-by-config-change
  labels:
    app: canary-by-config-change
spec:
  replicas: 2
  selector:
    matchLabels:
      app: canary-by-config-change
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: canary-by-config-change
        pipecd.dev/variant: primary
    spec:
      containers:
        - name: helloworld
          image: gcr.io/pipecd/helloworld:v0.5.0
          args:
            - server
          ports:
            - containerPort: 9085
          volumeMounts:
            - name: config
              mountPath: /etc/pipecd-config
              readOnly: true
      volumes:
        - name: config
          secret:
            secretName: canary-by-config-change
`,
			expected: []string{
				"canary-by-config-change",
			},
		},
		{
			name: "multiple secrets",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-by-config-change
  labels:
    app: canary-by-config-change
spec:
  replicas: 2
  selector:
    matchLabels:
      app: canary-by-config-change
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: canary-by-config-change
        pipecd.dev/variant: primary
    spec:
      initContainers:
        - name: init
          image: gcr.io/pipecd/helloworld:v0.5.0
          env:
            - name: env1
              valueFrom:
                secretKeyRef:
                  name: init-secret-1
                  key: key1
      containers:
        - name: helloworld
          image: gcr.io/pipecd/helloworld:v0.5.0
          args:
            - server
          ports:
            - containerPort: 9085
          env:
            - name: env1
              valueFrom:
                secretKeyRef:
                  name: secret-1
                  key: key1
            - name: env2
              valueFrom:
                secretKeyRef:
                  name: secret-2
                  key: key2
          volumeMounts:
            - name: config
              mountPath: /etc/pipecd-config
              readOnly: true
      volumes:
        - name: config
          secret:
            secretName: canary-by-config-change
        - name: config2
          secret:
            secretName: secret-2
`,
			expected: []string{
				"canary-by-config-change",
				"init-secret-1",
				"secret-1",
				"secret-2",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := ParseManifests(tc.manifest)
			require.NoError(t, err)
			require.Equal(t, 1, len(manifests))

			d := &appsv1.Deployment{}
			err = manifests[0].ConvertToStructuredObject(d)
			require.NoError(t, err)

			out := FindReferencingSecretsInDeployment(d)
			assert.Equal(t, tc.expected, out)
		})
	}
}
