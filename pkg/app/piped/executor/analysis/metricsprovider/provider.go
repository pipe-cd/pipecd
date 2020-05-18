package metricsprovider

import "errors"

var (
	ErrNoValuesFound = errors.New("no values found")
)

type Provider interface {
	Type() string
	RunQuery(query string) (float64, error)
}
