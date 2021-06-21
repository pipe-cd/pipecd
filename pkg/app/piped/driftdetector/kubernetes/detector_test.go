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

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
)

func TestGroupManifests(t *testing.T) {
	testcases := []struct {
		name               string
		heads              []provider.Manifest
		lives              []provider.Manifest
		expectedAdds       []provider.Manifest
		expectedDeletes    []provider.Manifest
		expectedHeadInters []provider.Manifest
		expectedLiveInters []provider.Manifest
	}{
		{
			name: "empty list",
		},
		{
			name: "only adds",
			lives: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "b"}},
				{Key: provider.ResourceKey{Name: "a"}},
			},
			expectedAdds: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
		},
		{
			name: "only deletes",
			heads: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "b"}},
				{Key: provider.ResourceKey{Name: "a"}},
			},
			expectedDeletes: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
		},
		{
			name: "only inters",
			heads: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "b"}},
				{Key: provider.ResourceKey{Name: "a"}},
			},
			lives: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
			expectedHeadInters: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
			expectedLiveInters: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
		},
		{
			name: "all kinds",
			heads: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "b"}},
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "c"}},
			},
			lives: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "d"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
			expectedAdds: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "d"}},
			},
			expectedDeletes: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "c"}},
			},
			expectedHeadInters: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
			expectedLiveInters: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			adds, deletes, headInters, liveInters := groupManifests(tc.heads, tc.lives)
			assert.Equal(t, tc.expectedAdds, adds)
			assert.Equal(t, tc.expectedDeletes, deletes)
			assert.Equal(t, tc.expectedHeadInters, headInters)
			assert.Equal(t, tc.expectedLiveInters, liveInters)
		})
	}
}
