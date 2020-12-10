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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicationGitPathValidate(t *testing.T) {
	testcases := []struct {
		name        string
		gitPath     *ApplicationGitPath
		expectedErr string
	}{
		{
			name:        "invalid: missing repo",
			gitPath:     &ApplicationGitPath{},
			expectedErr: "invalid ApplicationGitPath.Repo: value is required",
		},
		{
			name: "invalid: missing path",
			gitPath: &ApplicationGitPath{
				Repo: &ApplicationGitRepository{
					Id: "id",
				},
			},
			expectedErr: `invalid ApplicationGitPath.Path: value does not match regex pattern "^[^/].+$"`,
		},
		{
			name: "invalid: path must be relative",
			gitPath: &ApplicationGitPath{
				Repo: &ApplicationGitRepository{
					Id: "id",
				},
				Path: "/kubernetes/simple",
			},
			expectedErr: `invalid ApplicationGitPath.Path: value does not match regex pattern "^[^/].+$"`,
		},
		{
			name: "ok",
			gitPath: &ApplicationGitPath{
				Repo: &ApplicationGitRepository{
					Id: "id",
				},
				Path: "kuberntes/simple",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.gitPath.Validate()
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			assert.Equal(t, tc.expectedErr, errMsg)
		})
	}
}
