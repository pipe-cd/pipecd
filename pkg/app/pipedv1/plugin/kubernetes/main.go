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
	"log"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func main() {
	app := cli.NewApp(
		"pipecd-plugin-kubernetes",
		"Plugin component to deploy Kubernetes Application.",
	)
	app.AddCommands(
		NewPluginCommand(),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// TODO: use this after rewriting the plugin with the sdk
func _main() {
	sdk.RegisterDeploymentPlugin[sdk.ConfigNone, kubeconfig.KubernetesDeployTargetConfig](&plugin{})
	if err := sdk.Run(); err != nil {
		log.Fatalln(err)
	}
}
