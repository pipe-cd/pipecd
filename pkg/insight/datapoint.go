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

package insight

import (
	"fmt"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
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
		return fmt.Errorf("can not cast to DataPoint to DeployFrequency, %v", point)
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
		return fmt.Errorf("can not cast to DataPoint to ChangeFailureRate, %v", point)
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

// ToDataPoints converts a list of concrete points into the list of DataPoints
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

// UpdateDataPoint sets data point
func UpdateDataPoint(dp []DataPoint, point DataPoint, timestamp int64) ([]DataPoint, error) {
	latestData := dp[len(dp)-1]
	if timestamp < latestData.GetTimestamp() {
		return nil, fmt.Errorf("invalid timestamp")
	}

	if timestamp == latestData.GetTimestamp() {
		err := latestData.Merge(point)
		if err != nil {
			return nil, err
		}
		dp[len(dp)-1] = latestData
	} else {
		dp = append(dp, point)
	}
	return dp, nil
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
