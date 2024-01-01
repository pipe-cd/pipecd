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

package handler

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestMakeApplicationCounts(t *testing.T) {
	testcases := []struct {
		name           string
		counts         []model.InsightApplicationCount
		expectedTotal  int
		expectedGroups map[string]int
	}{
		{
			name:           "empty",
			expectedGroups: map[string]int{},
		},
		{
			name: "one count",
			counts: []model.InsightApplicationCount{
				{
					Labels: map[string]string{
						"KIND":            "KUBERNETES",
						"ACTIVITY_STATUS": "ENABLED",
					},
					Count: 5,
				},
			},
			expectedTotal: 5,
			expectedGroups: map[string]int{
				"KUBERNETES": 5,
			},
		},
		{
			name: "multiple counts",
			counts: []model.InsightApplicationCount{
				{
					Labels: map[string]string{
						"KIND":            "KUBERNETES",
						"ACTIVITY_STATUS": "ENABLED",
					},
					Count: 5,
				},
				{
					Labels: map[string]string{
						"KIND":            "KUBERNETES",
						"ACTIVITY_STATUS": "DISABLED",
					},
					Count: 3,
				},
				{
					Labels: map[string]string{
						"KIND":            "LAMBDA",
						"ACTIVITY_STATUS": "ENABLED",
					},
					Count: 2,
				},
			},
			expectedTotal: 10,
			expectedGroups: map[string]int{
				"KUBERNETES": 8,
				"LAMBDA":     2,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			total, groups := groupApplicationCounts(tc.counts)
			assert.Equal(t, tc.expectedTotal, total)
			assert.Equal(t, tc.expectedGroups, groups)
		})
	}
}

func TestApplicationCountsTmpl(t *testing.T) {
	testcases := []struct {
		name        string
		data        []map[string]interface{}
		expected    string
		expectedErr error
	}{
		{
			name: "ok",
			data: []map[string]interface{}{
				{
					"Project": "one-count",
					"Total":   5,
					"Counts": map[string]int{
						"KUBERNETES": 5,
					},
				},
				{
					"Project": "not-found",
					"Error":   "No data for this project",
				},
				{
					"Project": "multi-counts",
					"Total":   20,
					"Counts": map[string]int{
						"KUBERNETES": 10,
						"CLOUD_RUN":  2,
						"LAMBDA":     8,
					},
				},
			},
			expected: `<!DOCTYPE html>
<html>
<head>
<style>
table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(1) {
  background-color: #dddddd;
}
</style>
</head>
<body>

<h2 style="text-align: center;"><a href="/">Welcome to PipeCD Owner Page!</a></h2>

<h3>There are 3 registered projects</h3>
<h4>0. one-count</h4>

<table>
  <tr>
    <th>Application Kind</th>
    <th>Count</th>
  </tr>
  <tr>
    <td>KUBERNETES</td>
    <td>5</td>
  </tr>
  <tr>
    <td>TOTAL</td>
    <td>5</td>
  </tr>
</table>

<h4>1. not-found</h4>

Unable to fetch application counts (No data for this project).

<h4>2. multi-counts</h4>

<table>
  <tr>
    <th>Application Kind</th>
    <th>Count</th>
  </tr>
  <tr>
    <td>CLOUD_RUN</td>
    <td>2</td>
  </tr>
  <tr>
    <td>KUBERNETES</td>
    <td>10</td>
  </tr>
  <tr>
    <td>LAMBDA</td>
    <td>8</td>
  </tr>
  <tr>
    <td>TOTAL</td>
    <td>20</td>
  </tr>
</table>

</body>
</html>
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := applicationCountsTmpl.Execute(&buf, tc.data)
			assert.Equal(t, tc.expected, buf.String())
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
