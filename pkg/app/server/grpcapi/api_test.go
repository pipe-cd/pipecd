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

package grpcapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

func TestRequireAPIKey(t *testing.T) {
	testcases := []struct {
		name        string
		key         *model.APIKey
		requireRole model.APIKey_Role
		expectedKey *model.APIKey
		expectedErr string
	}{
		{
			name: "ok: using READ_ONLY to read",
			key: &model.APIKey{
				Role: model.APIKey_READ_ONLY,
			},
			requireRole: model.APIKey_READ_ONLY,
			expectedKey: &model.APIKey{
				Role: model.APIKey_READ_ONLY,
			},
		},
		{
			name: "ok: using READ_WRITE to read",
			key: &model.APIKey{
				Role: model.APIKey_READ_WRITE,
			},
			requireRole: model.APIKey_READ_ONLY,
			expectedKey: &model.APIKey{
				Role: model.APIKey_READ_WRITE,
			},
		},
		{
			name: "ok: using READ_WRITE to write",
			key: &model.APIKey{
				Role: model.APIKey_READ_WRITE,
			},
			requireRole: model.APIKey_READ_WRITE,
			expectedKey: &model.APIKey{
				Role: model.APIKey_READ_WRITE,
			},
		},
		{
			name: "invalid: using READ_ONLY to write",
			key: &model.APIKey{
				Role: model.APIKey_READ_ONLY,
			},
			requireRole: model.APIKey_READ_WRITE,
			expectedErr: "rpc error: code = PermissionDenied desc = Permission denied",
		},
		{
			name: "invalid: invalid role",
			key: &model.APIKey{
				Role: -1,
			},
			requireRole: model.APIKey_READ_ONLY,
			expectedErr: "rpc error: code = PermissionDenied desc = Invalid role",
		},
		{
			name:        "invalid: api key was not included",
			requireRole: model.APIKey_READ_ONLY,
			expectedErr: "rpc error: code = Unauthenticated desc = Unauthenticated",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			if tc.key != nil {
				ctx = rpcauth.ContextWithAPIKey(ctx, tc.key)
			}
			key, err := requireAPIKey(ctx, tc.requireRole, zap.NewNop())
			assert.Equal(t, tc.expectedKey, key)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err.Error())
			} else {
				assert.Equal(t, tc.expectedErr, "")
			}
		})
	}
}

func TestRequestPlanPreviewCachesRepositoryRoutes(t *testing.T) {
	pipedStore := &requestPlanPreviewPipedStore{
		pipeds: []*model.Piped{
			{
				Id:        "piped-1",
				ProjectId: "project-1",
				Repositories: []*model.ApplicationGitRepository{
					{
						Id:     "repo-1",
						Remote: "https://github.com/pipe-cd/example.git",
						Branch: "master",
					},
				},
			},
			{
				Id:        "piped-2",
				ProjectId: "project-1",
				Repositories: []*model.ApplicationGitRepository{
					{
						Id:     "repo-2",
						Remote: "https://github.com/pipe-cd/example.git",
						Branch: "master",
					},
				},
			},
			{
				Id:        "piped-3",
				ProjectId: "project-1",
				Repositories: []*model.ApplicationGitRepository{
					{
						Id:     "repo-3",
						Remote: "https://github.com/pipe-cd/other.git",
						Branch: "master",
					},
				},
			},
		},
	}
	commandStore := &requestPlanPreviewCommandStore{}
	api := &API{
		pipedStore:                 pipedStore,
		commandStore:               commandStore,
		planPreviewRepositoryCache: newPlanPreviewRepositoryCache(1024, time.Minute),
		logger:                     zap.NewNop(),
	}
	ctx := rpcauth.ContextWithAPIKey(context.TODO(), &model.APIKey{
		ProjectId: "project-1",
		Role:      model.APIKey_READ_WRITE,
	})
	req := &apiservice.RequestPlanPreviewRequest{
		RepoRemoteUrl: "https://github.com/pipe-cd/example.git",
		BaseBranch:    "master",
		HeadBranch:    "feature",
		HeadCommit:    "abcdef",
		Timeout:       60,
	}

	resp, err := api.RequestPlanPreview(ctx, req)
	require.NoError(t, err)
	assert.Len(t, resp.Commands, 2)

	resp, err = api.RequestPlanPreview(ctx, req)
	require.NoError(t, err)
	assert.Len(t, resp.Commands, 2)
	assert.Equal(t, 1, pipedStore.listCalls)
	assert.Len(t, commandStore.commands, 4)
}

func TestRequestPlanPreviewCachesEmptyRepositoryRoutes(t *testing.T) {
	pipedStore := &requestPlanPreviewPipedStore{
		pipeds: []*model.Piped{
			{
				Id:        "piped-1",
				ProjectId: "project-1",
				Repositories: []*model.ApplicationGitRepository{
					{
						Id:     "repo-1",
						Remote: "https://github.com/pipe-cd/other.git",
						Branch: "master",
					},
				},
			},
		},
	}
	api := &API{
		pipedStore:                 pipedStore,
		commandStore:               &requestPlanPreviewCommandStore{},
		planPreviewRepositoryCache: newPlanPreviewRepositoryCache(1024, time.Minute),
		logger:                     zap.NewNop(),
	}
	ctx := rpcauth.ContextWithAPIKey(context.TODO(), &model.APIKey{
		ProjectId: "project-1",
		Role:      model.APIKey_READ_WRITE,
	})
	req := &apiservice.RequestPlanPreviewRequest{
		RepoRemoteUrl: "https://github.com/pipe-cd/example.git",
		BaseBranch:    "master",
	}

	resp, err := api.RequestPlanPreview(ctx, req)
	require.NoError(t, err)
	assert.Empty(t, resp.Commands)

	resp, err = api.RequestPlanPreview(ctx, req)
	require.NoError(t, err)
	assert.Empty(t, resp.Commands)
	assert.Equal(t, 1, pipedStore.listCalls)
}

func BenchmarkRequestPlanPreviewRepositoryRouting(b *testing.B) {
	for _, pipedCount := range []int{1, 100, 10000} {
		b.Run(fmt.Sprintf("pipeds-%d", pipedCount), func(b *testing.B) {
			pipedStore := &requestPlanPreviewPipedStore{
				pipeds: makeRequestPlanPreviewPipeds(pipedCount),
			}
			api := &API{
				pipedStore:                 pipedStore,
				commandStore:               &requestPlanPreviewCommandStore{discard: true},
				planPreviewRepositoryCache: newPlanPreviewRepositoryCache(1024, time.Minute),
				logger:                     zap.NewNop(),
			}
			ctx := rpcauth.ContextWithAPIKey(context.TODO(), &model.APIKey{
				ProjectId: "project-1",
				Role:      model.APIKey_READ_WRITE,
			})
			req := &apiservice.RequestPlanPreviewRequest{
				RepoRemoteUrl: "https://github.com/pipe-cd/example.git",
				BaseBranch:    "master",
				HeadBranch:    "feature",
				HeadCommit:    "abcdef",
				Timeout:       60,
			}
			if _, err := api.RequestPlanPreview(ctx, req); err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := api.RequestPlanPreview(ctx, req); err != nil {
					b.Fatal(err)
				}
			}
			b.StopTimer()

			if pipedStore.listCalls != 1 {
				b.Fatalf("expected one piped list call, got %d", pipedStore.listCalls)
			}
		})
	}
}

func makeRequestPlanPreviewPipeds(n int) []*model.Piped {
	pipeds := make([]*model.Piped, 0, n)
	for i := 0; i < n; i++ {
		repoRemoteURL := fmt.Sprintf("https://github.com/pipe-cd/other-%d.git", i)
		if i == n-1 {
			repoRemoteURL = "https://github.com/pipe-cd/example.git"
		}
		pipeds = append(pipeds, &model.Piped{
			Id:        fmt.Sprintf("piped-%d", i),
			ProjectId: "project-1",
			Repositories: []*model.ApplicationGitRepository{
				{
					Id:     fmt.Sprintf("repo-%d", i),
					Remote: repoRemoteURL,
					Branch: "master",
				},
			},
		})
	}
	return pipeds
}

type requestPlanPreviewPipedStore struct {
	pipeds    []*model.Piped
	listCalls int
}

func (s *requestPlanPreviewPipedStore) Get(context.Context, string) (*model.Piped, error) {
	return nil, datastore.ErrNotFound
}

func (s *requestPlanPreviewPipedStore) List(context.Context, datastore.ListOptions) ([]*model.Piped, error) {
	s.listCalls++
	return s.pipeds, nil
}

func (s *requestPlanPreviewPipedStore) Add(context.Context, *model.Piped) error {
	return nil
}

func (s *requestPlanPreviewPipedStore) UpdateInfo(context.Context, string, string, string) error {
	return nil
}

func (s *requestPlanPreviewPipedStore) EnablePiped(context.Context, string) error {
	return nil
}

func (s *requestPlanPreviewPipedStore) DisablePiped(context.Context, string) error {
	return nil
}

type requestPlanPreviewCommandStore struct {
	commands []*model.Command
	discard  bool
}

func (s *requestPlanPreviewCommandStore) ListUnhandledCommands(context.Context, string) ([]*model.Command, error) {
	return nil, nil
}

func (s *requestPlanPreviewCommandStore) AddCommand(_ context.Context, command *model.Command) error {
	if s.discard {
		return nil
	}
	s.commands = append(s.commands, command)
	return nil
}

func (s *requestPlanPreviewCommandStore) GetCommand(context.Context, string) (*model.Command, error) {
	return nil, datastore.ErrNotFound
}

func (s *requestPlanPreviewCommandStore) UpdateCommandHandled(context.Context, string, model.CommandStatus, map[string]string, int64) error {
	return nil
}
