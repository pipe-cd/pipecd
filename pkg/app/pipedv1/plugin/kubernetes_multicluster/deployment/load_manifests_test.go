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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
)

// capturingLoader records the LoaderInput it receives so tests can assert on it.
type capturingLoader struct {
	capturedInput provider.LoaderInput
}

func (c *capturingLoader) LoadManifests(_ context.Context, input provider.LoaderInput) ([]provider.Manifest, error) {
	c.capturedInput = input
	return nil, nil
}

func TestLoadManifests_KustomizeOverrides(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		specInput        kubeconfig.KubernetesDeploymentInput
		multiTarget      *kubeconfig.KubernetesMultiTarget
		wantVersion      string
		wantOptions      map[string]string
	}{
		{
			name: "top-level version and options used when multiTarget has none",
			specInput: kubeconfig.KubernetesDeploymentInput{
				KustomizeVersion: "5.3.0",
				KustomizeOptions: map[string]string{"flag": "val"},
			},
			multiTarget: &kubeconfig.KubernetesMultiTarget{},
			wantVersion: "5.3.0",
			wantOptions: map[string]string{"flag": "val"},
		},
		{
			name: "multiTarget version overrides top-level version",
			specInput: kubeconfig.KubernetesDeploymentInput{
				KustomizeVersion: "5.3.0",
			},
			multiTarget: &kubeconfig.KubernetesMultiTarget{
				KustomizeVersion: "5.4.3",
			},
			wantVersion: "5.4.3",
			wantOptions: nil,
		},
		{
			name: "multiTarget options override top-level options",
			specInput: kubeconfig.KubernetesDeploymentInput{
				KustomizeOptions: map[string]string{"original": "val"},
			},
			multiTarget: &kubeconfig.KubernetesMultiTarget{
				KustomizeOptions: map[string]string{"load-restrictor": "LoadRestrictionsNone"},
			},
			wantVersion: "",
			wantOptions: map[string]string{"load-restrictor": "LoadRestrictionsNone"},
		},
		{
			name: "multiTarget version and options both override top-level",
			specInput: kubeconfig.KubernetesDeploymentInput{
				KustomizeVersion: "5.3.0",
				KustomizeOptions: map[string]string{"flag": "val"},
			},
			multiTarget: &kubeconfig.KubernetesMultiTarget{
				KustomizeVersion: "5.4.3",
				KustomizeOptions: map[string]string{"enable-helm": ""},
			},
			wantVersion: "5.4.3",
			wantOptions: map[string]string{"enable-helm": ""},
		},
		{
			name: "nil multiTarget uses top-level values",
			specInput: kubeconfig.KubernetesDeploymentInput{
				KustomizeVersion: "5.3.0",
				KustomizeOptions: map[string]string{"key": "val"},
			},
			multiTarget: nil,
			wantVersion: "5.3.0",
			wantOptions: map[string]string{"key": "val"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cl := &capturingLoader{}
			p := &Plugin{}

			spec := &kubeconfig.KubernetesApplicationSpec{
				Input: tt.specInput,
			}
			deploy := &sdk.Deployment{PipedID: "piped-id", ApplicationID: "app-id"}
			source := &sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{}

			_, err := p.loadManifests(context.Background(), deploy, spec, source, cl, zaptest.NewLogger(t), tt.multiTarget)
			require.NoError(t, err)

			assert.Equal(t, tt.wantVersion, cl.capturedInput.KustomizeVersion)
			assert.Equal(t, tt.wantOptions, cl.capturedInput.KustomizeOptions)
		})
	}
}
