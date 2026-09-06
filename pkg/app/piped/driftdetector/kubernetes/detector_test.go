// Copyright 2026 The PipeCD Authors.
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

package kubernetes

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/kubernetes"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/cache/memorycache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/git/gittest"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type fakeAppLister struct {
	apps []*model.Application
}

func (f *fakeAppLister) ListByPlatformProvider(name string) []*model.Application {
	return f.apps
}

type fakeStateGetter struct{}

func (fakeStateGetter) GetKubernetesAppLiveState(appID string) (kubernetes.AppState, bool) {
	return kubernetes.AppState{}, false
}

func (fakeStateGetter) NewEventIterator() kubernetes.EventIterator {
	return kubernetes.EventIterator{}
}

func (fakeStateGetter) GetWatchingResourceKinds() []provider.APIVersionKind {
	return nil
}

func (fakeStateGetter) GetAppLiveManifests(appID string) []provider.Manifest {
	return nil
}

func (fakeStateGetter) WaitForReady(ctx context.Context, timeout time.Duration) error {
	return nil
}

// TestCheck_ContextCanceled verifies that check() does not log the per-application
// failures as errors when the given context has already been canceled, e.g. during
// piped shutdown, while it still logs them as errors otherwise.
func TestCheck_ContextCanceled(t *testing.T) {
	testcases := []struct {
		name          string
		cancel        bool
		wantErrorLogs int
	}{
		{
			name:          "context not canceled: failures are logged as errors",
			cancel:        false,
			wantErrorLogs: 2,
		},
		{
			name:          "context canceled: failures are not logged as errors",
			cancel:        true,
			wantErrorLogs: 0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			core, logs := observer.New(zapcore.DebugLevel)
			logger := zap.New(core)

			repo := gittest.NewMockRepo(ctrl)
			repo.EXPECT().GetClonedBranch().Return("main").AnyTimes()
			repo.EXPECT().Pull(gomock.Any(), gomock.Any()).Return(nil)
			repo.EXPECT().GetLatestCommit(gomock.Any()).Return(git.Commit{Hash: "abc123"}, nil)
			// Point GetPath to an empty directory so loading the application
			// configuration fails quickly without touching any real repository.
			repo.EXPECT().GetPath().Return(t.TempDir()).AnyTimes()
			repo.EXPECT().CleanPath(gomock.Any(), gomock.Any()).Return(errors.New("clean failed"))

			app := &model.Application{
				Id:   "app-1",
				Kind: model.ApplicationKind_KUBERNETES,
				GitPath: &model.ApplicationGitPath{
					Repo: &model.ApplicationGitRepository{Id: "repo-1"},
					Path: "path/to/app",
				},
			}

			d := &detector{
				provider:          config.PipedPlatformProvider{Name: "kubernetes-default"},
				appLister:         &fakeAppLister{apps: []*model.Application{app}},
				stateGetter:       fakeStateGetter{},
				appManifestsCache: memorycache.NewCache(),
				config:            &config.PipedSpec{},
				logger:            logger,
				gitRepos:          map[string]git.Repo{"repo-1": repo},
				syncStates:        make(map[string]model.ApplicationSyncState),
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.cancel {
				cancel()
			}

			d.check(ctx)

			errorLogs := logs.FilterLevelExact(zapcore.ErrorLevel).All()
			assert.Len(t, errorLogs, tc.wantErrorLogs)
		})
	}
}
