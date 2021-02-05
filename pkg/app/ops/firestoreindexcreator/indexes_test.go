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
					field{
						FieldPath: "Deleted",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "CreatedAt",
						Order:     "ASCENDING",
					},
				},
			},
			{
				CollectionGroup: "Application",
				Fields: []field{
					field{
						FieldPath: "Disabled",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Application",
				Fields: []field{
					field{
						FieldPath: "EnvId",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Application",
				Fields: []field{
					field{
						FieldPath: "Kind",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Application",
				Fields: []field{
					field{
						FieldPath: "ProjectId",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Application",
				Fields: []field{
					field{
						FieldPath: "SyncState.Status",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Command",
				Fields: []field{
					field{
						FieldPath: "Status",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "CreatedAt",
						Order:     "ASCENDING",
					},
				},
			},
			{
				CollectionGroup: "Deployment",
				Fields: []field{
					field{
						FieldPath: "ApplicationId",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Deployment",
				Fields: []field{
					field{
						FieldPath: "EnvId",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Deployment",
				Fields: []field{
					field{
						FieldPath: "Kind",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Deployment",
				Fields: []field{
					field{
						FieldPath: "ProjectId",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Deployment",
				Fields: []field{
					field{
						FieldPath: "Status",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "UpdatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Event",
				Fields: []field{
					field{
						FieldPath: "EventKey",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "Name",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "ProjectId",
						Order:     "ASCENDING",
					},
					field{
						FieldPath: "CreatedAt",
						Order:     "DESCENDING",
					},
				},
			},
			{
				CollectionGroup: "Event",
				Fields: []field{
					field{
						FieldPath: "ProjectId",
						Order:     "ASCENDING",
					},
					field{
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
