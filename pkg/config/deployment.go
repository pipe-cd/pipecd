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

package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pipe-cd/pipe/pkg/model"
)

type Pipelineable interface {
	GetStage(index int32) (PipelineStage, bool)
}

// KubernetesDeploymentSpec represents a deployment configuration for Kubernetes application.
type KubernetesDeploymentSpec struct {
	Input         KubernetesDeploymentInput `json:"input"`
	CommitMatcher *DeploymentCommitMatcher  `json:"commitMatcher"`
	Pipeline      *DeploymentPipeline       `json:"pipeline"`

	PrimaryVariant  *PrimaryVariant  `json:"primaryVariant"`
	CanaryVariant   *CanaryVariant   `json:"canaryVariant"`
	BaselineVariant *BaselineVariant `json:"baselineVariant"`
	TrafficRouting  *TrafficRouting  `json:"trafficRouting"`
}

func (s *KubernetesDeploymentSpec) GetStage(index int32) (PipelineStage, bool) {
	if s.Pipeline == nil {
		return PipelineStage{}, false
	}
	if int(index) >= len(s.Pipeline.Stages) {
		return PipelineStage{}, false
	}
	return s.Pipeline.Stages[index], true
}

// Validate returns an error if any wrong configuration value was found.
func (s *KubernetesDeploymentSpec) Validate() error {
	return nil
}

// TerraformDeploymentSpec represents a deployment configuration for Terraform application.
type TerraformDeploymentSpec struct {
	Input    TerraformDeploymentInput `json:"input"`
	Pipeline *DeploymentPipeline      `json:"pipeline"`
}

func (s *TerraformDeploymentSpec) GetStage(index int32) (PipelineStage, bool) {
	if s.Pipeline == nil {
		return PipelineStage{}, false
	}
	if int(index) >= len(s.Pipeline.Stages) {
		return PipelineStage{}, false
	}
	return s.Pipeline.Stages[index], true
}

// Validate returns an error if any wrong configuration value was found.
func (s *TerraformDeploymentSpec) Validate() error {
	return nil
}

// DeploymentCommitMatcher provides a way to decide how to deploy.
type DeploymentCommitMatcher struct {
	// It makes sure to perform syncing if the commit message matches this regular expression.
	Sync string `json:"sync"`
	// It makes sure to perform pipeline if the commit message matches this regular expression.
	Pipeline string `json:"pipeline"`
}

// DeploymentPipeline represents the way to deploy the application.
// The pipeline is triggered by changes in any of the following objects:
// - Target PodSpec (Target can be Deployment, DaemonSet, StatefullSet)
// - ConfigMaps, Secrets that are mounted as volumes or envs in the deployment.
type DeploymentPipeline struct {
	Stages []PipelineStage `json:"stages"`
}

type PrimaryVariant struct {
	// Suffix that should be used when naming the PRIMARY variant's resources.
	// Default is "primary".
	Suffix  string            `json:"suffix"`
	Service K8sVariantService `json:"service"`
}

type CanaryVariant struct {
	// Suffix that should be used when naming the CANARY variant's resources.
	// Default is "canary".
	Suffix   string               `json:"suffix"`
	Service  K8sVariantService    `json:"service"`
	Workload K8sResourceReference `json:"workload"`
}

type BaselineVariant struct {
	// Suffix that should be used when naming the BASELINE variant's resources.
	// Default is "baseline".
	Suffix   string               `json:"suffix"`
	Service  K8sVariantService    `json:"service"`
	Workload K8sResourceReference `json:"workload"`
}

type K8sVariantService struct {
	K8sResourceReference
	Create bool `json:"create"`
}

type TrafficRoutingMethod string

const (
	TrafficRoutingMethodPod   TrafficRoutingMethod = "pod"
	TrafficRoutingMethodIstio TrafficRoutingMethod = "istio"
)

type TrafficRouting struct {
	Method TrafficRoutingMethod `json:"method"`
	Pod    *PodTrafficRouting   `json:"pod"`
	Istio  *IstioTrafficRouting `json:"istio"`
}

type PodTrafficRouting struct {
	Service K8sResourceReference `json:"service"`
}

type IstioTrafficRouting struct {
	EditableRoutes []string `json:"editableRoutes"`
	// TODO: Add a validate to ensure this was configured or using the default value by service name.
	Host           string               `json:"host"`
	VirtualService K8sResourceReference `json:"virtualService"`
}

type K8sResourceReference struct {
	Reference string `json:"reference"`
}

// ParseVariantResourceReference parses the given reference name
// and returns the resource Kind, Name.
// If the reference is malformed, empty kind, empty name and false will be returned.
// Reference format:
// - kind/name
// - name
func ParseVariantResourceReference(ref string) (kind, name string, ok bool) {
	parts := strings.Split(ref, "/")
	if len(parts) == 1 {
		return "", parts[0], true
	}
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}
	return "", "", false
}

// PiplineStage represents a single stage of a pipeline.
// This is used as a generic struct for all stage type.
type PipelineStage struct {
	Id      string
	Name    model.Stage
	Desc    string
	Timeout Duration

	WaitStageOptions               *WaitStageOptions
	WaitApprovalStageOptions       *WaitApprovalStageOptions
	AnalysisStageOptions           *AnalysisStageOptions
	K8sPrimaryRolloutStageOptions  *K8sPrimaryRolloutStageOptions
	K8sCanaryRolloutStageOptions   *K8sCanaryRolloutStageOptions
	K8sCanaryCleanStageOptions     *K8sCanaryCleanStageOptions
	K8sBaselineRolloutStageOptions *K8sBaselineRolloutStageOptions
	K8sBaselineCleanStageOptions   *K8sBaselineCleanStageOptions
	K8sTrafficRoutingStageOptions  *K8sTrafficRoutingStageOptions
	TerraformPlanStageOptions      *TerraformPlanStageOptions
	TerraformApplyStageOptions     *TerraformApplyStageOptions
}

type genericPipelineStage struct {
	Id      string          `json:"id"`
	Name    model.Stage     `json:"name"`
	Desc    string          `json:"desc,omitempty"`
	Timeout Duration        `json:"timeout"`
	With    json.RawMessage `json:"with"`
}

func (s *PipelineStage) UnmarshalJSON(data []byte) error {
	var err error
	gs := genericPipelineStage{}
	if err = json.Unmarshal(data, &gs); err != nil {
		return err
	}
	s.Id = gs.Id
	s.Name = gs.Name
	s.Desc = gs.Desc
	s.Timeout = gs.Timeout

	switch s.Name {
	case model.StageWait:
		s.WaitStageOptions = &WaitStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.WaitStageOptions)
		}
	case model.StageWaitApproval:
		s.WaitApprovalStageOptions = &WaitApprovalStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.WaitApprovalStageOptions)
		}
	case model.StageAnalysis:
		s.AnalysisStageOptions = &AnalysisStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.AnalysisStageOptions)
		}
	case model.StageK8sPrimaryRollout:
		s.K8sPrimaryRolloutStageOptions = &K8sPrimaryRolloutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sPrimaryRolloutStageOptions)
		}
	case model.StageK8sCanaryRollout:
		s.K8sCanaryRolloutStageOptions = &K8sCanaryRolloutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sCanaryRolloutStageOptions)
		}
	case model.StageK8sCanaryClean:
		s.K8sCanaryCleanStageOptions = &K8sCanaryCleanStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sCanaryCleanStageOptions)
		}
	case model.StageK8sBaselineRollout:
		s.K8sBaselineRolloutStageOptions = &K8sBaselineRolloutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sBaselineRolloutStageOptions)
		}
	case model.StageK8sBaselineClean:
		s.K8sBaselineCleanStageOptions = &K8sBaselineCleanStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sBaselineCleanStageOptions)
		}
	case model.StageK8sTrafficRouting:
		s.K8sTrafficRoutingStageOptions = &K8sTrafficRoutingStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sTrafficRoutingStageOptions)
		}
	case model.StageTerraformPlan:
		s.TerraformPlanStageOptions = &TerraformPlanStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.TerraformPlanStageOptions)
		}
	case model.StageTerraformApply:
		s.TerraformApplyStageOptions = &TerraformApplyStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.TerraformApplyStageOptions)
		}
	default:
		err = fmt.Errorf("unsupported stage name: %s", s.Name)
	}
	return err
}

// WaitStageOptions contains all configurable values for a WAIT stage.
type WaitStageOptions struct {
	Duration Duration `json:"duration"`
}

// WaitStageOptions contains all configurable values for a WAIT_APPROVAL stage.
type WaitApprovalStageOptions struct {
	Approvers []string `json:"approvers"`
}

// K8sPrimaryRolloutStageOptions contains all configurable values for a K8S_PRIMARY_ROLLOUT stage.
type K8sPrimaryRolloutStageOptions struct {
	Manifests []string `json:"manifests"`
	// Whether the resources that are no longer defined in Git will be removed.
	Prune bool `json:"prune"`
}

// K8sCanaryRolloutStageOptions contains all configurable values for a K8S_CANARY_ROLLOUT stage.
type K8sCanaryRolloutStageOptions struct {
	// How many pods for CANARY workloads.
	// An integer value can be specified to indicate an absolute value of pod number.
	// Or a string suffixed by "%" to indicate an percentage value compared to the pod number of PRIMARY.
	// Default is 1 pod.
	Replicas Replicas `json:"replicas"`
}

// K8sCanaryCleanStageOptions contains all configurable values for a K8S_CANARY_CLEAN stage.
type K8sCanaryCleanStageOptions struct {
}

// K8sBaselineRolloutStageOptions contains all configurable values for a K8S_BASELINE_ROLLOUT stage.
type K8sBaselineRolloutStageOptions struct {
	// How many pods for BASELINE workloads.
	// An integer value can be specified to indicate an absolute value of pod number.
	// Or a string suffixed by "%" to indicate an percentage value compared to the pod number of PRIMARY.
	// Default is 1 pod.
	Replicas Replicas `json:"replicas"`
}

// K8sBaselineCleanStageOptions contains all configurable values for a K8S_BASELINE_CLEAN stage.
type K8sBaselineCleanStageOptions struct {
}

// K8sTrafficRoutingStageOptions contains all configurable values for a K8S_TRAFFIC_ROUTING stage.
type K8sTrafficRoutingStageOptions struct {
	// Which variant should receive all traffic.
	All string `json:"all"`
	// The percentage of traffic should be routed to PRIMARY variant.
	Primary int `json:"primary"`
	// The percentage of traffic should be routed to CANARY variant.
	Canary int `json:"canary"`
	// The percentage of traffic should be routed to BASELINE variant.
	Baseline int `json:"baseline"`
}

func (opts K8sTrafficRoutingStageOptions) Percentages() (primary, canary, baseline int) {
	switch opts.All {
	case "primary":
		return 100, 0, 0
	case "canary":
		return 0, 100, 0
	case "baseline":
		return 0, 0, 100
	}
	return opts.Primary, opts.Canary, opts.Baseline
}

// TerraformPlanStageOptions contains all configurable values for a K8S_TERRAFORM_PLAN stage.
type TerraformPlanStageOptions struct {
}

// TerraformApplyStageOptions contains all configurable values for a K8S_TERRAFORM_APPLY stage.
type TerraformApplyStageOptions struct {
	// How many times to retry applying terraform changes.
	Retries int `json:"retries"`
}

// AnalysisStageOptions contains all configurable values for a K8S_ANALYSIS stage.
type AnalysisStageOptions struct {
	// How long the analysis process should be executed.
	Duration Duration `json:"duration"`
	// TODO: Consider about how to handle a pod restart
	// possible count of pod restarting
	RestartThreshold int                          `json:"restartThreshold"`
	Metrics          []TemplatableAnalysisMetrics `json:"metrics"`
	Logs             []TemplatableAnalysisLog     `json:"logs"`
	Https            []TemplatableAnalysisHTTP    `json:"https"`
	Dynamic          AnalysisDynamic              `json:"dynamic"`
}

type AnalysisTemplateRef struct {
	Name string            `json:"name"`
	Args map[string]string `json:"args"`
}

// TemplatableAnalysisMetrics wraps AnalysisMetrics to allow specify template to use.
type TemplatableAnalysisMetrics struct {
	AnalysisMetrics
	Template AnalysisTemplateRef `json:"template"`
}

// TemplatableAnalysisLog wraps AnalysisLog to allow specify template to use.
type TemplatableAnalysisLog struct {
	AnalysisLog
	Template AnalysisTemplateRef `json:"template"`
}

// TemplatableAnalysisHTTP wraps AnalysisHTTP to allow specify template to use.
type TemplatableAnalysisHTTP struct {
	AnalysisHTTP
	Template AnalysisTemplateRef `json:"template"`
}

type KubernetesDeploymentInput struct {
	Manifests      []string `json:"manifests"`
	KubectlVersion string   `json:"kubectlVersion"`

	KustomizeVersion string `json:"kustomizeVersion"`

	HelmChart   *InputHelmChart   `json:"helmChart"`
	HelmOptions *InputHelmOptions `json:"helmOptions"`
	HelmVersion string            `json:"helmVersion"`

	// The namespace where manifests will be applied.
	Namespace string `json:"namespace"`
	// Automatically reverts all changes from all stages when one of them failed.
	// Default is true.
	AutoRollback bool     `json:"autoRollback"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type TerraformDeploymentInput struct {
	Workspace        string `json:"workspace,omitempty"`
	TerraformVersion string `json:"terraformVersion,omitempty"`
	// Automatically reverts all changes from all stages when one of them failed.
	// Default is false.
	AutoRollback bool     `json:"autoRollback"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type InputHelmChart struct {
	// Empty means current repository.
	GitRemote string `json:"gitRemote"`
	// The commit SHA or tag for remote git.
	Ref string `json:"ref"`
	// Relative path from the repository root directory to the chart directory.
	Path string `json:"path"`

	// The name of an added Helm chart repository.
	Repository string `json:"repository"`
	Name       string `json:"name"`
	Version    string `json:"version"`
}

type InputHelmOptions struct {
	// By default the release name is equal to the application name.
	ReleaseName string `json:"releaseName"`
	// List of value files should be loaded.
	ValueFiles []string `json:"valueFiles"`
	SetFiles   map[string]string
}
