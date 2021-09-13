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

package metrics

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNoDataFound = errors.New("no data found")
)

// Provider represents a client for metrics provider which provides metrics for analysis.
type Provider interface {
	Type() string

	// Evaluate runs the given query against the metrics provider,
	// and then checks if the results are expected or not.
	// Returns the result reason if non-error occurred.
	// The first value "expected" must be false if err isn't nil.
	// TODO: Do not evaluate data points by Analysis providers
	//   Instead, the executor should do that by using QueryPoints().
	Evaluate(ctx context.Context, query string, queryRange QueryRange, evaluator Evaluator) (expected bool, reason string, err error)

	// QueryPoints gives back data points within the given range.
	QueryPoints(ctx context.Context, query string, queryRange QueryRange) (points []DataPoint, err error)
}

type DataPoint struct {
	Timestamp int64
	Value     float64
}

func (d *DataPoint) String() string {
	return fmt.Sprintf("timestamp: %v, value: %g", time.Unix(d.Timestamp, 0), d.Value)
}

// Evaluator evaluates the response from the metrics provider.
type Evaluator interface {
	// InRange checks if the value is expected one.
	InRange(value float64) bool
	String() string
}

// QueryRange represents a sliced time range.
type QueryRange struct {
	// Required: Start of the queried time period
	From time.Time
	// End of the queried time period. Defaults to the current time.
	To time.Time
}

func (q *QueryRange) Validate() error {
	if q.From.IsZero() {
		return fmt.Errorf("start of the query range is required")
	}
	if q.To.IsZero() {
		q.To = time.Now()
	}
	if q.From.After(q.To) {
		return fmt.Errorf("\"to\" should be after \"from\"")
	}
	return nil
}
