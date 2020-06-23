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

package kubernetes

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricsLabelTool    = "tool"
	metricsLabelVersion = "version"
	metricsLabelCommand = "command"
	metricsLabelStatus  = "status"

	metricsValueKubectl = "kubectl"
	metricsValueSuccess = "success"
	metricsValueFailure = "failure"
)

var (
	metricsToolCalls = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudprovider_kubernetes_tool_calls_total",
			Help: "Number of calls made to run the tool like kubectl, kustomize.",
		},
		[]string{
			metricsLabelTool,
			metricsLabelVersion,
			metricsLabelCommand,
			metricsLabelStatus,
		},
	)
)

func metricsKubectlCalled(version, command string, success bool) {
	status := metricsValueSuccess
	if !success {
		status = metricsValueFailure
	}
	metricsToolCalls.With(prometheus.Labels{
		metricsLabelTool:    metricsValueKubectl,
		metricsLabelVersion: version,
		metricsLabelCommand: command,
		metricsLabelStatus:  status,
	}).Inc()
}

func registerMetrics() {
	prometheus.MustRegister(metricsToolCalls)
}
