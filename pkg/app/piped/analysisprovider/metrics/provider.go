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

package metrics

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05 MST"

var (
	ErrNoDataFound = errors.New("no data found")
)

// Provider represents a client for metrics provider which provides metrics for analysis.
type Provider interface {
	Type() string
	// QueryPoints gives back data points within the given range.
	QueryPoints(ctx context.Context, query string, queryRange QueryRange) (points []DataPoint, err error)
}

type DataPoint struct {
	// Unix timestamp in seconds.
	Timestamp int64
	Value     float64
}

func (d *DataPoint) String() string {
	// Timestamp is shown in UTC.
	return fmt.Sprintf("timestamp: %q, value: %g", time.Unix(d.Timestamp, 0).UTC().Format(timeFormat), d.Value)
}

// QueryRange represents a sliced time range.
type QueryRange struct {
	// Required: Start of the queried time period
	From time.Time
	// End of the queried time period. Defaults to the current time.
	To time.Time
}

func (q *QueryRange) String() string {
	// Timestamps are shown in UTC.
	return fmt.Sprintf("from: %q, to: %q", q.From.UTC().Format(timeFormat), q.To.UTC().Format(timeFormat))
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
