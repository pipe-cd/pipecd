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
				SyncInterval: Duration(time.Minute),
				Git: PipedGit{
					Username:        "username",
					Email:           "username@email.com",
					SSHKeyFile:      "/etc/pipecd-piped/ssh.key",
					AccessTokenFile: "/etc/pipecd-piped/github.token",
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
				Destinations: []PipedDestination{
					{
						Name: "default",
						Kubernetes: &PipedDestinationKubernetes{
							AllowNamespaces: []string{"dev"},
						},
					},
					{
						Name: "gcp-terraform",
						Terraform: &PipedDestinationTerraform{
							GCP: &PipedTerraformGCP{
								Project:         "gcp-project",
								Region:          "us-central1",
								CredentialsFile: "/etc/pipecd-piped/terraform.serviceaccount",
							},
						},
					},
					{
						Name: "aws-terraform",
						Terraform: &PipedDestinationTerraform{
							AWS: &PipedTerraformAWS{
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
