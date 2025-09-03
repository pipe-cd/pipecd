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
	"net/url"
	"os"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	maskString = "******"
)

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
	// List of plugin configs
	Plugins []PipedPlugin `json:"plugins,omitempty"`
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
	s.Notifications.Mask()
	if s.SecretManagement != nil {
		s.SecretManagement.Mask()
	}

	// TODO: Mask plugin configs
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

func (s *PipedSpec) LoadPipedKey() ([]byte, error) {
	if s.PipedKeyData != "" {
		return base64.StdEncoding.DecodeString(s.PipedKeyData)
	}
	if s.PipedKeyFile != "" {
		return os.ReadFile(s.PipedKeyFile)
	}
	return nil, errors.New("either pipedKeyFile or pipedKeyData must be set")
}

type PipedGit = config.PipedGit

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

// PipedPlugin defines the plugin configuration for the piped.
type PipedPlugin struct {
	// The name of the plugin.
	Name string `json:"name"`
	// Source to download the plugin binary.
	URL string `json:"url"`
	// The port which the plugin listens to.
	Port int `json:"port"`
	// Configuration for the plugin.
	Config json.RawMessage `json:"config,omitempty"`
	// The deploy targets.
	DeployTargets []PipedDeployTarget `json:"deployTargets,omitempty"`
}

// PipedDeployTarget defines the deploy target configuration for the piped.
type PipedDeployTarget struct {
	// The name of the deploy target.
	Name string `json:"name"`
	// The labes of the deploy target.
	Labels map[string]string `json:"labels,omitempty"`
	// The configuration of the deploy target.
	Config json.RawMessage `json:"config"`
}

// ParsePluginConfig parses the given JSON string and returns the PipedPlugin.
func ParsePluginConfig(s string) (*PipedPlugin, error) {
	p := &PipedPlugin{}
	if err := json.Unmarshal([]byte(s), p); err != nil {
		return nil, err
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *PipedPlugin) Validate() error {
	if p.Name == "" {
		return errors.New("name must be set")
	}
	if p.URL == "" {
		return errors.New("url must be set")
	}
	u, err := url.Parse(p.URL)
	if err != nil {
		return fmt.Errorf("invalid plugin url: %w", err)
	}
	if u.Scheme != "file" && u.Scheme != "https" && u.Scheme != "oci" {
		return errors.New("only file, https and oci schemes are supported")
	}
	return nil
}

// FindDeployTarget finds the deploy target by the given name.
func (p *PipedPlugin) FindDeployTarget(name string) *PipedDeployTarget {
	for _, dt := range p.DeployTargets {
		if dt.Name == name {
			return &dt
		}
	}
	return nil
}
