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
)

// StageName represents predefined Stage that can be used in the pipeline.
type StageName string

const (
	StageNameWait                      StageName = "WAIT"
	StageNameWaitApproval              StageName = "WAIT_APPROVAL"
	StageNameAnalysis                  StageName = "ANALYSIS"
	StageNameK8sRollout                StageName = "K8S_ROLLOUT"
	StageNameK8sCanaryOut              StageName = "K8S_CANARY_OUT"
	StageNameK8sCanaryIn               StageName = "K8S_CANARY_IN"
	StageNameK8sBlueGreenOut           StageName = "K8S_BLUEGREEN_OUT"
	StageNameK8sBlueGreenSwitchTraffic StageName = "K8S_BLUEGREEN_SWITCH_TRAFFIC"
	StageNameK8sBlueGreenIn            StageName = "K8S_BLUEGREEN_IN"
	StageNameTerraformPlan             StageName = "TERRAFORM_PLAN"
	StageNameTerraformApply            StageName = "TERRAFORM_APPLY"
)

// StageName represents temporary desired state
// before reaching the final desired state.
// type StageName string

// const (
// 	StageNameWait            StageName = "WAIT"
// 	StageNameWaitApproval    StageName = "WAIT_APPROVAL"
// 	StageNameAnalysis        StageName = "ANALYSIS"
// 	StageNameK8sPrimaryOut   StageName = "K8S_PRIMARY_OUT"
// 	StageNameK8sStageOut     StageName = "K8S_STAGE_OUT"
// 	StageNameK8sStageIn      StageName = "K8S_STAGE_IN"
// 	StageNameK8sBaselineOut  StageName = "K8S_BASELINE_OUT"
// 	StageNameK8sBaselineIn   StageName = "K8S_BASELINE_IN"
// 	StageNameK8sTrafficRoute StageName = "K8S_TRAFFIC_ROUTE"
// 	StageNameTerraformPlan   StageName = "TERRAFORM_PLAN"
// 	StageNameTerraformApply  StageName = "TERRAFORM_APPLY"
// )

type K8sAppSpec struct {
	// Selector is a list of labels used to query all resources of this application.
	Selector    map[string]string `json:"selector"`
	Input       *K8sAppInput      `json:"input"`
	Pipeline    *AppPipeline      `json:"pipeline"`
	Destination string            `json:"destination"`
}

func (s *K8sAppSpec) Validate() error {
	return nil
}

type K8sKustomizationAppSpec struct {
	// Selector is a list of labels used to query all resources of this application.
	Selector    map[string]string         `json:"selector"`
	Input       *K8sKustomizationAppInput `json:"input"`
	Pipeline    *AppPipeline              `json:"pipeline"`
	Destination string                    `json:"destination"`
}

func (s *K8sKustomizationAppSpec) Validate() error {
	return nil
}

type K8sHelmAppSpec struct {
	// Selector is a list of labels used to query all resources of this application.
	Selector    map[string]string `json:"selector"`
	Input       *K8sHelmAppInput  `json:"input"`
	Pipeline    *AppPipeline      `json:"pipeline"`
	Destination string            `json:"destination"`
}

func (s *K8sHelmAppSpec) Validate() error {
	return nil
}

type TerraformAppSpec struct {
	Input       *TerraformAppInput `json:"input"`
	Pipeline    *AppPipeline       `json:"pipeline"`
	Destination string             `json:"destination"`
}

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

type PipelineStage struct {
	Name      StageName
	Desc      string
	Timeout   Duration
	PostDelay Duration

	ApprovalStageOptions                  *ApprovalStageOptions
	AnalysisStageOptions                  *AnalysisStageOptions
	K8sRolloutStageOptions                *K8sRolloutStageOptions
	K8sCanaryOutStageOptions              *K8sCanaryOutStageOptions
	K8sCanaryInStageOptions               *K8sCanaryInStageOptions
	K8sBlueGreenOutStageOptions           *K8sBlueGreenOutStageOptions
	K8sBlueGreenSwitchTrafficStageOptions *K8sBlueGreenSwitchTrafficStageOptions
	K8sBlueGreenInStageOptions            *K8sBlueGreenInStageOptions
	TerraformPlanStageOptions             *TerraformPlanStageOptions
	TerraformApplyStageOptions            *TerraformApplyStageOptions
}

type genericPipelineStage struct {
	Name      StageName       `json:"name"`
	Desc      string          `json:"desc,omitempty"`
	Timeout   Duration        `json:"timeout"`
	PostDelay Duration        `json:"postDelay"`
	With      json.RawMessage `json:"with"`
}

func (s *PipelineStage) UnmarshalJSON(data []byte) error {
	var err error
	gs := genericPipelineStage{}
	if err = json.Unmarshal(data, &gs); err != nil {
		return err
	}
	s.Name = gs.Name
	s.Desc = gs.Desc
	s.Timeout = gs.Timeout
	s.PostDelay = gs.PostDelay

	switch s.Name {
	case StageNameWaitApproval:
		s.ApprovalStageOptions = &ApprovalStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.ApprovalStageOptions)
		}
	case StageNameAnalysis:
		s.AnalysisStageOptions = &AnalysisStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.AnalysisStageOptions)
		}
	case StageNameK8sRollout:
		s.K8sRolloutStageOptions = &K8sRolloutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sRolloutStageOptions)
		}
	case StageNameK8sCanaryOut:
		s.K8sCanaryOutStageOptions = &K8sCanaryOutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sCanaryOutStageOptions)
		}
	case StageNameK8sCanaryIn:
		s.K8sCanaryInStageOptions = &K8sCanaryInStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sCanaryInStageOptions)
		}
	case StageNameK8sBlueGreenOut:
		s.K8sBlueGreenOutStageOptions = &K8sBlueGreenOutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sBlueGreenOutStageOptions)
		}
	case StageNameK8sBlueGreenSwitchTraffic:
		s.K8sBlueGreenSwitchTrafficStageOptions = &K8sBlueGreenSwitchTrafficStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sBlueGreenSwitchTrafficStageOptions)
		}
	case StageNameK8sBlueGreenIn:
		s.K8sBlueGreenInStageOptions = &K8sBlueGreenInStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.K8sBlueGreenInStageOptions)
		}
	case StageNameTerraformPlan:
		s.TerraformPlanStageOptions = &TerraformPlanStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.TerraformPlanStageOptions)
		}
	case StageNameTerraformApply:
		s.TerraformApplyStageOptions = &TerraformApplyStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.TerraformApplyStageOptions)
		}
	default:
		err = fmt.Errorf("unsupported stage name: %s", s.Name)
	}
	return err
}

type ApprovalStageOptions struct {
	Approvers []string `json:"approvers"`
}

type K8sRolloutStageOptions struct {
	Manifests []string `json:"manifests"`
}

type K8sCanaryOutStageOptions struct {
	// Percentage of canary traffic/pods after scale out.
	// Default is 10%.
	Weight        int             `json:"weight"`
	CanaryService K8sService      `json:"canaryService"`
	Target        K8sDeployTarget `json:"target"`
}

type K8sCanaryInStageOptions struct {
	// Percentage of canary traffic/pods after scale in.
	// Default is 0.
	Weight int
}

type K8sBlueGreenOutStageOptions struct {
	StageService string `json:"stageService"`
}

type K8sBlueGreenSwitchTrafficStageOptions struct {
}

type K8sBlueGreenInStageOptions struct {
}

type TerraformPlanStageOptions struct {
}

type TerraformApplyStageOptions struct {
}

type AnalysisStageOptions struct {
	Duration Duration
	// Maximum number of failed checks before the stage is considered as failure.
	Threshold int               `json:"threshold"`
	Metrics   []AnalysisMetrics `json:"metrics"`
	Logs      []AnalysisLog     `json:"logs"`
	Https     []AnalysisHTTP    `json:"https"`
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

type K8sAppInput struct {
	Manifests      []string
	KubectlVersion string
}

type K8sKustomizationAppInput struct {
	KubectlVersion string
}

type K8sHelmAppInput struct {
	Chart       string
	ValueFiles  []string
	Namespace   string
	HelmVersion string
}

type TerraformAppInput struct {
	Workspace        string
	TerraformVersion string
}

type K8sDeployTarget struct {
	Kind string
	Name string
}

type K8sService struct {
	Name string
}
