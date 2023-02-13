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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestPipedConfig(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/piped/piped-config.yaml",
			expectedKind:       KindPiped,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &PipedSpec{
				ProjectID:             "test-project",
				PipedID:               "test-piped",
				PipedKeyFile:          "etc/piped/key",
				APIAddress:            "your-pipecd.domain",
				WebAddress:            "https://your-pipecd.domain",
				SyncInterval:          Duration(time.Minute),
				AppConfigSyncInterval: Duration(time.Minute),
				Git: PipedGit{
					Username:   "username",
					Email:      "username@email.com",
					SSHKeyFile: "/etc/piped-secret/ssh-key",
				},
				Repositories: []PipedRepository{
					{
						RepoID: "repo1",
						Remote: "git@github.com:org/repo1.git",
						Branch: "master",
					},
					{
						RepoID: "repo2",
						Remote: "git@github.com:org/repo2.git",
						Branch: "master",
					},
				},
				ChartRepositories: []HelmChartRepository{
					{
						Type:    HTTPHelmChartRepository,
						Name:    "fantastic-charts",
						Address: "https://fantastic-charts.storage.googleapis.com",
					},
					{
						Type:     HTTPHelmChartRepository,
						Name:     "private-charts",
						Address:  "https://private-charts.com",
						Username: "basic-username",
						Password: "basic-password",
						Insecure: true,
					},
				},
				ChartRegistries: []HelmChartRegistry{
					{
						Type:     OCIHelmChartRegistry,
						Address:  "registry.example.com",
						Username: "sample-username",
						Password: "sample-password",
					},
				},
				PlatformProviders: []PipedPlatformProvider{
					{
						Name: "kubernetes-default",
						Type: model.PlatformProviderKubernetes,
						Labels: map[string]string{
							"group": "workload",
						},
						KubernetesConfig: &PlatformProviderKubernetesConfig{
							MasterURL:      "https://example.com",
							KubeConfigPath: "/etc/kube/config",
							AppStateInformer: KubernetesAppStateInformer{
								IncludeResources: []KubernetesResourceMatcher{
									{
										APIVersion: "pipecd.dev/v1beta1",
									},
									{
										APIVersion: "networking.gke.io/v1beta1",
										Kind:       "ManagedCertificate",
									},
								},
								ExcludeResources: []KubernetesResourceMatcher{
									{
										APIVersion: "v1",
										Kind:       "Endpoints",
									},
								},
							},
						},
					},
					{
						Name: "kubernetes-dev",
						Type: model.PlatformProviderKubernetes,
						Labels: map[string]string{
							"group": "config",
						},
						KubernetesConfig: &PlatformProviderKubernetesConfig{},
					},
					{
						Name: "terraform",
						Type: model.PlatformProviderTerraform,
						TerraformConfig: &PlatformProviderTerraformConfig{
							Vars: []string{
								"project=gcp-project",
								"region=us-centra1",
							},
						},
					},
					{
						Name: "cloudrun",
						Type: model.PlatformProviderCloudRun,
						CloudRunConfig: &PlatformProviderCloudRunConfig{
							Project:         "gcp-project-id",
							Region:          "cloud-run-region",
							CredentialsFile: "/etc/piped-secret/gcp-service-account.json",
						},
					},
					{
						Name: "lambda",
						Type: model.PlatformProviderLambda,
						LambdaConfig: &PlatformProviderLambdaConfig{
							Region: "us-east-1",
						},
					},
				},
				AnalysisProviders: []PipedAnalysisProvider{
					{
						Name: "prometheus-dev",
						Type: model.AnalysisProviderPrometheus,
						PrometheusConfig: &AnalysisProviderPrometheusConfig{
							Address: "https://your-prometheus.dev",
						},
					},
					{
						Name: "datadog-dev",
						Type: model.AnalysisProviderDatadog,
						DatadogConfig: &AnalysisProviderDatadogConfig{
							Address:            "https://your-datadog.dev",
							APIKeyFile:         "/etc/piped-secret/datadog-api-key",
							ApplicationKeyFile: "/etc/piped-secret/datadog-application-key",
						},
					},
					{
						Name: "stackdriver-dev",
						Type: model.AnalysisProviderStackdriver,
						StackdriverConfig: &AnalysisProviderStackdriverConfig{
							ServiceAccountFile: "/etc/piped-secret/gcp-service-account.json",
						},
					},
				},
				Notifications: Notifications{
					Routes: []NotificationRoute{
						{
							Name: "dev-slack",
							Labels: map[string]string{
								"env":  "dev",
								"team": "pipecd",
							},
							Receiver: "dev-slack-channel",
						},
						{
							Name: "prod-slack",
							Labels: map[string]string{
								"env": "dev",
							},
							Events:   []string{"DEPLOYMENT_TRIGGERED", "DEPLOYMENT_COMPLETED"},
							Receiver: "prod-slack-channel",
						},
						{
							Name:     "all-events-to-ci",
							Receiver: "ci-webhook",
						},
					},
					Receivers: []NotificationReceiver{
						{
							Name: "dev-slack-channel",
							Slack: &NotificationReceiverSlack{
								HookURL: "https://slack.com/dev",
							},
						},
						{
							Name: "prod-slack-channel",
							Slack: &NotificationReceiverSlack{
								HookURL: "https://slack.com/prod",
							},
						},
						{
							Name: "ci-webhook",
							Webhook: &NotificationReceiverWebhook{
								URL:            "https://pipecd.dev/dev-hook",
								SignatureKey:   "PipeCD-Signature",
								SignatureValue: "random-signature-string",
							},
						},
					},
				},
				SecretManagement: &SecretManagement{
					Type: model.SecretManagementTypeKeyPair,
					KeyPair: &SecretManagementKeyPair{
						PrivateKeyFile: "/etc/piped-secret/pair-private-key",
						PublicKeyFile:  "/etc/piped-secret/pair-public-key",
					},
				},
				EventWatcher: PipedEventWatcher{
					CheckInterval: Duration(10 * time.Minute),
					GitRepos: []PipedEventWatcherGitRepo{
						{
							RepoID:        "repo-1",
							CommitMessage: "Update values by Event watcher",
							Includes:      []string{"event-watcher-dev.yaml", "event-watcher-stg.yaml"},
						},
					},
				},
			},
			expectedError: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.fileName, func(t *testing.T) {
			cfg, err := LoadFromYAML(tc.fileName)
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedKind, cfg.Kind)
				assert.Equal(t, tc.expectedAPIVersion, cfg.APIVersion)
				assert.Equal(t, tc.expectedSpec, cfg.spec)
			}
		})
	}
}

func TestPipedEventWatcherValidate(t *testing.T) {
	testcases := []struct {
		name                  string
		eventWatcher          PipedEventWatcher
		wantErr               bool
		wantPipedEventWatcher PipedEventWatcher
	}{
		{
			name:    "missing repo id",
			wantErr: true,
			eventWatcher: PipedEventWatcher{
				GitRepos: []PipedEventWatcherGitRepo{
					{
						RepoID: "",
					},
				},
			},
			wantPipedEventWatcher: PipedEventWatcher{
				GitRepos: []PipedEventWatcherGitRepo{
					{
						RepoID: "",
					},
				},
			},
		},
		{
			name:    "duplicated repo exists",
			wantErr: true,
			eventWatcher: PipedEventWatcher{
				GitRepos: []PipedEventWatcherGitRepo{
					{
						RepoID: "foo",
					},
					{
						RepoID: "foo",
					},
				},
			},
			wantPipedEventWatcher: PipedEventWatcher{
				GitRepos: []PipedEventWatcherGitRepo{
					{
						RepoID: "foo",
					},
					{
						RepoID: "foo",
					},
				},
			},
		},
		{
			name:    "repos are unique",
			wantErr: false,
			eventWatcher: PipedEventWatcher{
				GitRepos: []PipedEventWatcherGitRepo{
					{
						RepoID: "foo",
					},
					{
						RepoID: "bar",
					},
				},
			},
			wantPipedEventWatcher: PipedEventWatcher{
				GitRepos: []PipedEventWatcherGitRepo{
					{
						RepoID: "foo",
					},
					{
						RepoID: "bar",
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.eventWatcher.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.wantPipedEventWatcher, tc.eventWatcher)
		})
	}
}

func TestNotificationReceiverWebhook_LoadSignatureValue(t *testing.T) {
	testcase := []struct {
		name    string
		webhook *NotificationReceiverWebhook
		want    string
		wantErr bool
	}{
		{
			name: "set signatureValue",
			webhook: &NotificationReceiverWebhook{
				URL:            "https://example.com",
				SignatureValue: "foo",
			},
			want:    "foo",
			wantErr: false,
		},
		{
			name: "set signatureValueFile",
			webhook: &NotificationReceiverWebhook{
				URL:                "https://example.com",
				SignatureValueFile: "testdata/piped/notification-receiver-webhook",
			},
			want:    "foo",
			wantErr: false,
		},
		{
			name: "set both of them",
			webhook: &NotificationReceiverWebhook{
				URL:                "https://example.com",
				SignatureValue:     "foo",
				SignatureValueFile: "testdata/piped/notification-receiver-webhook",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.webhook.LoadSignatureValue()
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPipedConfigMask(t *testing.T) {
	testcase := []struct {
		name    string
		spec    *PipedSpec
		want    *PipedSpec
		wantErr bool
	}{
		{
			name: "mask",
			spec: &PipedSpec{
				ProjectID:             "foo",
				PipedID:               "foo",
				PipedKeyFile:          "foo",
				PipedKeyData:          "foo",
				Name:                  "foo",
				APIAddress:            "foo",
				WebAddress:            "foo",
				SyncInterval:          Duration(time.Minute),
				AppConfigSyncInterval: Duration(time.Minute),
				Git: PipedGit{
					Username:          "foo",
					Email:             "foo",
					SSHConfigFilePath: "foo",
					Host:              "foo",
					HostName:          "foo",
					SSHKeyFile:        "foo",
					SSHKeyData:        "foo",
				},
				Repositories: []PipedRepository{
					{
						RepoID: "foo",
						Remote: "foo",
						Branch: "foo",
					},
				},
				ChartRepositories: []HelmChartRepository{
					{
						Type:       "foo",
						Name:       "foo",
						Address:    "foo",
						Username:   "foo",
						Password:   "foo",
						Insecure:   true,
						GitRemote:  "foo",
						SSHKeyFile: "foo",
					},
				},
				ChartRegistries: []HelmChartRegistry{
					{
						Type:     "foo",
						Address:  "foo",
						Username: "foo",
						Password: "foo",
					},
				},
				PlatformProviders: []PipedPlatformProvider{
					{
						Name: "foo",
						Type: model.PlatformProviderKubernetes,
						KubernetesConfig: &PlatformProviderKubernetesConfig{
							MasterURL:      "foo",
							KubeConfigPath: "foo",
							AppStateInformer: KubernetesAppStateInformer{
								Namespace: "",
								IncludeResources: []KubernetesResourceMatcher{
									{
										APIVersion: "foo",
										Kind:       "foo",
									},
								},
								ExcludeResources: []KubernetesResourceMatcher{
									{
										APIVersion: "foo",
										Kind:       "foo",
									},
								},
							},
						},
					},
					{
						Name: "bar",
						Type: model.PlatformProviderCloudRun,
						CloudRunConfig: &PlatformProviderCloudRunConfig{
							Project:         "bar",
							Region:          "bar",
							CredentialsFile: "/etc/cloudrun/credentials",
						},
					},
				},
				AnalysisProviders: []PipedAnalysisProvider{
					{
						Name: "foo",
						Type: "foo",
						PrometheusConfig: &AnalysisProviderPrometheusConfig{
							Address:      "foo",
							UsernameFile: "foo",
							PasswordFile: "foo",
						},
						DatadogConfig: &AnalysisProviderDatadogConfig{
							Address:            "foo",
							APIKeyFile:         "foo",
							ApplicationKeyFile: "foo",
						},
						StackdriverConfig: &AnalysisProviderStackdriverConfig{
							ServiceAccountFile: "foo",
						},
					},
				},
				Notifications: Notifications{
					Routes: []NotificationRoute{
						{
							Name:         "foo",
							Receiver:     "foo",
							Events:       []string{"foo"},
							IgnoreEvents: []string{"foo"},
							Groups:       []string{"foo"},
							IgnoreGroups: []string{"foo"},
							Apps:         []string{"foo"},
							IgnoreApps:   []string{"foo"},
							Labels:       map[string]string{"foo": "foo"},
							IgnoreLabels: map[string]string{"foo": "foo"},
						},
					},
					Receivers: []NotificationReceiver{
						{
							Name: "foo",
							Slack: &NotificationReceiverSlack{
								HookURL: "foo",
							},
							Webhook: &NotificationReceiverWebhook{
								URL:                "foo",
								SignatureKey:       "foo",
								SignatureValue:     "foo",
								SignatureValueFile: "foo",
							},
						},
					},
				},
				SecretManagement: &SecretManagement{
					Type: "foo",
					KeyPair: &SecretManagementKeyPair{
						PrivateKeyFile: "foo",
						PrivateKeyData: "foo",
						PublicKeyFile:  "foo",
						PublicKeyData:  "foo",
					},
					GCPKMS: &SecretManagementGCPKMS{
						KeyName:                   "foo",
						DecryptServiceAccountFile: "foo",
						EncryptServiceAccountFile: "foo",
					},
				},
				EventWatcher: PipedEventWatcher{
					CheckInterval: Duration(time.Minute),
					GitRepos: []PipedEventWatcherGitRepo{
						{
							RepoID:        "foo",
							CommitMessage: "foo",
							Includes:      []string{"foo"},
							Excludes:      []string{"foo"},
						},
					},
				},
				AppSelector: map[string]string{
					"foo": "foo",
				},
			},
			want: &PipedSpec{
				ProjectID:             "foo",
				PipedID:               "foo",
				PipedKeyFile:          maskString,
				PipedKeyData:          maskString,
				Name:                  "foo",
				APIAddress:            "foo",
				WebAddress:            "foo",
				SyncInterval:          Duration(time.Minute),
				AppConfigSyncInterval: Duration(time.Minute),
				Git: PipedGit{
					Username:          "foo",
					Email:             "foo",
					SSHConfigFilePath: maskString,
					Host:              "foo",
					HostName:          "foo",
					SSHKeyFile:        maskString,
					SSHKeyData:        maskString,
				},
				Repositories: []PipedRepository{
					{
						RepoID: "foo",
						Remote: "foo",
						Branch: "foo",
					},
				},
				ChartRepositories: []HelmChartRepository{
					{
						Type:       "foo",
						Name:       "foo",
						Address:    "foo",
						Username:   "foo",
						Password:   maskString,
						Insecure:   true,
						GitRemote:  "foo",
						SSHKeyFile: maskString,
					},
				},
				ChartRegistries: []HelmChartRegistry{
					{
						Type:     "foo",
						Address:  "foo",
						Username: "foo",
						Password: maskString,
					},
				},
				PlatformProviders: []PipedPlatformProvider{
					{
						Name: "foo",
						Type: model.PlatformProviderKubernetes,
						KubernetesConfig: &PlatformProviderKubernetesConfig{
							MasterURL:      "foo",
							KubeConfigPath: "foo",
							AppStateInformer: KubernetesAppStateInformer{
								Namespace: "",
								IncludeResources: []KubernetesResourceMatcher{
									{
										APIVersion: "foo",
										Kind:       "foo",
									},
								},
								ExcludeResources: []KubernetesResourceMatcher{
									{
										APIVersion: "foo",
										Kind:       "foo",
									},
								},
							},
						},
					},
					{
						Name: "bar",
						Type: model.PlatformProviderCloudRun,
						CloudRunConfig: &PlatformProviderCloudRunConfig{
							Project:         "bar",
							Region:          "bar",
							CredentialsFile: "******",
						},
					},
				},
				AnalysisProviders: []PipedAnalysisProvider{
					{
						Name: "foo",
						Type: "foo",
						PrometheusConfig: &AnalysisProviderPrometheusConfig{
							Address:      "foo",
							UsernameFile: "foo",
							PasswordFile: maskString,
						},
						DatadogConfig: &AnalysisProviderDatadogConfig{
							Address:            "foo",
							APIKeyFile:         maskString,
							ApplicationKeyFile: maskString,
						},
						StackdriverConfig: &AnalysisProviderStackdriverConfig{
							ServiceAccountFile: maskString,
						},
					},
				},
				Notifications: Notifications{
					Routes: []NotificationRoute{
						{
							Name:         "foo",
							Receiver:     "foo",
							Events:       []string{"foo"},
							IgnoreEvents: []string{"foo"},
							Groups:       []string{"foo"},
							IgnoreGroups: []string{"foo"},
							Apps:         []string{"foo"},
							IgnoreApps:   []string{"foo"},
							Labels:       map[string]string{"foo": "foo"},
							IgnoreLabels: map[string]string{"foo": "foo"},
						},
					},
					Receivers: []NotificationReceiver{
						{
							Name: "foo",
							Slack: &NotificationReceiverSlack{
								HookURL: maskString,
							},
							Webhook: &NotificationReceiverWebhook{
								URL:                maskString,
								SignatureKey:       maskString,
								SignatureValue:     maskString,
								SignatureValueFile: maskString,
							},
						},
					},
				},
				SecretManagement: &SecretManagement{
					Type: "foo",
					KeyPair: &SecretManagementKeyPair{
						PrivateKeyFile: maskString,
						PrivateKeyData: maskString,
						PublicKeyFile:  "foo",
						PublicKeyData:  "foo",
					},
					GCPKMS: &SecretManagementGCPKMS{
						KeyName:                   "foo",
						DecryptServiceAccountFile: maskString,
						EncryptServiceAccountFile: maskString,
					},
				},
				EventWatcher: PipedEventWatcher{
					CheckInterval: Duration(time.Minute),
					GitRepos: []PipedEventWatcherGitRepo{
						{
							RepoID:        "foo",
							CommitMessage: "foo",
							Includes:      []string{"foo"},
							Excludes:      []string{"foo"},
						},
					},
				},
				AppSelector: map[string]string{
					"foo": "foo",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			tc.spec.Mask()
			assert.Equal(t, tc.want, tc.spec)
		})
	}
}

func TestPipedSpecClone(t *testing.T) {
	testcases := []struct {
		name          string
		originalSpec  *PipedSpec
		expectedSpec  *PipedSpec
		expectedError error
	}{
		{
			name: "clone success",
			originalSpec: &PipedSpec{
				ProjectID:             "test-project",
				PipedID:               "test-piped",
				PipedKeyFile:          "etc/piped/key",
				APIAddress:            "your-pipecd.domain",
				WebAddress:            "https://your-pipecd.domain",
				SyncInterval:          Duration(time.Minute),
				AppConfigSyncInterval: Duration(time.Minute),
				Git: PipedGit{
					Username:   "username",
					Email:      "username@email.com",
					SSHKeyFile: "/etc/piped-secret/ssh-key",
				},
				Repositories: []PipedRepository{
					{
						RepoID: "repo1",
						Remote: "git@github.com:org/repo1.git",
						Branch: "master",
					},
					{
						RepoID: "repo2",
						Remote: "git@github.com:org/repo2.git",
						Branch: "master",
					},
				},
				ChartRepositories: []HelmChartRepository{
					{
						Type:    HTTPHelmChartRepository,
						Name:    "fantastic-charts",
						Address: "https://fantastic-charts.storage.googleapis.com",
					},
					{
						Type:     HTTPHelmChartRepository,
						Name:     "private-charts",
						Address:  "https://private-charts.com",
						Username: "basic-username",
						Password: "basic-password",
						Insecure: true,
					},
				},
				ChartRegistries: []HelmChartRegistry{
					{
						Type:     OCIHelmChartRegistry,
						Address:  "registry.example.com",
						Username: "sample-username",
						Password: "sample-password",
					},
				},
				PlatformProviders: []PipedPlatformProvider{
					{
						Name: "kubernetes-default",
						Type: model.PlatformProviderKubernetes,
						KubernetesConfig: &PlatformProviderKubernetesConfig{
							MasterURL:      "https://example.com",
							KubeConfigPath: "/etc/kube/config",
							AppStateInformer: KubernetesAppStateInformer{
								IncludeResources: []KubernetesResourceMatcher{
									{
										APIVersion: "pipecd.dev/v1beta1",
									},
									{
										APIVersion: "networking.gke.io/v1beta1",
										Kind:       "ManagedCertificate",
									},
								},
								ExcludeResources: []KubernetesResourceMatcher{
									{
										APIVersion: "v1",
										Kind:       "Endpoints",
									},
								},
							},
						},
					},
					{
						Name:             "kubernetes-dev",
						Type:             model.PlatformProviderKubernetes,
						KubernetesConfig: &PlatformProviderKubernetesConfig{},
					},
					{
						Name: "terraform",
						Type: model.PlatformProviderTerraform,
						TerraformConfig: &PlatformProviderTerraformConfig{
							Vars: []string{
								"project=gcp-project",
								"region=us-centra1",
							},
						},
					},
					{
						Name: "cloudrun",
						Type: model.PlatformProviderCloudRun,
						CloudRunConfig: &PlatformProviderCloudRunConfig{
							Project:         "gcp-project-id",
							Region:          "cloud-run-region",
							CredentialsFile: "/etc/piped-secret/gcp-service-account.json",
						},
					},
					{
						Name: "lambda",
						Type: model.PlatformProviderLambda,
						LambdaConfig: &PlatformProviderLambdaConfig{
							Region: "us-east-1",
						},
					},
				},
				AnalysisProviders: []PipedAnalysisProvider{
					{
						Name: "prometheus-dev",
						Type: model.AnalysisProviderPrometheus,
						PrometheusConfig: &AnalysisProviderPrometheusConfig{
							Address: "https://your-prometheus.dev",
						},
					},
					{
						Name: "datadog-dev",
						Type: model.AnalysisProviderDatadog,
						DatadogConfig: &AnalysisProviderDatadogConfig{
							Address:            "https://your-datadog.dev",
							APIKeyFile:         "/etc/piped-secret/datadog-api-key",
							ApplicationKeyFile: "/etc/piped-secret/datadog-application-key",
						},
					},
					{
						Name: "stackdriver-dev",
						Type: model.AnalysisProviderStackdriver,
						StackdriverConfig: &AnalysisProviderStackdriverConfig{
							ServiceAccountFile: "/etc/piped-secret/gcp-service-account.json",
						},
					},
				},
				Notifications: Notifications{
					Routes: []NotificationRoute{
						{
							Name: "dev-slack",
							Labels: map[string]string{
								"env":  "dev",
								"team": "pipecd",
							},
							Receiver: "dev-slack-channel",
						},
						{
							Name: "prod-slack",
							Labels: map[string]string{
								"env": "dev",
							},
							Events:   []string{"DEPLOYMENT_TRIGGERED", "DEPLOYMENT_COMPLETED"},
							Receiver: "prod-slack-channel",
						},
						{
							Name:     "all-events-to-ci",
							Receiver: "ci-webhook",
						},
					},
					Receivers: []NotificationReceiver{
						{
							Name: "dev-slack-channel",
							Slack: &NotificationReceiverSlack{
								HookURL: "https://slack.com/dev",
							},
						},
						{
							Name: "prod-slack-channel",
							Slack: &NotificationReceiverSlack{
								HookURL: "https://slack.com/prod",
							},
						},
						{
							Name: "ci-webhook",
							Webhook: &NotificationReceiverWebhook{
								URL:            "https://pipecd.dev/dev-hook",
								SignatureKey:   "PipeCD-Signature",
								SignatureValue: "random-signature-string",
							},
						},
					},
				},
				SecretManagement: &SecretManagement{
					Type: model.SecretManagementTypeKeyPair,
					KeyPair: &SecretManagementKeyPair{
						PrivateKeyFile: "/etc/piped-secret/pair-private-key",
						PublicKeyFile:  "/etc/piped-secret/pair-public-key",
					},
				},
				EventWatcher: PipedEventWatcher{
					CheckInterval: Duration(10 * time.Minute),
					GitRepos: []PipedEventWatcherGitRepo{
						{
							RepoID:        "repo-1",
							CommitMessage: "Update values by Event watcher",
							Includes:      []string{"event-watcher-dev.yaml", "event-watcher-stg.yaml"},
						},
					},
				},
			},
			expectedSpec: &PipedSpec{
				ProjectID:             "test-project",
				PipedID:               "test-piped",
				PipedKeyFile:          "etc/piped/key",
				APIAddress:            "your-pipecd.domain",
				WebAddress:            "https://your-pipecd.domain",
				SyncInterval:          Duration(time.Minute),
				AppConfigSyncInterval: Duration(time.Minute),
				Git: PipedGit{
					Username:   "username",
					Email:      "username@email.com",
					SSHKeyFile: "/etc/piped-secret/ssh-key",
				},
				Repositories: []PipedRepository{
					{
						RepoID: "repo1",
						Remote: "git@github.com:org/repo1.git",
						Branch: "master",
					},
					{
						RepoID: "repo2",
						Remote: "git@github.com:org/repo2.git",
						Branch: "master",
					},
				},
				ChartRepositories: []HelmChartRepository{
					{
						Type:    HTTPHelmChartRepository,
						Name:    "fantastic-charts",
						Address: "https://fantastic-charts.storage.googleapis.com",
					},
					{
						Type:     HTTPHelmChartRepository,
						Name:     "private-charts",
						Address:  "https://private-charts.com",
						Username: "basic-username",
						Password: "basic-password",
						Insecure: true,
					},
				},
				ChartRegistries: []HelmChartRegistry{
					{
						Type:     OCIHelmChartRegistry,
						Address:  "registry.example.com",
						Username: "sample-username",
						Password: "sample-password",
					},
				},
				PlatformProviders: []PipedPlatformProvider{
					{
						Name: "kubernetes-default",
						Type: model.PlatformProviderKubernetes,
						KubernetesConfig: &PlatformProviderKubernetesConfig{
							MasterURL:      "https://example.com",
							KubeConfigPath: "/etc/kube/config",
							AppStateInformer: KubernetesAppStateInformer{
								IncludeResources: []KubernetesResourceMatcher{
									{
										APIVersion: "pipecd.dev/v1beta1",
									},
									{
										APIVersion: "networking.gke.io/v1beta1",
										Kind:       "ManagedCertificate",
									},
								},
								ExcludeResources: []KubernetesResourceMatcher{
									{
										APIVersion: "v1",
										Kind:       "Endpoints",
									},
								},
							},
						},
					},
					{
						Name:             "kubernetes-dev",
						Type:             model.PlatformProviderKubernetes,
						KubernetesConfig: &PlatformProviderKubernetesConfig{},
					},
					{
						Name: "terraform",
						Type: model.PlatformProviderTerraform,
						TerraformConfig: &PlatformProviderTerraformConfig{
							Vars: []string{
								"project=gcp-project",
								"region=us-centra1",
							},
						},
					},
					{
						Name: "cloudrun",
						Type: model.PlatformProviderCloudRun,
						CloudRunConfig: &PlatformProviderCloudRunConfig{
							Project:         "gcp-project-id",
							Region:          "cloud-run-region",
							CredentialsFile: "/etc/piped-secret/gcp-service-account.json",
						},
					},
					{
						Name: "lambda",
						Type: model.PlatformProviderLambda,
						LambdaConfig: &PlatformProviderLambdaConfig{
							Region: "us-east-1",
						},
					},
				},
				AnalysisProviders: []PipedAnalysisProvider{
					{
						Name: "prometheus-dev",
						Type: model.AnalysisProviderPrometheus,
						PrometheusConfig: &AnalysisProviderPrometheusConfig{
							Address: "https://your-prometheus.dev",
						},
					},
					{
						Name: "datadog-dev",
						Type: model.AnalysisProviderDatadog,
						DatadogConfig: &AnalysisProviderDatadogConfig{
							Address:            "https://your-datadog.dev",
							APIKeyFile:         "/etc/piped-secret/datadog-api-key",
							ApplicationKeyFile: "/etc/piped-secret/datadog-application-key",
						},
					},
					{
						Name: "stackdriver-dev",
						Type: model.AnalysisProviderStackdriver,
						StackdriverConfig: &AnalysisProviderStackdriverConfig{
							ServiceAccountFile: "/etc/piped-secret/gcp-service-account.json",
						},
					},
				},
				Notifications: Notifications{
					Routes: []NotificationRoute{
						{
							Name: "dev-slack",
							Labels: map[string]string{
								"env":  "dev",
								"team": "pipecd",
							},
							Receiver: "dev-slack-channel",
						},
						{
							Name: "prod-slack",
							Labels: map[string]string{
								"env": "dev",
							},
							Events:   []string{"DEPLOYMENT_TRIGGERED", "DEPLOYMENT_COMPLETED"},
							Receiver: "prod-slack-channel",
						},
						{
							Name:     "all-events-to-ci",
							Receiver: "ci-webhook",
						},
					},
					Receivers: []NotificationReceiver{
						{
							Name: "dev-slack-channel",
							Slack: &NotificationReceiverSlack{
								HookURL: "https://slack.com/dev",
							},
						},
						{
							Name: "prod-slack-channel",
							Slack: &NotificationReceiverSlack{
								HookURL: "https://slack.com/prod",
							},
						},
						{
							Name: "ci-webhook",
							Webhook: &NotificationReceiverWebhook{
								URL:            "https://pipecd.dev/dev-hook",
								SignatureKey:   "PipeCD-Signature",
								SignatureValue: "random-signature-string",
							},
						},
					},
				},
				SecretManagement: &SecretManagement{
					Type: model.SecretManagementTypeKeyPair,
					KeyPair: &SecretManagementKeyPair{
						PrivateKeyFile: "/etc/piped-secret/pair-private-key",
						PublicKeyFile:  "/etc/piped-secret/pair-public-key",
					},
				},
				EventWatcher: PipedEventWatcher{
					CheckInterval: Duration(10 * time.Minute),
					GitRepos: []PipedEventWatcherGitRepo{
						{
							RepoID:        "repo-1",
							CommitMessage: "Update values by Event watcher",
							Includes:      []string{"event-watcher-dev.yaml", "event-watcher-stg.yaml"},
						},
					},
				},
			},
			expectedError: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cloned, err := tc.originalSpec.Clone()
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedSpec, cloned)
			}
		})
	}
}

func TestFindPlatformProvidersByLabel(t *testing.T) {
	pipedSpec := &PipedSpec{
		PlatformProviders: []PipedPlatformProvider{
			{
				Name: "provider-1",
				Type: model.PlatformProviderKubernetes,
				Labels: map[string]string{
					"group": "group-1",
					"foo":   "foo-1",
				},
			},
			{
				Name: "provider-2",
				Type: model.PlatformProviderKubernetes,
				Labels: map[string]string{
					"group": "group-2",
					"foo":   "foo-2",
				},
			},
			{
				Name: "provider-3",
				Type: model.PlatformProviderCloudRun,
				Labels: map[string]string{
					"group": "group-1",
					"foo":   "foo-3",
				},
			},
			{
				Name: "provider-4",
				Type: model.PlatformProviderKubernetes,
				Labels: map[string]string{
					"group": "group-2",
					"foo":   "foo-4",
				},
			},
		},
	}

	testcases := []struct {
		name   string
		labels map[string]string
		want   []PipedPlatformProvider
	}{
		{
			name: "empty due to missing label",
			labels: map[string]string{
				"group": "group-4",
			},
			want: []PipedPlatformProvider{},
		},
		{
			name: "found exactly one provider",
			labels: map[string]string{
				"group": "group-1",
			},
			want: []PipedPlatformProvider{
				{
					Name: "provider-1",
					Type: model.PlatformProviderKubernetes,
					Labels: map[string]string{
						"group": "group-1",
						"foo":   "foo-1",
					},
				},
			},
		},
		{
			name: "found multiple providers",
			labels: map[string]string{
				"group": "group-1",
			},
			want: []PipedPlatformProvider{
				{
					Name: "provider-1",
					Type: model.PlatformProviderKubernetes,
					Labels: map[string]string{
						"group": "group-1",
						"foo":   "foo-1",
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := pipedSpec.FindPlatformProvidersByLabels(tc.labels, model.ApplicationKind_KUBERNETES)
			assert.Equal(t, tc.want, got)
		})
	}
}
