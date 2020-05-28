package metrics

import (
	"context"
	"errors"

	"github.com/kapetaniosci/pipe/pkg/config"
)

var (
	ErrNoValuesFound = errors.New("no values found")
)

type Provider interface {
	Type() string
	// RunQuery runs the given query against the metrics provider,
	// and then checks if the results are expected or not.
	RunQuery(ctx context.Context, query string, expected config.AnalysisExpected) (result bool, err error)
}
