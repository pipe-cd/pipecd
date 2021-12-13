// Copyright 2021 The PipeCD Authors.
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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeApplicationURL(t *testing.T) {
	testcases := []struct {
		name          string
		baseURL       string
		applicationID string
		expected      string
	}{
		{
			name:          "baseURL has no suffix",
			baseURL:       "https://pipecd.dev",
			applicationID: "app-1",
			expected:      "https://pipecd.dev/applications/app-1",
		},
		{
			name:          "baseURL suffixed by /",
			baseURL:       "https://pipecd.dev/",
			applicationID: "app-2",
			expected:      "https://pipecd.dev/applications/app-2",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := MakeApplicationURL(tc.baseURL, tc.applicationID)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestApplication_ContainLabels(t *testing.T) {
	testcases := []struct {
		name   string
		app    *Application
		labels map[string]string
		want   bool
	}{
		{
			name: "all given tags aren't contained",
			app:  &Application{Labels: map[string]string{"key1": "value1"}},
			labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			want: false,
		},
		{
			name: "a label is contained",
			app: &Application{Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			}},
			labels: map[string]string{
				"key1": "value1",
			},
			want: true,
		},
		{
			name: "all tags are contained",
			app: &Application{Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			}},
			labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			want: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.app.ContainLabels(tc.labels)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestToApplicationKind(t *testing.T) {
	testcases := []struct {
		name         string
		kind         string
		expectedKind ApplicationKind
		valid        bool
	}{
		{
			name:         "KubernetesApp passed",
			kind:         "KubernetesApp",
			expectedKind: ApplicationKind_KUBERNETES,
			valid:        true,
		},
		{
			name:         "Kubernetes passed",
			kind:         "Kubernetes",
			expectedKind: ApplicationKind_KUBERNETES,
			valid:        true,
		},
		{
			name:         "TerraformApp passed",
			kind:         "TerraformApp",
			expectedKind: ApplicationKind_TERRAFORM,
			valid:        true,
		},
		{
			name:         "CloudrunApp passed",
			kind:         "CloudrunApp",
			expectedKind: ApplicationKind_CLOUDRUN,
			valid:        true,
		},
		{
			name:         "LambdaApp passed",
			kind:         "LambdaApp",
			expectedKind: ApplicationKind_LAMBDA,
			valid:        true,
		},
		{
			name:         "ECSApp passed",
			kind:         "ECSApp",
			expectedKind: ApplicationKind_ECS,
			valid:        true,
		},
		{
			name:         "Invalid app kind passed",
			kind:         "Kubernetest",
			expectedKind: -1,
			valid:        false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			oKind, valid := ToApplicationKind(tc.kind)
			require.Equal(t, tc.valid, valid)
			assert.Equal(t, tc.expectedKind, oKind)
		})
	}
}
