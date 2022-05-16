// Copyright 2022 The PipeCD Authors.
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

// GetResource returns *ProjectRBACResource and flag which represents wether the resource exists or not.
// If the resource exists in []*ProjectRBACPolicy but the action does not exist, this returns nil and false.
func (p *Role) GetResource(typ ProjectRBACResource_ResourceType, action ProjectRBACPolicy_Action) (*ProjectRBACResource, bool) {
	for _, v := range p.ProjectPolicies {
		res, ok := v.GetResource(typ, action)
		if ok {
			return res, true
		}
	}
	return nil, false
}
