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

package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	maskString = "******"
)

var defaultKubernetesPlatformProvider = PipedPlatformProvider{
	Name:             "kubernetes-default",
	Type:             model.PlatformProviderKubernetes,
	KubernetesConfig: &PlatformProviderKubernetesConfig{},
}

// PipedSpec contains configurable data used to while running Piped.
type PipedSpec struct {
	// The identifier of the PipeCD project where this piped belongs to.
	ProjectID string `json:"projectID"`
	// The unique identifier generated for this piped.
	PipedID string `json:"pipedID"`
	// The path to the file containing the generated Key string for this piped.
	PipedKeyFile string `json:"pipedKeyFile,omitempty"`
	// Base64 encoded string of Piped key.
	PipedKeyData string `json:"pipedKeyData,omitempty"`
	// The name of this piped.
	Name string `json:"name,omitempty"`
	// The address used to connect to the control-plane's API.
	APIAddress string `json:"apiAddress"`
	// The address to the control-plane's Web.
	WebAddress string `json:"webAddress,omitempty"`
	// How often to check whether an application should be synced.
	// Default is 1m.
	SyncInterval Duration `json:"syncInterval,omitempty" default:"1m"`
	// How often to check whether an application configuration file should be synced.
	// Default is 1m.
	AppConfigSyncInterval Duration `json:"appConfigSyncInterval,omitempty" default:"1m"`
	// Git configuration needed for git commands.
	Git PipedGit `json:"git,omitempty"`
	// List of git repositories this piped will handle.
	Repositories []PipedRepository `json:"repositories,omitempty"`
	// List of helm chart repositories that should be added while starting up.
	ChartRepositories []HelmChartRepository `json:"chartRepositories,omitempty"`
	// List of helm chart registries that should be logged in while starting up.
	ChartRegistries []HelmChartRegistry `json:"chartRegistries,omitempty"`
	// List of cloud providers can be used by this piped.
	// Deprecated: use PlatformProvider instead.
	CloudProviders []PipedPlatformProvider `json:"cloudProviders,omitempty"`
	// List of platform providers can be used by this piped.
	PlatformProviders []PipedPlatformProvider `json:"platformProviders,omitempty"`
	// List of plugiin configs
	Plugins []PipedPlugin `json:"plugins,omitempty"`
	// List of analysis providers can be used by this piped.
	AnalysisProviders []PipedAnalysisProvider `json:"analysisProviders,omitempty"`
	// Sending notification to Slack, Webhook…
	Notifications Notifications `json:"notifications"`
	// What secret management method should be used.
	SecretManagement *SecretManagement `json:"secretManagement,omitempty"`
	// Optional settings for event watcher.
	EventWatcher PipedEventWatcher `json:"eventWatcher"`
	// List of labels to filter all applications this piped will handle.
	AppSelector map[string]string `json:"appSelector,omitempty"`
}

func (s *PipedSpec) UnmarshalJSON(data []byte) error {
	type Alias PipedSpec
	ps := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &ps); err != nil {
		return err
	}

	// Add all CloudProviders configuration as PlatformProviders configuration.
	s.PlatformProviders = append(s.PlatformProviders, ps.CloudProviders...)
	s.CloudProviders = nil
	return nil
}

// Validate validates configured data of all fields.
func (s *PipedSpec) Validate() error {
	if s.ProjectID == "" {
		return errors.New("projectID must be set")
	}
	if s.PipedID == "" {
		return errors.New("pipedID must be set")
	}
	if s.PipedKeyData == "" && s.PipedKeyFile == "" {
		return errors.New("either pipedKeyFile or pipedKeyData must be set")
	}
	if s.PipedKeyData != "" && s.PipedKeyFile != "" {
		return errors.New("only pipedKeyFile or pipedKeyData can be set")
	}
	if s.APIAddress == "" {
		return errors.New("apiAddress must be set")
	}
	if s.SyncInterval < 0 {
		return errors.New("syncInterval must be greater than or equal to 0")
	}
	if err := s.Git.Validate(); err != nil {
		return err
	}
	for _, r := range s.ChartRepositories {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	for _, r := range s.ChartRegistries {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	if s.SecretManagement != nil {
		if err := s.SecretManagement.Validate(); err != nil {
			return err
		}
	}
	if err := s.EventWatcher.Validate(); err != nil {
		return err
	}
	for _, n := range s.Notifications.Receivers {
		if n.Slack != nil {
			if err := n.Slack.Validate(); err != nil {
				return err
			}
		}
	}
	for _, p := range s.AnalysisProviders {
		if err := p.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Clone generates a cloned PipedSpec object.
func (s *PipedSpec) Clone() (*PipedSpec, error) {
	js, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var clone PipedSpec
	if err = json.Unmarshal(js, &clone); err != nil {
		return nil, err
	}

	return &clone, nil
}

// Mask masks confidential fields.
func (s *PipedSpec) Mask() {
	if len(s.PipedKeyFile) != 0 {
		s.PipedKeyFile = maskString
	}
	if len(s.PipedKeyData) != 0 {
		s.PipedKeyData = maskString
	}
	s.Git.Mask()
	for i := 0; i < len(s.ChartRepositories); i++ {
		s.ChartRepositories[i].Mask()
	}
	for i := 0; i < len(s.ChartRegistries); i++ {
		s.ChartRegistries[i].Mask()
	}
	for _, p := range s.PlatformProviders {
		p.Mask()
	}
	for _, p := range s.AnalysisProviders {
		p.Mask()
	}
	s.Notifications.Mask()
	if s.SecretManagement != nil {
		s.SecretManagement.Mask()
	}
}

// EnableDefaultKubernetesPlatformProvider adds the default kubernetes cloud provider if it was not specified.
func (s *PipedSpec) EnableDefaultKubernetesPlatformProvider() {
	for _, cp := range s.PlatformProviders {
		if cp.Name == defaultKubernetesPlatformProvider.Name {
			return
		}
	}
	s.PlatformProviders = append(s.PlatformProviders, defaultKubernetesPlatformProvider)
}

// HasPlatformProvider checks whether the given provider is configured or not.
func (s *PipedSpec) HasPlatformProvider(name string, t model.ApplicationKind) bool {
	_, contains := s.FindPlatformProvider(name, t)
	return contains
}

// FindPlatformProvider finds and returns a Platform Provider by name and type.
func (s *PipedSpec) FindPlatformProvider(name string, t model.ApplicationKind) (PipedPlatformProvider, bool) {
	requiredProviderType := t.CompatiblePlatformProviderType()
	for _, p := range s.PlatformProviders {
		if p.Name != name {
			continue
		}
		if p.Type != requiredProviderType {
			continue
		}
		return p, true
	}
	return PipedPlatformProvider{}, false
}

// FindPlatformProvidersByLabels finds all PlatformProviders which match the provided labels.
func (s *PipedSpec) FindPlatformProvidersByLabels(labels map[string]string, t model.ApplicationKind) []PipedPlatformProvider {
	requiredProviderType := t.CompatiblePlatformProviderType()
	out := make([]PipedPlatformProvider, 0)

	labelMatch := func(providerLabels map[string]string) bool {
		if len(providerLabels) < len(labels) {
			return false
		}

		for k, v := range labels {
			if v != providerLabels[k] {
				return false
			}
		}
		return true
	}

	for _, p := range s.PlatformProviders {
		if p.Type != requiredProviderType {
			continue
		}
		if !labelMatch(p.Labels) {
			continue
		}
		out = append(out, p)
	}
	return out
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

func (s *PipedSpec) IsInsecureChartRepository(name string) bool {
	for _, cr := range s.ChartRepositories {
		if cr.Name == name {
			return cr.Insecure
		}
	}
	return false
}

func (s *PipedSpec) LoadPipedKey() ([]byte, error) {
	if s.PipedKeyData != "" {
		return base64.StdEncoding.DecodeString(s.PipedKeyData)
	}
	if s.PipedKeyFile != "" {
		return os.ReadFile(s.PipedKeyFile)
	}
	return nil, errors.New("either pipedKeyFile or pipedKeyData must be set")
}

type PipedGit struct {
	// The username that will be configured for `git` user.
	// Default is "piped".
	Username string `json:"username,omitempty"`
	// The email that will be configured for `git` user.
	// Default is "pipecd.dev@gmail.com".
	Email string `json:"email,omitempty"`
	// Where to write ssh config file.
	// Default is "$HOME/.ssh/config".
	SSHConfigFilePath string `json:"sshConfigFilePath,omitempty"`
	// The host name.
	// e.g. github.com, gitlab.com
	// Default is "github.com".
	Host string `json:"host,omitempty"`
	// The hostname or IP address of the remote git server.
	// e.g. github.com, gitlab.com
	// Default is the same value with Host.
	HostName string `json:"hostName,omitempty"`
	// The path to the private ssh key file.
	// This will be used to clone the source code of the specified git repositories.
	SSHKeyFile string `json:"sshKeyFile,omitempty"`
	// Base64 encoded string of ssh-key.
	SSHKeyData string `json:"sshKeyData,omitempty"`
	// Base64 encoded string of password.
	// This will be used to clone the source repo with https basic auth.
	Password string `json:"password,omitempty"`
}

func (g PipedGit) ShouldConfigureSSHConfig() bool {
	return g.SSHKeyData != "" || g.SSHKeyFile != ""
}

func (g PipedGit) LoadSSHKey() ([]byte, error) {
	if g.SSHKeyData != "" && g.SSHKeyFile != "" {
		return nil, errors.New("only either sshKeyFile or sshKeyData can be set")
	}
	if g.SSHKeyData != "" {
		return base64.StdEncoding.DecodeString(g.SSHKeyData)
	}
	if g.SSHKeyFile != "" {
		return os.ReadFile(g.SSHKeyFile)
	}
	return nil, errors.New("either sshKeyFile or sshKeyData must be set")
}

func (g *PipedGit) Validate() error {
	isPassword := g.Password != ""
	isSSH := g.ShouldConfigureSSHConfig()
	if isSSH && isPassword {
		return errors.New("cannot configure both sshKeyData or sshKeyFile and password authentication")
	}
	if isSSH && (g.SSHKeyData != "" && g.SSHKeyFile != "") {
		return errors.New("only either sshKeyFile or sshKeyData can be set")
	}
	if isPassword && (g.Username == "" || g.Password == "") {
		return errors.New("both username and password must be set")
	}
	return nil
}

func (g *PipedGit) Mask() {
	if len(g.SSHConfigFilePath) != 0 {
		g.SSHConfigFilePath = maskString
	}
	if len(g.SSHKeyFile) != 0 {
		g.SSHKeyFile = maskString
	}
	if len(g.SSHKeyData) != 0 {
		g.SSHKeyData = maskString
	}
	if len(g.Password) != 0 {
		g.Password = maskString
	}
}

func (g *PipedGit) DecodedPassword() (string, error) {
	if len(g.Password) == 0 {
		return "", nil
	}
	decoded, err := base64.StdEncoding.DecodeString(g.Password)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

type PipedRepository struct {
	// Unique identifier for this repository.
	// This must be unique in the piped scope.
	RepoID string `json:"repoId"`
	// Remote address of the repository used to clone the source code.
	// e.g. git@github.com:org/repo.git
	Remote string `json:"remote"`
	// The branch will be handled.
	Branch string `json:"branch"`
}

type HelmChartRepositoryType string

const (
	HTTPHelmChartRepository HelmChartRepositoryType = "HTTP"
	GITHelmChartRepository  HelmChartRepositoryType = "GIT"
)

type HelmChartRepository struct {
	// The repository type. Currently, HTTP and GIT are supported.
	// Default is HTTP.
	Type HelmChartRepositoryType `json:"type" default:"HTTP"`

	// Configuration for HTTP type.
	// The name of the Helm chart repository.
	Name string `json:"name,omitempty"`
	// The address to the Helm chart repository.
	Address string `json:"address,omitempty"`
	// Username used for the repository backed by HTTP basic authentication.
	Username string `json:"username,omitempty"`
	// Password used for the repository backed by HTTP basic authentication.
	Password string `json:"password,omitempty"`
	// Whether to skip TLS certificate checks for the repository or not.
	Insecure bool `json:"insecure"`

	// Configuration for GIT type.
	// Remote address of the Git repository used to clone Helm charts.
	// e.g. git@github.com:org/repo.git
	GitRemote string `json:"gitRemote,omitempty"`
	// The path to the private ssh key file used while cloning Helm charts from above Git repository.
	SSHKeyFile string `json:"sshKeyFile,omitempty"`
}

func (r *HelmChartRepository) IsHTTPRepository() bool {
	return r.Type == HTTPHelmChartRepository
}

func (r *HelmChartRepository) IsGitRepository() bool {
	return r.Type == GITHelmChartRepository
}

func (r *HelmChartRepository) Validate() error {
	if r.IsHTTPRepository() {
		if r.Name == "" {
			return errors.New("name must be set")
		}
		if r.Address == "" {
			return errors.New("address must be set")
		}
		return nil
	}

	if r.IsGitRepository() {
		if r.GitRemote == "" {
			return errors.New("gitRemote must be set")
		}
		return nil
	}

	return fmt.Errorf("either %s repository or %s repository must be configured", HTTPHelmChartRepository, GITHelmChartRepository)
}

func (r *HelmChartRepository) Mask() {
	if len(r.Password) != 0 {
		r.Password = maskString
	}
	if len(r.SSHKeyFile) != 0 {
		r.SSHKeyFile = maskString
	}
}

func (s *PipedSpec) HTTPHelmChartRepositories() []HelmChartRepository {
	repos := make([]HelmChartRepository, 0, len(s.ChartRepositories))
	for _, r := range s.ChartRepositories {
		if r.IsHTTPRepository() {
			repos = append(repos, r)
		}
	}
	return repos
}

func (s *PipedSpec) GitHelmChartRepositories() []HelmChartRepository {
	repos := make([]HelmChartRepository, 0, len(s.ChartRepositories))
	for _, r := range s.ChartRepositories {
		if r.IsGitRepository() {
			repos = append(repos, r)
		}
	}
	return repos
}

type HelmChartRegistryType string

// The registry types that hosts Helm charts.
const (
	OCIHelmChartRegistry HelmChartRegistryType = "OCI"
)

type HelmChartRegistry struct {
	// The registry type. Currently, only OCI is supported.
	Type HelmChartRegistryType `json:"type" default:"OCI"`

	// The address to the Helm chart registry.
	Address string `json:"address"`
	// Username used for the registry authentication.
	Username string `json:"username,omitempty"`
	// Password used for the registry authentication.
	Password string `json:"password,omitempty"`
}

func (r *HelmChartRegistry) IsOCI() bool {
	return r.Type == OCIHelmChartRegistry
}

func (r *HelmChartRegistry) Validate() error {
	if r.IsOCI() {
		if r.Address == "" {
			return errors.New("address must be set")
		}
		return nil
	}

	return fmt.Errorf("%s registry must be configured", OCIHelmChartRegistry)
}

func (r *HelmChartRegistry) Mask() {
	if len(r.Password) != 0 {
		r.Password = maskString
	}
}

type PipedPlatformProvider struct {
	Name   string                     `json:"name"`
	Type   model.PlatformProviderType `json:"type"`
	Labels map[string]string          `json:"labels,omitempty"`

	KubernetesConfig *PlatformProviderKubernetesConfig
	TerraformConfig  *PlatformProviderTerraformConfig
	CloudRunConfig   *PlatformProviderCloudRunConfig
	LambdaConfig     *PlatformProviderLambdaConfig
	ECSConfig        *PlatformProviderECSConfig
}

type genericPipedPlatformProvider struct {
	Name   string                     `json:"name"`
	Type   model.PlatformProviderType `json:"type"`
	Labels map[string]string          `json:"labels,omitempty"`
	Config json.RawMessage            `json:"config"`
}

func (p *PipedPlatformProvider) MarshalJSON() ([]byte, error) {
	var (
		err    error
		config json.RawMessage
	)

	switch p.Type {
	case model.PlatformProviderKubernetes:
		config, err = json.Marshal(p.KubernetesConfig)
	case model.PlatformProviderTerraform:
		config, err = json.Marshal(p.TerraformConfig)
	case model.PlatformProviderCloudRun:
		config, err = json.Marshal(p.CloudRunConfig)
	case model.PlatformProviderLambda:
		config, err = json.Marshal(p.LambdaConfig)
	case model.PlatformProviderECS:
		config, err = json.Marshal(p.ECSConfig)
	default:
		err = fmt.Errorf("unsupported platform provider type: %s", p.Name)
	}

	if err != nil {
		return nil, err
	}

	return json.Marshal(&genericPipedPlatformProvider{
		Name:   p.Name,
		Type:   p.Type,
		Labels: p.Labels,
		Config: config,
	})
}

func (p *PipedPlatformProvider) UnmarshalJSON(data []byte) error {
	var err error
	gp := genericPipedPlatformProvider{}
	if err = json.Unmarshal(data, &gp); err != nil {
		return err
	}
	p.Name = gp.Name
	p.Type = gp.Type
	p.Labels = gp.Labels

	switch p.Type {
	case model.PlatformProviderKubernetes:
		p.KubernetesConfig = &PlatformProviderKubernetesConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.KubernetesConfig)
		}
	case model.PlatformProviderTerraform:
		p.TerraformConfig = &PlatformProviderTerraformConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.TerraformConfig)
		}
	case model.PlatformProviderCloudRun:
		p.CloudRunConfig = &PlatformProviderCloudRunConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.CloudRunConfig)
		}
	case model.PlatformProviderLambda:
		p.LambdaConfig = &PlatformProviderLambdaConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.LambdaConfig)
		}
	case model.PlatformProviderECS:
		p.ECSConfig = &PlatformProviderECSConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.ECSConfig)
		}
	default:
		err = fmt.Errorf("unsupported platform provider type: %s", p.Name)
	}
	return err
}

func (p *PipedPlatformProvider) Mask() {
	if p.CloudRunConfig != nil {
		p.CloudRunConfig.Mask()
	}
	if p.LambdaConfig != nil {
		p.LambdaConfig.Mask()
	}
	if p.ECSConfig != nil {
		p.ECSConfig.Mask()
	}
}

type PlatformProviderKubernetesConfig struct {
	// The master URL of the kubernetes cluster.
	// Empty means in-cluster.
	MasterURL string `json:"masterURL,omitempty"`
	// The path to the kubeconfig file.
	// Empty means in-cluster.
	KubeConfigPath string `json:"kubeConfigPath,omitempty"`
	// Configuration for application resource informer.
	AppStateInformer KubernetesAppStateInformer `json:"appStateInformer"`
	// Version of kubectl will be used.
	KubectlVersion string `json:"kubectlVersion"`
}

type KubernetesAppStateInformer struct {
	// Only watches the specified namespace.
	// Empty means watching all namespaces.
	Namespace string `json:"namespace,omitempty"`
	// List of resources that should be added to the watching targets.
	IncludeResources []KubernetesResourceMatcher `json:"includeResources,omitempty"`
	// List of resources that should be ignored from the watching targets.
	ExcludeResources []KubernetesResourceMatcher `json:"excludeResources,omitempty"`
}

type KubernetesResourceMatcher struct {
	// The APIVersion of the kubernetes resource.
	APIVersion string `json:"apiVersion,omitempty"`
	// The kind name of the kubernetes resource.
	// Empty means all kinds are matching.
	Kind string `json:"kind,omitempty"`
}

type PlatformProviderTerraformConfig struct {
	// List of variables that will be set directly on terraform commands with "-var" flag.
	// The variable must be formatted by "key=value" as below:
	// "image_id=ami-abc123"
	// 'image_id_list=["ami-abc123","ami-def456"]'
	// 'image_id_map={"us-east-1":"ami-abc123","us-east-2":"ami-def456"}'
	Vars []string `json:"vars,omitempty"`
	// Enable drift detection.
	// TODO: This is a temporary option because Terraform drift detection is buggy and has performance issues. This will be possibly removed in the future release.
	DriftDetectionEnabled *bool `json:"driftDetectionEnabled" default:"true"`
}

type PlatformProviderCloudRunConfig struct {
	// The GCP project hosting the CloudRun service.
	Project string `json:"project"`
	// The region of running CloudRun service.
	Region string `json:"region"`
	// The path to the service account file for accessing CloudRun service.
	CredentialsFile string `json:"credentialsFile,omitempty"`
}

func (c *PlatformProviderCloudRunConfig) Mask() {
	if len(c.CredentialsFile) != 0 {
		c.CredentialsFile = maskString
	}
}

type PlatformProviderLambdaConfig struct {
	// The region to send requests to. This parameter is required.
	// e.g. "us-west-2"
	// A full list of regions is: https://docs.aws.amazon.com/general/latest/gr/rande.html
	Region string `json:"region"`
	// Path to the shared credentials file.
	CredentialsFile string `json:"credentialsFile,omitempty"`
	// The IAM role arn to use when assuming an role.
	RoleARN string `json:"roleARN,omitempty"`
	// Path to the WebIdentity token the SDK should use to assume a role with.
	TokenFile string `json:"tokenFile,omitempty"`
	// AWS Profile to extract credentials from the shared credentials file.
	// If empty, the environment variable "AWS_PROFILE" is used.
	// "default" is populated if the environment variable is also not set.
	Profile string `json:"profile,omitempty"`
}

func (c *PlatformProviderLambdaConfig) Mask() {
	if len(c.CredentialsFile) != 0 {
		c.CredentialsFile = maskString
	}
	if len(c.RoleARN) != 0 {
		c.RoleARN = maskString
	}
	if len(c.TokenFile) != 0 {
		c.TokenFile = maskString
	}
}

type PlatformProviderECSConfig struct {
	// The region to send requests to. This parameter is required.
	// e.g. "us-west-2"
	// A full list of regions is: https://docs.aws.amazon.com/general/latest/gr/rande.html
	Region string `json:"region"`
	// Path to the shared credentials file.
	CredentialsFile string `json:"credentialsFile,omitempty"`
	// The IAM role arn to use when assuming an role.
	RoleARN string `json:"roleARN,omitempty"`
	// Path to the WebIdentity token the SDK should use to assume a role with.
	TokenFile string `json:"tokenFile,omitempty"`
	// AWS Profile to extract credentials from the shared credentials file.
	// If empty, the environment variable "AWS_PROFILE" is used.
	// "default" is populated if the environment variable is also not set.
	Profile string `json:"profile,omitempty"`
}

func (c *PlatformProviderECSConfig) Mask() {
	if len(c.CredentialsFile) != 0 {
		c.CredentialsFile = maskString
	}
	if len(c.RoleARN) != 0 {
		c.RoleARN = maskString
	}
	if len(c.TokenFile) != 0 {
		c.TokenFile = maskString
	}
}

type PipedAnalysisProvider struct {
	Name string                     `json:"name"`
	Type model.AnalysisProviderType `json:"type"`

	PrometheusConfig  *AnalysisProviderPrometheusConfig
	DatadogConfig     *AnalysisProviderDatadogConfig
	StackdriverConfig *AnalysisProviderStackdriverConfig
}

func (p *PipedAnalysisProvider) Mask() {
	if p.PrometheusConfig != nil {
		p.PrometheusConfig.Mask()
	}
	if p.DatadogConfig != nil {
		p.DatadogConfig.Mask()
	}
	if p.StackdriverConfig != nil {
		p.StackdriverConfig.Mask()
	}
}

type genericPipedAnalysisProvider struct {
	Name   string                     `json:"name"`
	Type   model.AnalysisProviderType `json:"type"`
	Config json.RawMessage            `json:"config"`
}

func (p *PipedAnalysisProvider) MarshalJSON() ([]byte, error) {
	var (
		err    error
		config json.RawMessage
	)

	switch p.Type {
	case model.AnalysisProviderDatadog:
		config, err = json.Marshal(p.DatadogConfig)
	case model.AnalysisProviderPrometheus:
		config, err = json.Marshal(p.PrometheusConfig)
	case model.AnalysisProviderStackdriver:
		config, err = json.Marshal(p.StackdriverConfig)
	default:
		err = fmt.Errorf("unsupported analysis provider type: %s", p.Name)
	}

	if err != nil {
		return nil, err
	}

	return json.Marshal(&genericPipedAnalysisProvider{
		Name:   p.Name,
		Type:   p.Type,
		Config: config,
	})
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

func (p *PipedAnalysisProvider) Validate() error {
	switch p.Type {
	case model.AnalysisProviderPrometheus:
		return p.PrometheusConfig.Validate()
	case model.AnalysisProviderDatadog:
		return p.DatadogConfig.Validate()
	case model.AnalysisProviderStackdriver:
		return p.StackdriverConfig.Validate()
	default:
		return fmt.Errorf("unknow provider type: %s", p.Type)
	}
}

type AnalysisProviderPrometheusConfig struct {
	Address string `json:"address"`
	// The path to the username file.
	UsernameFile string `json:"usernameFile,omitempty"`
	// The path to the password file.
	PasswordFile string `json:"passwordFile,omitempty"`
}

func (a *AnalysisProviderPrometheusConfig) Validate() error {
	if a.Address == "" {
		return fmt.Errorf("prometheus analysis provider requires the address")
	}
	return nil
}

func (a *AnalysisProviderPrometheusConfig) Mask() {
	if len(a.PasswordFile) != 0 {
		a.PasswordFile = maskString
	}
}

type AnalysisProviderDatadogConfig struct {
	// The address of Datadog API server.
	// Only "datadoghq.com", "us3.datadoghq.com", "datadoghq.eu", "ddog-gov.com" are available.
	// Defaults to "datadoghq.com"
	Address string `json:"address,omitempty"`
	// Required: The path to the api key file.
	APIKeyFile string `json:"apiKeyFile"`
	// Required: The path to the application key file.
	ApplicationKeyFile string `json:"applicationKeyFile"`
	// Base64 API Key for Datadog API server.
	APIKeyData string `json:"apiKeyData,omitempty"`
	// Base64 Application Key for Datadog API server.
	ApplicationKeyData string `json:"applicationKeyData,omitempty"`
}

func (a *AnalysisProviderDatadogConfig) Validate() error {
	if a.APIKeyFile == "" && a.APIKeyData == "" {
		return fmt.Errorf("either datadog APIKeyFile or APIKeyData must be set")
	}
	if a.ApplicationKeyFile == "" && a.ApplicationKeyData == "" {
		return fmt.Errorf("either datadog ApplicationKeyFile or ApplicationKeyData must be set")
	}
	if a.APIKeyData != "" && a.APIKeyFile != "" {
		return fmt.Errorf("only datadog APIKeyFile or APIKeyData can be set")
	}
	if a.ApplicationKeyData != "" && a.ApplicationKeyFile != "" {
		return fmt.Errorf("only datadog ApplicationKeyFile or ApplicationKeyData can be set")
	}
	return nil
}

func (a *AnalysisProviderDatadogConfig) Mask() {
	if len(a.APIKeyFile) != 0 {
		a.APIKeyFile = maskString
	}
	if len(a.ApplicationKeyFile) != 0 {
		a.ApplicationKeyFile = maskString
	}
	if len(a.APIKeyData) != 0 {
		a.APIKeyData = maskString
	}
	if len(a.ApplicationKeyData) != 0 {
		a.ApplicationKeyData = maskString
	}
}

// func(a *AnalysisProviderDatadogConfig)

type AnalysisProviderStackdriverConfig struct {
	// The path to the service account file.
	ServiceAccountFile string `json:"serviceAccountFile"`
}

func (a *AnalysisProviderStackdriverConfig) Mask() {
	if len(a.ServiceAccountFile) != 0 {
		a.ServiceAccountFile = maskString
	}
}

func (a *AnalysisProviderStackdriverConfig) Validate() error {
	return nil
}

type Notifications struct {
	// List of notification routes.
	Routes []NotificationRoute `json:"routes,omitempty"`
	// List of notification receivers.
	Receivers []NotificationReceiver `json:"receivers,omitempty"`
}

func (n *Notifications) Mask() {
	for _, r := range n.Receivers {
		r.Mask()
	}
}

type NotificationRoute struct {
	Name         string            `json:"name"`
	Receiver     string            `json:"receiver"`
	Events       []string          `json:"events,omitempty"`
	IgnoreEvents []string          `json:"ignoreEvents,omitempty"`
	Groups       []string          `json:"groups,omitempty"`
	IgnoreGroups []string          `json:"ignoreGroups,omitempty"`
	Apps         []string          `json:"apps,omitempty"`
	IgnoreApps   []string          `json:"ignoreApps,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	IgnoreLabels map[string]string `json:"ignoreLabels,omitempty"`
}

type NotificationReceiver struct {
	Name    string                       `json:"name"`
	Slack   *NotificationReceiverSlack   `json:"slack,omitempty"`
	Webhook *NotificationReceiverWebhook `json:"webhook,omitempty"`
}

func (n *NotificationReceiver) Mask() {
	if n.Slack != nil {
		n.Slack.Mask()
	}
	if n.Webhook != nil {
		n.Webhook.Mask()
	}
}

type NotificationReceiverSlack struct {
	HookURL           string   `json:"hookURL"`
	OAuthToken        string   `json:"oauthToken"` // Deprecated: use OAuthTokenData instead.
	OAuthTokenData    string   `json:"oauthTokenData"`
	OAuthTokenFile    string   `json:"oauthTokenFile"`
	ChannelID         string   `json:"channelID"`
	MentionedAccounts []string `json:"mentionedAccounts,omitempty"`
	MentionedGroups   []string `json:"mentionedGroups,omitempty"`
}

func (n *NotificationReceiverSlack) Mask() {
	if len(n.HookURL) != 0 {
		n.HookURL = maskString
	}
	if len(n.OAuthToken) != 0 {
		n.OAuthToken = maskString
	}
	if len(n.OAuthTokenData) != 0 {
		n.OAuthTokenData = maskString
	}
}

func (n *NotificationReceiverSlack) Validate() error {
	mentionedAccounts := make([]string, 0, len(n.MentionedAccounts))
	for _, mentionedAccount := range n.MentionedAccounts {
		formatMentionedAccount := strings.TrimPrefix(mentionedAccount, "@")
		mentionedAccounts = append(mentionedAccounts, formatMentionedAccount)
	}
	mentionedGroups := make([]string, 0, len(n.MentionedGroups))
	for _, mentionedGroup := range n.MentionedGroups {
		if !strings.Contains(mentionedGroup, "!subteam^") {
			formatMentionedGroup := fmt.Sprintf("<!subteam^%s>", mentionedGroup)
			mentionedGroups = append(mentionedGroups, formatMentionedGroup)
		} else {
			mentionedGroups = append(mentionedGroups, mentionedGroup)
		}
	}
	if len(mentionedGroups) > 0 {
		n.MentionedGroups = mentionedGroups
	}
	if len(mentionedAccounts) > 0 {
		n.MentionedAccounts = mentionedAccounts
	}
	if n.HookURL != "" && (n.OAuthToken != "" || n.OAuthTokenFile != "" || n.OAuthTokenData != "" || n.ChannelID != "") {
		return errors.New("only one of sending via hook URL or API should be used")
	}
	if n.HookURL != "" {
		return nil
	}
	if n.ChannelID == "" || (n.OAuthToken == "" && n.OAuthTokenFile == "" && n.OAuthTokenData == "") {
		return errors.New("missing channelID or OAuth token configuration")
	}
	if (n.OAuthToken != "" && n.OAuthTokenFile != "") || (n.OAuthToken != "" && n.OAuthTokenData != "") || (n.OAuthTokenFile != "" && n.OAuthTokenData != "") {
		return errors.New("only one of OAuthToken, OAuthTokenData and OAuthTokenFile should be set")
	}
	return nil
}

type NotificationReceiverWebhook struct {
	URL                string `json:"url"`
	SignatureKey       string `json:"signatureKey,omitempty" default:"PipeCD-Signature"`
	SignatureValue     string `json:"signatureValue,omitempty"`
	SignatureValueFile string `json:"signatureValueFile,omitempty"`
}

func (n *NotificationReceiverWebhook) Mask() {
	if len(n.URL) != 0 {
		n.URL = maskString
	}
	if len(n.SignatureKey) != 0 {
		n.SignatureKey = maskString
	}
	if len(n.SignatureValue) != 0 {
		n.SignatureValue = maskString
	}
	if len(n.SignatureValueFile) != 0 {
		n.SignatureValueFile = maskString
	}
}

func (n *NotificationReceiverWebhook) LoadSignatureValue() (string, error) {
	if n.SignatureValue != "" && n.SignatureValueFile != "" {
		return "", errors.New("only either signatureValue or signatureValueFile can be set")
	}
	if n.SignatureValue != "" {
		return n.SignatureValue, nil
	}
	if n.SignatureValueFile != "" {
		val, err := os.ReadFile(n.SignatureValueFile)
		if err != nil {
			return "", err
		}
		return strings.TrimSuffix(string(val), "\n"), nil
	}
	return "", nil
}

type SecretManagement struct {
	// Which management service should be used.
	// Available values: KEY_PAIR, GCP_KMS, AWS_KMS
	Type model.SecretManagementType `json:"type"`

	KeyPair *SecretManagementKeyPair
	GCPKMS  *SecretManagementGCPKMS
}

type genericSecretManagement struct {
	Type   model.SecretManagementType `json:"type"`
	Config json.RawMessage            `json:"config"`
}

func (s *SecretManagement) MarshalJSON() ([]byte, error) {
	var (
		err    error
		config json.RawMessage
	)

	switch s.Type {
	case model.SecretManagementTypeKeyPair:
		config, err = json.Marshal(s.KeyPair)
	case model.SecretManagementTypeGCPKMS:
		config, err = json.Marshal(s.GCPKMS)
	default:
		err = fmt.Errorf("unsupported secret management type: %s", s.Type)
	}

	if err != nil {
		return nil, err
	}

	return json.Marshal(&genericSecretManagement{
		Type:   s.Type,
		Config: config,
	})
}

func (s *SecretManagement) UnmarshalJSON(data []byte) error {
	var err error
	g := genericSecretManagement{}
	if err = json.Unmarshal(data, &g); err != nil {
		return err
	}

	switch g.Type {
	case model.SecretManagementTypeKeyPair:
		s.Type = model.SecretManagementTypeKeyPair
		s.KeyPair = &SecretManagementKeyPair{}
		if len(g.Config) > 0 {
			err = json.Unmarshal(g.Config, s.KeyPair)
		}
	case model.SecretManagementTypeGCPKMS:
		s.Type = model.SecretManagementTypeGCPKMS
		s.GCPKMS = &SecretManagementGCPKMS{}
		if len(g.Config) > 0 {
			err = json.Unmarshal(g.Config, s.GCPKMS)
		}
	default:
		err = fmt.Errorf("unsupported secret management type: %s", s.Type)
	}
	return err
}

func (s *SecretManagement) Mask() {
	if s.KeyPair != nil {
		s.KeyPair.Mask()
	}
	if s.GCPKMS != nil {
		s.GCPKMS.Mask()
	}
}

func (s *SecretManagement) Validate() error {
	switch s.Type {
	case model.SecretManagementTypeKeyPair:
		return s.KeyPair.Validate()
	case model.SecretManagementTypeGCPKMS:
		return s.GCPKMS.Validate()
	default:
		return fmt.Errorf("unsupported sealed secret management type: %s", s.Type)
	}
}

type SecretManagementKeyPair struct {
	// The path to the private RSA key file.
	PrivateKeyFile string `json:"privateKeyFile"`
	// Base64 encoded string of private key.
	PrivateKeyData string `json:"privateKeyData,omitempty"`
	// The path to the public RSA key file.
	PublicKeyFile string `json:"publicKeyFile"`
	// Base64 encoded string of public key.
	PublicKeyData string `json:"publicKeyData,omitempty"`
}

func (s *SecretManagementKeyPair) Validate() error {
	if s.PrivateKeyFile == "" && s.PrivateKeyData == "" {
		return errors.New("either privateKeyFile or privateKeyData must be set")
	}
	if s.PrivateKeyFile != "" && s.PrivateKeyData != "" {
		return errors.New("only privateKeyFile or privateKeyData can be set")
	}
	if s.PublicKeyFile == "" && s.PublicKeyData == "" {
		return errors.New("either publicKeyFile or publicKeyData must be set")
	}
	if s.PublicKeyFile != "" && s.PublicKeyData != "" {
		return errors.New("only publicKeyFile or publicKeyData can be set")
	}
	return nil
}

func (s *SecretManagementKeyPair) Mask() {
	if len(s.PrivateKeyFile) != 0 {
		s.PrivateKeyFile = maskString
	}
	if len(s.PrivateKeyData) != 0 {
		s.PrivateKeyData = maskString
	}
}

func (s *SecretManagementKeyPair) LoadPrivateKey() ([]byte, error) {
	if s.PrivateKeyData != "" {
		return base64.StdEncoding.DecodeString(s.PrivateKeyData)
	}
	if s.PrivateKeyFile != "" {
		return os.ReadFile(s.PrivateKeyFile)
	}
	return nil, errors.New("either privateKeyFile or privateKeyData must be set")
}

func (s *SecretManagementKeyPair) LoadPublicKey() ([]byte, error) {
	if s.PublicKeyData != "" {
		return base64.StdEncoding.DecodeString(s.PublicKeyData)
	}
	if s.PublicKeyFile != "" {
		return os.ReadFile(s.PublicKeyFile)
	}
	return nil, errors.New("either publicKeyFile or publicKeyData must be set")
}

type SecretManagementGCPKMS struct {
	// Configurable fields when using Google Cloud KMS.
	// The key name used for decrypting the sealed secret.
	KeyName string `json:"keyName"`
	// The path to the service account used to decrypt secret.
	DecryptServiceAccountFile string `json:"decryptServiceAccountFile"`
	// The path to the service account used to encrypt secret.
	EncryptServiceAccountFile string `json:"encryptServiceAccountFile"`
}

func (s *SecretManagementGCPKMS) Validate() error {
	if s.KeyName == "" {
		return fmt.Errorf("keyName must be set")
	}
	if s.DecryptServiceAccountFile == "" {
		return fmt.Errorf("decryptServiceAccountFile must be set")
	}
	if s.EncryptServiceAccountFile == "" {
		return fmt.Errorf("encryptServiceAccountFile must be set")
	}
	return nil
}

func (s *SecretManagementGCPKMS) Mask() {
	if len(s.DecryptServiceAccountFile) != 0 {
		s.DecryptServiceAccountFile = maskString
	}
	if len(s.EncryptServiceAccountFile) != 0 {
		s.EncryptServiceAccountFile = maskString
	}
}

type PipedEventWatcher struct {
	// Interval to fetch the latest event and compare it with one defined in EventWatcher config files
	CheckInterval Duration `json:"checkInterval,omitempty"`
	// The configuration list of git repositories to be observed.
	// Only the repositories in this list will be observed by Piped.
	GitRepos []PipedEventWatcherGitRepo `json:"gitRepos,omitempty"`
}

func (p *PipedEventWatcher) Validate() error {
	seen := make(map[string]struct{}, len(p.GitRepos))
	for i, repo := range p.GitRepos {
		// Validate the existence of repo ID.
		if repo.RepoID == "" {
			return fmt.Errorf("missing repoID at index %d", i)
		}
		// Validate if duplicated repository settings exist.
		if _, ok := seen[repo.RepoID]; ok {
			return fmt.Errorf("duplicated repo id (%s) found in the eventWatcher directive", repo.RepoID)
		}
		seen[repo.RepoID] = struct{}{}
	}
	return nil
}

type PipedEventWatcherGitRepo struct {
	// Id of the git repository. This must be unique within
	// the repos' elements.
	RepoID string `json:"repoId,omitempty"`
	// The commit message used to push after replacing values.
	// Default message is used if not given.
	CommitMessage string `json:"commitMessage,omitempty"`
	// The file path patterns to be included.
	// Patterns can be used like "foo/*.yaml".
	Includes []string `json:"includes,omitempty"`
	// The file path patterns to be excluded.
	// Patterns can be used like "foo/*.yaml".
	// This is prioritized if both includes and this one are given.
	Excludes []string `json:"excludes,omitempty"`
}

type PipedPlugin struct {
	Name          string              `json:"name"`
	Port          int                 `json:"port"`
	DeployTargets []PipedDeployTarget `json:"deployTargets,omitempty"`
}

type PipedDeployTarget struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
	Config json.RawMessage   `json:"config"`
}
