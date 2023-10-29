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

package config

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestControlPlaneConfig(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       *ControlPlaneSpec
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
				SharedSSOConfigs: []SharedSSOConfig{
					{
						Name: "github",
						ProjectSSOConfig: model.ProjectSSOConfig{
							Provider: model.ProjectSSOConfig_GITHUB,
							Github: &model.ProjectSSOConfig_GitHub{
								ClientId:     "client-id",
								ClientSecret: "client-secret",
								BaseUrl:      "base-url",
								UploadUrl:    "upload-url",
							},
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
				InsightCollector: ControlPlaneInsightCollector{
					Application: InsightCollectorApplication{
						Enabled:  newBoolPointer(true),
						Schedule: "0 * * * *",
					},
					Deployment: InsightCollectorDeployment{
						Enabled:       newBoolPointer(true),
						Schedule:      "0 10 * * *",
						ChunkMaxCount: 1000,
					},
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
				require.Equal(t, 1, len(tc.expectedSpec.SharedSSOConfigs))
				require.Equal(t, 1, len(cfg.ControlPlaneSpec.SharedSSOConfigs))
				// Why don't we use assert.Equal to compare?
				// https://github.com/stretchr/testify/issues/758
				assert.True(t, proto.Equal(&tc.expectedSpec.SharedSSOConfigs[0].ProjectSSOConfig, &cfg.ControlPlaneSpec.SharedSSOConfigs[0].ProjectSSOConfig))

				tc.expectedSpec.SharedSSOConfigs = nil
				cfg.ControlPlaneSpec.SharedSSOConfigs = nil
				assert.Equal(t, tc.expectedSpec, cfg.ControlPlaneSpec)
			}
		})
	}
}
