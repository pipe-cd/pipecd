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

package appconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineIndent(t *testing.T) {
	testcases := []struct {
		name   string
		line   string
		indent string
		found  bool
	}{
		{
			name:   "empty string",
			line:   "",
			indent: "",
			found:  false,
		},
		{
			name: "contains only space",
			line: " 	",
			indent: "",
			found:  false,
		},
		{
			name:   "comment",
			line:   "  # This is a comment",
			indent: "",
			found:  false,
		},
		{
			name:   "indent just contains spaces",
			line:   "  foo:",
			indent: "  ",
			found:  true,
		},
		{
			name: "indent just contains tab",
			line: "	foo:",
			indent: "	",
			found: true,
		},
		{
			name: "indent contains both space and tab",
			line: "  	foo:",
			indent: "  	",
			found: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, found := determineIndent(tc.line)
			assert.Equal(t, tc.indent, got)
			assert.Equal(t, tc.found, found)
		})
	}
}

func TestConvert(t *testing.T) {
	testcases := []struct {
		name           string
		deployConfig   string
		appName        string
		appEnv         string
		appDescription string
		appConfig      string
	}{
		{
			name: "has description",
			deployConfig: `
# Deploy plain-yaml manifests in the application directory
# without using pipeline.
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  input:
    manifests:
      - deployment.yaml
      - service.yaml
    kubectlVersion: 1.18.5
`,
			appName:        "test-name",
			appEnv:         "test-env",
			appDescription: "test-description",
			appConfig: `
# Deploy plain-yaml manifests in the application directory
# without using pipeline.
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: test-name
  labels:
    env: test-env
  description: |
    test-description
  input:
    manifests:
      - deployment.yaml
      - service.yaml
    kubectlVersion: 1.18.5
`,
		},
		{
			name: "no description",
			deployConfig: `
# Deploy plain-yaml manifests in the application directory
# without using pipeline.
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  input:
    manifests:
      - deployment.yaml
      - service.yaml
    kubectlVersion: 1.18.5
`,
			appName: "test-name",
			appEnv:  "test-env",
			appConfig: `
# Deploy plain-yaml manifests in the application directory
# without using pipeline.
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: test-name
  labels:
    env: test-env
  input:
    manifests:
      - deployment.yaml
      - service.yaml
    kubectlVersion: 1.18.5
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convert([]byte(tc.deployConfig), tc.appName, tc.appEnv, tc.appDescription)
			assert.Equal(t, tc.appConfig, string(got))
			require.NoError(t, err)
		})
	}
}
