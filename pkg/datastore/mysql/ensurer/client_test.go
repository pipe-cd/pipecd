// Copyright 2023 The PipeCD Authors.
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

package ensurer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeCreateIndexStatements(t *testing.T) {
	testcases := []struct {
		name               string
		rawIndexes         string
		expectedStatements []string
	}{
		{
			name:       "Only one CREATE INDEX statement",
			rawIndexes: "CREATE INDEX application_created_at_desc ON Application (created_at DESC);\n",
			expectedStatements: []string{
				"CREATE INDEX application_created_at_desc ON Application (created_at DESC)",
			},
		},
		{
			name:       "Only one CREATE INDEX statement without `;`",
			rawIndexes: "CREATE INDEX application_created_at_desc ON Application (created_at DESC)\n",
			expectedStatements: []string{
				"CREATE INDEX application_created_at_desc ON Application (created_at DESC)",
			},
		},
		{
			name:       "Multi CREATE INDEX statements",
			rawIndexes: "CREATE INDEX application_updated_at_asc ON Application (updated_at ASC);\n\nCREATE INDEX application_created_at_desc ON Application (created_at DESC);\n",
			expectedStatements: []string{
				"CREATE INDEX application_updated_at_asc ON Application (updated_at ASC)",
				"CREATE INDEX application_created_at_desc ON Application (created_at DESC)",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			statements := makeCreateIndexStatements(tc.rawIndexes)
			assert.Equal(t, tc.expectedStatements, statements)
		})
	}
}
