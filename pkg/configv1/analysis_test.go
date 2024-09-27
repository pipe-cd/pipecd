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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func floatPointer(v float64) *float64 {
	return &v
}

func TestAnalysisExpectedString(t *testing.T) {
	testcases := []struct {
		name string
		Min  *float64
		Max  *float64
		want string
	}{
		{
			name: "only min given",
			Min:  floatPointer(1.5),
			want: "1.5 <=",
		},
		{
			name: "only max given",
			Max:  floatPointer(1.5),
			want: "<= 1.5",
		},
		{
			name: "both min and max given",
			Min:  floatPointer(1.5),
			Max:  floatPointer(2.5),
			want: "1.5 <= 2.5",
		},
		{
			name: "invalid range",
			want: "",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			e := &AnalysisExpected{
				Min: tc.Min,
				Max: tc.Max,
			}
			got := e.String()
			assert.Equal(t, tc.want, got)
		})
	}
}
