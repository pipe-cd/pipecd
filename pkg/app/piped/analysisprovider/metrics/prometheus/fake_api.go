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

func (m fakeAPI) Query(_ context.Context, _ string, _ time.Time) (model.Value, v1.Warnings, error) {
	if m.err != nil {
		return nil, m.warnings, m.err
	}
	return m.value, m.warnings, nil
}

func (m fakeAPI) QueryRange(_ context.Context, _ string, _ v1.Range) (model.Value, v1.Warnings, error) {
	if m.err != nil {
		return nil, m.warnings, m.err
	}
	return m.value, m.warnings, nil
}

// Below methods are required to meet the interface.

func (y fakeAPI) Alerts(ctx context.Context) (v1.AlertsResult, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) AlertManagers(ctx context.Context) (v1.AlertManagersResult, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) CleanTombstones(ctx context.Context) error {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) Config(ctx context.Context) (v1.ConfigResult, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) DeleteSeries(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) error {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) Flags(ctx context.Context) (v1.FlagsResult, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) LabelNames(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) ([]string, v1.Warnings, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) LabelValues(ctx context.Context, label string, matches []string, startTime time.Time, endTime time.Time) (model.LabelValues, v1.Warnings, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) QueryExemplars(ctx context.Context, query string, startTime time.Time, endTime time.Time) ([]v1.ExemplarQueryResult, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) Buildinfo(ctx context.Context) (v1.BuildinfoResult, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) Runtimeinfo(ctx context.Context) (v1.RuntimeinfoResult, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) Series(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) ([]model.LabelSet, v1.Warnings, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) Snapshot(ctx context.Context, skipHead bool) (v1.SnapshotResult, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) Rules(ctx context.Context) (v1.RulesResult, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) Targets(ctx context.Context) (v1.TargetsResult, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) TargetsMetadata(ctx context.Context, matchTarget string, metric string, limit string) ([]v1.MetricMetadata, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) Metadata(ctx context.Context, metric string, limit string) (map[string][]v1.Metadata, error) {
	panic("this method doesn't expect to be called")
}

func (y fakeAPI) TSDB(ctx context.Context) (v1.TSDBResult, error) {
	panic("this method doesn't expect to be called")
}
