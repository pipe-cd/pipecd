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

package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestUnmarshalConfig(t *testing.T) {
	testcases := []struct {
		name     string
		data     string
		wantSpec interface{}
		wantErr  bool
	}{
		{
			name: "correct config for KubernetesApp",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "KubernetesApp",
  "spec": {
		"input": {
		  "namespace": "default"
		}
  }
}`,
			wantSpec: &KubernetesApplicationSpec{
				Input: KubernetesDeploymentInput{
					Namespace: "default",
				},
			},
			wantErr: false,
		},
		{
			name: "config for KubernetesApp with unknown field",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "KubernetesApp",
  "spec": {
		"input": {
		  "namespace": "default"
		},
		"unknown": {}
  }
}`,
			wantSpec: &KubernetesApplicationSpec{
				Input: KubernetesDeploymentInput{
					Namespace: "default",
				},
			},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var got Config
			err := json.Unmarshal([]byte(tc.data), &got)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.wantSpec, got.spec)
		})
	}
}

func newBoolPointer(v bool) *bool {
	return &v
}

func TestKind_ToApplicationKind(t *testing.T) {
	testcases := []struct {
		name   string
		k      Kind
		want   model.ApplicationKind
		wantOk bool
	}{
		{
			name:   "App config",
			k:      KindKubernetesApp,
			want:   model.ApplicationKind_KUBERNETES,
			wantOk: true,
		},
		{
			name:   "Not an app config",
			k:      KindPiped,
			want:   model.ApplicationKind_KUBERNETES,
			wantOk: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, gotOk := tc.k.ToApplicationKind()
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantOk, gotOk)
		})
	}
}
