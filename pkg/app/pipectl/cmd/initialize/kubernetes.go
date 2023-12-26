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

package initialize

import (
	"io"

	"github.com/pipe-cd/pipecd/pkg/config"
)

type genericKubernetesApplicationSpec struct {
	Name           string                           `json:"name"`
	Input          config.KubernetesDeploymentInput `json:"input,omitempty"`
	QuickSync      *config.K8sSyncStageOptions      `json:"quickSync,omitempty"`
	Service        config.K8sResourceReference      `json:"service,omitempty"`
	Workloads      []config.K8sResourceReference    `json:"workloads,omitempty"`
	TrafficRouting *config.KubernetesTrafficRouting `json:"trafficRouting,omitempty"`
	VariantLabel   config.KubernetesVariantLabel    `json:"variantLabel,omitempty"`
	ResourceRoutes []config.KubernetesResourceRoute `json:"resourceRoutes,omitempty"`
	Description    string                           `json:"description,omitempty"`
}

func generateKubernetesConfig(in io.Reader) (*genericConfig, error) {
	spec, e := generateKubernetesSpec(in)
	if e != nil {
		return nil, e
	}

	return &genericConfig{
		Kind:            config.KindKubernetesApp,
		APIVersion:      config.VersionV1Beta1,
		ApplicationSpec: spec,
	}, nil
}

func generateKubernetesSpec(in io.Reader) (*genericKubernetesApplicationSpec, error) {
	appName := promptStringRequired("Name of the application: ", in)

	cfg := &genericKubernetesApplicationSpec{
		Name:  appName,
		Input: config.KubernetesDeploymentInput{},
		// QuickSync:      &config.K8sSyncStageOptions{},
		// Service:        config.K8sResourceReference{},
		// Workloads:      []config.K8sResourceReference{},
		// TrafficRouting: &config.KubernetesTrafficRouting{},
		// VariantLabel:   config.KubernetesVariantLabel{},
		// ResourceRoutes: []config.KubernetesResourceRoute{},
		Description: "Generated by `pipectl init`. See https://pipecd.dev/docs/user-guide/configuration-reference/ for more.",
	}

	return cfg, nil
}
