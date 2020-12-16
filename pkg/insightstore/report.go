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
	"strconv"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

var ErrValueNotFound = errors.New("value not found")

// deploy frequency

// DeployFrequencyChunk satisfy the interface `Chunk`.
type DeployFrequencyChunk struct {
	AccumulatedTo int64                    `json:"accumulated_to"`
	DataPoints    DeployFrequencyDataPoint `json:"datapoints"`
	FilePath      string
}

type DeployFrequencyDataPoint struct {
	Daily   map[string]DeployFrequency `json:"daily"`
	Weekly  map[string]DeployFrequency `json:"weekly"`
	Monthly map[string]DeployFrequency `json:"monthly"`
	Yearly  map[string]DeployFrequency `json:"yearly"`
}

type DeployFrequency struct {
	DeployCount float32 `json:"deploy_count"`
}

func (d *DeployFrequencyChunk) GetFilePath() string {
	return d.FilePath
}

func (d *DeployFrequencyChunk) PutFilePath(path string) {
	d.FilePath = path
}

func (d *DeployFrequencyChunk) Value(step model.InsightStep, key string) (float32, error) {
	switch step {
	case model.InsightStep_YEARLY:
		yearly, ok := d.DataPoints.Yearly[key]
		if d.DataPoints.Yearly == nil {
			return 0, fmt.Errorf("get value failed, because the chunk does not have Yearly field's value")
		}
		if !ok {
			return 0, ErrValueNotFound
		}
		return yearly.DeployCount, nil
	case model.InsightStep_MONTHLY:
		monthly, ok := d.DataPoints.Monthly[key]
		if d.DataPoints.Monthly == nil {
			return 0, fmt.Errorf("get value failed, because the chunk does not have Yearly field's value")
		}
		if !ok {
			return 0, ErrValueNotFound
		}
		return monthly.DeployCount, nil
	case model.InsightStep_WEEKLY:
		weekly, ok := d.DataPoints.Weekly[key]
		if d.DataPoints.Weekly == nil {
			return 0, fmt.Errorf("get value failed, because the chunk does not have Weekly field's value")
		}
		if !ok {
			return 0, ErrValueNotFound
		}
		return weekly.DeployCount, nil
	case model.InsightStep_DAILY:
		daily, ok := d.DataPoints.Daily[key]
		if d.DataPoints.Daily == nil {
			return 0, fmt.Errorf("get value failed, because the chunk does not have Daily field's value", key)
		}
		if !ok {
			return 0, ErrValueNotFound
		}
		return daily.DeployCount, nil
	}
	return 0, fmt.Errorf("value not found. step: %d, key: %s", step, key)
}

func (d *DeployFrequencyChunk) DataCount(step model.InsightStep) int {
	switch step {
	case model.InsightStep_YEARLY:
		return len(d.DataPoints.Yearly)
	case model.InsightStep_MONTHLY:
		return len(d.DataPoints.Monthly)
	case model.InsightStep_WEEKLY:
		return len(d.DataPoints.Weekly)
	case model.InsightStep_DAILY:
		return len(d.DataPoints.Daily)
	}
	return 0
}

// change failure rate

// ChangeFailureRateChunk satisfy the interface `Chunk`.
type ChangeFailureRateChunk struct {
	AccumulatedTo int64                      `json:"accumulated_to"`
	DataPoints    ChangeFailureRateDataPoint `json:"datapoints"`
	FilePath      string
}

type ChangeFailureRateDataPoint struct {
	Daily   map[string]ChangeFailureRate `json:"daily"`
	Weekly  map[string]ChangeFailureRate `json:"weekly"`
	Monthly map[string]ChangeFailureRate `json:"monthly"`
	Yearly  map[string]ChangeFailureRate `json:"yearly"`
}

type ChangeFailureRate struct {
	Rate         float32 `json:"rate"`
	SuccessCount int64   `json:"success_count"`
	FailureCount int64   `json:"failure_count"`
}

func (c *ChangeFailureRateChunk) GetFilePath() string {
	return c.FilePath
}

func (c *ChangeFailureRateChunk) PutFilePath(path string) {
	c.FilePath = path
}

func (c *ChangeFailureRateChunk) Value(step model.InsightStep, key string) (float32, error) {
	switch step {
	case model.InsightStep_YEARLY:
		yearly, ok := c.DataPoints.Yearly[key]
		if c.DataPoints.Yearly == nil {
			return 0, fmt.Errorf("get value failed, because the chunk does not have Yearly field's value")
		}
		if !ok {
			return 0, ErrValueNotFound
		}
		return yearly.Rate, nil
	case model.InsightStep_MONTHLY:
		monthly, ok := c.DataPoints.Monthly[key]
		if c.DataPoints.Monthly == nil {

			return 0, fmt.Errorf("get value failed, because the chunk does not have Monthly field's value")
		}
		if !ok {
			return 0, ErrValueNotFound
		}
		return monthly.Rate, nil
	case model.InsightStep_WEEKLY:
		weekly, ok := c.DataPoints.Weekly[key]
		if c.DataPoints.Weekly == nil {
			return 0, fmt.Errorf("get value failed, because the chunk does not have Weekly field's value")
		}
		if !ok {
			return 0, ErrValueNotFound
		}
		return weekly.Rate, nil
	case model.InsightStep_DAILY:
		daily, ok := c.DataPoints.Daily[key]
		if c.DataPoints.Daily == nil {
			return 0, fmt.Errorf("get value failed, because the chunk does not have Daily field's value")
		}
		if !ok {
			return 0, ErrValueNotFound
		}
		return daily.Rate, nil
	}
	return 0, fmt.Errorf("value not found. step: %d, key: %s", step, key)
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

type Chunk interface {
	// GetFilePath gets filepath
	GetFilePath() string
	// PutFilePath updates filepath
	PutFilePath(path string)
	// Value gets data by step and key
	Value(step model.InsightStep, key string) (float32, error)
	// DataCount returns count of data in specify step
	DataCount(step model.InsightStep) int
}

// convert below types to chunk.
// - pointer of DeployFrequencyChunk
// - pointer of ChangeFailureRateChunk
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

type Chunks []Chunk

func (cs Chunks) ChunksToDataPoints(from time.Time, count int, step model.InsightStep) ([]*model.InsightDataPoint, error) {
	var idps []*model.InsightDataPoint
	for _, c := range cs {
		idp, err := chunkToDataPoints(c, from, count, step)
		if err != nil {
			return nil, err
		}

		idps = append(idps, idp...)

		count = count - c.DataCount(step)

		nextMonth := time.Date(from.Year(), from.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		from = normalizeTime(nextMonth, step)
		if step == model.InsightStep_WEEKLY && from.Month() != nextMonth.Month() {
			from = from.AddDate(0, 0, 7)
		}
	}

	return idps, nil
}

func chunkToDataPoints(chunk Chunk, from time.Time, count int, step model.InsightStep) ([]*model.InsightDataPoint, error) {
	var getKey func(t time.Time) string
	var nextTargetDate func(t time.Time) time.Time
	switch step {
	case model.InsightStep_YEARLY:
		getKey = func(t time.Time) string {
			return strconv.Itoa(t.Year())
		}
		nextTargetDate = func(t time.Time) time.Time {
			return t.AddDate(1, 0, 0)
		}
	case model.InsightStep_MONTHLY:
		getKey = func(t time.Time) string {
			return t.Format("2006-01")
		}
		nextTargetDate = func(t time.Time) time.Time {
			return t.AddDate(0, 1, 0)
		}
	case model.InsightStep_WEEKLY:
		getKey = func(t time.Time) string {
			// This day must be a Sunday, otherwise it will fail to get the value from the map.
			return t.Format("2006-01-02")
		}
		nextTargetDate = func(t time.Time) time.Time {
			return t.AddDate(0, 0, 7)
		}
	case model.InsightStep_DAILY:
		getKey = func(t time.Time) string {
			return t.Format("2006-01-02")
		}
		nextTargetDate = func(t time.Time) time.Time {
			return t.AddDate(0, 0, 1)
		}
	}

	var idps []*model.InsightDataPoint
	targetDate := from
	for i := 0; i < count; i++ {
		key := getKey(targetDate)
		value, err := chunk.Value(step, key)
		if err != nil {
			if err == ErrValueNotFound {
				return idps, nil
			}
			return nil, err
		}

		idps = append(idps, &model.InsightDataPoint{
			Value:     value,
			Timestamp: targetDate.Unix(),
		})

		targetDate = nextTargetDate(targetDate)
	}

	return idps, nil
}
