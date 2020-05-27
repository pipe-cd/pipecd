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

	"github.com/kapetaniosci/pipe/pkg/model"
)

var (
	KubernetesDefaultPipeline = &DeploymentPipeline{
		Stages: []PipelineStage{
			{
				Name: model.StageK8sPrimaryUpdate,
				Desc: "Update primary to new version",
			},
		},
	}
	TerraformDefaultPipeline = &DeploymentPipeline{
		Stages: []PipelineStage{
			{
				Name: model.StageTerraformPlan,
				Desc: "Terraform Plan",
			},
			{
				Name: model.StageWaitApproval,
				Desc: "Wait for an approval",
			},
			{
				Name: model.StageTerraformApply,
				Desc: "Terraform Apply",
			},
		},
	}
)

type Pipelineable interface {
	GetPipeline() *DeploymentPipeline
	GetStage(index int32) (PipelineStage, bool)
}

// KubernetesDeploymentSpec represents a deployment configuration for Kubernetes application.
type KubernetesDeploymentSpec struct {
	// Selector is a list of labels used to query all resources of this application.
	Selector        map[string]string         `json:"selector"`
	Input           KubernetesDeploymentInput `json:"input"`
	Pipeline        *DeploymentPipeline       `json:"pipeline"`
	StageVariant    *StageVariant             `json:"stageVariant"`
	BaselineVariant *BaselineVariant          `json:"baselineVariant"`
	TrafficSplit    TrafficSplit              `json:"trafficSplit"`
	Destination     string                    `json:"destination"`
}

func (s *KubernetesDeploymentSpec) GetPipeline() *DeploymentPipeline {
	if s.Pipeline != nil {
		return s.Pipeline
	}
	return KubernetesDefaultPipeline
}

func (s *KubernetesDeploymentSpec) GetStage(index int32) (PipelineStage, bool) {
	p := s.GetPipeline()
	if int(index) >= len(p.Stages) {
		return PipelineStage{}, false
	}
	return p.Stages[index], true
}

// Validate returns an error if any wrong configuration value was found.
func (s *KubernetesDeploymentSpec) Validate() error {
	return nil
}

// TerraformDeploymentSpec represents a deployment configuration for Terraform application.
type TerraformDeploymentSpec struct {
	Input       TerraformDeploymentInput `json:"input"`
	Pipeline    *DeploymentPipeline      `json:"pipeline"`
	Destination string                   `json:"destination"`
}

func (s *TerraformDeploymentSpec) GetPipeline() *DeploymentPipeline {
	if s.Pipeline != nil {
		return s.Pipeline
	}
	return TerraformDefaultPipeline
}

func (s *TerraformDeploymentSpec) GetStage(index int32) (PipelineStage, bool) {
	p := s.GetPipeline()
	if int(index) >= len(p.Stages) {
		return PipelineStage{}, false
	}
	return p.Stages[index], true
}

// Validate returns an error if any wrong configuration value was found.
func (s *TerraformDeploymentSpec) Validate() error {
	if s.Destination == "" {
		return fmt.Errorf("spec.destination for terraform application is required")
	}
	return nil
}

// DeploymentPipeline represents the way to deploy the application.
// The pipeline is triggered by changes in any of the following objects:
// - Target PodSpec (Target can be Deployment, DaemonSet, StatefullSet)
// - ConfigMaps, Secrets that are mounted as volumes or envs in the deployment.
type DeploymentPipeline struct {
	Stages []PipelineStage `json:"stages"`
}

type StageVariant struct {
	Workload K8sWorkload `json:"workload"`
	Service  K8sService  `json:"service"`
	// Suffix that should be used when naming the STAGE variant's resources.
	// Default is "stage".
	Suffix string `json:"suffix"`
}

type BaselineVariant struct {
	Workload K8sWorkload `json:"workload"`
	Service  K8sService  `json:"service"`
	// Suffix that should be used when naming the BASELINE variant's resources.
	// Default is "baseline".
	Suffix string `json:"suffix"`
}

type TrafficSplitMethod string

const (
	TrafficSplitMethodPod   TrafficSplitMethod = "pod"
	TrafficSplitMethodIstio TrafficSplitMethod = "istio"
)

type TrafficSplit struct {
	Method TrafficSplitMethod `json:"method"`
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
	K8sPrimaryUpdateStageOptions   *K8sPrimaryUpdateStageOptions
	K8sStageRolloutStageOptions    *K8sStageRolloutStageOptions
	K8sStageCleanStageOptions      *K8sStageCleanStageOptions
	K8sBaselineRolloutStageOptions *K8sBaselineRolloutStageOptions
	K8sBaselineCleanStageOptions   *K8sBaselineCleanStageOptions
	K8sTrafficSplitStageOptions    *K8sTrafficSplitStageOptions
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
	case model.StageK8sPrimaryUpdate:
		s.K8sPrimaryUpdateStageOptions = &K8sPrimaryUpdateStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sPrimaryUpdateStageOptions)
		}
	case model.StageK8sStageRollout:
		s.K8sStageRolloutStageOptions = &K8sStageRolloutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sStageRolloutStageOptions)
		}
	case model.StageK8sStageClean:
		s.K8sStageCleanStageOptions = &K8sStageCleanStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sStageCleanStageOptions)
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
	case model.StageK8sTrafficSplit:
		s.K8sTrafficSplitStageOptions = &K8sTrafficSplitStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sTrafficSplitStageOptions)
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

// WaitStageOptions contains all configurable values for a K8S_PRIMARY_UPDATE stage.
type K8sPrimaryUpdateStageOptions struct {
	Manifests []string `json:"manifests"`
}

// K8sStageRolloutStageOptions contains all configurable values for a K8S_STAGE_ROLLOUT stage.
type K8sStageRolloutStageOptions struct {
	// How many pods for STAGE workloads.
	// An integer value can be specified to indicate an absolute value of pod number.
	// Or a string suffixed by "%" to indicate an percantage value compared to the pod number of PRIMARY.
	// Default is 1 pod.
	Replicas Replicas `json:"replicas"`
}

// K8sStageCleanStageOptions contains all configurable values for a K8S_STAGE_CLEAN stage.
type K8sStageCleanStageOptions struct {
}

// K8sBaselineRolloutStageOptions contains all configurable values for a K8S_BASELINE_ROLLOUT stage.
type K8sBaselineRolloutStageOptions struct {
	// How many pods for BASELINE workloads.
	// An integer value can be specified to indicate an absolute value of pod number.
	// Or a string suffixed by "%" to indicate an percantage value compared to the pod number of PRIMARY.
	// Default is 1 pod.
	Replicas Replicas `json:"replicas"`
}

// K8sBaselineCleanStageOptions contains all configurable values for a K8S_BASELINE_CLEAN stage.
type K8sBaselineCleanStageOptions struct {
}

// K8sTrafficSplitStageOptions contains all configurable values for a K8S_TRAFFIC_SPLIT stage.
type K8sTrafficSplitStageOptions struct {
	// Which variant should receive all traffic.
	All string `json:"all"`
	// The percentage of traffic should be routed to PRIMARY variant.
	Primary int `json:"primary"`
	// The percentage of traffic should be routed to STAGE variant.
	Stage int `json:"stage"`
	// The percentage of traffic should be routed to BASELINE variant.
	Baseline int `json:"baseline"`
}

// TerraformPlanStageOptions contains all configurable values for a K8S_TERRAFORM_PLAN stage.
type TerraformPlanStageOptions struct {
}

// TerraformApplyStageOptions contains all configurable values for a K8S_TERRAFORM_APPLY stage.
type TerraformApplyStageOptions struct {
}

// AnalysisStageOptions contains all configurable values for a K8S_ANALYSIS stage.
type AnalysisStageOptions struct {
	// How long the analysis process should be executed.
	Duration Duration `json:"duration"`
	// Maximum number of failed checks before the stage is considered as failure.
	// TODO: Remove
	FailureLimit     int               `json:"failureLimit"`
	RestartThreshold int               `json:"restartThreshold"`
	Metrics          []AnalysisMetrics `json:"metrics"`
	Logs             []AnalysisLog     `json:"logs"`
	Https            []AnalysisHTTP    `json:"https"`
	Dynamic          AnalysisDynamic   `json:"dynamic"`
}

type AnalysisMetrics struct {
	Query    string           `json:"query"`
	Expected AnalysisExpected `json:"expected"`
	Interval Duration         `json:"interval"`
	// Maximum number of failed checks before the query result is considered as failure.
	FailureLimit int `json:"failureLimit"`
	// How long after which the query times out.
	Timeout     Duration `json:"timeout"`
	Provider    string   `json:"provider"`
	UseTemplate string   `json:"useTemplate"`
}

// AnalysisExpected defines the range used for metrics analysis.
type AnalysisExpected struct {
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`
}

type AnalysisLog struct {
	Query    string   `json:"query"`
	Interval Duration `json:"interval"`
	// Maximum number of failed checks before the query result is considered as failure.
	FailureLimit int `json:"failureLimit"`
	// How long after which the query times out.
	Timeout     Duration `json:"timeout"`
	Provider    string   `json:"provider"`
	UseTemplate string   `json:"useTemplate"`
}

type AnalysisHTTP struct {
	URL    string
	Method string
	// Custom headers to set in the request. HTTP allows repeated headers.
	Headers          []string
	ExpectedStatus   string
	ExpectedResponse string
	Interval         Duration
	// Maximum number of failed checks before the response is considered as failure.
	FailureLimit int      `json:"failureLimit"`
	Timeout      Duration `json:"timeout"`
	UseTemplate  string   `json:"useTemplate"`
}

// AnalysisDynamic contains settings for analysis by comparing  with dynamic data.
type AnalysisDynamic struct {
	Metrics []AnalysisDynamicMetrics `json:"metrics"`
	Logs    []AnalysisDynamicLog     `json:"logs"`
	Https   []AnalysisDynamicHTTP    `json:"https"`
}

type AnalysisDynamicMetrics struct {
	Query    string   `json:"query"`
	Provider string   `json:"provider"`
	Timeout  Duration `json:"timeout"`
}

type AnalysisDynamicLog struct {
	Query    string   `json:"query"`
	Provider string   `json:"provider"`
	Timeout  Duration `json:"timeout"`
}

type AnalysisDynamicHTTP struct {
	URL              string
	Method           string
	Headers          []string
	ExpectedStatus   string
	ExpectedResponse string
	Interval         Duration
	Timeout          Duration `json:"timeout"`
}

type KubernetesDeploymentInput struct {
	Manifests        []string        `json:"manifests"`
	KubectlVersion   string          `json:"kubectlVersion"`
	KustomizeVersion string          `json:"kustomizeVersion"`
	HelmChart        *InputHelmChart `json:"helmChart"`
	HelmValueFiles   []string        `json:"helmValueFiles"`
	HelmVersion      string          `json:"helmVersion"`
	Dependencies     []string        `json:"dependencies,omitempty"`
}

type TerraformDeploymentInput struct {
	Workspace        string   `json:"workspace,omitempty"`
	TerraformVersion string   `json:"terraformVersion,omitempty"`
	Dependencies     []string `json:"dependencies,omitempty"`
}

type InputHelmChart struct {
	Git        string `json:"git"`
	Path       string `json:"path"`
	Ref        string `json:"ref"`
	Repository string `json:"repository"`
	Name       string `json:"name"`
	Version    string `json:"version"`
}

type K8sWorkload struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type K8sService struct {
	Name string `json:"name"`
}
