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
	"fmt"
	"path/filepath"
	"strings"
)

const DefaultDeploymentConfigFileName = ".pipe.yaml"

// GetDeploymentConfigFilePath returns the path to deployment configuration file.
func (p ApplicationGitPath) GetDeploymentConfigFilePath() string {
	filename := DefaultDeploymentConfigFileName
	if n := p.ConfigFilename; n != "" {
		filename = n
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

func MakeApplicationURL(baseURL, applicationID string) string {
	return fmt.Sprintf("%s/applications/%s", strings.TrimSuffix(baseURL, "/"), applicationID)
}

// ContainTags checks if it has all the given tags.
func (a *Application) ContainTags(tagIDs []string) bool {
	if len(a.Tags) < len(tagIDs) {
		return false
	}
	tagMap := make(map[string]struct{}, len(a.Tags))
	for i := range a.Tags {
		tagMap[a.Tags[i].Id] = struct{}{}
	}
	for _, tag := range tagIDs {
		if _, ok := tagMap[tag]; !ok {
			return false
		}
	}
	return true
}
