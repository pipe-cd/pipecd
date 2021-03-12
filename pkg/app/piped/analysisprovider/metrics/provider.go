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
	ErrNoValuesFound = errors.New("no values found")
)

// Provider represents a client for metrics provider which provides metrics for analysis.
type Provider interface {
	Type() string
	// RunQuery runs the given query against the metrics provider,
	// and then checks if the results are expected or not.
	// TODO: Give back the reason of the result
	RunQuery(ctx context.Context, query string, evaluator Evaluator, queryRange QueryRange) (result bool, err error)
}

// Evaluator evaluates the response from the metrics provider.
type Evaluator interface {
	// InRange checks if the value is expected one.
	InRange(value float64) bool
	// Validates ensures its own configuration has no problem.
	Validate() error
}

// QueryRange represents a sliced time range.
type QueryRange struct {
	// Required: Start of the queried time period
	From time.Time
	// End of the queried time period. Defaults to the current time.
	To time.Time
	// Query resolution step width. Defaults to 1m.
	Step time.Duration
}

func (q *QueryRange) Validate() error {
	if q.From.IsZero() {
		return fmt.Errorf("start of the query range is required")
	}
	if q.To.IsZero() {
		q.To = time.Now()
	}
	// TODO: Look into the appropriate default value of a Step
	if q.Step == 0 {
		q.Step = time.Minute
	}
	return nil
}
