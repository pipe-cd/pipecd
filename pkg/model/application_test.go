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
