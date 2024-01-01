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

package model

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	DefaultApplicationConfigFilename    = "app.pipecd.yaml"
	oldDefaultApplicationConfigFilename = ".pipe.yaml"
	applicationConfigFileExtention      = ".pipecd.yaml"
)

// GetApplicationConfigFilePath returns the path to application configuration file.
func (p ApplicationGitPath) GetApplicationConfigFilePath() string {
	return filepath.Join(p.Path, p.GetApplicationConfigFilename())
}

func (p ApplicationGitPath) GetApplicationConfigFilename() string {
	// The config file name used to allow to be empty until the default name got changed.
	// So empty means the old default name.
	filename := oldDefaultApplicationConfigFilename
	if n := p.ConfigFilename; n != "" {
		filename = n
	}
	return filename
}

// HasChanged checks whether the content of sync state has been changed.
// This ignores the timestamp value.
func (s ApplicationSyncState) HasChanged(next ApplicationSyncState) bool {
	if s.Status != next.Status {
		return true
	}
	if s.ShortReason != next.ShortReason {
		return true
	}
	if s.Reason != next.Reason {
		return true
	}
	return false
}

func MakeApplicationURL(baseURL, applicationID string) string {
	return fmt.Sprintf("%s/applications/%s", strings.TrimSuffix(baseURL, "/"), applicationID)
}

// ContainLabels checks if it has all the given labels.
func (a *Application) ContainLabels(labels map[string]string) bool {
	if len(a.Labels) < len(labels) {
		return false
	}

	for k, v := range labels {
		value, ok := a.Labels[k]
		if !ok {
			return false
		}
		if value != v {
			return false
		}
	}
	return true
}

func (a *Application) IsOutOfSync() bool {
	if a.SyncState == nil {
		return false
	}
	return a.SyncState.Status == ApplicationSyncStatus_OUT_OF_SYNC
}

func (a *Application) SetUpdatedAt(t int64) {
	a.UpdatedAt = t
}

func (ak ApplicationKind) CompatiblePlatformProviderType() PlatformProviderType {
	switch ak {
	case ApplicationKind_KUBERNETES:
		return PlatformProviderKubernetes
	case ApplicationKind_TERRAFORM:
		return PlatformProviderTerraform
	case ApplicationKind_LAMBDA:
		return PlatformProviderLambda
	case ApplicationKind_CLOUDRUN:
		return PlatformProviderCloudRun
	case ApplicationKind_ECS:
		return PlatformProviderECS
	default:
		return PlatformProviderKubernetes
	}
}

func IsApplicationConfigFile(filename string) bool {
	return filename == DefaultApplicationConfigFilename || strings.HasSuffix(filename, applicationConfigFileExtention) || filename == oldDefaultApplicationConfigFilename
}
