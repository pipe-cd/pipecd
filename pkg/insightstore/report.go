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
	"strconv"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

// deploy frequency

// DeployFrequencyReport satisfy the interface `Report`.
type DeployFrequencyReport struct {
	AccumulatedTo int64                    `json:"accumulated_to"`
	Datapoints    DeployFrequencyDataPoint `json:"datapoints"`
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

func (d *DeployFrequencyReport) GetFilePath() string {
	return d.FilePath
}

func (d *DeployFrequencyReport) PutFilePath(path string) {
	d.FilePath = path
}

func (d *DeployFrequencyReport) Value(step model.InsightStep, key string) (float32, error) {
	switch step {
	case model.InsightStep_YEARLY:
		return d.Datapoints.Yearly[key].DeployCount, nil
	case model.InsightStep_MONTHLY:
		return d.Datapoints.Monthly[key].DeployCount, nil
	case model.InsightStep_WEEKLY:
		return d.Datapoints.Weekly[key].DeployCount, nil
	case model.InsightStep_DAILY:
		return d.Datapoints.Daily[key].DeployCount, nil
	}
	return 0, fmt.Errorf("value not found. step: %d, key: %s", step, key)
}

func (d *DeployFrequencyReport) DataCount(step model.InsightStep) int {
	switch step {
	case model.InsightStep_YEARLY:
		return len(d.Datapoints.Yearly)
	case model.InsightStep_MONTHLY:
		return len(d.Datapoints.Monthly)
	case model.InsightStep_WEEKLY:
		return len(d.Datapoints.Weekly)
	case model.InsightStep_DAILY:
		return len(d.Datapoints.Daily)
	}
	return 0
}

// change failure rate

// ChangeFailureRateReport satisfy the interface `Report`.
type ChangeFailureRateReport struct {
	AccumulatedTo int64                      `json:"accumulated_to"`
	Datapoints    ChangeFailureRateDataPoint `json:"datapoints"`
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

func (c *ChangeFailureRateReport) GetFilePath() string {
	return c.FilePath
}

func (c *ChangeFailureRateReport) PutFilePath(path string) {
	c.FilePath = path
}

func (c *ChangeFailureRateReport) Value(step model.InsightStep, key string) (float32, error) {
	switch step {
	case model.InsightStep_YEARLY:
		return c.Datapoints.Yearly[key].Rate, nil
	case model.InsightStep_MONTHLY:
		return c.Datapoints.Monthly[key].Rate, nil
	case model.InsightStep_WEEKLY:
		return c.Datapoints.Weekly[key].Rate, nil
	case model.InsightStep_DAILY:
		return c.Datapoints.Daily[key].Rate, nil
	}
	return 0, fmt.Errorf("value not found. step: %d, key: %s", step, key)
}

func (c *ChangeFailureRateReport) DataCount(step model.InsightStep) int {
	switch step {
	case model.InsightStep_YEARLY:
		return len(c.Datapoints.Yearly)
	case model.InsightStep_MONTHLY:
		return len(c.Datapoints.Monthly)
	case model.InsightStep_WEEKLY:
		return len(c.Datapoints.Weekly)
	case model.InsightStep_DAILY:
		return len(c.Datapoints.Daily)
	}
	return 0
}

type Report interface {
	// GetFilePath gets filepath
	GetFilePath() string
	// PutFilePath updates filepath
	PutFilePath(path string)
	// Value gets data by step and key
	Value(step model.InsightStep, key string) (float32, error)
	// DataCount returns count of data in specify step
	DataCount(step model.InsightStep) int
}

// convert below types to report.
// - pointer of DeployFrequencyReport
// - pointer of ChangeFailureRateReport
func toReport(i interface{}) (Report, error) {
	switch p := i.(type) {
	case *DeployFrequencyReport:
		return p, nil
	case *ChangeFailureRateReport:
		return p, nil
	default:
		return nil, fmt.Errorf("cannot convert to Report: %v", p)
	}

}

func convertToInsightDataPoints(report Report, from time.Time, dataPointCount int, step model.InsightStep) ([]*model.InsightDataPoint, error) {
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

	idps := make([]*model.InsightDataPoint, dataPointCount)
	targetDate := from
	for i := 0; i < dataPointCount; i++ {
		key := getKey(targetDate)
		value, err := report.Value(step, key)
		if err != nil {
			return nil, err
		}

		idps[i] = &model.InsightDataPoint{
			Value:     value,
			Timestamp: targetDate.Unix(),
		}

		targetDate = nextTargetDate(targetDate)
	}

	return idps, nil
}
