package insightstore

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

type Report interface {
	// GetFilePath get filepath
	GetFilePath() string
	// PutFilePath update filepath
	PutFilePath(path string)
	// Value get data by step and key
	Value(step model.InsightStep, key string) (float32, error)
}

// convert below types to report
// - pointer of deployFrequencyReport
// - pointer of changeFailureRateReport
func toReport(i interface{}) (Report, error) {
	switch p := i.(type) {
	case *deployFrequencyReport:
		return p, nil
	case *changeFailureRateReport:
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

// deploy frequency

// deployFrequencyReport satisfy the interface `Report`
type deployFrequencyReport struct {
	AccumulatedTo int64                    `json:"accumulated_to"`
	Datapoints    deployFrequencyDataPoint `json:"datapoints"`
	FilePath      string
}

type deployFrequencyDataPoint struct {
	Daily   map[string]deployFrequency `json:"daily"`
	Weekly  map[string]deployFrequency `json:"weekly"`
	Monthly map[string]deployFrequency `json:"monthly"`
	Yearly  map[string]deployFrequency `json:"yearly"`
}

type deployFrequency struct {
	DeployCount float32 `json:"deploy_count"`
}

func (d *deployFrequencyReport) GetFilePath() string {
	return d.FilePath
}

func (d *deployFrequencyReport) PutFilePath(path string) {
	d.FilePath = path
}

func (d *deployFrequencyReport) Value(step model.InsightStep, key string) (float32, error) {
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

// change failure rate

// changeFailureRateReport satisfy the interface `Report`
type changeFailureRateReport struct {
	AccumulatedTo int64                      `json:"accumulated_to"`
	Datapoints    changeFailureRateDataPoint `json:"datapoints"`
	FilePath      string
}

type changeFailureRateDataPoint struct {
	Daily   map[string]changeFailureRate `json:"daily"`
	Weekly  map[string]changeFailureRate `json:"weekly"`
	Monthly map[string]changeFailureRate `json:"monthly"`
	Yearly  map[string]changeFailureRate `json:"yearly"`
}

type changeFailureRate struct {
	Rate         float32 `json:"rate"`
	SuccessCount int64   `json:"success_count"`
	FailureCount int64   `json:"failure_count"`
}

func (c *changeFailureRateReport) GetFilePath() string {
	return c.FilePath
}

func (c *changeFailureRateReport) PutFilePath(path string) {
	c.FilePath = path
}

func (c *changeFailureRateReport) Value(step model.InsightStep, key string) (float32, error) {
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
