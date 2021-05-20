// Copyright 2020 The PipeCD Authors.
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

package grpcapi

import (
	"github.com/pipe-cd/pipe/pkg/app/api/service/webservice"
	"github.com/pipe-cd/pipe/pkg/model"
)

var (
	k8sDeploymentConfigTemplates = []*webservice.DeploymentConfigTemplate{
		{
			ApplicationKind: model.ApplicationKind_KUBERNETES,
			Name:            "Simple",
			Labels:          []webservice.DeploymentConfigTemplateLabel{},
			Content:         DeploymentConfigTemplates["KubernetesSimple"],
		},
		{
			ApplicationKind: model.ApplicationKind_KUBERNETES,
			Name:            "Canary",
			Labels:          []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY},
			Content:         DeploymentConfigTemplates["KubernetesCanary"],
		},
		{
			ApplicationKind: model.ApplicationKind_KUBERNETES,
			Name:            "Blue/Green",
			Labels:          []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
			Content:         DeploymentConfigTemplates["KubernetesBlueGreen"],
		},
		{
			ApplicationKind: model.ApplicationKind_KUBERNETES,
			Name:            "Kustomize",
			Labels:          []webservice.DeploymentConfigTemplateLabel{},
			Content:         DeploymentConfigTemplates["KubernetesKustomize"],
		},
	}

	terraformDeploymentConfigTemplates  = []*webservice.DeploymentConfigTemplate{}
	crossplaneDeploymentConfigTemplates = []*webservice.DeploymentConfigTemplate{}
	lambdaDeploymentConfigTemplates     = []*webservice.DeploymentConfigTemplate{}
	cloudrunDeploymentConfigTemplates   = []*webservice.DeploymentConfigTemplate{}
	ecsDeploymentConfigTemplates        = []*webservice.DeploymentConfigTemplate{}
)
