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

package grpcapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/datastore/datastoretest"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

// fakePipedTokenVerifier accepts every piped token.
type fakePipedTokenVerifier struct{}

func (fakePipedTokenVerifier) Verify(_ context.Context, _, _, _ string) error { return nil }

// pipedAuthContext returns a context carrying a valid piped token, built by
// running the real PipedTokenUnaryServerInterceptor over a stub handler so the
// test exercises only exported rpcauth APIs.
func pipedAuthContext(t *testing.T) context.Context {
	t.Helper()

	token := rpcauth.MakePipedToken("projectID", "pipedID", "pipedKey")
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %s", rpcauth.PipedTokenCredentials, token))
	ctx := metadata.NewIncomingContext(context.Background(), md)

	verifier := fakePipedTokenVerifier{}
	var authed context.Context
	interceptor := rpcauth.PipedTokenUnaryServerInterceptor(verifier, nil)
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "TestListApplications"}, func(c context.Context, _ interface{}) (interface{}, error) {
		authed = c
		return nil, nil
	})
	if err != nil {
		t.Fatalf("failed to build authorized context: %v", err)
	}
	return authed
}

func TestListApplicationsPagination(t *testing.T) {
	ctx := pipedAuthContext(t)

	t.Run("single page", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		apps := []*model.Application{
			{Id: "app-1", PipedId: "pipedID"},
			{Id: "app-2", PipedId: "pipedID"},
		}
		s := datastoretest.NewMockApplicationStore(ctrl)
		s.EXPECT().
			List(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, opts datastore.ListOptions) ([]*model.Application, string, error) {
				assert.Equal(t, listApplicationsPageSize, opts.Limit)
				assert.NotEmpty(t, opts.Orders, "cursor paging requires a stable order")
				return apps, "", nil
			})

		api := &PipedAPI{applicationStore: s}
		resp, err := api.ListApplications(ctx, &pipedservice.ListApplicationsRequest{})
		assert.NoError(t, err)
		assert.Equal(t, apps, resp.Applications)
	})

	t.Run("multiple pages are aggregated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		page1 := []*model.Application{{Id: "app-1", PipedId: "pipedID"}}
		page2 := []*model.Application{{Id: "app-2", PipedId: "pipedID"}}

		s := datastoretest.NewMockApplicationStore(ctrl)
		firstCall := s.EXPECT().
			List(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, opts datastore.ListOptions) ([]*model.Application, string, error) {
				assert.Empty(t, opts.Cursor, "first call must not carry a cursor")
				return page1, "cursor-1", nil
			})
		secondCall := s.EXPECT().
			List(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, opts datastore.ListOptions) ([]*model.Application, string, error) {
				assert.Equal(t, "cursor-1", opts.Cursor, "second call must continue from the first cursor")
				return page2, "", nil
			})
		gomock.InOrder(firstCall, secondCall)

		api := &PipedAPI{applicationStore: s}
		resp, err := api.ListApplications(ctx, &pipedservice.ListApplicationsRequest{})
		assert.NoError(t, err)
		want := make([]*model.Application, 0, len(page1)+len(page2))
		want = append(want, page1...)
		want = append(want, page2...)
		assert.Equal(t, want, resp.Applications)
	})

	t.Run("store error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := datastoretest.NewMockApplicationStore(ctrl)
		s.EXPECT().
			List(gomock.Any(), gomock.Any()).
			Return(nil, "", assert.AnError)

		api := &PipedAPI{applicationStore: s}
		resp, err := api.ListApplications(ctx, &pipedservice.ListApplicationsRequest{})
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}
