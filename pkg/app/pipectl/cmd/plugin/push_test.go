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

package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/oci"
)

func TestPush_parseFilePaths(t *testing.T) {
	t.Parallel()

	p := &push{}
	tests := []struct {
		name    string
		input   []string
		want    map[oci.Platform]string
		wantErr bool
	}{
		{
			name:  "single valid entry",
			input: []string{"linux/amd64=/path/to/file1"},
			want: map[oci.Platform]string{
				{OS: "linux", Arch: "amd64"}: "/path/to/file1",
			},
			wantErr: false,
		},
		{
			name:  "multiple valid entries",
			input: []string{"linux/amd64=/path/to/file1", "linux/arm64=/path/to/file2"},
			want: map[oci.Platform]string{
				{OS: "linux", Arch: "amd64"}: "/path/to/file1",
				{OS: "linux", Arch: "arm64"}: "/path/to/file2",
			},
			wantErr: false,
		},
		{
			name:    "invalid format (missing =)",
			input:   []string{"linux/amd64:/path/to/file1"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid platform (missing arch)",
			input:   []string{"linux=/path/to/file1"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid platform (extra part)",
			input:   []string{"linux/amd64/extra=/path/to/file1"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := p.parseFilePaths(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
