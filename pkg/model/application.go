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

package model

import (
	"path/filepath"
)

// GetDeploymentConfigFilePath returns the path to deployment configuration directory.
func (p ApplicationGitPath) GetDeploymentConfigFilePath(filename string) string {
	if path := p.ConfigPath; path != "" {
		return path
	}
	return filepath.Join(p.Path, filename)
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
