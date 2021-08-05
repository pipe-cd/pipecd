package prometheus

import (
	"context"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type fakeAPI struct {
	value    model.Value
	err      error
	warnings v1.Warnings
}

func (f fakeAPI) Query(_ context.Context, _ string, _ time.Time) (model.Value, v1.Warnings, error) {
	if f.err != nil {
		return nil, f.warnings, f.err
	}
	return f.value, f.warnings, nil
}

func (f fakeAPI) QueryRange(_ context.Context, _ string, _ v1.Range) (model.Value, v1.Warnings, error) {
	if f.err != nil {
		return nil, f.warnings, f.err
	}
	return f.value, f.warnings, nil
}

// Below methods are required to meet the interface.

func (f fakeAPI) Alerts(ctx context.Context) (v1.AlertsResult, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) AlertManagers(ctx context.Context) (v1.AlertManagersResult, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) CleanTombstones(ctx context.Context) error {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) Config(ctx context.Context) (v1.ConfigResult, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) DeleteSeries(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) error {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) Flags(ctx context.Context) (v1.FlagsResult, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) LabelNames(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) ([]string, v1.Warnings, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) LabelValues(ctx context.Context, label string, matches []string, startTime time.Time, endTime time.Time) (model.LabelValues, v1.Warnings, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) QueryExemplars(ctx context.Context, query string, startTime time.Time, endTime time.Time) ([]v1.ExemplarQueryResult, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) Buildinfo(ctx context.Context) (v1.BuildinfoResult, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) Runtimeinfo(ctx context.Context) (v1.RuntimeinfoResult, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) Series(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) ([]model.LabelSet, v1.Warnings, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) Snapshot(ctx context.Context, skipHead bool) (v1.SnapshotResult, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) Rules(ctx context.Context) (v1.RulesResult, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) Targets(ctx context.Context) (v1.TargetsResult, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) TargetsMetadata(ctx context.Context, matchTarget string, metric string, limit string) ([]v1.MetricMetadata, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) Metadata(ctx context.Context, metric string, limit string) (map[string][]v1.Metadata, error) {
	panic("this method doesn't expect to be called")
}

func (f fakeAPI) TSDB(ctx context.Context) (v1.TSDBResult, error) {
	panic("this method doesn't expect to be called")
}
