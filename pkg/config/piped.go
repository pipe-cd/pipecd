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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/pipe-cd/pipe/pkg/model"
)

var DefaultKubernetesCloudProvider = PipedCloudProvider{
	Name:             "kubernetes-default",
	Type:             model.CloudProviderKubernetes,
	KubernetesConfig: &CloudProviderKubernetesConfig{},
}

// PipedSpec contains configurable data used to while running Piped.
type PipedSpec struct {
	// The identifier of the PipeCD project where this piped belongs to.
	ProjectID string
	// The unique identifier generated for this piped.
	PipedID string
	// The path to the file containing the generated Key string for this piped.
	PipedKeyFile string
	// Base64 encoded string of Piped key.
	PipedKeyData string
	// The name of this piped.
	Name string
	// The address used to connect to the control-plane's API.
	APIAddress string `json:"apiAddress"`
	// The address to the control-plane's Web.
	WebAddress string `json:"webAddress"`
	// How often to check whether an application should be synced.
	// Default is 1m.
	SyncInterval Duration `json:"syncInterval" default:"1m"`
	// How often to check whether an application configuration file should be synced.
	// Default is 5m.
	AppConfigSyncInterval Duration `json:"appConfigSyncInterval" default:"5m"`
	// Git configuration needed for git commands.
	Git PipedGit `json:"git"`
	// List of git repositories this piped will handle.
	Repositories []PipedRepository `json:"repositories"`
	// List of helm chart repositories that should be added while starting up.
	ChartRepositories []HelmChartRepository `json:"chartRepositories"`
	// List of cloud providers can be used by this piped.
	CloudProviders []PipedCloudProvider `json:"cloudProviders"`
	// List of analysis providers can be used by this piped.
	AnalysisProviders []PipedAnalysisProvider `json:"analysisProviders"`
	// Sending notification to Slack, Webhook…
	Notifications Notifications `json:"notifications"`
	// What secret management method should be used.
	SecretManagement *SecretManagement `json:"secretManagement"`
	// Optional settings for event watcher.
	EventWatcher PipedEventWatcher `json:"eventWatcher"`
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
	for _, r := range s.ChartRepositories {
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
	for _, p := range s.AnalysisProviders {
		if err := p.Validate(); err != nil {
			return err
		}
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
	Username string `json:"username"`
	// The email that will be configured for `git` user.
	// Default is "pipecd.dev@gmail.com".
	Email string `json:"email"`
	// Where to write ssh config file.
	// Default is "$HOME/.ssh/config".
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
	// This will be used to clone the source code of the specified git repositories.
	SSHKeyFile string `json:"sshKeyFile"`
	// Base64 encoded string of ssh-key.
	SSHKeyData string `json:"sshKeyData"`
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

type HelmChartRepository struct {
	// The name of the Helm chart repository.
	Name string `json:"name"`
	// The address to the Helm chart repository.
	Address string `json:"address"`
	// Username used for the repository backed by HTTP basic authentication.
	Username string `json:"username"`
	// Password used for the repository backed by HTTP basic authentication.
	Password string `json:"password"`
	// Whether to skip TLS certificate checks for the repository or not.
	Insecure bool `json:"insecure"`

	// Remote address of the Git repository used to clone Helm charts.
	// e.g. git@github.com:org/repo.git
	GitRemote string `json:"gitRemote"`
	// The path to the private ssh key file used while cloning Helm charts from above Git repository.
	SSHKeyFile string `json:"sshKeyFile"`
}

func (r *HelmChartRepository) IsHTTPRepository() bool {
	return r.Name != "" && r.Address != ""
}

func (r *HelmChartRepository) IsGitRepository() bool {
	return r.GitRemote != ""
}

func (r *HelmChartRepository) Validate() error {
	if r.IsHTTPRepository() {
		if r.Name == "" {
			return errors.New("name must be set")
		}
		if r.Address == "" {
			return errors.New("address must be set")
		}
	}

	if r.IsGitRepository() {
		if r.GitRemote == "" {
			return errors.New("gitRemote must be set")
		}
	}
	return nil
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

type PipedCloudProvider struct {
	Name string
	Type model.CloudProviderType

	KubernetesConfig *CloudProviderKubernetesConfig
	TerraformConfig  *CloudProviderTerraformConfig
	CloudRunConfig   *CloudProviderCloudRunConfig
	LambdaConfig     *CloudProviderLambdaConfig
	ECSConfig        *CloudProviderECSConfig
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
	case model.CloudProviderECS:
		p.ECSConfig = &CloudProviderECSConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.ECSConfig)
		}
	default:
		err = fmt.Errorf("unsupported cloud provider type: %s", p.Name)
	}
	return err
}

type CloudProviderKubernetesConfig struct {
	// The master URL of the kubernetes cluster.
	// Empty means in-cluster.
	MasterURL string `json:"masterURL"`
	// The path to the kubeconfig file.
	// Empty means in-cluster.
	KubeConfigPath string `json:"kubeConfigPath"`
	// Configuration for application resource informer.
	AppStateInformer KubernetesAppStateInformer `json:"appStateInformer"`
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
	// The APIVersion of the kubernetes resource.
	APIVersion string `json:"apiVersion"`
	// The kind name of the kubernetes resource.
	// Empty means all kinds are matching.
	Kind string `json:"kind"`
}

type CloudProviderTerraformConfig struct {
	// List of variables that will be set directly on terraform commands with "-var" flag.
	// The variable must be formatted by "key=value" as below:
	// "image_id=ami-abc123"
	// 'image_id_list=["ami-abc123","ami-def456"]'
	// 'image_id_map={"us-east-1":"ami-abc123","us-east-2":"ami-def456"}'
	Vars []string `json:"vars"`
}

type CloudProviderCloudRunConfig struct {
	// The GCP project hosting the CloudRun service.
	Project string `json:"project"`
	// The region of running CloudRun service.
	Region string `json:"region"`
	// The path to the service account file for accessing CloudRun service.
	CredentialsFile string `json:"credentialsFile"`
}

type CloudProviderLambdaConfig struct {
	// The region to send requests to. This parameter is required.
	// e.g. "us-west-2"
	// A full list of regions is: https://docs.aws.amazon.com/general/latest/gr/rande.html
	Region string `json:"region"`
	// Path to the shared credentials file.
	CredentialsFile string `json:"credentialsFile"`
	// The IAM role arn to use when assuming an role.
	RoleARN string `json:"roleARN"`
	// Path to the WebIdentity token the SDK should use to assume a role with.
	TokenFile string `json:"tokenFile"`
	// AWS Profile to extract credentials from the shared credentials file.
	// If empty, the environment variable "AWS_PROFILE" is used.
	// "default" is populated if the environment variable is also not set.
	Profile string `json:"profile"`
}

type CloudProviderECSConfig struct {
	// The region to send requests to. This parameter is required.
	// e.g. "us-west-2"
	// A full list of regions is: https://docs.aws.amazon.com/general/latest/gr/rande.html
	Region string `json:"region"`
	// Path to the shared credentials file.
	CredentialsFile string `json:"credentialsFile"`
	// The IAM role arn to use when assuming an role.
	RoleARN string `json:"roleARN"`
	// Path to the WebIdentity token the SDK should use to assume a role with.
	TokenFile string `json:"tokenFile"`
	// AWS Profile to extract credentials from the shared credentials file.
	// If empty, the environment variable "AWS_PROFILE" is used.
	// "default" is populated if the environment variable is also not set.
	Profile string `json:"profile"`
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
	UsernameFile string `json:"usernameFile"`
	// The path to the password file.
	PasswordFile string `json:"passwordFile"`
}

func (a *AnalysisProviderPrometheusConfig) Validate() error {
	if a.Address == "" {
		return fmt.Errorf("prometheus analysis provider requires the address")
	}
	return nil
}

type AnalysisProviderDatadogConfig struct {
	// The address of Datadog API server.
	// Only "datadoghq.com", "us3.datadoghq.com", "datadoghq.eu", "ddog-gov.com" are available.
	// Defaults to "datadoghq.com"
	Address string `json:"address"`
	// Required: The path to the api key file.
	APIKeyFile string `json:"apiKeyFile"`
	// Required: The path to the application key file.
	ApplicationKeyFile string `json:"applicationKeyFile"`
}

func (a *AnalysisProviderDatadogConfig) Validate() error {
	if a.APIKeyFile == "" {
		return fmt.Errorf("datadog analysis provider requires the api key file")
	}
	if a.ApplicationKeyFile == "" {
		return fmt.Errorf("datadog analysis provider requires the application key file")
	}
	return nil
}

type AnalysisProviderStackdriverConfig struct {
	// The path to the service account file.
	ServiceAccountFile string `json:"serviceAccountFile"`
}

func (a *AnalysisProviderStackdriverConfig) Validate() error {
	return nil
}

type Notifications struct {
	// List of notification routes.
	Routes []NotificationRoute `json:"routes"`
	// List of notification receivers.
	Receivers []NotificationReceiver `json:"receivers"`
}

type NotificationRoute struct {
	Name         string   `json:"name"`
	Receiver     string   `json:"receiver"`
	Events       []string `json:"events"`
	IgnoreEvents []string `json:"ignoreEvents"`
	Groups       []string `json:"groups"`
	IgnoreGroups []string `json:"ignoreGroups"`
	Apps         []string `json:"apps"`
	IgnoreApps   []string `json:"ignoreApps"`
	Envs         []string `json:"envs"`
	IgnoreEnvs   []string `json:"ignoreEnvs"`
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
	URL            string `json:"url"`
	SignatureKey   string `json:"signatureKey" default:"PipeCD-Signature"`
	SignatureValue string `json:"signatureValue"`
}

type SecretManagement struct {
	// Which management service should be used.
	// Available values: KEY_PAIR, SEALING_KEY, GCP_KMS, AWS_KMS
	Type model.SecretManagementType `json:"type"`

	KeyPair *SecretManagementKeyPair
	GCPKMS  *SecretManagementGCPKMS
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
	// Configurable fields for SEALING_KEY.
	// The path to the private RSA key file.
	PrivateKeyFile string `json:"privateKeyFile"`
	// Base64 encoded string of private key.
	PrivateKeyData string `json:"privateKeyData"`
	// The path to the public RSA key file.
	PublicKeyFile string `json:"publicKeyFile"`
	// Base64 encoded string of public key.
	PublicKeyData string `json:"publicKeyData"`
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

type genericSecretManagement struct {
	Type   model.SecretManagementType `json:"type"`
	Config json.RawMessage            `json:"config"`
}

func (s *SecretManagement) UnmarshalJSON(data []byte) error {
	var err error
	g := genericSecretManagement{}
	if err = json.Unmarshal(data, &g); err != nil {
		return err
	}

	switch g.Type {
	case model.SecretManagementTypeSealingKey:
		fallthrough
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

type PipedEventWatcher struct {
	// Interval to fetch the latest event and compare it with one defined in EventWatcher config files
	CheckInterval Duration `json:"checkInterval"`
	// The configuration list of git repositories to be observed.
	// Only the repositories in this list will be observed by Piped.
	GitRepos []PipedEventWatcherGitRepo `json:"gitRepos"`
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
	RepoID string `json:"repoId"`
	// The commit message used to push after replacing values.
	// Default message is used if not given.
	CommitMessage string `json:"commitMessage"`
	// The file path patterns to be included.
	// Patterns can be used like "foo/*.yaml".
	Includes []string `json:"includes"`
	// The file path patterns to be excluded.
	// Patterns can be used like "foo/*.yaml".
	// This is prioritized if both includes and this one are given.
	Excludes []string `json:"excludes"`
}
