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

package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

func TestCheckVariantSelectorInWorkload(t *testing.T) {
	t.Parallel()

	const (
		variantLabel   = "pipecd.dev/variant"
		primaryVariant = "primary"
	)
	testcases := []struct {
		name      string
		manifest  string
		expectErr bool
	}{
		{
			name: "missing variant in selector",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			expectErr: true,
		},
		{
			name: "missing variant in template labels",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
`,
			expectErr: true,
		},
		{
			name: "wrong variant in selector",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: canary
  template:
    metadata:
      labels:
        app: simple
`,
			expectErr: true,
		},
		{
			name: "wrong variant in temlate labels",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: canary
`,
			expectErr: true,
		},
		{
			name: "ok",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: primary
`,
		},
	}

	expected := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: primary
`
	generatedManifests, err := provider.ParseManifests(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(generatedManifests))

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.ParseManifests(tc.manifest)
			require.NoError(t, err)
			require.Equal(t, 1, len(manifests))

			err = checkVariantSelectorInWorkload(manifests[0], variantLabel, primaryVariant)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			err = ensureVariantSelectorInWorkload(manifests[0], variantLabel, primaryVariant)
			assert.NoError(t, err)
			assert.Equal(t, generatedManifests[0], manifests[0])
		})
	}

}

func TestGenerateVariantServiceManifests(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name         string
		inputYAML    string
		variantLabel string
		variant      string
		nameSuffix   string
		expectYAML   string
	}{
		{
			name: "basic service variant",
			inputYAML: `
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: my-app
  type: NodePort
  ports:
    - port: 80
      targetPort: 8080
  externalIPs:
    - 1.2.3.4
  loadBalancerIP: 5.6.7.8
  loadBalancerSourceRanges:
    - 0.0.0.0/0
`,
			variantLabel: "pipecd.dev/variant",
			variant:      "canary",
			nameSuffix:   "canary",
			expectYAML: `
apiVersion: v1
kind: Service
metadata:
  name: my-service-canary
spec:
  selector:
    app: my-app
    pipecd.dev/variant: canary
  type: ClusterIP
  ports:
    - port: 80
      targetPort: 8080
`,
		},
		{
			name: "service with no selector",
			inputYAML: `
apiVersion: v1
kind: Service
metadata:
  name: test-svc
spec:
  ports:
    - port: 443
      targetPort: 8443
`,
			variantLabel: "pipecd.dev/variant",
			variant:      "primary",
			nameSuffix:   "primary",
			expectYAML: `
apiVersion: v1
kind: Service
metadata:
  name: test-svc-primary
spec:
  selector:
    pipecd.dev/variant: primary
  type: ClusterIP
  ports:
    - port: 443
      targetPort: 8443
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			services, err := provider.ParseManifests(tc.inputYAML)
			require.NoError(t, err)
			got, err := generateVariantServiceManifests(services, tc.variantLabel, tc.variant, tc.nameSuffix)
			require.NoError(t, err)
			expects, err := provider.ParseManifests(tc.expectYAML)
			require.NoError(t, err)
			require.Equal(t, len(expects), len(got))

			for i := range expects {
				var wantSvc, gotSvc corev1.Service
				err := expects[i].ConvertToStructuredObject(&wantSvc)
				require.NoError(t, err)
				err = got[i].ConvertToStructuredObject(&gotSvc)
				require.NoError(t, err)

				assert.Equal(t, wantSvc, gotSvc)
			}
		})
	}
}
