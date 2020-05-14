package metricsprovider

import "errors"

var (
	ErrNoValuesFound = errors.New("no values found")
)

type Provider interface {
	// RunQuery executes the query and converts the first result to float64
	RunQuery(query string) (float64, error)
}
