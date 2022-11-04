// Copyright 2022 The PipeCD Authors.
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

package insightstore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/filestore/filestoretest"
	"github.com/pipe-cd/pipecd/pkg/insight"
)

func TestLoadMilestone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		content string

		expected    *insight.Milestone
		readerErr   error
		expectedErr error
	}{
		{
			name:        "file not found in filestore",
			content:     "",
			readerErr:   filestore.ErrNotFound,
			expectedErr: filestore.ErrNotFound,
		},
		{
			name: "file found in filestore",
			content: `{
				"deployment_created_at_milestone": 1234,
				"deployment_completed_at_milestone": 1234
			}`,
			expected: &insight.Milestone{
				DeploymentCreatedAtMilestone:   1234,
				DeploymentCompletedAtMilestone: 1234,
			},
			readerErr:   nil,
			expectedErr: nil,
		},
	}

	fs := filestoretest.NewMockStore(ctrl)
	s := &store{filestore: fs}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			obj := []byte(tc.content)
			fs.EXPECT().Get(context.TODO(), milestonePath).Return(obj, tc.readerErr)
			state, err := s.LoadMilestone(context.TODO())
			if err != nil {
				if tc.expectedErr == nil {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tc.expectedErr)
				return
			}
			assert.Equal(t, tc.expected, state)
		})
	}
}
