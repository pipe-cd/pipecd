// Copyright 2022 The PipeCD Authors.
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

package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestOverlap(t *testing.T) {
	t.Parallel()

	const (
		testDateUnix = 1647820800 // 2022/03/21:00:00:00 UTC
		dayUnix      = int64(time.Hour * 24 / time.Second)
	)

	testcases := []struct {
		name           string
		lhsFrom, lhsTo int64
		rhsFrom, rhsTo int64
		expected       bool
	}{
		{
			name:     "No overlap 1",
			lhsFrom:  testDateUnix - dayUnix,
			lhsTo:    testDateUnix,
			rhsFrom:  testDateUnix,
			rhsTo:    testDateUnix + dayUnix,
			expected: false,
		},
		{
			name:     "No overlap 2",
			lhsFrom:  testDateUnix,
			lhsTo:    testDateUnix + dayUnix,
			rhsFrom:  testDateUnix - dayUnix,
			rhsTo:    testDateUnix,
			expected: false,
		},
		{
			name:     "Overlap same day",
			lhsFrom:  testDateUnix,
			lhsTo:    testDateUnix + dayUnix,
			rhsFrom:  testDateUnix,
			rhsTo:    testDateUnix + dayUnix,
			expected: true,
		},
		{
			name:     "Overlap",
			lhsFrom:  testDateUnix - dayUnix,
			lhsTo:    testDateUnix + dayUnix,
			rhsFrom:  testDateUnix,
			rhsTo:    testDateUnix + dayUnix,
			expected: true,
		},
		{
			name:     "Overlap contain",
			lhsFrom:  testDateUnix - dayUnix,
			lhsTo:    testDateUnix + dayUnix,
			rhsFrom:  testDateUnix,
			rhsTo:    testDateUnix,
			expected: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := overlap(tc.lhsFrom, tc.lhsTo, tc.rhsFrom, tc.rhsTo)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestExtractDailyDeploymentFromChunk(t *testing.T) {
	t.Parallel()

	const (
		testDateUnix = 1647820800 // 2022/03/21:00:00:00 UTC
		dayUnix      = int64(time.Hour * 24 / time.Second)
	)

	testcases := []struct {
		name                     string
		chunk                    *InsightDeploymentChunk
		from, to                 int64
		expectedDailyDeployments []*InsightDeployment
	}{
		{
			name: "SingleDeployments-SingleDeployments",
			chunk: &InsightDeploymentChunk{
				From: testDateUnix, To: testDateUnix + 2,
				Deployments: []*InsightDeployment{
					{CompletedAt: testDateUnix + 1},
					{CompletedAt: testDateUnix + 2},
				},
			},
			from: testDateUnix, to: testDateUnix + 2,
			expectedDailyDeployments: []*InsightDeployment{
				{CompletedAt: testDateUnix + 1},
				{CompletedAt: testDateUnix + 2},
			},
		},
		{
			name: "MultipleDeployments-SingleDeployments",
			chunk: &InsightDeploymentChunk{
				From: testDateUnix - dayUnix, To: testDateUnix + dayUnix + 2,
				Deployments: []*InsightDeployment{
					{CompletedAt: testDateUnix - dayUnix + 1},
					{CompletedAt: testDateUnix - dayUnix + 2},
					{CompletedAt: testDateUnix + 1},
					{CompletedAt: testDateUnix + 2},
					{CompletedAt: testDateUnix + dayUnix + 1},
					{CompletedAt: testDateUnix + dayUnix + 2},
				},
			},
			from: testDateUnix, to: testDateUnix + dayUnix,
			expectedDailyDeployments: []*InsightDeployment{
				{CompletedAt: testDateUnix + 1},
				{CompletedAt: testDateUnix + 2},
			},
		},
		{
			name: "MultipleDeployments-MultipleDeployments",
			chunk: &InsightDeploymentChunk{
				From: testDateUnix - dayUnix, To: testDateUnix + dayUnix + 2,
				Deployments: []*InsightDeployment{
					{CompletedAt: testDateUnix - dayUnix + 1},
					{CompletedAt: testDateUnix - dayUnix + 2},
					{CompletedAt: testDateUnix + 1},
					{CompletedAt: testDateUnix + 2},
					{CompletedAt: testDateUnix + dayUnix + 1},
					{CompletedAt: testDateUnix + dayUnix + 2},
				},
			},
			from: testDateUnix - dayUnix, to: testDateUnix + 2*dayUnix,
			expectedDailyDeployments: []*InsightDeployment{
				{CompletedAt: testDateUnix - dayUnix + 1},
				{CompletedAt: testDateUnix - dayUnix + 2},
				{CompletedAt: testDateUnix + 1},
				{CompletedAt: testDateUnix + 2},
				{CompletedAt: testDateUnix + dayUnix + 1},
				{CompletedAt: testDateUnix + dayUnix + 2},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.chunk.ExtractDeployments(tc.from, tc.to)
			assert.True(t, len(tc.expectedDailyDeployments) == len(got))
			for i := range tc.expectedDailyDeployments {
				assert.True(t, proto.Equal(tc.expectedDailyDeployments[i], got[i]))
			}
		})
	}
}
