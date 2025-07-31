// Copyright 2025 The PipeCD Authors.
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

package livestate

import (
	"fmt"
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/provider"
)

func TestMakeSyncState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		result  provider.PlanResult
		commit  string
		want    sdk.ApplicationSyncState
		wantErr bool
	}{
		{
			name: "no changes",
			result: provider.PlanResult{
				Imports:  0,
				Adds:     0,
				Changes:  0,
				Destroys: 0,
			},
			commit: "commit1",
			want: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateSynced,
				ShortReason: "",
				Reason:      "",
			},
			wantErr: false,
		},
		{
			name: "has changes",
			result: provider.PlanResult{
				Imports:    1,
				Adds:       2,
				Changes:    3,
				Destroys:   4,
				PlanOutput: "Terraform will perform the following actions:\nsome changes\nPlan: 1 to import, 2 to add, 3 to change, 4 to destroy.",
			},
			commit: "1234567890abcdef",
			want: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateOutOfSync,
				ShortReason: "There are 10 manifests that are not synced (1 imports, 2 adds, 4 deletes, 3 changes)",
				Reason: fmt.Sprintf("Diff between the defined state in Git at commit 1234567 and actual live state:\n\n" +
					"--- Actual   (LiveState)\n" +
					"+++ Expected (Git)\n\n" +
					"some changes\nPlan: 1 to import, 2 to add, 3 to change, 4 to destroy.\n"),
			},
			wantErr: false,
		},
		{
			name: "invalid plan output",
			result: provider.PlanResult{
				Imports:    1,
				Adds:       1,
				Changes:    1,
				Destroys:   1,
				PlanOutput: "<invalid plan output> Terraform will perform the following actions:",
			},
			commit:  "1234567890abcdef",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := makeSyncState(tt.result, tt.commit)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
