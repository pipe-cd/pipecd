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

package config

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	AnalysisStrategyThreshold      = "THRESHOLD"
	AnalysisStrategyPrevious       = "PREVIOUS"
	AnalysisStrategyCanaryBaseline = "CANARY_BASELINE"
	AnalysisStrategyCanaryPrimary  = "CANARY_PRIMARY"

	AnalysisDeviationEither = "EITHER"
	AnalysisDeviationHigh   = "HIGH"
	AnalysisDeviationLow    = "LOW"
)

// AnalysisMetrics contains common configurable values for deployment analysis with metrics.
type AnalysisMetrics struct {
	// The strategy name. One of THRESHOLD or PREVIOUS or CANARY_BASELINE or CANARY_PRIMARY is available.
	// Defaults to THRESHOLD.
	Strategy string `json:"strategy" default:"THRESHOLD"`
	// The unique name of provider defined in the Piped Configuration.
	// Required field.
	Provider string `json:"provider"`
	// A query performed against the Analysis Provider.
	// Required field.
	Query string `json:"query"`
	// The expected query result.
	// Required field for the THRESHOLD strategy.
	Expected AnalysisExpected `json:"expected"`
	// Run a query at this intervals.
	// Required field.
	Interval Duration `json:"interval"`
	// Acceptable number of failures. For instance, If 1 is set,
	// the analysis will be considered a failure after 2 failures.
	// Default is 0.
	FailureLimit int `json:"failureLimit"`
	// If true, it considers as a success when no data returned from the analysis provider.
	// Default is false.
	SkipOnNoData bool `json:"skipOnNoData"`
	// How long after which the query times out.
	// Default is 30s.
	Timeout Duration `json:"timeout"`

	// The stage fails on deviation in the specified direction. One of LOW or HIGH or EITHER is available.
	// This can be used only for PREVIOUS, CANARY_BASELINE or CANARY_PRIMARY. Defaults to EITHER.
	Deviation string `json:"deviation" default:"EITHER"`
	// The custom arguments to be populated for the Canary query.
	// They can be referred as {{ .VariantArgs.xxx }}.
	CanaryArgs map[string]string `json:"canaryArgs"`
	// The custom arguments to be populated for the Baseline query.
	// They can be referred as {{ .VariantArgs.xxx }}.
	BaselineArgs map[string]string `json:"baselineArgs"`
	// The custom arguments to be populated for the Primary query.
	// They can be referred as {{ .VariantArgs.xxx }}.
	PrimaryArgs map[string]string `json:"primaryArgs"`
}

func (m *AnalysisMetrics) Validate() error {
	if m.Provider == "" {
		return fmt.Errorf("missing \"provider\" field")
	}
	if m.Query == "" {
		return fmt.Errorf("missing \"query\" field")
	}
	if m.Interval == 0 {
		return fmt.Errorf("missing \"interval\" field")
	}
	if m.Deviation != AnalysisDeviationEither && m.Deviation != AnalysisDeviationHigh && m.Deviation != AnalysisDeviationLow {
		return fmt.Errorf("\"deviation\" have to be one of %s, %s or %s", AnalysisDeviationEither, AnalysisDeviationHigh, AnalysisDeviationLow)
	}
	return nil
}

// AnalysisExpected defines the range used for metrics analysis.
type AnalysisExpected struct {
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`
}

func (e *AnalysisExpected) Validate() error {
	if e.Min == nil && e.Max == nil {
		return fmt.Errorf("expected range is undefined")
	}
	return nil
}

// InRange returns true if the given value is within the range.
func (e *AnalysisExpected) InRange(value float64) bool {
	if e.Min != nil && *e.Min > value {
		return false
	}
	if e.Max != nil && *e.Max < value {
		return false
	}
	return true
}

func (e *AnalysisExpected) String() string {
	if e.Min == nil && e.Max == nil {
		return ""
	}

	var b strings.Builder
	if e.Min != nil {
		min := strconv.FormatFloat(*e.Min, 'f', -1, 64)
		b.WriteString(min + " ")
	}

	b.WriteString("<=")

	if e.Max != nil {
		max := strconv.FormatFloat(*e.Max, 'f', -1, 64)
		b.WriteString(" " + max)
	}
	return b.String()
}

// AnalysisLog contains common configurable values for deployment analysis with log.
type AnalysisLog struct {
	Query    string   `json:"query"`
	Interval Duration `json:"interval"`
	// Maximum number of failed checks before the query result is considered as failure.
	FailureLimit int `json:"failureLimit"`
	// If true, it considers as success when no data returned from the analysis provider.
	// Default is false.
	SkipOnNoData bool `json:"skipOnNoData"`
	// How long after which the query times out.
	Timeout  Duration `json:"timeout"`
	Provider string   `json:"provider"`
}

func (a *AnalysisLog) Validate() error {
	return nil
}

// AnalysisHTTP contains common configurable values for deployment analysis with http.
type AnalysisHTTP struct {
	URL    string `json:"url"`
	Method string `json:"method"`
	// Custom headers to set in the request. HTTP allows repeated headers.
	Headers          []AnalysisHTTPHeader `json:"headers"`
	ExpectedCode     int                  `json:"expectedCode"`
	ExpectedResponse string               `json:"expectedResponse"`
	Interval         Duration             `json:"interval"`
	// Maximum number of failed checks before the response is considered as failure.
	FailureLimit int `json:"failureLimit"`
	// If true, it considers as success when no data returned from the analysis provider.
	// Default is false.
	SkipOnNoData bool     `json:"skipOnNoData"`
	Timeout      Duration `json:"timeout"`
}

func (a *AnalysisHTTP) Validate() error {
	return nil
}

type AnalysisHTTPHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
