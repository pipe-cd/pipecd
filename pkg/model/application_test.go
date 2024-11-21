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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestGetKindString(t *testing.T) {
	tests := []struct {
		name     string
		app      *Application
		expected string
	}{
		{
			name:     "Kubernetes Application",
			app:      &Application{Kind: ApplicationKind_KUBERNETES},
			expected: "KUBERNETES",
		},
		{
			name:     "Terraform Application",
			app:      &Application{Kind: ApplicationKind_TERRAFORM},
			expected: "TERRAFORM",
		},
		{
			name:     "Lambda Application",
			app:      &Application{Kind: ApplicationKind_LAMBDA},
			expected: "LAMBDA",
		},
		{
			name:     "CloudRun Application",
			app:      &Application{Kind: ApplicationKind_CLOUDRUN},
			expected: "CLOUDRUN",
		},
		{
			name:     "ECS Application",
			app:      &Application{Kind: ApplicationKind_ECS},
			expected: "ECS",
		},
		{
			name:     "Application",
			app:      &Application{Kind: ApplicationKind_APPLICATION, Labels: map[string]string{"kind": "KUBERNETES"}},
			expected: "KUBERNETES",
		},
		{
			name:     "Application with no kind label",
			app:      &Application{Kind: ApplicationKind_APPLICATION, Labels: map[string]string{}},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.app.GetKindString()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCompatiblePlatformProviderType(t *testing.T) {
	tests := []struct {
		name     string
		kind     ApplicationKind
		expected PlatformProviderType
	}{
		{
			name:     "Kubernetes Application",
			kind:     ApplicationKind_KUBERNETES,
			expected: PlatformProviderKubernetes,
		},
		{
			name:     "Terraform Application",
			kind:     ApplicationKind_TERRAFORM,
			expected: PlatformProviderTerraform,
		},
		{
			name:     "Lambda Application",
			kind:     ApplicationKind_LAMBDA,
			expected: PlatformProviderLambda,
		},
		{
			name:     "CloudRun Application",
			kind:     ApplicationKind_CLOUDRUN,
			expected: PlatformProviderCloudRun,
		},
		{
			name:     "ECS Application",
			kind:     ApplicationKind_ECS,
			expected: PlatformProviderECS,
		},
		{
			name:     "Default Case (assumed non-defined ApplicationKind)",
			kind:     ApplicationKind(9999), // assuming this isn't a defined kind
			expected: PlatformProviderKubernetes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.kind.CompatiblePlatformProviderType()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestIsOutOfSync(t *testing.T) {
	tests := []struct {
		name     string
		app      *Application
		expected bool
	}{
		{
			name:     "Nil SyncState",
			app:      &Application{},
			expected: false,
		},
		{
			name:     "SyncState IN_SYNC",
			app:      &Application{SyncState: &ApplicationSyncState{Status: ApplicationSyncStatus_SYNCED}},
			expected: false,
		},
		{
			name:     "SyncState OUT_OF_SYNC",
			app:      &Application{SyncState: &ApplicationSyncState{Status: ApplicationSyncStatus_OUT_OF_SYNC}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.app.IsOutOfSync()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestHasChanged(t *testing.T) {
	tests := []struct {
		name     string
		current  ApplicationSyncState
		next     ApplicationSyncState
		expected bool
	}{
		{
			name:     "No Change",
			current:  ApplicationSyncState{Status: ApplicationSyncStatus_SYNCED, ShortReason: "B", Reason: "C"},
			next:     ApplicationSyncState{Status: ApplicationSyncStatus_SYNCED, ShortReason: "B", Reason: "C"},
			expected: false,
		},
		{
			name:     "Status Changed",
			current:  ApplicationSyncState{Status: ApplicationSyncStatus_SYNCED, ShortReason: "B", Reason: "C"},
			next:     ApplicationSyncState{Status: ApplicationSyncStatus_DEPLOYING, ShortReason: "B", Reason: "C"},
			expected: true,
		},
		{
			name:     "ShortReason Changed",
			current:  ApplicationSyncState{Status: ApplicationSyncStatus_SYNCED, ShortReason: "B", Reason: "C"},
			next:     ApplicationSyncState{Status: ApplicationSyncStatus_SYNCED, ShortReason: "Y", Reason: "C"},
			expected: true,
		},
		{
			name:     "Reason Changed",
			current:  ApplicationSyncState{Status: ApplicationSyncStatus_SYNCED, ShortReason: "B", Reason: "C"},
			next:     ApplicationSyncState{Status: ApplicationSyncStatus_SYNCED, ShortReason: "B", Reason: "Z"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.current.HasChanged(tt.next)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGetApplicationConfigFilename(t *testing.T) {
	tests := []struct {
		name         string
		gitPath      ApplicationGitPath
		expectedName string
	}{
		{
			name:         "Default Filename",
			gitPath:      ApplicationGitPath{ConfigFilename: ""},
			expectedName: oldDefaultApplicationConfigFilename,
		},
		{
			name:         "Custom Filename",
			gitPath:      ApplicationGitPath{ConfigFilename: "customConfig.yaml"},
			expectedName: "customConfig.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualName := tt.gitPath.GetApplicationConfigFilename()
			assert.Equal(t, tt.expectedName, actualName)
		})
	}
}
