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

// RunnerSpec contains configurable data used to while running Runner.
type RunnerSpec struct {
	Git               RunnerGit           `json:"git"`
	Repositories      []RunnerRepository  `json:"repositories"`
	Destinations      []RunnerDestination `json:"destinations"`
	AnalysisProviders []AnalysisProvider  `json:"analysisProviders"`
}

// Validate validates configured data of all fields.
func (s *RunnerSpec) Validate() error {
	return nil
}

type RunnerGit struct {
	// The path to the private ssh key file.
	// This will be used to clone the source code of the git repositories.
	SSHKeyFile string `json:"sshKeyFile"`
	// The path to the GitHub/GitLab access token file.
	// This will be used to authenticate while creating pull request...
	AccessTokenFile string `json:"accessTokenFile"`
}

type RunnerRepository struct {
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	// How often to check the new commit merged into the branch.
	PollInterval Duration `json:"pollInterval"`
}

type DestinationType string

const (
	KubernetesDestination DestinationType = "Kubernetes"
	TerraformDestination  DestinationType = "Terraform"
)

type RunnerDestination struct {
	Name       string                       `json:"name"`
	Type       DestinationType              `json:"-"`
	Kubernetes *RunnerDestinationKubernetes `json:"kubernetes"`
	Terraform  *RunnerDestinationTerraform  `json:"terraform"`
}

type RunnerDestinationKubernetes struct {
	AllowNamespaces []string `json:"allowNamespaces"`
}

type RunnerDestinationTerraform struct {
	GCP *RunnerTerraformGCP `json:"gcp"`
	AWS *RunnerTerraformAWS `json:"aws"`
}

type RunnerTerraformGCP struct {
	Project         string `json:"project"`
	Region          string `json:"region"`
	CredentialsFile string `json:"credentialsFile"`
}

type RunnerTerraformAWS struct {
	Region string `json:"region"`
}

type AnalysisProvider struct {
	Name string `json:"name"`
}
