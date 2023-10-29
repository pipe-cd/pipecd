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

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/model"
)

const allEventsSymbol = "*"

type GenericApplicationSpec struct {
	// The application name.
	// This is required if you set the application through the application configuration file.
	Name string `json:"name"`
	// Additional attributes to identify applications.
	Labels map[string]string `json:"labels"`
	// Notes on the Application.
	Description string `json:"description"`

	// Configuration used while planning deployment.
	Planner DeploymentPlanner `json:"planner"`
	// Forcibly use QuickSync or Pipeline when commit message matched the specified pattern.
	CommitMatcher DeploymentCommitMatcher `json:"commitMatcher"`
	// Pipeline for deploying progressively.
	Pipeline *DeploymentPipeline `json:"pipeline"`
	// The trigger configuration use to determine trigger logic.
	Trigger Trigger `json:"trigger"`
	// Configuration to be used once the deployment is triggered successfully.
	PostSync *PostSync `json:"postSync"`
	// The maximum length of time to execute deployment before giving up.
	// Default is 6h.
	Timeout Duration `json:"timeout,omitempty" default:"6h"`
	// List of encrypted secrets and targets that should be decoded before using.
	Encryption *SecretEncryption `json:"encryption"`
	// List of files that should be attached to application manifests before using.
	Attachment *Attachment `json:"attachment"`
	// Additional configuration used while sending notification to external services.
	DeploymentNotification *DeploymentNotification `json:"notification"`
	// List of the configuration for event watcher.
	EventWatcher []EventWatcherConfig `json:"eventWatcher"`
	// Configuration for drift detection
	DriftDetection *DriftDetection `json:"driftDetection"`
}

type DeploymentPlanner struct {
	// Disable auto-detecting to use QUICK_SYNC or PROGRESSIVE_SYNC.
	// Always use the speficied pipeline for all deployments.
	AlwaysUsePipeline bool `json:"alwaysUsePipeline"`
}

type Trigger struct {
	// Configurable fields used while deciding the application
	// should be triggered or not based on commit changes.
	OnCommit OnCommit `json:"onCommit"`
	// Configurable fields used while deciding the application
	// should be triggered or not based on received SYNC command.
	OnCommand OnCommand `json:"onCommand"`
	// Configurable fields used while deciding the application
	// should be triggered or not based on OUT_OF_SYNC state.
	OnOutOfSync OnOutOfSync `json:"onOutOfSync"`
	// Configurable fields used while deciding the application
	// should be triggered based on received CHAIN_SYNC command.
	OnChain OnChain `json:"onChain"`
}

type OnCommit struct {
	// Whether to exclude application from triggering target
	// when a new commit touched the application.
	// Default is false.
	Disabled bool `json:"disabled,omitempty"`
	// List of directories or files where their changes will trigger the deployment.
	// Regular expression can be used.
	Paths []string `json:"paths,omitempty"`
	// List of directories or files where their changes will be ignored.
	// Regular expression can be used.
	Ignores []string `json:"ignores,omitempty"`
}

type OnCommand struct {
	// Whether to exclude application from triggering target
	// when received a new SYNC command.
	// Default is false.
	Disabled bool `json:"disabled,omitempty"`
}

type OnOutOfSync struct {
	// Whether to exclude application from triggering target
	// when application is at OUT_OF_SYNC state.
	// Default is true.
	Disabled *bool `json:"disabled,omitempty" default:"true"`
	// Minimum amount of time must be elapsed since the last deployment.
	// This can be used to avoid triggering unnecessary continuous deployments based on OUT_OF_SYNC status.
	MinWindow Duration `json:"minWindow,omitempty" default:"5m"`
}

type OnChain struct {
	// Whether to exclude application from triggering target
	// when received a new CHAIN_SYNC command.
	// Default is true.
	Disabled *bool `json:"disabled,omitempty" default:"true"`
}

func (s *GenericApplicationSpec) Validate() error {
	if s.Pipeline != nil {
		for _, stage := range s.Pipeline.Stages {
			if stage.AnalysisStageOptions != nil {
				if err := stage.AnalysisStageOptions.Validate(); err != nil {
					return err
				}
			}
			if stage.WaitApprovalStageOptions != nil {
				if err := stage.WaitApprovalStageOptions.Validate(); err != nil {
					return err
				}
			}
			if stage.CustomSyncOptions != nil {
				if err := stage.CustomSyncOptions.Validate(); err != nil {
					return err
				}
			}
		}
	}

	if ps := s.PostSync; ps != nil {
		if err := ps.Validate(); err != nil {
			return err
		}
	}

	if e := s.Encryption; e != nil {
		if err := e.Validate(); err != nil {
			return err
		}
	}

	if am := s.Attachment; am != nil {
		if err := am.Validate(); err != nil {
			return err
		}
	}

	if s.DeploymentNotification != nil {
		for _, m := range s.DeploymentNotification.Mentions {
			if err := m.Validate(); err != nil {
				return err
			}
		}
	}

	if dd := s.DriftDetection; dd != nil {
		if err := dd.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (s GenericApplicationSpec) GetStage(index int32) (PipelineStage, bool) {
	if s.Pipeline == nil {
		return PipelineStage{}, false
	}
	if int(index) >= len(s.Pipeline.Stages) {
		return PipelineStage{}, false
	}
	return s.Pipeline.Stages[index], true
}

// HasStage checks if the given stage is included in the pipeline.
func (s GenericApplicationSpec) HasStage(stage model.Stage) bool {
	if s.Pipeline == nil {
		return false
	}
	for _, s := range s.Pipeline.Stages {
		if s.Name == stage {
			return true
		}
	}
	return false
}

// DeploymentCommitMatcher provides a way to decide how to deploy.
type DeploymentCommitMatcher struct {
	// It makes sure to perform syncing if the commit message matches this regular expression.
	QuickSync string `json:"quickSync"`
	// It makes sure to perform pipeline if the commit message matches this regular expression.
	Pipeline string `json:"pipeline"`
}

// DeploymentPipeline represents the way to deploy the application.
// The pipeline is triggered by changes in any of the following objects:
// - Target PodSpec (Target can be Deployment, DaemonSet, StatefulSet)
// - ConfigMaps, Secrets that are mounted as volumes or envs in the deployment.
type DeploymentPipeline struct {
	Stages []PipelineStage `json:"stages"`
}

// PipelineStage represents a single stage of a pipeline.
// This is used as a generic struct for all stage type.
type PipelineStage struct {
	ID      string
	Name    model.Stage
	Desc    string
	Timeout Duration

	CustomSyncOptions        *CustomSyncOptions
	WaitStageOptions         *WaitStageOptions
	WaitApprovalStageOptions *WaitApprovalStageOptions
	AnalysisStageOptions     *AnalysisStageOptions

	K8sPrimaryRolloutStageOptions  *K8sPrimaryRolloutStageOptions
	K8sCanaryRolloutStageOptions   *K8sCanaryRolloutStageOptions
	K8sCanaryCleanStageOptions     *K8sCanaryCleanStageOptions
	K8sBaselineRolloutStageOptions *K8sBaselineRolloutStageOptions
	K8sBaselineCleanStageOptions   *K8sBaselineCleanStageOptions
	K8sTrafficRoutingStageOptions  *K8sTrafficRoutingStageOptions

	TerraformSyncStageOptions  *TerraformSyncStageOptions
	TerraformPlanStageOptions  *TerraformPlanStageOptions
	TerraformApplyStageOptions *TerraformApplyStageOptions

	CloudRunSyncStageOptions    *CloudRunSyncStageOptions
	CloudRunPromoteStageOptions *CloudRunPromoteStageOptions

	LambdaSyncStageOptions          *LambdaSyncStageOptions
	LambdaCanaryRolloutStageOptions *LambdaCanaryRolloutStageOptions
	LambdaPromoteStageOptions       *LambdaPromoteStageOptions

	ECSSyncStageOptions           *ECSSyncStageOptions
	ECSCanaryRolloutStageOptions  *ECSCanaryRolloutStageOptions
	ECSPrimaryRolloutStageOptions *ECSPrimaryRolloutStageOptions
	ECSCanaryCleanStageOptions    *ECSCanaryCleanStageOptions
	ECSTrafficRoutingStageOptions *ECSTrafficRoutingStageOptions
}

type genericPipelineStage struct {
	ID      string          `json:"id"`
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
	s.ID = gs.ID
	s.Name = gs.Name
	s.Desc = gs.Desc
	s.Timeout = gs.Timeout

	switch s.Name {
	case model.StageCustomSync:
		s.CustomSyncOptions = &CustomSyncOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.CustomSyncOptions)
		}
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

	case model.StageTerraformSync:
		s.TerraformSyncStageOptions = &TerraformSyncStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.TerraformSyncStageOptions)
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

	case model.StageCloudRunSync:
		s.CloudRunSyncStageOptions = &CloudRunSyncStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.CloudRunSyncStageOptions)
		}
	case model.StageCloudRunPromote:
		s.CloudRunPromoteStageOptions = &CloudRunPromoteStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.CloudRunPromoteStageOptions)
		}

	case model.StageLambdaSync:
		s.LambdaSyncStageOptions = &LambdaSyncStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.LambdaSyncStageOptions)
		}
	case model.StageLambdaPromote:
		s.LambdaPromoteStageOptions = &LambdaPromoteStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.LambdaPromoteStageOptions)
		}
	case model.StageLambdaCanaryRollout:
		s.LambdaCanaryRolloutStageOptions = &LambdaCanaryRolloutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.LambdaCanaryRolloutStageOptions)
		}

	case model.StageECSSync:
		s.ECSSyncStageOptions = &ECSSyncStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.ECSSyncStageOptions)
		}
	case model.StageECSCanaryRollout:
		s.ECSCanaryRolloutStageOptions = &ECSCanaryRolloutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.ECSCanaryRolloutStageOptions)
		}
	case model.StageECSPrimaryRollout:
		s.ECSPrimaryRolloutStageOptions = &ECSPrimaryRolloutStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.ECSPrimaryRolloutStageOptions)
		}
	case model.StageECSCanaryClean:
		s.ECSCanaryCleanStageOptions = &ECSCanaryCleanStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.ECSCanaryCleanStageOptions)
		}
	case model.StageECSTrafficRouting:
		s.ECSTrafficRoutingStageOptions = &ECSTrafficRoutingStageOptions{}
		if len(gs.With) > 0 {
			err = json.Unmarshal(gs.With, s.ECSTrafficRoutingStageOptions)
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
	// The maximum length of time to wait before giving up.
	// Defaults to 6h.
	Timeout        Duration `json:"timeout" default:"6h"`
	Approvers      []string `json:"approvers"`
	MinApproverNum int      `json:"minApproverNum" default:"1"`
}

func (w *WaitApprovalStageOptions) Validate() error {
	if w.MinApproverNum < 1 {
		return fmt.Errorf("minApproverNum %d should be greater than 0", w.MinApproverNum)
	}
	return nil
}

type CustomSyncOptions struct {
	Timeout Duration          `json:"timeout" default:"6h"`
	Envs    map[string]string `json:"envs"`
	Run     string            `json:"run"`
}

func (c *CustomSyncOptions) Validate() error {
	if c.Run == "" {
		return fmt.Errorf("the CUSTOM_SYNC stage requires run field")
	}
	return nil
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
	HTTPS            []TemplatableAnalysisHTTP    `json:"https"`
}

func (a *AnalysisStageOptions) Validate() error {
	if a.Duration == 0 {
		return fmt.Errorf("the ANALYSIS stage requires duration field")
	}

	for _, m := range a.Metrics {
		if m.Template.Name != "" {
			if err := m.Template.Validate(); err != nil {
				return fmt.Errorf("one of metrics configurations of ANALYSIS stage is invalid: %w", err)
			}
			continue
		}
		if err := m.AnalysisMetrics.Validate(); err != nil {
			return fmt.Errorf("one of metrics configurations of ANALYSIS stage is invalid: %w", err)
		}
	}

	for _, l := range a.Logs {
		if l.Template.Name != "" {
			if err := l.Template.Validate(); err != nil {
				return fmt.Errorf("one of log configurations of ANALYSIS stage is invalid: %w", err)
			}
			continue
		}
		if err := l.AnalysisLog.Validate(); err != nil {
			return fmt.Errorf("one of log configurations of ANALYSIS stage is invalid: %w", err)
		}
	}
	for _, h := range a.HTTPS {
		if h.Template.Name != "" {
			if err := h.Template.Validate(); err != nil {
				return fmt.Errorf("one of http configurations of ANALYSIS stage is invalid: %w", err)
			}
			continue
		}
		if err := h.AnalysisHTTP.Validate(); err != nil {
			return fmt.Errorf("one of http configurations of ANALYSIS stage is invalid: %w", err)
		}
	}
	return nil
}

type AnalysisTemplateRef struct {
	Name    string            `json:"name"`
	AppArgs map[string]string `json:"appArgs"`
}

func (a *AnalysisTemplateRef) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("the reference of analysis template name is empty")
	}
	return nil
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

type SecretEncryption struct {
	// List of encrypted secrets.
	EncryptedSecrets map[string]string `json:"encryptedSecrets"`
	// List of files to be decrypted before using.
	DecryptionTargets []string `json:"decryptionTargets"`
}

func (e *SecretEncryption) Validate() error {
	for k, v := range e.EncryptedSecrets {
		if k == "" {
			return fmt.Errorf("key field in encryptedSecrets must not be empty")
		}
		if v == "" {
			return fmt.Errorf("value field of %s in encryptedSecrets must not be empty", k)
		}
	}
	return nil
}

type Attachment struct {
	// Map of name to refer with the file path which contain embedding source data.
	Sources map[string]string `json:"sources"`
	// List of files to be embedded before using.
	Targets []string `json:"targets"`
}

func (e *Attachment) Validate() error {
	for k, v := range e.Sources {
		if k == "" {
			return fmt.Errorf("key field in sources must not be empty")
		}
		if v == "" {
			return fmt.Errorf("value field in sources must not be empty")
		}
	}
	return nil
}

// DeploymentNotification represents the way to send to users.
type DeploymentNotification struct {
	// List of users to be notified for each event.
	Mentions []NotificationMention `json:"mentions"`
}

func (n *DeploymentNotification) FindSlackAccounts(event model.NotificationEventType) []string {
	as := make(map[string]struct{})
	for _, m := range n.Mentions {
		if m.Event != allEventsSymbol && "EVENT_"+m.Event != event.String() {
			continue
		}
		for _, s := range m.Slack {
			as[s] = struct{}{}
		}
	}

	approvers := make([]string, 0, len(as))
	for a := range as {
		approvers = append(approvers, a)
	}
	return approvers
}

type NotificationMention struct {
	// The event to be notified to users.
	Event string `json:"event"`
	// List of user IDs for mentioning in Slack.
	// See https://api.slack.com/reference/surfaces/formatting#mentioning-users
	// for more information on how to check them.
	Slack []string `json:"slack"`
	// TODO: Support for email notification
	// The email for notification.
	Email []string `json:"email"`
}

func (n *NotificationMention) Validate() error {
	if n.Event == allEventsSymbol {
		return nil
	}

	e := "EVENT_" + n.Event
	for k := range model.NotificationEventType_value {
		if e == k {
			return nil
		}
	}
	return fmt.Errorf("event %q is incorrect as NotificationEventType", n.Event)
}

// PostSync provides all configurations to be used once the current deployment
// is triggered successfully.
type PostSync struct {
	DeploymentChain *DeploymentChain `json:"chain"`
}

func (p *PostSync) Validate() error {
	if dc := p.DeploymentChain; dc != nil {
		return dc.Validate()
	}
	return nil
}

// DeploymentChain provides all configurations used to trigger a chain of deployments.
type DeploymentChain struct {
	// ApplicationMatchers provides list of ChainApplicationMatcher which contain filters to be used
	// to find applications to deploy as chain node. It's required to not empty.
	ApplicationMatchers []ChainApplicationMatcher `json:"applications"`
	// Conditions provides configuration used to determine should the piped in charge in
	// the first applications in the chain trigger a whole new deployment chain or not.
	// If this field is not set, always trigger a whole new deployment chain when the current
	// application is triggered.
	// TODO: Add conditions to deployment chain configuration.
	// Conditions *DeploymentChainTriggerCondition `json:"conditions,omitempty"`
}

func (dc *DeploymentChain) Validate() error {
	if len(dc.ApplicationMatchers) == 0 {
		return fmt.Errorf("missing specified applications that will be triggered on this chain of deployment")
	}

	for _, m := range dc.ApplicationMatchers {
		if err := m.Validate(); err != nil {
			return err
		}
	}

	// if cc := dc.Conditions; cc != nil {
	// 	if err := cc.Validate(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// ChainApplicationMatcher provides filters used to find the right applications to trigger
// as a part of the deployment chain.
type ChainApplicationMatcher struct {
	Name   string            `json:"name"`
	Kind   string            `json:"kind"`
	Labels map[string]string `json:"labels"`
}

func (m *ChainApplicationMatcher) Validate() error {
	hasFilterCond := m.Name != "" || m.Kind != "" || len(m.Labels) != 0

	if !hasFilterCond {
		return fmt.Errorf("at least one of \"name\", \"kind\" or \"labels\" must be set to find applications to deploy")
	}
	return nil
}

type DeploymentChainTriggerCondition struct {
	CommitPrefix string `json:"commitPrefix"`
}

func (c *DeploymentChainTriggerCondition) Validate() error {
	hasCond := c.CommitPrefix != ""
	if !hasCond {
		return fmt.Errorf("missing commitPrefix configration as deployment chain trigger condition")
	}
	return nil
}

type DriftDetection struct {
	// IgnoreFields are a list of 'apiVersion:kind:namespace:name#fieldPath'
	IgnoreFields []string `json:"ignoreFields"`
}

func (dd *DriftDetection) Validate() error {
	for _, ignoreField := range dd.IgnoreFields {
		splited := strings.Split(ignoreField, "#")
		if len(splited) != 2 {
			return fmt.Errorf("ignoreFields must be in the form of 'apiVersion:kind:namespace:name#fieldPath'")
		}
	}
	return nil
}

func LoadApplication(repoPath, configRelPath string, appKind model.ApplicationKind) (*GenericApplicationSpec, error) {
	var absPath = filepath.Join(repoPath, configRelPath)

	cfg, err := LoadFromYAML(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("application config file %s was not found in Git", configRelPath)
		}
		return nil, err
	}
	if kind, ok := cfg.Kind.ToApplicationKind(); !ok || kind != appKind {
		return nil, fmt.Errorf("invalid application kind in the application config file, got: %s, expected: %s", kind, appKind)
	}

	spec, ok := cfg.GetGenericApplication()
	if !ok {
		return nil, fmt.Errorf("unsupported application kind: %s", appKind)
	}

	return &spec, nil
}
