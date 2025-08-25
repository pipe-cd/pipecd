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

package main

import (
	"context"
	"log"

	"go.uber.org/zap"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/deployment"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/livestate"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
)

type initializer struct{}

func (i *initializer) Initialize(ctx context.Context, input *sdk.InitializeInput[config.KubernetesPluginConfig, config.KubernetesDeployTargetConfig]) error {
	// Initialize the plugin with the given context and input.

	toolregistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	// TODO: get helm version from piped config
	helmPath, err := toolregistry.Helm(ctx, "")
	if err != nil {
		return err
	}

	helm := provider.NewHelm(helmPath, input.Logger)

	// Login to OCI registries
	for _, registry := range input.Config.ChartRegistries {
		if !registry.IsOCI() {
			continue
		}

		if err := helm.LoginToOCIRegistry(ctx, registry.Address, registry.Username, registry.Password); err != nil {
			input.Logger.Error("failed to login to helm oci registry", zap.String("address", registry.Address), zap.Error(err))
			return err
		}
	}

	return nil
}

func main() {
	plugin, err := sdk.NewPlugin(
		"0.0.1",
		sdk.WithInitializer[config.KubernetesApplicationSpec](&initializer{}),
		sdk.WithDeploymentPlugin(&deployment.Plugin{}),
		sdk.WithLivestatePlugin(&livestate.Plugin{}),
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err := plugin.Run(); err != nil {
		log.Fatalln(err)
	}
}
