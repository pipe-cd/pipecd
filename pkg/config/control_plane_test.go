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

package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipe/pkg/model"
)

func TestControlPlaneConfig(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/control-plane/control-plane-config.yaml",
			expectedKind:       KindControlPlane,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &ControlPlaneSpec{
				Projects: []ControlPlaneProject{
					{
						ID: "abc",
						StaticAdmin: ProjectStaticUser{
							Username:     "test-user",
							PasswordHash: "test-password",
						},
					},
				},
				SharedSSO: []SharedSingleSignOn{
					{
						Name:     "default",
						Provider: Github,
						Github: SharedSingleSignOnGitHub{
							ClientID:     "client-id",
							ClientSecret: "client-secret",
							BaseUrl:      "base-url",
							UploadUrl:    "upload-url",
						},
					},
				},
				Datastore: ControlPlaneDataStore{
					Type: model.DataStoreFirestore,
					FirestoreConfig: &DataStoreFireStoreConfig{
						Namespace:       "pipecd-test",
						Environment:     "unit-test",
						Project:         "project",
						CredentialsFile: "datastore-credentials-file.json",
					},
				},
				Filestore: ControlPlaneFileStore{
					Type: model.FileStoreGCS,
					GCSConfig: &FileStoreGCSConfig{
						Bucket:          "bucket",
						CredentialsFile: "filestore-credentials-file.json",
					},
				},
				Cache: ControlPlaneCache{
					TTL: Duration(5 * time.Minute),
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.fileName, func(t *testing.T) {
			cfg, err := LoadFromYAML(tc.fileName)
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedKind, cfg.Kind)
				assert.Equal(t, tc.expectedAPIVersion, cfg.APIVersion)
				assert.Equal(t, tc.expectedSpec, cfg.spec)
			}
		})
	}
}
