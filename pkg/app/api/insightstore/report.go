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

// deployFrequencyDataPoint satisfy the interface `Report`
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

// changeFailureRateDataPoint satisfy the interface `Report`
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

type Report interface {
	// Value get data by step and key
	Value(step model.InsightStep, key string) (float32, error)
}

func toDatapoint(i interface{}) (Report, error) {
	switch p := i.(type) {
	case deployFrequencyReport:
		return p, nil
	case changeFailureRateReport:
		return p, nil
	default:
		return nil, fmt.Errorf("cannot convert to Report: %v", p)
	}

}

func (d deployFrequencyReport) Value(step model.InsightStep, key string) (float32, error) {
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

func (c changeFailureRateReport) Value(step model.InsightStep, key string) (float32, error) {
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
