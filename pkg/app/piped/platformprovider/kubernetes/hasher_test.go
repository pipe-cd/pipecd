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

package kubernetes

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashManifests(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		manifests     string
		expected      string
		expectedError error
	}{
		{
			name:          "no manifests",
			expectedError: errors.New("no manifest to hash"),
		},
		{
			name: "configmap: emptydata",
			manifests: `
apiVersion: v1
kind: ConfigMap
data: {}
binaryData: {}
`,
			expected: "42745tchd9",
		},
		{
			name: "configmap: one key",
			manifests: `
apiVersion: v1
kind: ConfigMap
data:
  one: ""
binaryData: {}
`,
			expected: "9g67k2htb6",
		},
		{
			name: "configmap: there keys for checking order",
			manifests: `
apiVersion: v1
kind: ConfigMap
data:
  two: "2"
  one: ""
  three: "3"
binaryData: {}
`,
			expected: "f5h7t85m9b",
		},
		{
			name: "secret: emptydata",
			manifests: `
apiVersion: v1
kind: Secret
type: my-type
data: {}
`,
			expected: "t75bgf6ctb",
		},
		{
			name: "secret: one key",
			manifests: `
apiVersion: v1
kind: Secret
type: my-type
data:
  "one": ""
`,
			expected: "74bd68bm66",
		},
		{
			name: "secret: there keys for checking order",
			manifests: `
apiVersion: v1
kind: Secret
type: my-type
data:
  two: Mg==
  one: ""
  three: Mw==
`,
			expected: "dgcb6h9tmk",
		},
		{
			name: "multiple configs",
			manifests: `
apiVersion: v1
kind: ConfigMap
data:
  two: "2"
  three: "3"
binaryData: {}
---
apiVersion: v1
kind: Secret
type: my-type
data:
  one: ""
  three: Mw==
`,
			expected: "57hhd7795k",
		},
		{
			name: "not config manifest",
			manifests: `
apiVersion: apps/v1
kind: Foo
metadata:
  name: simple
  labels:
    app: simple
    pipecd.dev/managed-by: piped
spec:
  replicas: 2
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
        component: foo
    spec:
      containers:
      - name: helloworld
        image: gcr.io/pipecd/helloworld:v1.0.0
        args:
          - hi
          - hello
        ports:
        - containerPort: 9085
`,
			expected: "db48kd6689",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := ParseManifests(tc.manifests)
			require.NoError(t, err)

			out, err := HashManifests(manifests)
			assert.Equal(t, tc.expected, out)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
