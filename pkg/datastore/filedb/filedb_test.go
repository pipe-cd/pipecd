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

package filedb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestDataTo(t *testing.T) {
	testcases := []struct {
		name      string
		src       interface{}
		dst       interface{}
		expectErr bool
	}{
		{
			name:      "returns error on type miss match",
			src:       &model.Command{},
			dst:       &model.Event{},
			expectErr: true,
		},
		{
			name:      "map data successfully",
			src:       &model.Command{Id: "command-id"},
			dst:       &model.Command{},
			expectErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := dataTo(tc.src, tc.dst)
			require.Equal(t, tc.expectErr, err != nil)

			if err == nil {
				assert.Equal(t, tc.src, tc.dst)
			}
		})
	}
}
