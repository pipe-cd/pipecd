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

package transfer

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// fakeAPIClient records the requests that restoreApplications sends.
// Only the methods restoreApplications calls are implemented; the embedded
// interface is nil, so any other call fails loudly.
type fakeAPIClient struct {
	apiservice.Client

	addErr   error
	added    []*apiservice.AddApplicationRequest
	disabled []string
}

func (f *fakeAPIClient) AddApplication(_ context.Context, req *apiservice.AddApplicationRequest, _ ...grpc.CallOption) (*apiservice.AddApplicationResponse, error) {
	if f.addErr != nil {
		return nil, f.addErr
	}
	f.added = append(f.added, req)
	return &apiservice.AddApplicationResponse{ApplicationId: fmt.Sprintf("new-%d", len(f.added))}, nil
}

func (f *fakeAPIClient) DisableApplication(_ context.Context, req *apiservice.DisableApplicationRequest, _ ...grpc.CallOption) (*apiservice.DisableApplicationResponse, error) {
	f.disabled = append(f.disabled, req.ApplicationId)
	return &apiservice.DisableApplicationResponse{}, nil
}

func TestRestoreApplications(t *testing.T) {
	t.Parallel()

	pipedIDMap := map[string]string{"old-piped": "new-piped"}
	app := &model.Application{
		Id:          "old-app",
		Name:        "backend",
		PipedId:     "old-piped",
		GitPath:     &model.ApplicationGitPath{Repo: &model.ApplicationGitRepository{Id: "repo"}, Path: "apps/backend"},
		Description: "backend service",
		Labels:      map[string]string{"team": "backend", "env": "prod"},
	}

	t.Run("restored application keeps the backed-up labels", func(t *testing.T) {
		t.Parallel()
		cli := &fakeAPIClient{}

		restored, failed := restoreApplications(context.Background(), cli, []*model.Application{app}, pipedIDMap, zap.NewNop())

		assert.Equal(t, 1, restored)
		assert.Equal(t, 0, failed)
		require.Len(t, cli.added, 1)
		req := cli.added[0]
		assert.Equal(t, "backend", req.Name)
		assert.Equal(t, "new-piped", req.PipedId)
		assert.Equal(t, app.GitPath, req.GitPath)
		assert.Equal(t, "backend service", req.Description)
		assert.Equal(t, map[string]string{"team": "backend", "env": "prod"}, req.Labels)
		assert.Empty(t, cli.disabled)
	})

	t.Run("disabled application is re-disabled after creation", func(t *testing.T) {
		t.Parallel()
		cli := &fakeAPIClient{}
		disabledApp := &model.Application{Id: "old-app", Name: "backend", PipedId: "old-piped", Disabled: true}

		restored, failed := restoreApplications(context.Background(), cli, []*model.Application{disabledApp}, pipedIDMap, zap.NewNop())

		assert.Equal(t, 1, restored)
		assert.Equal(t, 0, failed)
		assert.Equal(t, []string{"new-1"}, cli.disabled)
	})

	t.Run("application without a piped mapping is skipped", func(t *testing.T) {
		t.Parallel()
		cli := &fakeAPIClient{}
		unmapped := &model.Application{Id: "old-app", Name: "backend", PipedId: "unknown-piped"}

		restored, failed := restoreApplications(context.Background(), cli, []*model.Application{unmapped}, pipedIDMap, zap.NewNop())

		assert.Equal(t, 0, restored)
		assert.Equal(t, 1, failed)
		assert.Empty(t, cli.added)
	})

	t.Run("failed creation is counted", func(t *testing.T) {
		t.Parallel()
		cli := &fakeAPIClient{addErr: errors.New("boom")}

		restored, failed := restoreApplications(context.Background(), cli, []*model.Application{app}, pipedIDMap, zap.NewNop())

		assert.Equal(t, 0, restored)
		assert.Equal(t, 1, failed)
	})
}
