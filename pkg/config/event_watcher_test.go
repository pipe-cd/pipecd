// Copyright 2021 The PipeCD Authors.
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

func TestLoadEventWatcher(t *testing.T) {
	want := &EventWatcherSpec{Events: []EventWatcherEvent{
		{
			Name: "app1-image-update",
			Replacements: []EventWatcherReplacement{
				{
					File:      "app1/deployment.yaml",
					YAMLField: "$.spec.template.spec.containers[0].image",
				},
			},
		},
		{
			Name: "app2-helm-release",
			Labels: map[string]string{
				"repoId": "repo-1",
			},
			Replacements: []EventWatcherReplacement{
				{
					File:      "app2/.pipe.yaml",
					YAMLField: "$.spec.input.helmChart.version",
				},
			},
		},
	}}

	t.Run("valid config files given", func(t *testing.T) {
		got, err := LoadEventWatcher("testdata", nil, []string{"README.md"})
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestEventWatcherValidate(t *testing.T) {
	testcases := []struct {
		name             string
		eventWatcherSpec EventWatcherSpec
		wantErr          bool
	}{
		{
			name: "no name given",
			eventWatcherSpec: EventWatcherSpec{
				Events: []EventWatcherEvent{
					{
						Replacements: []EventWatcherReplacement{
							{
								File:      "file",
								YAMLField: "$.foo",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no replacements given",
			eventWatcherSpec: EventWatcherSpec{
				Events: []EventWatcherEvent{
					{
						Name: "event-a",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no replacement file given",
			eventWatcherSpec: EventWatcherSpec{
				Events: []EventWatcherEvent{
					{
						Replacements: []EventWatcherReplacement{
							{
								YAMLField: "$.foo",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no replacement field given",
			eventWatcherSpec: EventWatcherSpec{
				Events: []EventWatcherEvent{
					{
						Replacements: []EventWatcherReplacement{
							{
								File: "file",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "both yaml and json given",
			eventWatcherSpec: EventWatcherSpec{
				Events: []EventWatcherEvent{
					{
						Replacements: []EventWatcherReplacement{
							{
								File:      "file",
								YAMLField: "$.foo",
								JSONField: "$.foo",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "both yaml and hcl given",
			eventWatcherSpec: EventWatcherSpec{
				Events: []EventWatcherEvent{
					{
						Replacements: []EventWatcherReplacement{
							{
								File:      "file",
								YAMLField: "$.foo",
								HCLField:  "$.foo",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "both json and hcl given",
			eventWatcherSpec: EventWatcherSpec{
				Events: []EventWatcherEvent{
					{
						Replacements: []EventWatcherReplacement{
							{
								File:      "file",
								JSONField: "$.foo",
								HCLField:  "$.foo",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid config given",
			eventWatcherSpec: EventWatcherSpec{
				Events: []EventWatcherEvent{
					{
						Name: "event-a",
						Replacements: []EventWatcherReplacement{
							{
								File:      "file",
								YAMLField: "$.foo",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.eventWatcherSpec.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

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

func TestFilterEventWatcherFiles(t *testing.T) {
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
			got, err := filterEventWatcherFiles(tc.files, tc.includes, tc.excludes)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
