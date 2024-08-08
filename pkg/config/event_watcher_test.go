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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterEventWatcherFiles(t *testing.T) {
	testcases := []struct {
		name     string
		files    []string
		includes []string
		excludes []string
		want     []string
		wantErr  bool
	}{
		{
			name:    "both includes and excludes aren't given",
			files:   []string{"file-1"},
			want:    []string{"file-1"},
			wantErr: false,
		},
		{
			name:     "both includes and excludes are given",
			files:    []string{"file-1"},
			want:     []string{},
			includes: []string{"file-1"},
			excludes: []string{"file-1"},
			wantErr:  false,
		},
		{
			name:     "includes given",
			files:    []string{"file-1", "file-2", "file-3"},
			includes: []string{"file-1", "file-3"},
			want:     []string{"file-1", "file-3"},
			wantErr:  false,
		},
		{
			name:     "excludes given",
			files:    []string{"file-1", "file-2", "file-3"},
			excludes: []string{"file-1", "file-3"},
			want:     []string{"file-2"},
			wantErr:  false,
		},
		{
			name:     "includes with pattern given",
			files:    []string{"dir/file-1.yaml", "dir/file-2.yaml", "dir/file-3.yaml"},
			includes: []string{"dir/*.yaml"},
			want:     []string{"dir/file-1.yaml", "dir/file-2.yaml", "dir/file-3.yaml"},
			wantErr:  false,
		},
		{
			name:     "excludes with pattern given",
			files:    []string{"dir/file-1.yaml", "dir/file-2.yaml", "dir/file-3.yaml", "dir-2/file-1.yaml"},
			excludes: []string{"dir/*.yaml"},
			want:     []string{"dir-2/file-1.yaml"},
			wantErr:  false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := filterEventWatcherFiles(tc.files, tc.includes, tc.excludes)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
