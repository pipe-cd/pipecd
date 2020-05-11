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

// PipedSpec contains configurable data used to while running Piped.
type PipedSpec struct {
	// How often to check whether an application should be synced.
	SyncInterval Duration `json:"syncInterval"`
	// Git configuration needed for git commands.
	Git PipedGit `json:"git"`
	// List of git repositories this piped will handle.
	Repositories      []PipedRepository  `json:"repositories"`
	Destinations      []PipedDestination `json:"destinations"`
	AnalysisProviders []AnalysisProvider `json:"analysisProviders"`
}

// Validate validates configured data of all fields.
func (s *PipedSpec) Validate() error {
	return nil
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

type PipedGit struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	// The path to the private ssh key file.
	// This will be used to clone the source code of the git repositories.
	SSHKeyFile string `json:"sshKeyFile"`
	// The path to the GitHub/GitLab access token file.
	// This will be used to authenticate while creating pull request...
	AccessTokenFile string `json:"accessTokenFile"`
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

type DestinationType string

const (
	KubernetesDestination DestinationType = "Kubernetes"
	TerraformDestination  DestinationType = "Terraform"
)

type PipedDestination struct {
	Name       string                      `json:"name"`
	Type       DestinationType             `json:"-"`
	Kubernetes *PipedDestinationKubernetes `json:"kubernetes"`
	Terraform  *PipedDestinationTerraform  `json:"terraform"`
}

type PipedDestinationKubernetes struct {
	AllowNamespaces []string `json:"allowNamespaces"`
}

type PipedDestinationTerraform struct {
	GCP *PipedTerraformGCP `json:"gcp"`
	AWS *PipedTerraformAWS `json:"aws"`
}

type PipedTerraformGCP struct {
	Project         string `json:"project"`
	Region          string `json:"region"`
	CredentialsFile string `json:"credentialsFile"`
}

type PipedTerraformAWS struct {
	Region string `json:"region"`
}

type AnalysisProvider struct {
	Name string `json:"name"`
}
