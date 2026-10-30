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

package provider

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
)

func TestMergeVars(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		deployTargetVars []string
		appVars          []string
		want             []string
	}{
		{
			name:             "empty vars",
			deployTargetVars: []string{},
			appVars:          []string{},
			want:             []string{},
		},
		{
			name:             "only deploy target vars",
			deployTargetVars: []string{"key1=value1", "key2=value2"},
			appVars:          []string{},
			want:             []string{"key1=value1", "key2=value2"},
		},
		{
			name:             "only app vars",
			deployTargetVars: []string{},
			appVars:          []string{"key3=value3", "key4=value4"},
			want:             []string{"key3=value3", "key4=value4"},
		},
		{
			name:             "both deploy target and app vars",
			deployTargetVars: []string{"key1=value1", "key2=value2"},
			appVars:          []string{"key3=value3", "key4=value4"},
			want:             []string{"key1=value1", "key2=value2", "key3=value3", "key4=value4"},
		},
		{
			// TODO: Validate duplication ans the want should not contain "key2=valueX"
			name:             "duplicate vars",
			deployTargetVars: []string{"key1=value1", "key2=value2"},
			appVars:          []string{"key2=valueX", "key3=value3"},
			want:             []string{"key1=value1", "key2=value2", "key2=valueX", "key3=value3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := mergeVars(tt.deployTargetVars, tt.appVars)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewTerraformCommand_Options(t *testing.T) {
	t.Parallel()

	appDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(appDir, "main.tf"), []byte(`variable "env" {
  type = string
}
`), 0o644))

	newCommand := func(t *testing.T, opts ...Option) *Terraform {
		t.Helper()
		// Called outside a deployment, so there is no stage log persister.
		client := sdk.NewClient(nil, "terraform", "app-id", "", nil, toolregistrytest.NewTestToolRegistry(t))
		cmd, err := NewTerraformCommand(context.Background(), client, sdk.DeploymentSource[config.ApplicationConfigSpec]{
			ApplicationDirectory: appDir,
			ApplicationConfig: &sdk.ApplicationConfig[config.ApplicationConfigSpec]{
				Spec: &config.ApplicationConfigSpec{TerraformVersion: "1.12.0", Vars: []string{"env=app"}},
			},
		}, &sdk.DeployTarget[config.DeployTargetConfig]{
			Name:   "dt-1",
			Config: config.DeployTargetConfig{Vars: []string{"region=dt"}},
		}, opts...)
		require.NoError(t, err)
		return cmd
	}

	t.Run("without options", func(t *testing.T) {
		t.Parallel()
		cmd := newCommand(t)
		assert.False(t, cmd.options.noColor)
		assert.Equal(t, []string{"region=dt", "env=app"}, cmd.options.vars)
	})

	t.Run("with WithoutColor", func(t *testing.T) {
		t.Parallel()
		cmd := newCommand(t, WithoutColor())
		assert.True(t, cmd.options.noColor)
		assert.Equal(t, []string{"region=dt", "env=app"}, cmd.options.vars)
	})

	t.Run("given options are applied after the application config", func(t *testing.T) {
		t.Parallel()
		cmd := newCommand(t, WithVars([]string{"env=override"}))
		assert.Equal(t, []string{"env=override"}, cmd.options.vars)
	})
}
