// Copyright 2020 The Dianomi Authors.
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
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestUnaryServerInterceptor(t *testing.T) {
	in := UnaryServerInterceptor()
	testcases := []struct {
		name                    string
		ctx                     context.Context
		expectedCredentials     string
		expectedCredentialsType CredentialsType
		failed                  bool
	}{
		{
			name:                    "missing credentials",
			ctx:                     context.TODO(),
			expectedCredentials:     "",
			expectedCredentialsType: UnknownCredentials,
			failed:                  true,
		},
		{
			name: "should be ok with IDToken",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"ID-TOKEN token"},
			}),
			expectedCredentials:     "token",
			expectedCredentialsType: IDTokenCredentials,
			failed:                  false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := in(tc.ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
				creds, err := ExtractCredentials(ctx)
				if err != nil {
					return nil, err
				}
				if creds.Type != tc.expectedCredentialsType {
					return nil, errors.New("different credentials type")
				}
				if creds.Data != tc.expectedCredentials {
					return nil, errors.New("different credentials data")
				}
				return nil, nil
			})
			assert.Equal(t, tc.failed, err != nil)
		})
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	in := StreamServerInterceptor()
	testcases := []struct {
		name                    string
		ctx                     context.Context
		expectedCredentials     string
		expectedCredentialsType CredentialsType
		failed                  bool
	}{
		{
			name:                    "missing credentials",
			ctx:                     context.TODO(),
			expectedCredentials:     "",
			expectedCredentialsType: UnknownCredentials,
			failed:                  true,
		},
		{
			name: "should be ok with IDToken",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{"ID-TOKEN token"},
			}),
			expectedCredentials:     "token",
			expectedCredentialsType: IDTokenCredentials,
			failed:                  false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			stream := &fakeServerStream{
				ctx: tc.ctx,
			}
			err := in(nil, stream, nil, func(srv interface{}, stream grpc.ServerStream) error {
				ctx := stream.Context()
				creds, err := ExtractCredentials(ctx)
				if err != nil {
					return err
				}
				if creds.Type != tc.expectedCredentialsType {
					return errors.New("different credentials type")
				}
				if creds.Data != tc.expectedCredentials {
					return errors.New("different credentials data")
				}
				return nil
			})
			assert.Equal(t, tc.failed, err != nil)
		})
	}
}
