package metric

import (
	"context"
	"errors"
)

var (
	ErrNoValuesFound = errors.New("no values found")
)

type Provider interface {
	Type() string
	// RunQuery runs the given query against the metrics provider,
	// and then checks if the results are expected or not.
	RunQuery(ctx context.Context, query, expected string) (result bool, err error)
}
