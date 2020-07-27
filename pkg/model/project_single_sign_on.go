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
)

// CreateAuthURL creates a auth url.
func (p *ProjectSingleSignOn) CreateAuthURL(project, apiURL, callbackPath, state string) (string, error) {
	switch p.Provider {
	case ProjectSingleSignOnProvider_GITHUB:
		if p.Github == nil {
			return "", fmt.Errorf("there are no date in the project github oauth configurations")
		}
		return p.Github.CreateAuthURL(project, apiURL, callbackPath, state)
	default:
		return "", fmt.Errorf("not implemented")
	}
}
