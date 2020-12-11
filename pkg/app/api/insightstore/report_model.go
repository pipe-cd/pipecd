package insightstore

import (
	"encoding/json"
	"fmt"
)

type commonReport struct {
	AccumulatedTo int64      `json:"accumulated_to"`
	Datapoints    datapoints `json:"datapoints"`
}

type datapoints struct {
	Daily   json.RawMessage `json:"daily"`
	Weekly  json.RawMessage `json:"weekly"`
	Monthly json.RawMessage `json:"monthly"`
	Yearly  json.RawMessage `json:"yearly"`
}

type datapoint interface {
	Value() float32
}

func toDatapoint(i interface{}) (map[string]datapoint, error) {
	datapoints := map[string]datapoint{}
	switch p := i.(type) {
	case map[string]deployFrequency:
		for k, v := range p {
			datapoints[k] = v
		}
		return datapoints, nil
	case map[string]changeFailureRate:
		for k, v := range p {
			datapoints[k] = v
		}
		return datapoints, nil
	default:
		return nil, fmt.Errorf("cannot convert to datapoint: %v", p)
	}

}

// deployFrequency satisfy the interface `datapoint`
type deployFrequency struct {
	DeployCount float32 `json:"deploy_count"`
}

func (d deployFrequency) Value() float32 {
	return d.DeployCount
}

// changeFailureRate satisfy the interface `datapoint`
type changeFailureRate struct {
	Rate         float32 `json:"rate"`
	SuccessCount int64   `json:"success_count"`
	FailureCount int64   `json:"failure_count"`
}

func (c changeFailureRate) Value() float32 {
	return c.Rate
}
