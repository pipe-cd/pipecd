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

package scriptrun

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ContextInfo_BuildEnv(t *testing.T) {
	tests := []struct {
		name    string
		ci      *ContextInfo
		want    map[string]string
		wantErr bool
	}{
		{
			name: "success",
			ci: &ContextInfo{
				DeploymentID:        "deployment-id",
				ApplicationID:       "application-id",
				ApplicationName:     "application-name",
				TriggeredAt:         1234567890,
				TriggeredCommitHash: "commit-hash",
				RepositoryURL:       "repo-url",
				Labels: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				Summary: "summary",
			},
			want: map[string]string{
				"SR_DEPLOYMENT_ID":         "deployment-id",
				"SR_APPLICATION_ID":        "application-id",
				"SR_APPLICATION_NAME":      "application-name",
				"SR_TRIGGERED_AT":          "1234567890",
				"SR_TRIGGERED_COMMIT_HASH": "commit-hash",
				"SR_REPOSITORY_URL":        "repo-url",
				"SR_SUMMARY":               "summary",
				"SR_LABELS_KEY1":           "value1",
				"SR_LABELS_KEY2":           "value2",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ci.BuildEnv()
			assert.Equal(t, tt.wantErr, err != nil)

			for k, v := range got {
				if k == "SR_CONTEXT_RAW" {
					continue
				}
				assert.Equal(t, tt.want[k], v)
			}

			var gotRaw ContextInfo
			err = json.Unmarshal([]byte(got["SR_CONTEXT_RAW"]), &gotRaw)
			assert.Nil(t, err)
			assert.Equal(t, tt.ci, &gotRaw)
		})
	}
}
