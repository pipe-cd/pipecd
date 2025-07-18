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

package cloudrun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/run/v1"
)

func TestMakeCloudRunParent(t *testing.T) {
	t.Parallel()

	const projectID = "projectID"
	got := makeCloudRunParent(projectID)
	want := "namespaces/projectID"
	assert.Equal(t, want, got)
}

func TestMakeCloudRunServiceName(t *testing.T) {
	t.Parallel()

	const (
		projectID = "projectID"
		serviceID = "serviceID"
	)
	got := makeCloudRunServiceName(projectID, serviceID)
	want := "namespaces/projectID/services/serviceID"
	assert.Equal(t, want, got)
}

func TestMakeCloudRunRevisionName(t *testing.T) {
	t.Parallel()

	const (
		projectID  = "projectID"
		revisionID = "revisionID"
	)
	got := makeCloudRunRevisionName(projectID, revisionID)
	want := "namespaces/projectID/revisions/revisionID"
	assert.Equal(t, want, got)
}

func TestPreserveRevisionTags(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		currentService *run.Service
		newSvcCfg      *run.Service
		expected       map[string]struct {
			tag     string
			percent int64
		}
		expectedLen int
	}{
		{
			name: "preserve tags and add missing revisions",
			currentService: &run.Service{
				Spec: &run.ServiceSpec{
					Traffic: []*run.TrafficTarget{
						{RevisionName: "rev1", Tag: "tag1", Percent: 50},
						{RevisionName: "rev2", Tag: "tag2", Percent: 50},
					},
				},
			},
			newSvcCfg: &run.Service{
				Spec: &run.ServiceSpec{
					Traffic: []*run.TrafficTarget{
						{RevisionName: "rev1", Percent: 80},
						{RevisionName: "rev3", Percent: 20},
					},
				},
			},
			expected: map[string]struct {
				tag     string
				percent int64
			}{
				"rev1": {tag: "tag1", percent: 80},
				"rev2": {tag: "tag2", percent: 0},
				"rev3": {tag: "", percent: 20},
			},
			expectedLen: 3,
		},
		{
			name: "no tags to preserve",
			currentService: &run.Service{
				Spec: &run.ServiceSpec{
					Traffic: []*run.TrafficTarget{
						{RevisionName: "old-rev", Percent: 100},
					},
				},
			},
			newSvcCfg: &run.Service{
				Spec: &run.ServiceSpec{
					Traffic: []*run.TrafficTarget{
						{RevisionName: "new-rev", Percent: 100},
					},
				},
			},
			expected: map[string]struct {
				tag     string
				percent int64
			}{
				"new-rev": {tag: "", percent: 100},
			},
			expectedLen: 1,
		},
		{
			name: "revisions with empty tags should not be preserved",
			currentService: &run.Service{
				Spec: &run.ServiceSpec{
					Traffic: []*run.TrafficTarget{
						{RevisionName: "rev1", Tag: "tag1", Percent: 50},
						{RevisionName: "rev2", Tag: "", Percent: 50},
					},
				},
			},
			newSvcCfg: &run.Service{
				Spec: &run.ServiceSpec{
					Traffic: []*run.TrafficTarget{
						{RevisionName: "rev3", Percent: 100},
					},
				},
			},
			expected: map[string]struct {
				tag     string
				percent int64
			}{
				"rev1": {tag: "tag1", percent: 0},
				"rev3": {tag: "", percent: 100},
			},
			expectedLen: 2,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			preserveRevisionTags(tc.currentService, tc.newSvcCfg)

			assert.Len(t, tc.newSvcCfg.Spec.Traffic, tc.expectedLen)

			revisionMap := make(map[string]*run.TrafficTarget)
			for _, traffic := range tc.newSvcCfg.Spec.Traffic {
				revisionMap[traffic.RevisionName] = traffic
			}

			for revName, expected := range tc.expected {
				assert.Contains(t, revisionMap, revName)

				actual := revisionMap[revName]
				assert.Equal(t, expected.tag, actual.Tag)
				assert.Equal(t, expected.percent, actual.Percent)
			}

			assert.Equal(t, len(tc.expected), len(revisionMap))

			var activeTrafficTotal int64
			for _, traffic := range tc.newSvcCfg.Spec.Traffic {
				if traffic.Percent > 0 {
					activeTrafficTotal += traffic.Percent
				}
			}

			assert.Equal(t, int64(100), activeTrafficTotal)
		})
	}
}
