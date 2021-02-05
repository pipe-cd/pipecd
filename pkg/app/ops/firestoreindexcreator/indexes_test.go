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

package firestoreindexcreator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseIndexes(t *testing.T) {
	want := &indexes{
		Indexes: []index{
			{
				CollectionGroup: "Application",
				Fields: []field{
					{
						FieldPath: "Deleted",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "CreatedAt",
						Order:     "ASCENDING",
					},
				},
			},
			{
				CollectionGroup: "Application",
				Fields: []field{
					{
						FieldPath: "Disabled",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Application",
				Fields: []field{
					{
						FieldPath: "EnvId",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Application",
				Fields: []field{
					{
						FieldPath: "Kind",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Application",
				Fields: []field{
					{
						FieldPath: "ProjectId",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Application",
				Fields: []field{
					{
						FieldPath: "SyncState.Status",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Command",
				Fields: []field{
					{
						FieldPath: "Status",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "CreatedAt",
						Order:     "ASCENDING",
					},
				},
			},
			{
				CollectionGroup: "Deployment",
				Fields: []field{
					{
						FieldPath: "ApplicationId",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Deployment",
				Fields: []field{
					{
						FieldPath: "EnvId",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Deployment",
				Fields: []field{
					{
						FieldPath: "Kind",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Deployment",
				Fields: []field{
					{
						FieldPath: "ProjectId",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Deployment",
				Fields: []field{
					{
						FieldPath: "Status",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Event",
				Fields: []field{
					{
						FieldPath: "EventKey",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "Name",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "ProjectId",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "CreatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Event",
				Fields: []field{
					{
						FieldPath: "ProjectId",
						Order:     "ASCENDING",
					},
					{
						FieldPath: "CreatedAt",
						Order:     "ASCENDING",
					},
				},
			},
		},
	}

	got, err := parseIndexes()
	assert.Equal(t, want, got)
	require.NoError(t, err)
}
