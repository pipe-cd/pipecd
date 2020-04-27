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
)

func TestRunnerConfig(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/runner/runner-config.yaml",
			expectedKind:       KindRunner,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &RunnerSpec{
				Git: RunnerGit{
					SSHKeyFile:      "/etc/pipecd-runner/ssh.key",
					AccessTokenFile: "/etc/pipecd-runner/github.token",
				},
				Repositories: []RunnerRepository{
					RunnerRepository{
						Repo:         "git@github.com:org/repo1",
						Branch:       "master",
						PollInterval: Duration(time.Minute),
					},
					RunnerRepository{
						Repo:         "git@github.com:org/repo2",
						Branch:       "master",
						PollInterval: Duration(2 * time.Minute),
					},
				},
				Destinations: []RunnerDestination{
					RunnerDestination{
						Name: "default",
						Kubernetes: &RunnerDestinationKubernetes{
							AllowNamespaces: []string{"dev"},
						},
					},
					RunnerDestination{
						Name: "gcp-terraform",
						Terraform: &RunnerDestinationTerraform{
							GCP: &RunnerTerraformGCP{
								Project:         "gcp-project",
								Region:          "us-central1",
								CredentialsFile: "/etc/pipecd-runner/terraform.serviceaccount",
							},
						},
					},
					RunnerDestination{
						Name: "aws-terraform",
						Terraform: &RunnerDestinationTerraform{
							AWS: &RunnerTerraformAWS{
								Region: "us-east-1",
							},
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
