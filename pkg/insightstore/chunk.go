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
	"fmt"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

// deploy frequency

// DeployFrequencyChunk represents a chunk of DeployFrequency data points.
type DeployFrequencyChunk struct {
	AccumulatedTo int64                    `json:"accumulated_to"`
	DataPoints    DeployFrequencyDataPoint `json:"data_points"`
	FilePath      string
}

type DeployFrequencyDataPoint struct {
	Daily   []DeployFrequency `json:"daily"`
	Weekly  []DeployFrequency `json:"weekly"`
	Monthly []DeployFrequency `json:"monthly"`
	Yearly  []DeployFrequency `json:"yearly"`
}

// DeployFrequency represents a data point that shows the deployment frequency metrics.
type DeployFrequency struct {
	Timestamp   int64   `json:"timestamp"`
	DeployCount float32 `json:"deploy_count"`
}

func (c *DeployFrequencyChunk) GetFilePath() string {
	return c.FilePath
}

func (c *DeployFrequencyChunk) SetFilePath(path string) {
	c.FilePath = path
}

func (c *DeployFrequencyChunk) GetAccumulatedTo() int64 {
	return c.AccumulatedTo
}

func (c *DeployFrequencyChunk) SetAccumulatedTo(a int64) {
	c.AccumulatedTo = a
}

func (c *DeployFrequencyChunk) DataCount(step model.InsightStep) int {
	switch step {
	case model.InsightStep_YEARLY:
		return len(c.DataPoints.Yearly)
	case model.InsightStep_MONTHLY:
		return len(c.DataPoints.Monthly)
	case model.InsightStep_WEEKLY:
		return len(c.DataPoints.Weekly)
	case model.InsightStep_DAILY:
		return len(c.DataPoints.Daily)
	}
	return 0
}

func (c *DeployFrequencyChunk) GetDataPoints(step model.InsightStep) ([]DataPoint, error) {
	switch step {
	case model.InsightStep_YEARLY:
		return toDataPoints(c.DataPoints.Yearly)
	case model.InsightStep_MONTHLY:
		return toDataPoints(c.DataPoints.Monthly)
	case model.InsightStep_WEEKLY:
		return toDataPoints(c.DataPoints.Weekly)
	case model.InsightStep_DAILY:
		return toDataPoints(c.DataPoints.Daily)
	}
	return []DataPoint{}, fmt.Errorf("invalid step: %v", step)
}

func (d DeployFrequency) GetTimestamp() int64 {
	return d.Timestamp
}

func (d DeployFrequency) Value() float32 {
	return d.DeployCount
}

// change failure rate

// ChangeFailureRateChunk represents a chunk of ChangeFailureRate data points.
type ChangeFailureRateChunk struct {
	AccumulatedTo int64                      `json:"accumulated_to"`
	DataPoints    ChangeFailureRateDataPoint `json:"data_points"`
	FilePath      string
}

type ChangeFailureRateDataPoint struct {
	Daily   []ChangeFailureRate `json:"daily"`
	Weekly  []ChangeFailureRate `json:"weekly"`
	Monthly []ChangeFailureRate `json:"monthly"`
	Yearly  []ChangeFailureRate `json:"yearly"`
}

// ChangeFailureRate represents a data point that shows the change failure rate metrics.
type ChangeFailureRate struct {
	Timestamp    int64   `json:"timestamp"`
	Rate         float32 `json:"rate"`
	SuccessCount int64   `json:"success_count"`
	FailureCount int64   `json:"failure_count"`
}

func (c *ChangeFailureRateChunk) GetFilePath() string {
	return c.FilePath
}

func (c *ChangeFailureRateChunk) SetFilePath(path string) {
	c.FilePath = path
}

func (c *ChangeFailureRateChunk) GetAccumulatedTo() int64 {
	return c.AccumulatedTo
}

func (c *ChangeFailureRateChunk) SetAccumulatedTo(a int64) {
	c.AccumulatedTo = a
}

func (c *ChangeFailureRateChunk) GetDataPoints(step model.InsightStep) ([]DataPoint, error) {
	switch step {
	case model.InsightStep_YEARLY:
		return toDataPoints(c.DataPoints.Yearly)
	case model.InsightStep_MONTHLY:
		return toDataPoints(c.DataPoints.Monthly)
	case model.InsightStep_WEEKLY:
		return toDataPoints(c.DataPoints.Weekly)
	case model.InsightStep_DAILY:
		return toDataPoints(c.DataPoints.Daily)
	}
	return []DataPoint{}, fmt.Errorf("invalid step: %v", step)
}

func (c *ChangeFailureRateChunk) DataCount(step model.InsightStep) int {
	switch step {
	case model.InsightStep_YEARLY:
		return len(c.DataPoints.Yearly)
	case model.InsightStep_MONTHLY:
		return len(c.DataPoints.Monthly)
	case model.InsightStep_WEEKLY:
		return len(c.DataPoints.Weekly)
	case model.InsightStep_DAILY:
		return len(c.DataPoints.Daily)
	}
	return 0
}

func (c ChangeFailureRate) GetTimestamp() int64 {
	return c.Timestamp
}

func (c ChangeFailureRate) Value() float32 {
	return c.Rate
}

type Chunk interface {
	// GetFilePath gets filepath
	GetFilePath() string
	// SetFilePath sets filepath
	SetFilePath(path string)
	// GetAccumulatedTo gets AccumulatedTo
	GetAccumulatedTo() int64
	// SetAccumulatedTo sets AccumulatedTo
	SetAccumulatedTo(a int64)
	// GetDataPoint gets list of data points of specify step
	GetDataPoints(step model.InsightStep) ([]DataPoint, error)
	// DataCount returns count of data in specify step
	DataCount(step model.InsightStep) int
}

// convert types to Chunk.
func toChunk(i interface{}) (Chunk, error) {
	switch p := i.(type) {
	case *DeployFrequencyChunk:
		return p, nil
	case *ChangeFailureRateChunk:
		return p, nil
	default:
		return nil, fmt.Errorf("cannot convert to Chunk: %v", p)
	}
}

type DataPoint interface {
	// Value gets data for model.InsightDataPoint
	Value() float32
	// Timestamp gets timestamp
	GetTimestamp() int64
}

// convert types to list of DataPoint.
func toDataPoints(i interface{}) ([]DataPoint, error) {
	switch dps := i.(type) {
	case []DeployFrequency:
		dataPoints := make([]DataPoint, len(dps))
		for j, dp := range dps {
			dataPoints[j] = dp
		}
		return dataPoints, nil
	case []ChangeFailureRate:
		dataPoints := make([]DataPoint, len(dps))
		for j, dp := range dps {
			dataPoints[j] = dp
		}
		return dataPoints, nil
	default:
		return nil, fmt.Errorf("cannot convert to DataPoints: %v", dps)
	}
}

type Chunks []Chunk

func (cs Chunks) ExtractDataPoints(step model.InsightStep, from time.Time, count int) ([]*model.InsightDataPoint, error) {
	var out []*model.InsightDataPoint
	var to time.Time
	switch step {
	case model.InsightStep_YEARLY:
		to = from.AddDate(count-1, 0, 0)
	case model.InsightStep_MONTHLY:
		to = from.AddDate(0, count-1, 0)
	case model.InsightStep_WEEKLY:
		to = from.AddDate(0, 0, (count-1)*7)
	case model.InsightStep_DAILY:
		to = from.AddDate(0, 0, count-1)
	}

	for _, c := range cs {
		idp, err := extractDataPoints(c, step, from, to)
		if err != nil {
			return nil, err
		}

		out = append(out, idp...)
	}

	return out, nil
}

func extractDataPoints(chunk Chunk, step model.InsightStep, from, to time.Time) ([]*model.InsightDataPoint, error) {
	target, err := chunk.GetDataPoints(step)
	if err != nil {
		return []*model.InsightDataPoint{}, nil
	}

	var result []*model.InsightDataPoint
	for _, d := range target {
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
