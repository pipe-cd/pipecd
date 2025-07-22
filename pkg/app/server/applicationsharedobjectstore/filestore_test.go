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

package applicationsharedobjectstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/filestore/filestoretest"
)

func TestBuildPath(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name       string
		appID      string
		pluginName string
		key        string
		want       string
	}{
		{
			name:       "simple path",
			appID:      "app-1",
			pluginName: "plugin-1",
			key:        "key-1",
			want:       "application-shared-objects/app-1/plugin-1/key-1.json",
		},
		{
			name:       "path with special characters",
			appID:      "app/1",
			pluginName: "plugin.1",
			key:        "key_1",
			want:       "application-shared-objects/app/1/plugin.1/key_1.json",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildPath(tc.appID, tc.pluginName, tc.key)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storeMock := filestoretest.NewMockStore(ctrl)

	testcases := []struct {
		name          string
		applicationID string
		pluginName    string
		key           string
		content       []byte
		readerErr     error

		expectedObject []byte
		expectedErr    error
	}{
		{
			name:          "success",
			applicationID: "app-1",
			pluginName:    "plugin-1",
			key:           "key-1",
			content: []byte(`{
				"key1": "value1",
			}`),
			readerErr: nil,

			expectedObject: []byte(`{
				"key1": "value1",
			}`),
			expectedErr: nil,
		},
		{
			name:          "file not found",
			applicationID: "app-1",
			pluginName:    "plugin-1",
			key:           "key-1",
			readerErr:     filestore.ErrNotFound,

			expectedObject: nil,
			expectedErr:    filestore.ErrNotFound,
		},
	}

	fs := &appObjectFileStore{
		backend: storeMock,
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path := buildPath(tc.applicationID, tc.pluginName, tc.key)
			storeMock.EXPECT().Get(t.Context(), path).Return(tc.content, tc.readerErr)

			obj, err := fs.Get(t.Context(), path)
			require.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedObject, obj)
		})
	}
}
