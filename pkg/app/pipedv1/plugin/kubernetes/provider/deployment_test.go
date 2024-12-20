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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNestedStringSlice(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		obj      any
		fields   []string
		expected []string
	}{
		{
			name:     "simple string slice",
			obj:      []string{"a", "b", "c"},
			fields:   []string{},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "nested string slice",
			obj:      map[string]any{"key": []string{"a", "b", "c"}},
			fields:   []string{"key"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "nested string slice with map",
			obj:      map[string]any{"key": map[string]any{"innerKey": []string{"a", "b", "c"}}},
			fields:   []string{"key", "innerKey"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "nested string slice with mixed types",
			obj:      map[string]any{"key": []any{"a", "b", "c"}},
			fields:   []string{"key"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "nested string slice with non-string types",
			obj:      map[string]any{"key": []any{1, 2, 3}},
			fields:   []string{"key"},
			expected: nil,
		},
		{
			name:     "nested string slice with missing field",
			obj:      map[string]any{"key": []string{"a", "b", "c"}},
			fields:   []string{"missingKey"},
			expected: nil,
		},
		{
			name:     "nested string slice with empty fields",
			obj:      "singleString",
			fields:   []string{},
			expected: []string{"singleString"},
		},
		{
			name:     "nested string slice with nil object",
			obj:      nil,
			fields:   []string{"key"},
			expected: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := nestedStringSlice(tc.obj, tc.fields...)
			assert.Equal(t, tc.expected, out)
		})
	}
}

func TestFindReferencingConfigMaps(t *testing.T) {
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
			t.Parallel()

			manifests, err := ParseManifests(tc.manifest)
			require.NoError(t, err)
			require.Equal(t, 1, len(manifests))

			out := FindReferencingConfigMaps(manifests[0])
			assert.Equal(t, tc.expected, out)
		})
	}
}

func TestFindReferencingSecrets(t *testing.T) {
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
			t.Parallel()

			manifests, err := ParseManifests(tc.manifest)
			require.NoError(t, err)
			require.Equal(t, 1, len(manifests))

			out := FindReferencingSecrets(manifests[0])
			assert.Equal(t, tc.expected, out)
		})
	}
}
func TestFindContainerImages(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		manifest string
		expected []string
	}{
		{
			name: "no container image",
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
          args:
            - server
          ports:
            - containerPort: 9085
`,
			expected: nil,
		},
		{
			name: "one container image",
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
    spec:
      containers:
        - name: helloworld
          image: gcr.io/pipecd/helloworld:v0.5.0
          args:
            - server
          ports:
            - containerPort: 9085
`,
			expected: []string{"gcr.io/pipecd/helloworld:v0.5.0"},
		},
		{
			name: "multiple container images",
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
    spec:
      containers:
        - name: helloworld
          image: gcr.io/pipecd/helloworld:v0.5.0
          args:
            - server
          ports:
            - containerPort: 9085
        - name: sidecar
          image: gcr.io/pipecd/sidecar:v0.5.0
          args:
            - proxy
          ports:
            - containerPort: 9086
`,
			expected: []string{"gcr.io/pipecd/helloworld:v0.5.0", "gcr.io/pipecd/sidecar:v0.5.0"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			manifests, err := ParseManifests(tc.manifest)
			require.NoError(t, err)
			require.Equal(t, 1, len(manifests))

			out := FindContainerImages(manifests[0])
			assert.Equal(t, tc.expected, out)
		})
	}
}
