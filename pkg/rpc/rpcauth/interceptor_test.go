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

package rpcauth

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type fakeServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *fakeServerStream) Context() context.Context {
	return s.ctx
}

type testRunnerTokenVerifier struct {
	runnerKey string
}

func (v testRunnerTokenVerifier) Verify(projectID, runnerID, runnerKey string) error {
	if runnerKey != v.runnerKey {
		return fmt.Errorf("invalid runner key, want: %s, got: %s", v.runnerKey, runnerKey)
	}
	return nil
}

func TestRunnerTokenUnaryServerInterceptor(t *testing.T) {
	verifier := testRunnerTokenVerifier{"test-runner-key"}
	in := RunnerTokenUnaryServerInterceptor(verifier, zap.NewNop())
	testcases := []struct {
		name              string
		ctx               context.Context
		expectedRunnerKey string
		failed            bool
	}{
		{
			name:              "missing credentials",
			ctx:               context.TODO(),
			expectedRunnerKey: "",
			failed:            true,
		},
		{
			name: "wrong credentials type",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"ID-TOKEN test-project-id,test-runner-id,test-runner-key"},
			}),
			expectedRunnerKey: "",
			failed:            true,
		},
		{
			name: "malformed runner token",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"RUNNER-TOKEN test-runner-key"},
			}),
			expectedRunnerKey: "",
			failed:            true,
		},
		{
			name: "should be ok with RunnerToken",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"RUNNER-TOKEN test-project-id,test-runner-id,test-runner-key"},
			}),
			expectedRunnerKey: "test-runner-key",
			failed:            false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := in(tc.ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
				_, _, runnerKey, err := ExtractRunnerToken(ctx)
				if err != nil {
					return nil, err
				}
				if runnerKey != tc.expectedRunnerKey {
					return nil, errors.New("invalid runner key")
				}
				return nil, nil
			})
			assert.Equal(t, tc.failed, err != nil)
		})
	}
}

func TestRunnerTokenStreamServerInterceptor(t *testing.T) {
	verifier := testRunnerTokenVerifier{"test-runner-key"}
	in := RunnerTokenStreamServerInterceptor(verifier, zap.NewNop())
	testcases := []struct {
		name              string
		ctx               context.Context
		expectedRunnerKey string
		failed            bool
	}{
		{
			name:              "missing credentials",
			ctx:               context.TODO(),
			expectedRunnerKey: "",
			failed:            true,
		},
		{
			name: "wrong credentials type",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"ID-TOKEN test-project-id,test-runner-id,test-runner-key"},
			}),
			expectedRunnerKey: "",
			failed:            true,
		},
		{
			name: "malformed runner token",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"RUNNER-TOKEN test-runner-key"},
			}),
			expectedRunnerKey: "",
			failed:            true,
		},
		{
			name: "should be ok with RunnerToken",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"RUNNER-TOKEN test-project-id,test-runner-id,test-runner-key"},
			}),
			expectedRunnerKey: "test-runner-key",
			failed:            false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			stream := &fakeServerStream{
				ctx: tc.ctx,
			}
			err := in(nil, stream, nil, func(srv interface{}, stream grpc.ServerStream) error {
				ctx := stream.Context()
				_, _, runnerKey, err := ExtractRunnerToken(ctx)
				if err != nil {
					return err
				}
				if runnerKey != tc.expectedRunnerKey {
					return errors.New("invalid runner key")
				}
				return nil
			})
			assert.Equal(t, tc.failed, err != nil)
		})
	}
}
