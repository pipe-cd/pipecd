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

package kubernetesmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	toolKey          = "tool"
	versionKey       = "version"
	toolCommandKey   = "command"
	commandOutputKey = "status"
)

type Tool string

const (
	LabelToolKubectl Tool = "kubectl"
)

type ToolCommand string

const (
	LabelApplyCommand   ToolCommand = "apply"
	LabelCreateCommand  ToolCommand = "create"
	LabelReplaceCommand ToolCommand = "replace"
	LabelDeleteCommand  ToolCommand = "delete"
	LabelGetCommand     ToolCommand = "get"
)

type CommandOutput string

const (
	LabelOutputSuccess CommandOutput = "success"
	LabelOutputFailre  CommandOutput = "failure"
)

var (
	toolCallsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudprovider_kubernetes_tool_calls_total",
			Help: "Number of calls made to run the tool like kubectl, kustomize.",
		},
		[]string{
			toolKey,
			versionKey,
			toolCommandKey,
			commandOutputKey,
		},
	)
)

func IncKubectlCallsCounter(version string, command ToolCommand, success bool) {
	status := LabelOutputSuccess
	if !success {
		status = LabelOutputFailre
	}
	toolCallsCounter.With(prometheus.Labels{
		toolKey:          string(LabelToolKubectl),
		versionKey:       version,
		toolCommandKey:   string(command),
		commandOutputKey: string(status),
	}).Inc()
}

func Register(r prometheus.Registerer) {
	r.MustRegister(toolCallsCounter)
}
