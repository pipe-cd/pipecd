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

package filedb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

func TestNext(t *testing.T) {
	type fakeModel struct {
		Name string
	}
	testcases := []struct {
		name          string
		iter          Iterator
		expectDst     interface{}
		expectCurrent int
		expectErr     error
	}{
		{
			name: "pass in normal case",
			iter: Iterator{
				data: []interface{}{
					&fakeModel{Name: "1"},
					&fakeModel{Name: "2"},
				},
				current: 0,
			},
			expectDst:     &fakeModel{Name: "1"},
			expectCurrent: 1,
			expectErr:     nil,
		},
		{
			name: "iterator done",
			iter: Iterator{
				data: []interface{}{
					&fakeModel{Name: "1"},
					&fakeModel{Name: "2"},
				},
				current: 2,
			},
			expectDst:     fakeModel{},
			expectCurrent: 2,
			expectErr:     datastore.ErrIteratorDone,
		},
		{
			name: "data type miss match",
			iter: Iterator{
				data: []interface{}{
					"1",
				},
			},
			expectErr: fmt.Errorf("data type miss match"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			dst := &fakeModel{}
			err := tc.iter.Next(dst)

			require.Equal(t, tc.expectErr, err)
			if err == nil {
				assert.Equal(t, tc.expectCurrent, tc.iter.current)
				assert.Equal(t, tc.expectDst, dst)
			}
		})
	}
}
