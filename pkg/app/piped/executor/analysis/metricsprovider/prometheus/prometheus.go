package prometheus

import (
	"net/url"
	"time"
)

type Provider struct {
	timeout  time.Duration
	url      url.URL
	username string
	password string
}

// response represents a response from prometheus server.
type response struct {
	Data struct {
		Result []struct {
			Metric struct {
				Name string `json:"name"`
			}
			Value []interface{} `json:"value"`
		}
	}
}

func NewPrometheusProvider() (*Provider, error) {
	return &Provider{}, nil
}

// RunQuery executes the promQL query and returns the the first result as float64.
func (p *Provider) RunQuery(query string) (float64, error) {
	return 0, nil
}
