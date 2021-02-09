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

package firestoreindexensurer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseIndexes(t *testing.T) {
	want := []index{
		{
			CollectionGroup: "Application",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "Disabled",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Application",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "EnvId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Application",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "Deleted",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "CreatedAt",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Application",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "Kind",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Application",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "SyncState.Status",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Application",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "ProjectId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Command",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "Status",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "CreatedAt",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Deployment",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "ApplicationId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Deployment",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "ProjectId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Deployment",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "EnvId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Deployment",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "Kind",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Deployment",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "Status",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Event",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "ProjectId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "CreatedAt",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "Event",
			QueryScope:      "COLLECTION",
			Fields: []field{
				{
					FieldPath:   "EventKey",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Name",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "ProjectId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "CreatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "__name__",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
			},
		},
	}

	got, err := parseIndexes()
	assert.Equal(t, want, got)
	require.NoError(t, err)
}

func TestFilterIndexes(t *testing.T) {
	testcases := []struct {
		name     string
		indexes  []index
		excludes []index
		want     []index
	}{
		{
			name: "no excludes given",
			indexes: []index{
				{
					CollectionGroup: "collection-group",
					QueryScope:      "COLLECTION",
					Fields: []field{
						{
							FieldPath: "field-path",
							Order:     "ASCENDING",
						},
					},
				},
			},
			excludes: []index{},
			want: []index{
				{
					CollectionGroup: "collection-group",
					QueryScope:      "COLLECTION",
					Fields: []field{
						{
							FieldPath: "field-path",
							Order:     "ASCENDING",
						},
					},
				},
			},
		},
		{
			name: "exclude an index",
			indexes: []index{
				{
					CollectionGroup: "collection-group",
					QueryScope:      "COLLECTION",
					Fields: []field{
						{
							FieldPath: "field-path",
							Order:     "ASCENDING",
						},
					},
				},
			},
			excludes: []index{
				{
					CollectionGroup: "collection-group",
					QueryScope:      "COLLECTION",
					Fields: []field{
						{
							FieldPath: "field-path",
							Order:     "ASCENDING",
						},
					},
				},
			},
			want: []index{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := filterIndexes(tc.indexes, tc.excludes)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIndexID(t *testing.T) {
	testcases := []struct {
		name string
		idx  index
		want string
	}{
		{
			name: "single field",
			idx: index{
				CollectionGroup: "collection-group",
				QueryScope:      "COLLECTION",
				Fields: []field{
					{
						FieldPath: "field-path",
						Order:     "ASCENDING",
					},
				},
			},
			want: "collection-group/COLLECTION/field-path:field-path/order:ASCENDING",
		},
		{
			name: "two fields",
			idx: index{
				CollectionGroup: "collection-group",
				QueryScope:      "COLLECTION",
				Fields: []field{
					{
						FieldPath: "field-path1",
						Order:     "ASCENDING",
					},
					{
						FieldPath:   "field-path2",
						ArrayConfig: "contains",
					},
				},
			},
			want: "collection-group/COLLECTION/field-path:field-path1/order:ASCENDING/field-path:field-path2/array-config:contains",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.idx.id()
			assert.Equal(t, tc.want, got)
		})
	}
}
