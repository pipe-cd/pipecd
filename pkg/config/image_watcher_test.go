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

	"github.com/stretchr/testify/assert"
)

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
			name: "both includes and excludes are given",
			files: []os.FileInfo{
				&fakeFileInfo{
					name: "file-1",
				},
			},
			want:     []os.FileInfo{},
			includes: []string{"file-1"},
			excludes: []string{"file-1"},
			wantErr:  false,
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

func TestImageWatcherValidate(t *testing.T) {
	testcases := []struct {
		name             string
		imageWatcherSpec ImageWatcherSpec
		wantErr          bool
	}{
		{
			name: "empty provider given",
			imageWatcherSpec: ImageWatcherSpec{Targets: []ImageWatcherTarget{
				{
					Image:    "image",
					FilePath: "filePath",
					Field:    "field",
				},
			}},
			wantErr: true,
		},
		{
			name: "empty image given",
			imageWatcherSpec: ImageWatcherSpec{Targets: []ImageWatcherTarget{
				{
					Provider: "provider",
					FilePath: "filePath",
					Field:    "field",
				},
			}},
			wantErr: true,
		},
		{
			name: "empty file path given",
			imageWatcherSpec: ImageWatcherSpec{Targets: []ImageWatcherTarget{
				{
					Provider: "provider",
					Image:    "image",
					Field:    "field",
				},
			}},
			wantErr: true,
		},
		{
			name: "empty field given",
			imageWatcherSpec: ImageWatcherSpec{Targets: []ImageWatcherTarget{
				{
					Provider: "provider",
					Image:    "image",
					FilePath: "filePath",
				},
			}},
			wantErr: true,
		},
		{
			name: "all fields given",
			imageWatcherSpec: ImageWatcherSpec{Targets: []ImageWatcherTarget{
				{
					Provider: "provider",
					Image:    "image",
					FilePath: "filePath",
					Field:    "field",
				},
			}},
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.imageWatcherSpec.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
