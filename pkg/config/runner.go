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

type RunnerSpec struct {
	Git               RunnerGit           `json:"git"`
	Repositories      []RunnerRepository  `json:"repositories"`
	Destinations      []RunnerDestination `json:"destinations"`
	AnalysisProviders []AnalysisProvider  `json:"analysisProviders"`
}

func (s *RunnerSpec) Validate() error {
	return nil
}

type RunnerGit struct {
	SSSKeyFile      string `json:"sshKeyFile"`
	AccessTokenFile string `json:"accessTokenFile"`
}

type RunnerRepository struct {
	Repo         string   `json:"repo"`
	Branch       string   `json:"branch"`
	SyncInterval Duration `json:"syncInterval"`
	PollInterval Duration `json:"pollInterval"`
}

type DestinationType string

const (
	KubernetesDestination DestinationType = "Kubernetes"
	TerraformDestination  DestinationType = "Terraform"
)

type RunnerDestination struct {
	Name string          `json:"name"`
	Type DestinationType `json:"type"`
}

type AnalysisProvider struct {
	Name string `json:"name"`
}
