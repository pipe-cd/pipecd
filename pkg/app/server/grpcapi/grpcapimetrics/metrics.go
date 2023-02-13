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

package grpcapimetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	projectKey = "project"
)

var (
	deploymentCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpcapi_create_deployment_total",
			Help: "Number of successful CreateDeployment RPC with project label",
		},
		[]string{
			projectKey,
		},
	)
)

func Register(r prometheus.Registerer) {
	r.MustRegister(deploymentCounter)
}

func IncDeploymentCounter(project string) {
	deploymentCounter.With(prometheus.Labels{
		projectKey: project,
	}).Inc()
}
