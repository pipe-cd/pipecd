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

type InputHelmOptions struct {
	// The release name of helm deployment.
	// By default the release name is equal to the application name.
	ReleaseName string `json:"releaseName,omitempty"`
	// List of values.
	SetValues map[string]string `json:"setValues,omitempty"`
	// List of value files should be loaded.
	ValueFiles []string `json:"valueFiles,omitempty"`
	// List of file path for values.
	SetFiles map[string]string `json:"setFiles,omitempty"`
	// Set of supported Kubernetes API versions.
	APIVersions []string `json:"apiVersions,omitempty"`
	// Kubernetes version used for Capabilities.KubeVersion
	KubeVersion string `json:"kubeVersion,omitempty"`
}

type InputHelmChart struct {
	// Git remote address where the chart is placing.
	// Empty means the same repository.
	GitRemote string `json:"gitRemote,omitempty"`
	// The commit SHA or tag for remote git.
	Ref string `json:"ref,omitempty"`
	// Relative path from the repository root directory to the chart directory.
	Path string `json:"path,omitempty"`

	// The name of an added Helm Chart Repository.
	Repository string `json:"repository,omitempty"`
	Name       string `json:"name,omitempty"`
	Version    string `json:"version,omitempty"`
	// Whether to skip TLS certificate checks for the repository or not.
	// This option will automatically set the value of HelmChartRepository.Insecure.
	Insecure bool `json:"-"`
}
