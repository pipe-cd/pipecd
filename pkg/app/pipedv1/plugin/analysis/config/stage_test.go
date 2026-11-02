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

	"github.com/pipe-cd/piped-plugin-sdk-go/unit"
	"github.com/stretchr/testify/assert"
)

func TestAnalysisStageOptions_ValidateMetricsDuration(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		duration time.Duration
		interval time.Duration
		wantErr  string
	}{
		{name: "shorter interval", duration: time.Minute, interval: time.Second},
		{name: "equal interval", duration: time.Minute, interval: time.Minute},
		{name: "longer interval", duration: time.Second, interval: time.Minute, wantErr: "exceeds the analysis duration"},
		{name: "missing duration", interval: time.Second, wantErr: "requires duration"},
		{name: "negative duration", duration: -time.Minute, interval: time.Second, wantErr: "duration must be positive"},
		{name: "missing interval", duration: time.Minute, wantErr: "missing \"interval\""},
		{name: "negative interval", duration: time.Minute, interval: -time.Second, wantErr: "interval must be positive"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			opts := AnalysisStageOptions{
				Duration: unit.Duration(tc.duration),
				Metrics: []TemplatableAnalysisMetrics{{AnalysisMetrics: AnalysisMetrics{
					Provider:  "prometheus",
					Query:     "query",
					Deviation: AnalysisDeviationEither,
					Interval:  unit.Duration(tc.interval),
				}}},
			}
			err := opts.Validate()
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestValidateAnalysisTemplateRef(t *testing.T) {
	testcases := []struct {
		name    string
		tplName string
		wantErr bool
	}{
		{
			name:    "valid",
			tplName: "name",
			wantErr: false,
		},
		{
			name:    "invalid due to empty template name",
			tplName: "",
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			a := &AnalysisTemplateRef{
				Name: tc.tplName,
			}
			err := a.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
