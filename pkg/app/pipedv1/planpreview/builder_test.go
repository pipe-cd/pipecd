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
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	planpreviewapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/planpreview"
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

type fakePlanPreviewPlugin struct {
	pluginapi.PluginClient
	name    string
	results []*planpreviewapi.PlanPreviewResult
	err     error
}

func (p *fakePlanPreviewPlugin) Name() string {
	return p.name
}

func (p *fakePlanPreviewPlugin) GetPlanPreview(context.Context, *planpreviewapi.GetPlanPreviewRequest, ...grpc.CallOption) (*planpreviewapi.GetPlanPreviewResponse, error) {
	if p.err != nil {
		return nil, p.err
	}
	return &planpreviewapi.GetPlanPreviewResponse{Results: p.results}, nil
}

func TestBuildPluginPlanResultsNoChange(t *testing.T) {
	t.Parallel()

	unchanged := &planpreviewapi.PlanPreviewResult{DeployTarget: "dt", NoChange: true}
	changed := &planpreviewapi.PlanPreviewResult{DeployTarget: "dt", NoChange: false}
	unimplemented := &fakePlanPreviewPlugin{name: "unimplemented", err: status.Error(codes.Unimplemented, "")}

	testcases := []struct {
		name     string
		plugins  []pluginapi.PluginClient
		expected bool
	}{
		{
			name: "all results have no change",
			plugins: []pluginapi.PluginClient{
				&fakePlanPreviewPlugin{name: "p1", results: []*planpreviewapi.PlanPreviewResult{unchanged, unchanged}},
				&fakePlanPreviewPlugin{name: "p2", results: []*planpreviewapi.PlanPreviewResult{unchanged}},
				unimplemented,
			},
			expected: true,
		},
		{
			name: "one result has changes",
			plugins: []pluginapi.PluginClient{
				&fakePlanPreviewPlugin{name: "p1", results: []*planpreviewapi.PlanPreviewResult{unchanged}},
				&fakePlanPreviewPlugin{name: "p2", results: []*planpreviewapi.PlanPreviewResult{unchanged, changed}},
			},
			expected: false,
		},
		{
			name:     "no plugin supports plan preview",
			plugins:  []pluginapi.PluginClient{unimplemented},
			expected: false,
		},
		{
			name: "a plugin fails",
			plugins: []pluginapi.PluginClient{
				&fakePlanPreviewPlugin{name: "p1", results: []*planpreviewapi.PlanPreviewResult{unchanged}},
				&fakePlanPreviewPlugin{name: "p2", err: errors.New("failed")},
			},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			b := &builder{pipedCfg: &config.PipedSpec{}}
			result := &model.ApplicationPlanPreviewResult{}
			b.buildPluginPlanResults(context.Background(), zap.NewNop(), &model.Application{}, tc.plugins, nil, nil, result)
			assert.Equal(t, tc.expected, result.NoChange)
		})
	}
}
