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

package grpcapi

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/filestore"
)

func TestGRPCStoreError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		inputErr error
		inputMsg string
		expected error
	}{
		{
			name:     "datastore not found error",
			inputErr: datastore.ErrNotFound,
			inputMsg: "datastore",
			expected: status.Error(codes.NotFound, "Entity was not found to datastore"),
		},
		{
			name:     "filestore not found error",
			inputErr: filestore.ErrNotFound,
			inputMsg: "filestore",
			expected: status.Error(codes.NotFound, "Entity was not found to filestore"),
		},
		{
			name:     "stagelogstore not found error",
			inputErr: filestore.ErrNotFound,
			inputMsg: "stagelogstore",
			expected: status.Error(codes.NotFound, "Entity was not found to stagelogstore"),
		},
		{
			name:     "datastore invalid argument error",
			inputErr: datastore.ErrInvalidArgument,
			inputMsg: "datastore",
			expected: status.Error(codes.InvalidArgument, "Invalid argument to datastore"),
		},
		{
			name:     "datastore already exists error",
			inputErr: datastore.ErrAlreadyExists,
			inputMsg: "datastore",
			expected: status.Error(codes.AlreadyExists, "Entity already exists to datastore"),
		},
		{
			name:     "user defined error",
			inputErr: datastore.ErrUserDefined,
			expected: status.Error(codes.FailedPrecondition, "user defined error"),
		},
		{
			name:     "user defined error with message",
			inputErr: fmt.Errorf("%w: %s", datastore.ErrUserDefined, "test"),
			expected: status.Error(codes.FailedPrecondition, "user defined error: test"),
		},
		{
			name:     "internal error",
			inputErr: errors.New("internal error"),
			inputMsg: "test",
			expected: status.Error(codes.Internal, "Failed to test"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := gRPCStoreError(tt.inputErr, tt.inputMsg)
			assert.Equal(t, tt.expected, err)
		})
	}
}
