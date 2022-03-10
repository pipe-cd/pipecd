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

package insight

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestMakeApplicationCounts(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testcases := []struct {
		name     string
		apps     []*model.Application
		expected ApplicationCounts
	}{
		{
			name: "empty",
			apps: nil,
			expected: ApplicationCounts{
				UpdatedAt: now.Unix(),
			},
		},
		{
			name: "multiple applications",
			apps: []*model.Application{
				{
					Id:       "1",
					Kind:     model.ApplicationKind_KUBERNETES,
					Disabled: false,
				},
				{
					Id:       "2",
					Kind:     model.ApplicationKind_KUBERNETES,
					Disabled: false,
				},
				{
					Id:       "3",
					Kind:     model.ApplicationKind_KUBERNETES,
					Disabled: true,
				},
				{
					Id:       "4",
					Kind:     model.ApplicationKind_CLOUDRUN,
					Disabled: true,
				},
			},
			expected: ApplicationCounts{
				Counts: []model.InsightApplicationCount{
					{
						Labels: map[string]string{
							"KIND":          "KUBERNETES",
							"ACTIVE_STATUS": "ENABLED",
						},
						Count: 2,
					},
					{
						Labels: map[string]string{
							"KIND":          "KUBERNETES",
							"ACTIVE_STATUS": "DISABLED",
						},
						Count: 1,
					},
					{
						Labels: map[string]string{
							"KIND":          "CLOUDRUN",
							"ACTIVE_STATUS": "DISABLED",
						},
						Count: 1,
					},
				},
				UpdatedAt: now.Unix(),
			},
		},
	}

	for _, tc := range testcases {
		c := MakeApplicationCounts(tc.apps, now)
		// We can use fmt to sort by Labels because maps are printed in key-sorted order.
		sort.Slice(c.Counts, func(i, j int) bool {
			return fmt.Sprint(c.Counts[i].Labels) > fmt.Sprint(c.Counts[j].Labels)
		})
		sort.Slice(tc.expected.Counts, func(i, j int) bool {
			return fmt.Sprint(tc.expected.Counts[i].Labels) > fmt.Sprint(tc.expected.Counts[j].Labels)
		})
		assert.Equal(t, tc.expected, c)
	}
}
