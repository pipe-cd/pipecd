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
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/model"
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := gRPCStoreError(tt.inputErr, tt.inputMsg)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestMakeGitPath(t *testing.T) {
	t.Parallel()

	piped := &model.Piped{
		Repositories: []*model.ApplicationGitRepository{
			{Id: "repo-1", Remote: "git@github.com:org/repo.git", Branch: "main"},
		},
	}

	tests := []struct {
		name         string
		repoID       string
		path         string
		cfgFilename  string
		expectedCode codes.Code
	}{
		{name: "valid", repoID: "repo-1", path: "apps/foo", cfgFilename: "app.pipecd.yaml", expectedCode: codes.OK},
		{name: "path escapes repository", repoID: "repo-1", path: "../../etc", expectedCode: codes.InvalidArgument},
		{name: "nested path escapes repository", repoID: "repo-1", path: "apps/../../etc", expectedCode: codes.InvalidArgument},
		{name: "config filename is parent dir", repoID: "repo-1", path: "apps/foo", cfgFilename: "..", expectedCode: codes.InvalidArgument},
		{name: "config filename is a path", repoID: "repo-1", path: "apps/foo", cfgFilename: "../app.pipecd.yaml", expectedCode: codes.InvalidArgument},
		{name: "unknown repository", repoID: "repo-2", path: "apps/foo", expectedCode: codes.NotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gitPath, err := makeGitPath(tt.repoID, tt.path, tt.cfgFilename, piped, zap.NewNop())
			assert.Equal(t, tt.expectedCode, status.Code(err))
			if tt.expectedCode == codes.OK {
				assert.Equal(t, tt.path, gitPath.Path)
				assert.Equal(t, tt.cfgFilename, gitPath.ConfigFilename)
			}
		})
	}
}
