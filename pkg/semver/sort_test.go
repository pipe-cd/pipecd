// Copyright (C) 2014-2019, Matt Butcher and Matt Farina

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

package semver

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollection(t *testing.T) {
	cases := []struct {
		name     string
		versions []string
		want     []string
	}{
		{
			name: "only major versions",
			versions: []string{
				"3",
				"1",
				"2",
			},
			want: []string{
				"3.0.0",
				"2.0.0",
				"1.0.0",
			},
		},
		{
			name: "various versions",
			versions: []string{
				"1.2.3",
				"1.0",
				"1.3",
				"2",
				"0.4.2",
			},
			want: []string{
				"2.0.0",
				"1.3.0",
				"1.2.3",
				"1.0.0",
				"0.4.2",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			versions := make([]*Version, 0, len(tc.versions))
			for _, s := range tc.versions {
				v, err := NewVersion(s)
				require.NoError(t, err)
				versions = append(versions, v)
			}
			vs := ByNewer(versions)
			sort.Sort(vs)

			got := make([]string, 0, len(vs))
			for _, v := range vs {
				got = append(got, v.String())
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
