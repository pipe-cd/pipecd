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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRole_GetResource(t *testing.T) {
	type args struct {
		resourceType ProjectRBACResource_ResourceType
		action       ProjectRBACPolicy_Action
	}
	testcases := []struct {
		name   string
		role   *Role
		want   *ProjectRBACResource
		exists bool
		args   args
	}{
		{
			name: "role has the permission",
			role: &Role{
				ProjectPolicies: []*ProjectRBACPolicy{
					{
						Resources: []*ProjectRBACResource{
							{
								Type: ProjectRBACResource_APPLICATION,
							},
						},
						Actions: []ProjectRBACPolicy_Action{
							ProjectRBACPolicy_ALL,
						},
					},
				},
			},
			want:   &ProjectRBACResource{Type: ProjectRBACResource_APPLICATION},
			exists: true,
			args: args{
				resourceType: ProjectRBACResource_APPLICATION,
				action:       ProjectRBACPolicy_CREATE,
			},
		},
		{
			name: "role does not have the permission",
			role: &Role{
				ProjectPolicies: []*ProjectRBACPolicy{
					{
						Resources: []*ProjectRBACResource{
							{
								Type: ProjectRBACResource_APPLICATION,
							},
						},
						Actions: []ProjectRBACPolicy_Action{
							ProjectRBACPolicy_GET,
						},
					},
				},
			},
			want:   nil,
			exists: false,
			args: args{
				resourceType: ProjectRBACResource_APPLICATION,
				action:       ProjectRBACPolicy_CREATE,
			},
		},
		{
			name: "role does not have the resource permission",
			role: &Role{
				ProjectPolicies: []*ProjectRBACPolicy{
					{
						Resources: []*ProjectRBACResource{
							{
								Type: ProjectRBACResource_APPLICATION,
							},
						},
						Actions: []ProjectRBACPolicy_Action{
							ProjectRBACPolicy_GET,
						},
					},
				},
			},
			want:   nil,
			exists: false,
			args: args{
				resourceType: ProjectRBACResource_APPLICATION,
				action:       ProjectRBACPolicy_CREATE,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ret, ok := tc.role.GetResource(tc.args.resourceType, tc.args.action)
			assert.Equal(t, tc.exists, ok)
			assert.Equal(t, tc.want, ret)
		})
	}

}
