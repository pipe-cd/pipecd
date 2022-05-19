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
					{
						Type:     OCIHelmChartRepository,
						Address:  "oci://private-charts.com",
						Username: "sample-username",
						Password: "sample-password",
					},
				},
				CloudProviders: []PipedCloudProvider{
					{
						Name: "kubernetes-default",
						Type: model.ApplicationKind_KUBERNETES,
						KubernetesConfig: &CloudProviderKubernetesConfig{
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
						Type:             model.ApplicationKind_KUBERNETES,
						KubernetesConfig: &CloudProviderKubernetesConfig{},
					},
					{
						Name: "terraform",
						Type: model.ApplicationKind_TERRAFORM,
						TerraformConfig: &CloudProviderTerraformConfig{
							Vars: []string{
								"project=gcp-project",
								"region=us-centra1",
							},
						},
					},
					{
						Name: "cloudrun",
						Type: model.ApplicationKind_CLOUDRUN,
						CloudRunConfig: &CloudProviderCloudRunConfig{
							Project:         "gcp-project-id",
							Region:          "cloud-run-region",
							CredentialsFile: "/etc/piped-secret/gcp-service-account.json",
						},
					},
					{
						Name: "lambda",
						Type: model.ApplicationKind_LAMBDA,
						LambdaConfig: &CloudProviderLambdaConfig{
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
							Events:   []string{"DEPLOYMENT_STARTED", "DEPLOYMENT_COMPLETED"},
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
