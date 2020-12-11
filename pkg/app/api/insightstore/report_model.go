package insightstore

import (
	"fmt"

	"github.com/pipe-cd/pipe/pkg/model"
)

// deploy frequency

type deployFrequencyReport struct {
	AccumulatedTo int64                    `json:"accumulated_to"`
	Datapoints    deployFrequencyDataPoint `json:"datapoints"`
}

// deployFrequencyDataPoint satisfy the interface `datapoint`
type deployFrequencyDataPoint struct {
	Daily   map[string]deployFrequency `json:"daily"`
	Weekly  map[string]deployFrequency `json:"weekly"`
	Monthly map[string]deployFrequency `json:"monthly"`
	Yearly  map[string]deployFrequency `json:"yearly"`
}

type deployFrequency struct {
	DeployCount float32 `json:"deploy_count"`
}

// change failure rate
type changeFailureRateReport struct {
	AccumulatedTo int64                      `json:"accumulated_to"`
	Datapoints    changeFailureRateDataPoint `json:"datapoints"`
}

// changeFailureRateDataPoint satisfy the interface `datapoint`
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

type datapoint interface {
	// Value get data by step and key
	Value(step model.InsightStep, key string) (float32, error)
}

func toDatapoint(i interface{}) (datapoint, error) {
	switch p := i.(type) {
	case deployFrequencyDataPoint:
		return p, nil
	case changeFailureRateDataPoint:
		return p, nil
	default:
		return nil, fmt.Errorf("cannot convert to datapoint: %v", p)
	}

}

func (d deployFrequencyDataPoint) Value(step model.InsightStep, key string) (float32, error) {
	switch step {
	case model.InsightStep_YEARLY:
		return d.Yearly[key].DeployCount, nil
	case model.InsightStep_MONTHLY:
		return d.Monthly[key].DeployCount, nil
	case model.InsightStep_WEEKLY:
		return d.Weekly[key].DeployCount, nil
	case model.InsightStep_DAILY:
		return d.Daily[key].DeployCount, nil
	}
	return 0, fmt.Errorf("value not found. step: %d, key: %s", step, key)
}

func (c changeFailureRateDataPoint) Value(step model.InsightStep, key string) (float32, error) {
	switch step {
	case model.InsightStep_YEARLY:
		return c.Yearly[key].Rate, nil
	case model.InsightStep_MONTHLY:
		return c.Monthly[key].Rate, nil
	case model.InsightStep_WEEKLY:
		return c.Weekly[key].Rate, nil
	case model.InsightStep_DAILY:
		return c.Daily[key].Rate, nil
	}
	return 0, fmt.Errorf("value not found. step: %d, key: %s", step, key)
}
