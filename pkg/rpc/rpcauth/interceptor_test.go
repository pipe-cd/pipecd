// Copyright 2023 The PipeCD Authors.
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

package rpcauth

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type fakeServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *fakeServerStream) Context() context.Context {
	return s.ctx
}

type testPipedTokenVerifier struct {
	pipedKey string
}

func (v testPipedTokenVerifier) Verify(ctx context.Context, projectID, pipedID, pipedKey string) error {
	if pipedKey != v.pipedKey {
		return fmt.Errorf("invalid piped key, want: %s, got: %s", v.pipedKey, pipedKey)
	}
	return nil
}

func TestPipedTokenUnaryServerInterceptor(t *testing.T) {
	verifier := testPipedTokenVerifier{"test-piped-key"}
	in := PipedTokenUnaryServerInterceptor(verifier, zap.NewNop())
	testcases := []struct {
		name             string
		ctx              context.Context
		expectedPipedKey string
		failed           bool
	}{
		{
			name:             "missing credentials",
			ctx:              context.TODO(),
			expectedPipedKey: "",
			failed:           true,
		},
		{
			name: "wrong credentials type",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"ID-TOKEN test-project-id,test-piped-id,test-piped-key"},
			}),
			expectedPipedKey: "",
			failed:           true,
		},
		{
			name: "malformed piped token",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"PIPED-TOKEN test-piped-key"},
			}),
			expectedPipedKey: "",
			failed:           true,
		},
		{
			name: "should be ok with PipedToken",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"PIPED-TOKEN test-project-id,test-piped-id,test-piped-key"},
			}),
			expectedPipedKey: "test-piped-key",
			failed:           false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := in(tc.ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
				_, _, pipedKey, err := ExtractPipedToken(ctx)
				if err != nil {
					return nil, err
				}
				if pipedKey != tc.expectedPipedKey {
					return nil, errors.New("invalid piped key")
				}
				return nil, nil
			})
			assert.Equal(t, tc.failed, err != nil)
		})
	}
}

func TestPipedTokenStreamServerInterceptor(t *testing.T) {
	verifier := testPipedTokenVerifier{"test-piped-key"}
	in := PipedTokenStreamServerInterceptor(verifier, zap.NewNop())
	testcases := []struct {
		name             string
		ctx              context.Context
		expectedPipedKey string
		failed           bool
	}{
		{
			name:             "missing credentials",
			ctx:              context.TODO(),
			expectedPipedKey: "",
			failed:           true,
		},
		{
			name: "wrong credentials type",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"ID-TOKEN test-project-id,test-piped-id,test-piped-key"},
			}),
			expectedPipedKey: "",
			failed:           true,
		},
		{
			name: "malformed piped token",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"PIPED-TOKEN test-piped-key"},
			}),
			expectedPipedKey: "",
			failed:           true,
		},
		{
			name: "should be ok with PipedToken",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"PIPED-TOKEN test-project-id,test-piped-id,test-piped-key"},
			}),
			expectedPipedKey: "test-piped-key",
			failed:           false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			stream := &fakeServerStream{
				ctx: tc.ctx,
			}
			err := in(nil, stream, nil, func(srv interface{}, stream grpc.ServerStream) error {
				ctx := stream.Context()
				_, _, pipedKey, err := ExtractPipedToken(ctx)
				if err != nil {
					return err
				}
				if pipedKey != tc.expectedPipedKey {
					return errors.New("invalid piped key")
				}
				return nil
			})
			assert.Equal(t, tc.failed, err != nil)
		})
	}
}

type testAPIKeyVerifier struct {
	keyString string
	key       *model.APIKey
}

func (v testAPIKeyVerifier) Verify(_ context.Context, key string) (*model.APIKey, error) {
	if key != v.keyString {
		return nil, fmt.Errorf("invalid API key, want: %s, got: %s", v.keyString, key)
	}
	return v.key, nil
}

func TestAPIKeyUnaryServerInterceptor(t *testing.T) {
	verifier := testAPIKeyVerifier{
		keyString: "test-api-key",
		key: &model.APIKey{
			Id: "test-api-key",
		},
	}
	in := APIKeyUnaryServerInterceptor(verifier, zap.NewNop())
	testcases := []struct {
		name        string
		ctx         context.Context
		expectedKey *model.APIKey
		errString   string
	}{
		{
			name:      "missing credentials",
			ctx:       context.TODO(),
			errString: "rpc error: code = Unauthenticated desc = missing credentials",
		},
		{
			name: "wrong credentials type",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"PIPED-TOKEN test-project-id,test-piped-id,test-piped-key"},
			}),
			errString: "rpc error: code = Unauthenticated desc = Unauthenticated",
		},
		{
			name: "ok",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"API-KEY test-api-key"},
			}),
			expectedKey: &model.APIKey{
				Id: "test-api-key",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := in(tc.ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
				apiKey, err := ExtractAPIKey(ctx)
				if err != nil {
					return nil, err
				}
				if apiKey.Id != tc.expectedKey.Id {
					return nil, errors.New("invalid api key")
				}
				return nil, nil
			})
			if tc.errString != "" {
				require.NotNil(t, err)
				assert.Equal(t, tc.errString, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
