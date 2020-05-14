package datadog

import "time"

type Provider struct {
	metricsQueryEndpoint     string
	apiKeyValidationEndpoint string

	timeout        time.Duration
	apiKey         string
	applicationKey string
	fromDelta      int64
}

func NewProvider() (*Provider, error) {
	return &Provider{}, nil
}

// response represents a response from datadog server.
type response struct {
	Series []struct {
		Pointlist [][]float64 `json:"pointlist"`
	}
}

// RunQuery executes the datadog query against datadog endpoint
// and returns the the first result as float64.
func (p *Provider) RunQuery(query string) (float64, error) {
	return 0, nil
}
