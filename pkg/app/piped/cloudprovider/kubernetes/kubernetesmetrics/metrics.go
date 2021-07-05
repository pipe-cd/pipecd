// Copyright 2021 The PipeCD Authors.
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

package kubernetesmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricsLabelTool    = "tool"
	metricsLabelVersion = "version"
	metricsLabelCommand = "command"
	metricsLabelStatus  = "status"
)

type Tool string

const (
	LabelToolKubectl Tool = "kubectl"
)

type ToolCommand string

const (
	LabelApplyCommand  ToolCommand = "apply"
	LabelDeleteCommand ToolCommand = "delete"
)

type CommandOutput string

const (
	LabelOutputSuccess CommandOutput = "success"
	LabelOutputFailre  CommandOutput = "failure"
)

var (
	toolCallsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "piped_cloudprovider_kubernetes_tool_calls_total",
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

func IncKubectlCallsCounter(version string, command ToolCommand, success bool) {
	status := LabelOutputSuccess
	if !success {
		status = LabelOutputFailre
	}
	toolCallsCounter.With(prometheus.Labels{
		metricsLabelTool:    string(LabelToolKubectl),
		metricsLabelVersion: version,
		metricsLabelCommand: string(command),
		metricsLabelStatus:  string(status),
	}).Inc()
}

func Register(r prometheus.Registerer) {
	r.MustRegister(toolCallsCounter)
}
