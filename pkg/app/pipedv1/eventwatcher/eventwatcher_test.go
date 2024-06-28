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

package eventwatcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStr(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name    string
		value   interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "string",
			value:   "value",
			want:    "value",
			wantErr: false,
		},
		{
			name:    "int",
			value:   1,
			want:    "1",
			wantErr: false,
		},
		{
			name:    "int64",
			value:   int64(1),
			want:    "1",
			wantErr: false,
		},
		{
			name:    "uint64",
			value:   uint64(1),
			want:    "1",
			wantErr: false,
		},
		{
			name:    "float64",
			value:   1.1,
			want:    "1.1",
			wantErr: false,
		},
		{
			name:    "bool",
			value:   true,
			want:    "true",
			wantErr: false,
		},
		{
			name:    "map",
			value:   make(map[string]interface{}),
			want:    "",
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convertStr(tc.value)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestModifyYAML(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name         string
		path         string
		field        string
		newValue     string
		wantNewYml   []byte
		wantUpToDate bool
		wantErr      bool
	}{
		{
			name:         "different between defined one and given one",
			path:         "testdata/a.yaml",
			field:        "$.foo",
			newValue:     "2",
			wantNewYml:   []byte("foo: 2\n"),
			wantUpToDate: false,
			wantErr:      false,
		},
		{
			name:         "already up-to-date",
			path:         "testdata/a.yaml",
			field:        "$.foo",
			newValue:     "1",
			wantNewYml:   nil,
			wantUpToDate: true,
			wantErr:      false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotNewYml, gotUpToDate, err := modifyYAML(tc.path, tc.field, tc.newValue)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.wantNewYml, gotNewYml)
			assert.Equal(t, tc.wantUpToDate, gotUpToDate)
		})
	}
}

func TestModifyText(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name         string
		path         string
		regex        string
		newValue     string
		want         []byte
		wantUpToDate bool
		wantErr      bool
	}{
		{
			name:         "invalid regex given",
			path:         "testdata/with-template.yaml",
			regex:        "[",
			newValue:     "v0.2.0",
			want:         nil,
			wantUpToDate: false,
			wantErr:      true,
		},
		{
			name:         "no capturing group given",
			path:         "testdata/with-template.yaml",
			regex:        "image: gcr.io/pipecd/foo:v[0-9].[0-9].[0-9]",
			newValue:     "v0.2.0",
			want:         nil,
			wantUpToDate: false,
			wantErr:      true,
		},
		{
			name:         "invalid capturing group given",
			path:         "testdata/with-template.yaml",
			regex:        "image: gcr.io/pipecd/foo:([)",
			newValue:     "v0.2.0",
			want:         nil,
			wantUpToDate: false,
			wantErr:      true,
		},
		{
			name:         "the file doesn't match regex",
			path:         "testdata/with-template.yaml",
			regex:        "abcdefg",
			newValue:     "v0.1.0",
			want:         nil,
			wantUpToDate: false,
			wantErr:      true,
		},
		{
			name:         "the file is up-to-date",
			path:         "testdata/with-template.yaml",
			regex:        "image: gcr.io/pipecd/foo:(v[0-9].[0-9].[0-9])",
			newValue:     "v0.1.0",
			want:         nil,
			wantUpToDate: true,
			wantErr:      false,
		},
		{
			name:     "replace a part of text",
			path:     "testdata/with-template.yaml",
			regex:    "image: gcr.io/pipecd/foo:(v[0-9].[0-9].[0-9])",
			newValue: "v0.2.0",
			want: []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
spec:
  template:
    spec:
      containers:
      - name: foo
        image: gcr.io/pipecd/foo:v0.2.0
        ports:
        - containerPort: 9085
        env:
        - name: FOO
          value: {{ .encryptedSecrets.foo }}

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: bar
spec:
  template:
    spec:
      containers:
      - name: bar
        image: gcr.io/pipecd/bar:v0.1.0
        ports:
        - containerPort: 9085
        env:
        - name: BAR
          value: {{ .encryptedSecrets.bar }}
`),
			wantUpToDate: false,
			wantErr:      false,
		},
		{
			name:     "replace text",
			path:     "testdata/kustomization.yaml",
			regex:    "newTag: (v[0-9].[0-9].[0-9])",
			newValue: "v0.2.0",
			want: []byte(`images:
- name: gcr.io/pipecd/foo
  newTag: v0.2.0
`),
			wantUpToDate: false,
			wantErr:      false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, gotUpToDate, err := modifyText(tc.path, tc.regex, tc.newValue)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantUpToDate, gotUpToDate)
		})
	}
}

func TestGetBranchName(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name      string
		newBranch bool
		eventName string
		branch    string
		want      string
	}{
		{
			name:      "create new branch",
			newBranch: true,
			eventName: "event",
			branch:    "main",
		},
		{
			name:      "return existing branch",
			newBranch: false,
			eventName: "event",
			branch:    "main",
			want:      "main",
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := makeBranchName(tc.newBranch, tc.eventName, tc.branch)
			if tc.newBranch {
				assert.NotEqual(t, tc.branch, got)
			} else {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
