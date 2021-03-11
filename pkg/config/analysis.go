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

package config

import "fmt"

// AnalysisMetrics contains common configurable values for deployment analysis with metrics.
type AnalysisMetrics struct {
	Query    string           `json:"query"`
	Expected AnalysisExpected `json:"expected"`
	Interval Duration         `json:"interval"`
	// Maximum number of failed checks before the query result is considered as failure.
	// For instance, If 1 is set, the analysis will be considered a failure after 2 failures.
	FailureLimit int `json:"failureLimit"`
	// How long after which the query times out.
	Timeout  Duration `json:"timeout"`
	Provider string   `json:"provider"`
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
	if min := e.Min; min != nil && *min > value {
		return false
	}
	if max := e.Max; max != nil && *max < value {
		return false
	}
	return true
}

// AnalysisLog contains common configurable values for deployment analysis with log.
type AnalysisLog struct {
	Query    string   `json:"query"`
	Interval Duration `json:"interval"`
	// Maximum number of failed checks before the query result is considered as failure.
	FailureLimit int `json:"failureLimit"`
	// How long after which the query times out.
	Timeout  Duration `json:"timeout"`
	Provider string   `json:"provider"`
}

// AnalysisHTTP contains common configurable values for deployment analysis with http.
type AnalysisHTTP struct {
	URL    string `json:"url"`
	Method string `json:"method"`
	// Custom headers to set in the request. HTTP allows repeated headers.
	Headers          []AnalysisHeader `json:"headers"`
	ExpectedCode     int              `json:"expectedCode"`
	ExpectedResponse string           `json:"expectedResponse"`
	Interval         Duration         `json:"interval"`
	// Maximum number of failed checks before the response is considered as failure.
	FailureLimit int      `json:"failureLimit"`
	Timeout      Duration `json:"timeout"`
}

type AnalysisHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AnalysisDynamic contains settings for analysis by comparing  with dynamic data.
type AnalysisDynamic struct {
	Metrics []AnalysisDynamicMetrics `json:"metrics"`
	Logs    []AnalysisDynamicLog     `json:"logs"`
	Https   []AnalysisDynamicHTTP    `json:"https"`
}

type AnalysisDynamicMetrics struct {
	Query    string   `json:"query"`
	Provider string   `json:"provider"`
	Timeout  Duration `json:"timeout"`
}

type AnalysisDynamicLog struct {
	Query    string   `json:"query"`
	Provider string   `json:"provider"`
	Timeout  Duration `json:"timeout"`
}

type AnalysisDynamicHTTP struct {
	URL              string           `json:"url"`
	Method           string           `json:"method"`
	Headers          []AnalysisHeader `json:"headers"`
	ExpectedCode     int              `json:"expectedCode"`
	ExpectedResponse string           `json:"expectedResponse"`
	Interval         Duration         `json:"interval"`
	Timeout          Duration         `json:"timeout"`
}
