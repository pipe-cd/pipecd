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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockPipedStore struct {
	apiPipedStore
	getFunc func(ctx context.Context, id string) (*model.Piped, error)
}

func (m *mockPipedStore) Get(ctx context.Context, id string) (*model.Piped, error) {
	return m.getFunc(ctx, id)
}

type mockApplicationStore struct {
	apiApplicationStore
	listFunc func(ctx context.Context, opts datastore.ListOptions) ([]*model.Application, string, error)
	addFunc  func(ctx context.Context, app *model.Application) error
}

func (m *mockApplicationStore) List(ctx context.Context, opts datastore.ListOptions) ([]*model.Application, string, error) {
	return m.listFunc(ctx, opts)
}

func (m *mockApplicationStore) Add(ctx context.Context, app *model.Application) error {
	return m.addFunc(ctx, app)
}

func TestAddApplication(t *testing.T) {
	testcases := []struct {
		name         string
		req          *apiservice.AddApplicationRequest
		piped        *model.Piped
		existingApps []*model.Application
		expectErr    bool
		expectedCode codes.Code
	}{
		{
			name: "ok: no duplicates",
			req: &apiservice.AddApplicationRequest{
				Name:    "app-1",
				PipedId: "piped-1",
				GitPath: &model.ApplicationGitPath{
					Repo:           &model.ApplicationGitRepository{Id: "repo-1"},
					Path:           "path-1",
					ConfigFilename: "app.yaml",
				},
				Kind:             model.ApplicationKind_KUBERNETES,
				PlatformProvider: "provider-1",
			},
			piped: &model.Piped{
				Id:        "piped-1",
				ProjectId: "project-1",
				Repositories: []*model.ApplicationGitRepository{
					{
						Id:     "repo-1",
						Remote: "https://github.com/org/repo",
						Branch: "master",
					},
				},
			},
			existingApps: []*model.Application{},
			expectErr:    false,
		},
		{
			name: "error: duplicate found",
			req: &apiservice.AddApplicationRequest{
				Name:    "app-1",
				PipedId: "piped-1",
				GitPath: &model.ApplicationGitPath{
					Repo:           &model.ApplicationGitRepository{Id: "repo-1"},
					Path:           "path-1",
					ConfigFilename: "app.yaml",
				},
				Kind:             model.ApplicationKind_KUBERNETES,
				PlatformProvider: "provider-1",
			},
			piped: &model.Piped{
				Id:        "piped-1",
				ProjectId: "project-1",
				Repositories: []*model.ApplicationGitRepository{
					{
						Id:     "repo-1",
						Remote: "https://github.com/org/repo",
						Branch: "master",
					},
				},
			},
			existingApps: []*model.Application{
				{
					Id:      "existing-1",
					Name:    "app-1",
					PipedId: "piped-1",
					GitPath: &model.ApplicationGitPath{
						Repo:           &model.ApplicationGitRepository{Id: "repo-1"},
						Path:           "path-1",
						ConfigFilename: "app.yaml",
					},
				},
			},
			expectErr:    true,
			expectedCode: codes.AlreadyExists,
		},
		{
			name: "ok: same name different path",
			req: &apiservice.AddApplicationRequest{
				Name:    "app-1",
				PipedId: "piped-1",
				GitPath: &model.ApplicationGitPath{
					Repo:           &model.ApplicationGitRepository{Id: "repo-1"},
					Path:           "path-2",
					ConfigFilename: "app.yaml",
				},
				Kind:             model.ApplicationKind_KUBERNETES,
				PlatformProvider: "provider-1",
			},
			piped: &model.Piped{
				Id:        "piped-1",
				ProjectId: "project-1",
				Repositories: []*model.ApplicationGitRepository{
					{
						Id:     "repo-1",
						Remote: "https://github.com/org/repo",
						Branch: "master",
					},
				},
			},
			existingApps: []*model.Application{},
			expectErr:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ps := &mockPipedStore{
				getFunc: func(ctx context.Context, id string) (*model.Piped, error) {
					return tc.piped, nil
				},
			}
			as := &mockApplicationStore{
				listFunc: func(ctx context.Context, opts datastore.ListOptions) ([]*model.Application, string, error) {
					// Verify filters
					filters := make(map[string]interface{})
					for _, f := range opts.Filters {
						filters[f.Field] = f.Value
					}
					assert.Equal(t, "project-1", filters["ProjectId"])
					assert.Equal(t, tc.req.GitPath.Repo.Id, filters["GitPath.Repo.Id"])
					assert.Equal(t, tc.req.GitPath.Path, filters["GitPath.Path"])
					assert.Equal(t, tc.req.GitPath.ConfigFilename, filters["GitPath.ConfigFilename"])

					return tc.existingApps, "", nil
				},
				addFunc: func(ctx context.Context, app *model.Application) error {
					return nil
				},
			}

			api := &API{
				pipedStore:       ps,
				applicationStore: as,
				logger:           zap.NewNop(),
			}

			ctx := rpcauth.ContextWithAPIKey(context.Background(), &model.APIKey{
				ProjectId: "project-1",
				Role:      model.APIKey_READ_WRITE,
			})

			_, err := api.AddApplication(ctx, tc.req)
			if tc.expectErr {
				assert.Error(t, err)
				s, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCode, s.Code())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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
