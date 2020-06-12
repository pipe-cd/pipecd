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

package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiffResultListFind(t *testing.T) {
	testcases := []struct {
		name        string
		list        DiffResultList
		query       string
		expected    *DiffResult
		expectedErr error
	}{
		{
			name: "not found",
			list: DiffResultList{
				{
					PathString: "spec.template.spec.containers.[0].image",
					Before:     "gcr.io/pipecd/helloworld:v1.0.0",
					After:      "gcr.io/pipecd/helloworld:v2.0.0",
				},
			},
			query: `spec.template2.spec.containers.\[\d+\].image`,
		},
		{
			name: "found one",
			list: DiffResultList{
				{
					PathString: "spec.template.spec.containers.[0].image",
					Before:     "gcr.io/pipecd/helloworld:v1.0.0",
					After:      "gcr.io/pipecd/helloworld:v2.0.0",
				},
			},
			query: `spec.template.spec.containers.\[\d+\].image`,
			expected: &DiffResult{
				PathString: "spec.template.spec.containers.[0].image",
				Before:     "gcr.io/pipecd/helloworld:v1.0.0",
				After:      "gcr.io/pipecd/helloworld:v2.0.0",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for i := range tc.list {
				path, err := parseDiffPath(tc.list[i].PathString)
				require.NoError(t, err)
				tc.list[i].Path = path
			}

			dr, ok, err := tc.list.Find(tc.query)

			assert.Equal(t, tc.expectedErr, err)
			if tc.expected == nil {
				assert.Equal(t, false, ok)
				assert.Equal(t, DiffResult{}, dr)
			} else {
				assert.Equal(t, true, ok)
				assert.Equal(t, tc.expected.PathString, dr.PathString)
				assert.Equal(t, tc.expected.Before, dr.Before)
				assert.Equal(t, tc.expected.After, dr.After)
			}
		})
	}
}

func TestDiffResultListFindAll(t *testing.T) {
	testcases := []struct {
		name     string
		list     DiffResultList
		query    string
		expected []DiffResult
	}{
		{
			name: "not found",
			list: DiffResultList{
				{
					PathString: "spec.template.spec.containers.[0].image",
					Before:     "gcr.io/pipecd/helloworld:v1.0.0",
					After:      "gcr.io/pipecd/helloworld:v2.0.0",
				},
			},
			query: `spec.template2.spec.containers.\[\d+\].image`,
		},
		{
			name: "found two objects",
			list: DiffResultList{
				{
					PathString: "spec.template.spec.containers.[0].image",
					Before:     "gcr.io/pipecd/helloworld:v1.0.0",
					After:      "gcr.io/pipecd/helloworld:v2.0.0",
				},
				{
					PathString: "spec.template.spec.containers.[1].image",
					Before:     "envoy:v1.0.0",
					After:      "envoy:v2.0.0",
				},
				{
					PathString: "spec.template.spec.containers.[1].name",
					Before:     "foo",
					After:      "bar",
				},
			},
			query: `^spec.template.spec.containers.\[\d+\].image$`,
			expected: []DiffResult{
				{
					PathString: "spec.template.spec.containers.[0].image",
					Before:     "gcr.io/pipecd/helloworld:v1.0.0",
					After:      "gcr.io/pipecd/helloworld:v2.0.0",
				},
				{
					PathString: "spec.template.spec.containers.[1].image",
					Before:     "envoy:v1.0.0",
					After:      "envoy:v2.0.0",
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for i := range tc.list {
				path, err := parseDiffPath(tc.list[i].PathString)
				require.NoError(t, err)
				tc.list[i].Path = path
			}

			list := tc.list.FindAll(tc.query)
			for i := range list {
				list[i].Path = nil
			}

			assert.Equal(t, tc.expected, list)
		})
	}
}

func TestDiff(t *testing.T) {
	testcases := []struct {
		name     string
		yamlFile string
		options  []DiffOption
		result   DiffResultList
	}{
		{
			name:     "no diff",
			yamlFile: "testdata/diff_no_diff.yaml",
		},
		{
			name:     "no diff with ignore order",
			yamlFile: "testdata/diff_ignore_order_no_diff.yaml",
			options: []DiffOption{
				WithDiffIgnoreOrder(),
			},
		},
		{
			name:     "has some diffs",
			yamlFile: "testdata/diff_multi_diffs.yaml",
			result: []DiffResult{
				{
					Path: []PathStep{
						{
							Type: MapKeyPathStep,
							Key:  "metadata",
						},
						{
							Type: MapKeyPathStep,
							Key:  "labels",
						},
						{
							Type: MapKeyPathStep,
							Key:  "change",
						},
					},
					PathString: "metadata.labels.change",
					Before:     "first",
					After:      "second",
				},
				{
					Path: []PathStep{
						{
							Type: MapKeyPathStep,
							Key:  "spec",
						},
						{
							Type: MapKeyPathStep,
							Key:  "template",
						},
						{
							Type: MapKeyPathStep,
							Key:  "spec",
						},
						{
							Type: MapKeyPathStep,
							Key:  "containers",
						},
						{
							Type:  SliceIndexPathStep,
							Index: 0,
						},
						{
							Type: MapKeyPathStep,
							Key:  "image",
						},
					},
					PathString: "spec.template.spec.containers.[0].image",
					Before:     "gcr.io/pipecd/helloworld:v1.0.0",
					After:      "gcr.io/pipecd/helloworld:v2.0.0",
				},
			},
		},
		{
			name:     "one filtered by path prefix",
			yamlFile: "testdata/diff_multi_diffs.yaml",
			options: []DiffOption{
				WithPathPrefix("spec.template"),
			},
			result: []DiffResult{
				{
					Path: []PathStep{
						{
							Type: MapKeyPathStep,
							Key:  "spec",
						},
						{
							Type: MapKeyPathStep,
							Key:  "template",
						},
						{
							Type: MapKeyPathStep,
							Key:  "spec",
						},
						{
							Type: MapKeyPathStep,
							Key:  "containers",
						},
						{
							Type:  SliceIndexPathStep,
							Index: 0,
						},
						{
							Type: MapKeyPathStep,
							Key:  "image",
						},
					},
					PathString: "spec.template.spec.containers.[0].image",
					Before:     "gcr.io/pipecd/helloworld:v1.0.0",
					After:      "gcr.io/pipecd/helloworld:v2.0.0",
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := LoadManifestsFromYAMLFile(tc.yamlFile)
			require.NoError(t, err)
			require.Equal(t, 2, len(manifests))

			dl := Diff(manifests[0], manifests[1], tc.options...)
			assert.Equal(t, tc.result, dl)
		})
	}
}
