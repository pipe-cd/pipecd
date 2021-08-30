package prometheus

import (
	"context"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type fakeClient struct {
	value    model.Value
	err      error
	warnings v1.Warnings
}

func (f fakeClient) QueryRange(_ context.Context, _ string, _ v1.Range) (model.Value, v1.Warnings, error) {
	if f.err != nil {
		return nil, f.warnings, f.err
	}
	return f.value, f.warnings, nil
}
