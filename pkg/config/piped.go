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
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

var DefaultKubernetesCloudProvider = PipedCloudProvider{
	Name:             "kubernetes-default",
	Type:             model.CloudProviderKubernetes,
	KubernetesConfig: &CloudProviderKubernetesConfig{},
}

// PipedSpec contains configurable data used to while running Piped.
type PipedSpec struct {
	// The identifier of the project which this piped belongs to.
	ProjectID string
	// The unique identifier generated for this piped.
	PipedID string
	// The path to the key generated for this piped.
	PipedKeyFile string
	WebURL       string `json:"webURL"`
	// How often to check whether an application should be synced.
	// Default is 1m.
	SyncInterval Duration `json:"syncInterval"`
	// Git configuration needed for git commands.
	Git PipedGit `json:"git"`
	// List of git repositories this piped will handle.
	Repositories []PipedRepository `json:"repositories"`
	// List of helm chart repositories that should be added while starting up.
	ChartRepositories []HelmChartRepository   `json:"chartRepositories"`
	CloudProviders    []PipedCloudProvider    `json:"cloudProviders"`
	AnalysisProviders []PipedAnalysisProvider `json:"analysisProviders"`
	Notifications     Notifications           `json:"notifications"`
}

// Validate validates configured data of all fields.
func (s *PipedSpec) Validate() error {
	if s.ProjectID == "" {
		return fmt.Errorf("projectID must be set")
	}
	if s.PipedID == "" {
		return fmt.Errorf("pipedID must be set")
	}
	if s.PipedKeyFile == "" {
		return fmt.Errorf("pipedKeyFile must be set")
	}

	if s.SyncInterval < 0 {
		s.SyncInterval = Duration(time.Minute)
	}
	return nil
}

// EnableDefaultKubernetesCloudProvider adds the default kubernetes cloud provider if it was not specified.
func (s *PipedSpec) EnableDefaultKubernetesCloudProvider() {
	for _, cp := range s.CloudProviders {
		if cp.Name == DefaultKubernetesCloudProvider.Name {
			return
		}
	}
	s.CloudProviders = append(s.CloudProviders, DefaultKubernetesCloudProvider)
}

// HasCloudProvider checks whether the given provider is configured or not.
func (s *PipedSpec) HasCloudProvider(name string, t model.CloudProviderType) bool {
	for _, cp := range s.CloudProviders {
		if cp.Name != name {
			continue
		}
		if cp.Type != t {
			continue
		}
		return true
	}
	return false
}

// FindCloudProvider finds and returns a Cloud Provider by name and type.
func (s *PipedSpec) FindCloudProvider(name string, t model.CloudProviderType) (PipedCloudProvider, bool) {
	for _, p := range s.CloudProviders {
		if p.Name != name {
			continue
		}
		if p.Type != t {
			continue
		}
		return p, true
	}
	return PipedCloudProvider{}, false
}

// GetRepositoryMap returns a map of repositories where key is repo id.
func (s *PipedSpec) GetRepositoryMap() map[string]PipedRepository {
	m := make(map[string]PipedRepository, len(s.Repositories))
	for _, repo := range s.Repositories {
		m[repo.RepoID] = repo
	}
	return m
}

// GetRepository finds a repository with the given ID from the configured list.
func (s *PipedSpec) GetRepository(id string) (PipedRepository, bool) {
	for _, repo := range s.Repositories {
		if repo.RepoID == id {
			return repo, true
		}
	}
	return PipedRepository{}, false
}

// GetAnalysisProvider finds and returns an Analysis Provider config whose name is the given string.
func (s *PipedSpec) GetAnalysisProvider(name string) (PipedAnalysisProvider, bool) {
	for _, p := range s.AnalysisProviders {
		if p.Name == name {
			return p, true
		}
	}
	return PipedAnalysisProvider{}, false
}

type PipedGit struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	// Where to write ssh config file.
	// Default is "/etc/ssh/ssh_config".
	SSHConfigFilePath string `json:"sshConfigFilePath"`
	// The host name.
	// e.g. github.com, gitlab.com
	// Default is "github.com".
	Host string `json:"host"`
	// The hostname or IP address of the remote git server.
	// e.g. github.com, gitlab.com
	// Default is the same value with Host.
	HostName string `json:"hostName"`
	// The path to the private ssh key file.
	// This will be used to clone the source code of the git repositories.
	SSHKeyFile string `json:"sshKeyFile"`
}

func (g PipedGit) ShouldConfigureSSHConfig() bool {
	return g.SSHKeyFile != ""
}

type PipedRepository struct {
	// Unique identifier for this repository.
	// This must be unique in the piped scope.
	RepoID string `json:"repoId"`
	// Remote address of the repository.
	// e.g. git@github.com:org/repo1.git
	Remote string `json:"remote"`
	// The branch should be tracked.
	Branch string `json:"branch"`
}

type HelmChartRepository struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type PipedCloudProvider struct {
	Name string
	Type model.CloudProviderType

	KubernetesConfig *CloudProviderKubernetesConfig
	TerraformConfig  *CloudProviderTerraformConfig
	CloudRunConfig   *CloudProviderCloudRunConfig
	LambdaConfig     *CloudProviderLambdaConfig
}

type genericPipedCloudProvider struct {
	Name   string                  `json:"name"`
	Type   model.CloudProviderType `json:"type"`
	Config json.RawMessage         `json:"config"`
}

func (p *PipedCloudProvider) UnmarshalJSON(data []byte) error {
	var err error
	gp := genericPipedCloudProvider{}
	if err = json.Unmarshal(data, &gp); err != nil {
		return err
	}
	p.Name = gp.Name
	p.Type = gp.Type

	switch p.Type {
	case model.CloudProviderKubernetes:
		p.KubernetesConfig = &CloudProviderKubernetesConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.KubernetesConfig)
		}
	case model.CloudProviderTerraform:
		p.TerraformConfig = &CloudProviderTerraformConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.TerraformConfig)
		}
	case model.CloudProviderCloudRun:
		p.CloudRunConfig = &CloudProviderCloudRunConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.CloudRunConfig)
		}
	case model.CloudProviderLambda:
		p.LambdaConfig = &CloudProviderLambdaConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.LambdaConfig)
		}
	default:
		err = fmt.Errorf("unsupported cloud provider type: %s", p.Name)
	}
	return err
}

type CloudProviderKubernetesConfig struct {
	AppStateInformer KubernetesAppStateInformer `json:"appStateInformer"`
	MasterURL        string                     `json:"masterURL"`
	KubeConfigPath   string                     `json:"kubeConfigPath"`
}

type KubernetesAppStateInformer struct {
	// Only watches the specified namespace.
	// Empty means watching all namespaces.
	Namespace string `json:"namespace"`
	// List of resources that should be added to the watching targets.
	IncludeResources []KubernetesResourceMatcher `json:"includeResources"`
	// List of resources that should be ignored from the watching targets.
	ExcludeResources []KubernetesResourceMatcher `json:"excludeResources"`
}

type KubernetesResourceMatcher struct {
	// APIVersion of kubernetes resource.
	APIVersion string `json:"apiVersion"`
	// Kind name of kubernetes resource.
	// Empty means all kinds are matching.
	Kind string `json:"kind"`
}

type CloudProviderTerraformConfig struct {
	GCP *CloudProviderTerraformGCP `json:"gcp"`
	AWS *CloudProviderTerraformAWS `json:"aws"`
}

type CloudProviderTerraformGCP struct {
	Project         string `json:"project"`
	Region          string `json:"region"`
	CredentialsFile string `json:"credentialsFile"`
}

type CloudProviderTerraformAWS struct {
	Region string `json:"region"`
}

type CloudProviderCloudRunConfig struct {
	Project         string `json:"project"`
	Region          string `json:"region"`
	CredentialsFile string `json:"credentialsFile"`
}

type CloudProviderLambdaConfig struct {
	Region string `json:"region"`
}

type PipedAnalysisProvider struct {
	Name string                     `json:"name"`
	Type model.AnalysisProviderType `json:"type"`

	PrometheusConfig  *AnalysisProviderPrometheusConfig  `json:"prometheus"`
	DatadogConfig     *AnalysisProviderDatadogConfig     `json:"datadog"`
	StackdriverConfig *AnalysisProviderStackdriverConfig `json:"stackdriver"`
}

type genericPipedAnalysisProvider struct {
	Name   string                     `json:"name"`
	Type   model.AnalysisProviderType `json:"type"`
	Config json.RawMessage            `json:"config"`
}

func (p *PipedAnalysisProvider) UnmarshalJSON(data []byte) error {
	var err error
	gp := genericPipedAnalysisProvider{}
	if err = json.Unmarshal(data, &gp); err != nil {
		return err
	}
	p.Name = gp.Name
	p.Type = gp.Type

	switch p.Type {
	case model.AnalysisProviderPrometheus:
		p.PrometheusConfig = &AnalysisProviderPrometheusConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.PrometheusConfig)
		}
	case model.AnalysisProviderDatadog:
		p.DatadogConfig = &AnalysisProviderDatadogConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.DatadogConfig)
		}
	case model.AnalysisProviderStackdriver:
		p.StackdriverConfig = &AnalysisProviderStackdriverConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.StackdriverConfig)
		}
	default:
		err = fmt.Errorf("unsupported analysis provider type: %s", p.Name)
	}
	return err
}

type AnalysisProviderPrometheusConfig struct {
	Address string `json:"address"`
	// The path to the username file.
	UsernameFile string `json:"usernameFile"`
	// The path to the password file.
	PasswordFile string `json:"passwordFile"`
}

type AnalysisProviderDatadogConfig struct {
	Address string `json:"address"`
	// The path to the api key file.
	APIKeyFile string `json:"apiKeyFile"`
	// The path to the application key file.
	ApplicationKeyFile string `json:"applicationKeyFile"`
}

type AnalysisProviderStackdriverConfig struct {
	// The path to the service account file.
	ServiceAccountFile string `json:"serviceAccountFile"`
}

type Notifications struct {
	Routes    []NotificationRoute    `json:"routes"`
	Receivers []NotificationReceiver `json:"receivers"`
}

type NotificationRoute struct {
	Name         string   `json:"name"`
	Events       []string `json:"events"`
	IgnoreEvents []string `json:"ignoreEvents"`
	Groups       []string `json:"groups"`
	IgnoreGroups []string `json:"ignoreGroups"`
	Apps         []string `json:"apps"`
	IgnoreApps   []string `json:"ignoreApps"`
	Envs         []string `json:"envs"`
	IgnoreEnvs   []string `json:"ignoreEnvs"`
	Receiver     string   `json:"receiver"`
}

type NotificationReceiver struct {
	Name    string                       `json:"name"`
	Slack   *NotificationReceiverSlack   `json:"slack"`
	Webhook *NotificationReceiverWebhook `json:"webhook"`
}

type NotificationReceiverSlack struct {
	HookURL string `json:"hookURL"`
}

type NotificationReceiverWebhook struct {
	URL string `json:"url"`
}
