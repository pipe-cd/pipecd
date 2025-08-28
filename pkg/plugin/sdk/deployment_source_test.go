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

package sdk

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

func TestApplicationConfig_HasStage(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		appConfig *ApplicationConfig[any]
		stageName string
		want      bool
	}{
		{
			name: "stage exists",
			appConfig: &ApplicationConfig[any]{
				commonSpec: &config.GenericApplicationSpec{
					Pipeline: &config.DeploymentPipeline{
						Stages: []config.PipelineStage{
							{Name: "stage1"},
						},
					},
				},
			},
			stageName: "stage1",
			want:      true,
		},
		{
			name: "stage exists with different name",
			appConfig: &ApplicationConfig[any]{
				commonSpec: &config.GenericApplicationSpec{
					Pipeline: &config.DeploymentPipeline{
						Stages: []config.PipelineStage{{Name: "stage1"}},
					},
				},
			},
			stageName: "stage2",
			want:      false,
		},
		{
			name: "multiple stages",
			appConfig: &ApplicationConfig[any]{
				commonSpec: &config.GenericApplicationSpec{
					Pipeline: &config.DeploymentPipeline{
						Stages: []config.PipelineStage{{Name: "stage1"}, {Name: "stage2"}},
					},
				},
			},
			stageName: "stage2",
			want:      true,
		},
		{
			name: "stage does not exist",
			appConfig: &ApplicationConfig[any]{
				commonSpec: &config.GenericApplicationSpec{
					Pipeline: &config.DeploymentPipeline{
						Stages: []config.PipelineStage{},
					},
				},
			},
			stageName: "stage1",
			want:      false,
		},
		{
			name: "pipeline is nil",
			appConfig: &ApplicationConfig[any]{
				commonSpec: &config.GenericApplicationSpec{},
			},
			stageName: "stage1",
			want:      false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := c.appConfig.HasStage(c.stageName)
			if got != c.want {
				t.Errorf("HasStage(%q) = %v, want %v", c.stageName, got, c.want)
			}
		})
	}
}

type testPluginSpec struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (s *testPluginSpec) Validate() error {
	if s.Value < 0 {
		return fmt.Errorf("value must not be a negative value")
	}
	return nil
}

func TestApplicationConfig_ParsePluginConfig(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		pluginName string
		config     *ApplicationConfig[testPluginSpec]
		wantSpec   *testPluginSpec
		wantErr    bool
	}{
		{
			name:       "no plugin config present",
			pluginName: "test-plugin",
			config: &ApplicationConfig[testPluginSpec]{
				pluginConfigs: nil,
			},
			wantSpec: &testPluginSpec{
				Name:  "",
				Value: 0,
			},
			wantErr: false,
		},
		{
			name:       "empty plugin configs map",
			pluginName: "test-plugin",
			config: &ApplicationConfig[testPluginSpec]{
				pluginConfigs: make(map[string]json.RawMessage),
			},
			wantSpec: &testPluginSpec{
				Name:  "",
				Value: 0,
			},
			wantErr: false,
		},
		{
			name:       "invalid json in plugin config",
			pluginName: "test-plugin",
			config: &ApplicationConfig[testPluginSpec]{
				pluginConfigs: map[string]json.RawMessage{
					"test-plugin": json.RawMessage(`{invalid-json`),
				},
			},
			wantSpec: nil,
			wantErr:  true,
		},
		{
			name:       "validation failure",
			pluginName: "test-plugin",
			config: &ApplicationConfig[testPluginSpec]{
				pluginConfigs: map[string]json.RawMessage{
					"test-plugin": json.RawMessage(`{"name": "test", "value": -1}`),
				},
			},
			wantSpec: nil,
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.config.parsePluginConfig(tc.pluginName)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantSpec, tc.config.Spec)
			assert.Nil(t, tc.config.pluginConfigs, "pluginConfigs should be cleared after successful parsing")
		})
	}
}
