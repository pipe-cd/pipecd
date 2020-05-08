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

// KubernetesAppSpec represents configuration for a Kubernetes application.
type KubernetesAppSpec struct {
	// Selector is a list of labels used to query all resources of this application.
	Selector    map[string]string   `json:"selector"`
	Input       *KubernetesAppInput `json:"input"`
	Pipeline    *AppPipeline        `json:"pipeline"`
	Destination string              `json:"destination"`
}

// Validate returns an error if any wrong configuration value was found.
func (s *KubernetesAppSpec) Validate() error {
	return nil
}

// TerraformAppSpec represents configuration for a Terraform application.
type TerraformAppSpec struct {
	Input       *TerraformAppInput `json:"input"`
	Pipeline    *AppPipeline       `json:"pipeline"`
	Destination string             `json:"destination"`
}

// Validate returns an error if any wrong configuration value was found.
func (s *TerraformAppSpec) Validate() error {
	if s.Destination == "" {
		return fmt.Errorf("spec.destination for terraform application is required")
	}
	return nil
}

// AppPipeline represents the way to deploy the application.
// The pipeline is triggered by changes in any of the following objects:
// - Target PodSpec (Target can be Deployment, DaemonSet, StatefullSet)
// - ConfigMaps, Secrets that are mounted as volumes or envs in the deployment.
type AppPipeline struct {
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

	WaitStageOptions            *WaitStageOptions
	WaitApprovalStageOptions    *WaitApprovalStageOptions
	AnalysisStageOptions        *AnalysisStageOptions
	K8sPrimaryOutStageOptions   *K8sPrimaryOutStageOptions
	K8sStageOutStageOptions     *K8sStageOutStageOptions
	K8sStageInStageOptions      *K8sStageInStageOptions
	K8sBaselineOutStageOptions  *K8sBaselineOutStageOptions
	K8sBaselineInStageOptions   *K8sBaselineInStageOptions
	K8sTrafficRouteStageOptions *K8sTrafficRouteStageOptions
	TerraformPlanStageOptions   *TerraformPlanStageOptions
	TerraformApplyStageOptions  *TerraformApplyStageOptions
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
	case model.StageK8sPrimaryOut:
		s.K8sPrimaryOutStageOptions = &K8sPrimaryOutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sPrimaryOutStageOptions)
		}
	case model.StageK8sStageOut:
		s.K8sStageOutStageOptions = &K8sStageOutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sStageOutStageOptions)
		}
	case model.StageK8sStageIn:
		s.K8sStageInStageOptions = &K8sStageInStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sStageInStageOptions)
		}
	case model.StageK8sBaselineOut:
		s.K8sBaselineOutStageOptions = &K8sBaselineOutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sBaselineOutStageOptions)
		}
	case model.StageK8sBaselineIn:
		s.K8sBaselineInStageOptions = &K8sBaselineInStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sBaselineInStageOptions)
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

// WaitStageOptions contains all configurable values for a K8S_PRIMARY_OUT stage.
type K8sPrimaryOutStageOptions struct {
	Manifests []string `json:"manifests"`
}

// K8sStageOutStageOptions contains all configurable values for a K8S_STAGE_OUT stage.
type K8sStageOutStageOptions struct {
	// Percentage of pods for STAGE workloads.
	// Default is 1 pod.
	Weight int `json:"weight"`
	// Suffix that should be used when naming the STAGE resources.
	// Default is "stage".
	Suffix string
	// If true the service resource for the STAGE will be created.
	WithService bool
}

// K8sStageInStageOptions contains all configurable values for a K8S_STAGE_IN stage.
type K8sStageInStageOptions struct {
}

// K8sBaselineOutStageOptions contains all configurable values for a K8S_BASELINE_OUT stage.
type K8sBaselineOutStageOptions struct {
	// Percentage of pods for BASELINE workloads.
	// Default is 1 pod.
	Weight int `json:"weight"`
	// Suffix that should be used when naming the STAGE resources.
	// Default is "baseline".
	Suffix string
	// If true the service resource for the BASELINE will be created.
	WithService bool
}

// K8sBaselineInStageOptions contains all configurable values for a K8S_BASELINE_IN stage.
type K8sBaselineInStageOptions struct {
}

// K8sTrafficRouteStageOptions contains all configurable values for a K8S_TRAFFIC_ROUTE stage.
type K8sTrafficRouteStageOptions struct {
	// Target can be "primary", "stage", "baseline".
	// If this field was configured, all of the traffic to the applicaiton
	// should be routing to the target.
	Target string `json:"target,omitempty"`
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

type KubernetesAppInput struct {
	Manifests      []string        `json:"manifests"`
	KubectlVersion string          `json:"kubectlVersion"`
	HelmChart      *InputHelmChart `json:"helmChart"`
	HelmValueFiles []string        `json:"helmValueFiles"`
	HelmVersion    string          `json:"helmVersion"`
	Dependencies   []string        `json:"dependencies,omitempty"`
}

type TerraformAppInput struct {
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
