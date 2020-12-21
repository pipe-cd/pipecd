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

package insightstore

import (
	"errors"
	"fmt"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

var (
	ErrNotFound = errors.New("data point not found")
)

// DeployFrequency represents a data point that shows the deployment frequency metrics.
type DeployFrequency struct {
	Timestamp   int64   `json:"timestamp"`
	DeployCount float32 `json:"deploy_count"`
}

func (d *DeployFrequency) GetTimestamp() int64 {
	return d.Timestamp
}

func (d *DeployFrequency) Value() float32 {
	return d.DeployCount
}

func (d *DeployFrequency) Merge(point DataPoint) error {
	if point == nil {
		return nil
	}

	df, ok := point.(*DeployFrequency)
	if !ok {
		return fmt.Errorf("cannot cast to DeployFrequency. DataPoint: %v", point)
	}

	if df.Timestamp != d.Timestamp {
		return fmt.Errorf("mismatch timestamp. want: %d, acutual: %d", d.Timestamp, df.Timestamp)
	}

	d.DeployCount += df.DeployCount
	return nil

}

// ChangeFailureRate represents a data point that shows the change failure rate metrics.
type ChangeFailureRate struct {
	Timestamp    int64   `json:"timestamp"`
	Rate         float32 `json:"rate"`
	SuccessCount int64   `json:"success_count"`
	FailureCount int64   `json:"failure_count"`
}

func (c *ChangeFailureRate) GetTimestamp() int64 {
	return c.Timestamp
}

func (c *ChangeFailureRate) Value() float32 {
	return c.Rate
}

func (c *ChangeFailureRate) Merge(point DataPoint) error {
	if point == nil {
		return nil
	}

	cfr, ok := point.(*ChangeFailureRate)
	if !ok {
		return fmt.Errorf("cannot cast to DeployFrequency. DataPoint: %v", point)
	}

	if cfr.Timestamp != c.Timestamp {
		return fmt.Errorf("mismatch timestamp. want: %d, acutual: %d", c.Timestamp, cfr.Timestamp)
	}

	c.FailureCount += cfr.FailureCount
	c.SuccessCount += cfr.SuccessCount
	c.Rate = float32(c.FailureCount) / float32(c.FailureCount+c.SuccessCount)
	return nil
}

type DataPoint interface {
	// Value gets data for model.InsightDataPoint.
	Value() float32
	// Timestamp gets timestamp.
	GetTimestamp() int64
	// Merge merges other DataPoint.
	Merge(point DataPoint) error
}

// convert types to list of DataPoint.
func ToDataPoints(i interface{}) ([]DataPoint, error) {
	switch dps := i.(type) {
	case []*DeployFrequency:
		dataPoints := make([]DataPoint, len(dps))
		for j, dp := range dps {
			dataPoints[j] = dp
		}
		return dataPoints, nil
	case []*ChangeFailureRate:
		dataPoints := make([]DataPoint, len(dps))
		for j, dp := range dps {
			dataPoints[j] = dp
		}
		return dataPoints, nil
	default:
		return nil, fmt.Errorf("cannot convert to DataPoints: %v", dps)
	}
}

// findDataPoint find key in the list of data points by timestamp
func findDataPoint(dp []DataPoint, timestamp int64) (int, error) {
	for i, d := range dp {
		ts := d.GetTimestamp()
		if ts == timestamp {
			return i, nil
		}
	}
	return 0, ErrNotFound
}

// GetDataPoint gets a data point by timestamp
func GetDataPoint(dp []DataPoint, timestamp int64) (DataPoint, error) {
	for _, d := range dp {
		ts := d.GetTimestamp()
		if ts == timestamp {
			return d, nil
		}
	}
	return nil, ErrNotFound
}

// SetDataPoint sets data point specified by timestamp
func SetDataPoint(dp []DataPoint, point DataPoint, timestamp int64) []DataPoint {
	k, err := findDataPoint(dp, timestamp)
	if err != nil {
		if err == ErrNotFound {
			return append(dp, point)
		}
	}
	dp[k] = point
	return dp
}

func extractDataPoints(dp []DataPoint, from, to time.Time) ([]*model.InsightDataPoint, error) {
	var result []*model.InsightDataPoint
	for _, d := range dp {
		ts := d.GetTimestamp()
		if ts > to.Unix() {
			break
		}

		if ts >= from.Unix() {
			result = append(result, &model.InsightDataPoint{
				Timestamp: ts,
				Value:     d.Value(),
			})
		}
	}
	return result, nil
}
