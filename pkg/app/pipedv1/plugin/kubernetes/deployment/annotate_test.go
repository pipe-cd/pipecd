// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

func TestAnnotateConfigHash(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		manifests     string
		expected      string
		expectedError error
	}{
		{
			name: "empty list",
		},
		{
			name: "one config",
			manifests: `
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
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: canary-by-config-change
data:
  two: "2"
`,
			expected: `
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
      annotations:
        pipecd.dev/config-hash: 75c9m2btb6
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
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: canary-by-config-change
data:
  two: "2"
`,
		},
		{
			name: "multiple configs",
			manifests: `
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
        - name: secret
          secret:
            secretName: secret-1
        - name: unmanaged-config
          configMap:
            name: unmanaged-config
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: canary-by-config-change
data:
  two: "2"
---
apiVersion: v1
kind: Secret
metadata:
  name: secret-1
type: my-type
data:
  "one": "Mg=="
`,
			expected: `
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
      annotations:
        pipecd.dev/config-hash: t7dtkdm455
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
        - name: secret
          secret:
            secretName: secret-1
        - name: unmanaged-config
          configMap:
            name: unmanaged-config
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: canary-by-config-change
data:
  two: "2"
---
apiVersion: v1
kind: Secret
metadata:
  name: secret-1
type: my-type
data:
  "one": "Mg=="
`,
		},
		{
			name: "one secret",
			manifests: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-by-secret-change
  labels:
    app: canary-by-secret-change
spec:
  replicas: 2
  selector:
    matchLabels:
      app: canary-by-secret-change
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: canary-by-secret-change
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
        - name: secret
          mountPath: /etc/pipecd-secret
          readOnly: true
      volumes:
      - name: secret
        secret:
          secretName: canary-by-secret-change
---
apiVersion: v1
kind: Secret
metadata:
  name: canary-by-secret-change
type: Opaque
data:
  one: "MQ=="
`,
			expected: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-by-secret-change
  labels:
    app: canary-by-secret-change
spec:
  replicas: 2
  selector:
    matchLabels:
      app: canary-by-secret-change
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: canary-by-secret-change
        pipecd.dev/variant: primary
      annotations:
        pipecd.dev/config-hash: t58h88cd4b
    spec:
      containers:
      - name: helloworld
        image: gcr.io/pipecd/helloworld:v0.5.0
        args:
        - server
        ports:
        - containerPort: 9085
        volumeMounts:
        - name: secret
          mountPath: /etc/pipecd-secret
          readOnly: true
      volumes:
      - name: secret
        secret:
          secretName: canary-by-secret-change
---
apiVersion: v1
kind: Secret
metadata:
  name: canary-by-secret-change
type: Opaque
data:
  one: "MQ=="
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.ParseManifests(tc.manifests)
			require.NoError(t, err)

			expected, err := provider.ParseManifests(tc.expected)
			require.NoError(t, err)

			err = annotateConfigHash(manifests)
			assert.Equal(t, expected, manifests)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
