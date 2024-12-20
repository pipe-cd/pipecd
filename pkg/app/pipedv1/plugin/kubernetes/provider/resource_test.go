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
)

func TestResourceKey_IsDeployment(t *testing.T) {
	tests := []struct {
		name string
		key  ResourceKey
		want bool
	}{
		{
			name: "Deployment with built-in API version",
			key: ResourceKey{
				apiVersion: "apps/v1",
				kind:       KindDeployment,
				namespace:  "default",
				name:       "test-deployment",
			},
			want: true,
		},
		{
			name: "Deployment with non built-in API version",
			key: ResourceKey{
				apiVersion: "custom/v1",
				kind:       KindDeployment,
				namespace:  "default",
				name:       "test-deployment",
			},
			want: false,
		},
		{
			name: "Non-Deployment kind",
			key: ResourceKey{
				apiVersion: "apps/v1",
				kind:       KindConfigMap,
				namespace:  "default",
				name:       "test-configmap",
			},
			want: false,
		},
		{
			name: "Empty ResourceKey",
			key:  ResourceKey{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.key.IsDeployment(); got != tt.want {
				t.Errorf("ResourceKey.IsDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}
