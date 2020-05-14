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

// KubernetesDeploymentSpec represents a deployment configuration for Kubernetes application.
type KubernetesDeploymentSpec struct {
	// Selector is a list of labels used to query all resources of this application.
	Selector    map[string]string          `json:"selector"`
	Input       *KubernetesDeploymentInput `json:"input"`
	Pipeline    *DeploymentPipeline        `json:"pipeline"`
	Destination string                     `json:"destination"`
}

// Validate returns an error if any wrong configuration value was found.
func (s *KubernetesDeploymentSpec) Validate() error {
	return nil
}

// TerraformDeploymentSpec represents a deployment configuration for Terraform application.
type TerraformDeploymentSpec struct {
	Input       *TerraformDeploymentInput `json:"input"`
	Pipeline    *DeploymentPipeline       `json:"pipeline"`
	Destination string                    `json:"destination"`
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

type StageTrackOptions struct {
	Target  *K8sDeployTarget
	Service *K8sService
}

type BaselineTrackOptions struct {
	Target  *K8sDeployTarget
	Service *K8sService
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
	K8sTrafficRouteStageOptions    *K8sTrafficRouteStageOptions
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
	case model.StageK8sTrafficRoute:
		s.K8sTrafficRouteStageOptions = &K8sTrafficRouteStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sTrafficRouteStageOptions)
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
	// Suffix that should be used when naming the STAGE resources.
	// Default is "stage".
	Suffix string
	// If true the service resource for the STAGE will be created.
	WithService bool
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
	// Suffix that should be used when naming the BASELINE resources.
	// Default is "baseline".
	Suffix string
	// If true the service resource for the BASELINE will be created.
	WithService bool
}

// K8sBaselineCleanStageOptions contains all configurable values for a K8S_BASELINE_CLEAN stage.
type K8sBaselineCleanStageOptions struct {
}

// K8sTrafficRouteStageOptions contains all configurable values for a K8S_TRAFFIC_ROUTE stage.
type K8sTrafficRouteStageOptions struct {
	// The percentage of traffic should be routed to PRIMARY.
	Primary int `json:"primary"`
	// The percentage of traffic should be routed to STAGE.
	Stage int `json:"stage"`
	// The percentage of traffic should be routed to BASELINE.
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
	Threshold int `json:"threshold"`
	// Maximum number of container restarts before the stage is considered as failure.
	RestartThreshold int               `json:"restartThreshold"`
	Metrics          []AnalysisMetrics `json:"metrics"`
	Logs             []AnalysisLog     `json:"logs"`
	Https            []AnalysisHTTP    `json:"https"`
}

type AnalysisMetrics struct {
	Query    string   `json:"query"`
	Expected string   `json:"expected"`
	Interval Duration `json:"interval"`
	// How long after which the query times out.
	Timeout     Duration `json:"timeout"`
	Provider    string   `json:"provider"`
	UseTemplate string   `json:"useTemplate"`
}

// Comprare the log entries between new version and old version.
type AnalysisLog struct {
	Query       string
	Threshold   int
	Provider    string
	UseTemplate string
}

type AnalysisHTTP struct {
	URL    string
	Method string
	// Custom headers to set in the request. HTTP allows repeated headers.
	Headers          []string
	ExpectedStatus   string
	ExpectedResponse string
	Interval         Duration
	Timeout          Duration
	UseTemplate      string
}

type KubernetesDeploymentInput struct {
	Manifests      []string        `json:"manifests"`
	KubectlVersion string          `json:"kubectlVersion"`
	HelmChart      *InputHelmChart `json:"helmChart"`
	HelmValueFiles []string        `json:"helmValueFiles"`
	HelmVersion    string          `json:"helmVersion"`
	Dependencies   []string        `json:"dependencies,omitempty"`
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

type K8sDeployTarget struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type K8sService struct {
	Name string `json:"name"`
}
