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

package controllermetrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	deploymentIDKey     = "deployment"
	applicationIDKey    = "application_id"
	applicationNameKey  = "application_name"
	applicationKindKey  = "application_kind"
	platformProviderKey = "platform_provider"
	deploymentStatusKey = "status"
)

var (
	deploymentStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deployment_status",
			Help: "The current status of deployment. 1 for current status, 0 for others.",
		},
		[]string{deploymentIDKey, applicationIDKey, applicationNameKey, applicationKindKey, platformProviderKey, deploymentStatusKey},
	)
)

func UpdateDeploymentStatus(d *model.Deployment, status model.DeploymentStatus) {
	for name, value := range model.DeploymentStatus_value {
		if model.DeploymentStatus(value) == status {
			deploymentStatus.WithLabelValues(d.Id, d.ApplicationId, d.ApplicationName, d.Kind.String(), d.PlatformProvider, name).Set(1)
		} else {
			deploymentStatus.WithLabelValues(d.Id, d.ApplicationId, d.ApplicationName, d.Kind.String(), d.PlatformProvider, name).Set(0)
		}
	}
}

func Register(r prometheus.Registerer) {
	r.MustRegister(
		deploymentStatus,
	)
}
