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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeFileInfo struct {
	name string
}

func (f *fakeFileInfo) Name() string { return f.name }

// Below methods are required to meet the interface.
func (f *fakeFileInfo) Size() int64        { return 0 }
func (f *fakeFileInfo) Mode() os.FileMode  { return 0 }
func (f *fakeFileInfo) ModTime() time.Time { return time.Now() }
func (f *fakeFileInfo) IsDir() bool        { return false }
func (f *fakeFileInfo) Sys() interface{}   { return nil }

func TestFilterImageWatcherFiles(t *testing.T) {
	testcases := []struct {
		name     string
		files    []os.FileInfo
		includes []string
		excludes []string
		want     []os.FileInfo
		wantErr  bool
	}{
		{
			name: "both includes and excludes aren't given",
			files: []os.FileInfo{
				&fakeFileInfo{
					name: "file-1",
				},
			},
			want: []os.FileInfo{
				&fakeFileInfo{
					name: "file-1",
				},
			},
			wantErr: false,
		},
		{
			name: "both includes and excludes aren given",
			files: []os.FileInfo{
				&fakeFileInfo{
					name: "file-1",
				},
			},
			includes: []string{"file-1"},
			excludes: []string{"file-1"},
			wantErr:  true,
		},
		{
			name: "includes given",
			files: []os.FileInfo{
				&fakeFileInfo{
					name: "file-1",
				},
				&fakeFileInfo{
					name: "file-2",
				},
				&fakeFileInfo{
					name: "file-3",
				},
			},
			includes: []string{"file-1", "file-3"},
			want: []os.FileInfo{
				&fakeFileInfo{
					name: "file-1",
				},
				&fakeFileInfo{
					name: "file-3",
				},
			},
			wantErr: false,
		},
		{
			name: "excludes given",
			files: []os.FileInfo{
				&fakeFileInfo{
					name: "file-1",
				},
				&fakeFileInfo{
					name: "file-2",
				},
				&fakeFileInfo{
					name: "file-3",
				},
			},
			excludes: []string{"file-1", "file-3"},
			want: []os.FileInfo{
				&fakeFileInfo{
					name: "file-2",
				},
			},
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := filterImageWatcherFiles(tc.files, tc.includes, tc.excludes)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
