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
					FieldPath:   "Id",
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
					FieldPath:   "Id",
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
					FieldPath:   "Name",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Id",
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
					FieldPath:   "Id",
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
					FieldPath:   "Id",
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
					FieldPath:   "Id",
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
					FieldPath:   "PipedId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Id",
					Order:       "ASCENDING",
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
					FieldPath:   "Id",
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
					FieldPath:   "ApplicationName",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Id",
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
					FieldPath:   "Id",
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
					FieldPath:   "Id",
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
					FieldPath:   "Id",
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
					FieldPath:   "DeploymentChainId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Id",
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
					FieldPath:   "CompletedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Id",
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
					FieldPath:   "CompletedAt",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Id",
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
					FieldPath:   "CreatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Id",
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
					FieldPath:   "ProjectId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "CreatedAt",
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
					FieldPath:   "Id",
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
					FieldPath:   "ProjectId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Status",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "CreatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Id",
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
					FieldPath:   "ProjectId",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
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
					FieldPath:   "Id",
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
					FieldPath:   "UpdatedAt",
					Order:       "DESCENDING",
					ArrayConfig: "",
				},
				{
					FieldPath:   "Id",
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
					FieldPath:   "Id",
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
					FieldPath:   "Id",
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
					FieldPath:   "Id",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "DeploymentChain",
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
					FieldPath:   "Id",
					Order:       "ASCENDING",
					ArrayConfig: "",
				},
			},
		},
		{
			CollectionGroup: "DeploymentChain",
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
					FieldPath:   "Id",
					Order:       "ASCENDING",
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
		{
			name: "no exclude a composite index in case the fields order is changed",
			indexes: []index{
				{
					CollectionGroup: "collection-group",
					QueryScope:      "COLLECTION",
					Fields: []field{
						{
							FieldPath: "field-path-1",
							Order:     "ASCENDING",
						},
						{
							FieldPath: "field-path-2",
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
							FieldPath: "field-path-2",
							Order:     "ASCENDING",
						},
						{
							FieldPath: "field-path-1",
							Order:     "ASCENDING",
						},
					},
				},
			},
			want: []index{
				{
					CollectionGroup: "collection-group",
					QueryScope:      "COLLECTION",
					Fields: []field{
						{
							FieldPath: "field-path-1",
							Order:     "ASCENDING",
						},
						{
							FieldPath: "field-path-2",
							Order:     "ASCENDING",
						},
					},
				},
			},
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
			name: "ensure the fields order is not changed",
			idx: index{
				CollectionGroup: "collection-group",
				QueryScope:      "COLLECTION",
				Fields: []field{
					{
						FieldPath:   "field-path2",
						ArrayConfig: "contains",
					},
					{
						FieldPath: "field-path3",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "field-path1",
						Order:     "ASCENDING",
					},
				},
			},
			want: "collection-group/COLLECTION/field-path:field-path2/array-config:contains/field-path:field-path3/order:ASCENDING/field-path:field-path1/order:ASCENDING",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.idx.id()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPrefixIndexes(t *testing.T) {
	testcases := []struct {
		name     string
		indexes  []index
		prefix   string
		expected []index
	}{
		{
			name:   "nil list",
			prefix: "TestPrefix",
		},
		{
			name:   "normal list",
			prefix: "TestPrefix",
			indexes: []index{
				{
					CollectionGroup: "CollectionA",
				},
				{
					CollectionGroup: "CollectionB",
				},
			},
			expected: []index{
				{
					CollectionGroup: "TestPrefixCollectionA",
				},
				{
					CollectionGroup: "TestPrefixCollectionB",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			indexes := tc.indexes
			prefixIndexes(indexes, tc.prefix)
			assert.Equal(t, tc.expected, indexes)
		})
	}
}
