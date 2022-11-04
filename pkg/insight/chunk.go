// Copyright 2022 The PipeCD Authors.
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

// import (
// 	"fmt"
// 	"time"

// 	"github.com/pipe-cd/pipecd/pkg/model"
// )

// // deploy frequency

// // DeployFrequencyChunk represents a chunk of DeployFrequency data points.
// type DeployFrequencyChunk struct {
// 	AccumulatedTo int64                    `json:"accumulated_to"`
// 	DataPoints    DeployFrequencyDataPoint `json:"data_points"`
// 	FilePath      string
// }

// type DeployFrequencyDataPoint struct {
// 	Daily   []*DeployFrequency `json:"daily"`
// 	Weekly  []*DeployFrequency `json:"weekly"`
// 	Monthly []*DeployFrequency `json:"monthly"`
// 	Yearly  []*DeployFrequency `json:"yearly"`
// }

// func (c *DeployFrequencyChunk) GetFilePath() string {
// 	return c.FilePath
// }

// func (c *DeployFrequencyChunk) SetFilePath(path string) {
// 	c.FilePath = path
// }

// func (c *DeployFrequencyChunk) GetAccumulatedTo() int64 {
// 	return c.AccumulatedTo
// }

// func (c *DeployFrequencyChunk) SetAccumulatedTo(a int64) {
// 	c.AccumulatedTo = a
// }

// func (c *DeployFrequencyChunk) GetDataPoints(step model.InsightStep) ([]DataPoint, error) {
// 	switch step {
// 	case model.InsightStep_YEARLY:
// 		return ToDataPoints(c.DataPoints.Yearly)
// 	case model.InsightStep_MONTHLY:
// 		return ToDataPoints(c.DataPoints.Monthly)
// 	case model.InsightStep_WEEKLY:
// 		return ToDataPoints(c.DataPoints.Weekly)
// 	case model.InsightStep_DAILY:
// 		return ToDataPoints(c.DataPoints.Daily)
// 	}
// 	return nil, fmt.Errorf("invalid step: %v", step)
// }

// func (c *DeployFrequencyChunk) SetDataPoints(step model.InsightStep, points []DataPoint) error {
// 	dfs := make([]*DeployFrequency, len(points))
// 	for i, p := range points {
// 		dfs[i] = p.(*DeployFrequency)
// 	}
// 	switch step {
// 	case model.InsightStep_YEARLY:
// 		c.DataPoints.Yearly = dfs
// 	case model.InsightStep_MONTHLY:
// 		c.DataPoints.Monthly = dfs
// 	case model.InsightStep_WEEKLY:
// 		c.DataPoints.Weekly = dfs
// 	case model.InsightStep_DAILY:
// 		c.DataPoints.Daily = dfs
// 	default:
// 		return fmt.Errorf("invalid step: %v", step)
// 	}
// 	return nil
// }

// // change failure rate

// // ChangeFailureRateChunk represents a chunk of ChangeFailureRate data points.
// type ChangeFailureRateChunk struct {
// 	AccumulatedTo int64                      `json:"accumulated_to"`
// 	DataPoints    ChangeFailureRateDataPoint `json:"data_points"`
// 	FilePath      string
// }

// type ChangeFailureRateDataPoint struct {
// 	Daily   []*ChangeFailureRate `json:"daily"`
// 	Weekly  []*ChangeFailureRate `json:"weekly"`
// 	Monthly []*ChangeFailureRate `json:"monthly"`
// 	Yearly  []*ChangeFailureRate `json:"yearly"`
// }

// func (c *ChangeFailureRateChunk) GetFilePath() string {
// 	return c.FilePath
// }

// func (c *ChangeFailureRateChunk) SetFilePath(path string) {
// 	c.FilePath = path
// }

// func (c *ChangeFailureRateChunk) GetAccumulatedTo() int64 {
// 	return c.AccumulatedTo
// }

// func (c *ChangeFailureRateChunk) SetAccumulatedTo(a int64) {
// 	c.AccumulatedTo = a
// }

// func (c *ChangeFailureRateChunk) GetDataPoints(step model.InsightStep) ([]DataPoint, error) {
// 	switch step {
// 	case model.InsightStep_YEARLY:
// 		return ToDataPoints(c.DataPoints.Yearly)
// 	case model.InsightStep_MONTHLY:
// 		return ToDataPoints(c.DataPoints.Monthly)
// 	case model.InsightStep_WEEKLY:
// 		return ToDataPoints(c.DataPoints.Weekly)
// 	case model.InsightStep_DAILY:
// 		return ToDataPoints(c.DataPoints.Daily)
// 	}
// 	return nil, fmt.Errorf("invalid step: %v", step)
// }

// func (c *ChangeFailureRateChunk) SetDataPoints(step model.InsightStep, points []DataPoint) error {
// 	cfs := make([]*ChangeFailureRate, len(points))
// 	for i, p := range points {
// 		cfs[i] = p.(*ChangeFailureRate)
// 	}
// 	switch step {
// 	case model.InsightStep_YEARLY:
// 		c.DataPoints.Yearly = cfs
// 	case model.InsightStep_MONTHLY:
// 		c.DataPoints.Monthly = cfs
// 	case model.InsightStep_WEEKLY:
// 		c.DataPoints.Weekly = cfs
// 	case model.InsightStep_DAILY:
// 		c.DataPoints.Daily = cfs
// 	default:
// 		return fmt.Errorf("invalid step: %v", step)
// 	}
// 	return nil
// }

// type Chunk interface {
// 	// GetFilePath gets filepath
// 	GetFilePath() string
// 	// SetFilePath sets filepath
// 	SetFilePath(path string)
// 	// GetAccumulatedTo gets AccumulatedTo
// 	GetAccumulatedTo() int64
// 	// SetAccumulatedTo sets AccumulatedTo
// 	SetAccumulatedTo(a int64)
// 	// GetDataPoints gets list of data points of specify step
// 	GetDataPoints(step model.InsightStep) ([]DataPoint, error)
// 	// SetDataPoints sets list of data points of specify step
// 	SetDataPoints(step model.InsightStep, points []DataPoint) error
// }

// func NewChunk(projectID string, metricsKind model.InsightMetricsKind, step model.InsightStep, appID string, timestamp time.Time) Chunk {
// 	paths := DetermineFilePaths(projectID, appID, metricsKind, step, timestamp, 1)
// 	path := paths[0]

// 	var chunk Chunk
// 	switch metricsKind {
// 	case model.InsightMetricsKind_DEPLOYMENT_FREQUENCY:
// 		chunk = &DeployFrequencyChunk{
// 			FilePath: path,
// 		}
// 	case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
// 		chunk = &ChangeFailureRateChunk{
// 			FilePath: path,
// 		}
// 	default:
// 		return nil
// 	}

// 	return chunk
// }

// // convert types to Chunk.
// func ToChunk(i interface{}) (Chunk, error) {
// 	switch p := i.(type) {
// 	case *DeployFrequencyChunk:
// 		return p, nil
// 	case *ChangeFailureRateChunk:
// 		return p, nil
// 	default:
// 		return nil, fmt.Errorf("cannot convert to Chunk: %v", p)
// 	}
// }

// type Chunks []Chunk

// func (cs Chunks) ExtractDataPoints(step model.InsightStep, from time.Time, count int) ([]*model.InsightDataPoint, error) {
// 	var out []*model.InsightDataPoint
// 	var to time.Time
// 	switch step {
// 	case model.InsightStep_YEARLY:
// 		to = from.AddDate(count-1, 0, 0)
// 	case model.InsightStep_MONTHLY:
// 		to = from.AddDate(0, count-1, 0)
// 	case model.InsightStep_WEEKLY:
// 		to = from.AddDate(0, 0, (count-1)*7)
// 	case model.InsightStep_DAILY:
// 		to = from.AddDate(0, 0, count-1)
// 	}

// 	for _, c := range cs {
// 		dp, err := c.GetDataPoints(step)
// 		if err != nil {
// 			return nil, err
// 		}

// 		idp, err := extractDataPoints(dp, from, to)
// 		if err != nil {
// 			return nil, err
// 		}

// 		out = append(out, idp...)
// 	}

// 	return out, nil
// }
