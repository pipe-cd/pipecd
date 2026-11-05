// Copyright 2025 The PipeCD Authors.
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
	"time"

	"github.com/stretchr/testify/assert"

	unit "github.com/pipe-cd/piped-plugin-sdk-go/unit"
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

func TestAnalysisLogValidate(t *testing.T) {
	testcases := []struct {
		name     string
		interval unit.Duration
		wantErr  bool
	}{
		{name: "valid", interval: unit.Duration(time.Minute), wantErr: false},
		{name: "missing interval", interval: 0, wantErr: true},
		{name: "negative interval", interval: -unit.Duration(time.Minute), wantErr: true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			l := &AnalysisLog{Query: "query", Interval: tc.interval}
			err := l.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestAnalysisHTTPValidate(t *testing.T) {
	testcases := []struct {
		name     string
		interval unit.Duration
		wantErr  bool
	}{
		{name: "valid", interval: unit.Duration(time.Minute), wantErr: false},
		{name: "missing interval", interval: 0, wantErr: true},
		{name: "negative interval", interval: -unit.Duration(time.Minute), wantErr: true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			h := &AnalysisHTTP{
				URL:          "http://example.com",
				Method:       "GET",
				ExpectedCode: 200,
				Interval:     tc.interval,
			}
			err := h.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
