// Copyright 2025 The PipeCD Authors.
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

package planpreview

import (
	"testing"

	"github.com/stretchr/testify/assert"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type fakeApplicationLister struct {
	apps []*model.Application
}

func (l *fakeApplicationLister) List() []*model.Application {
	return l.apps
}

func TestListApplications(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		repo     config.PipedRepository
		apps     []*model.Application
		expected []*model.Application
	}{
		{
			name: "no applications",
			repo: config.PipedRepository{
				RepoID: "repo-1",
				Remote: "git@github.com:org/repo-1.git",
				Branch: "main",
			},
			apps:     []*model.Application{},
			expected: []*model.Application{},
		},
		{
			name: "filter by repository configuration",
			repo: config.PipedRepository{
				RepoID: "repo-1",
				Remote: "git@github.com:org/repo-1.git",
				Branch: "main",
			},
			apps: []*model.Application{
				{
					Id: "app-1",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-1",
							Remote: "git@github.com:org/repo-1.git",
							Branch: "main",
						},
					},
				},
				{
					Id: "app-2",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-2", // Different repo ID
							Remote: "git@github.com:org/repo-1.git",
							Branch: "main",
						},
					},
				},
				{
					Id: "app-3",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-1",
							Remote: "git@github.com:org/repo-2.git", // Different remote
							Branch: "main",
						},
					},
				},
				{
					Id: "app-4",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-1",
							Remote: "git@github.com:org/repo-1.git",
							Branch: "develop", // Different branch
						},
					},
				},
			},
			expected: []*model.Application{
				{
					Id: "app-1",
					GitPath: &model.ApplicationGitPath{
						Repo: &model.ApplicationGitRepository{
							Id:     "repo-1",
							Remote: "git@github.com:org/repo-1.git",
							Branch: "main",
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			b := &builder{
				applicationLister: &fakeApplicationLister{
					apps: tc.apps,
				},
			}

			got := b.listApplications(tc.repo)
			assert.Equal(t, tc.expected, got)
		})
	}
}
