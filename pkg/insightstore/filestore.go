// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package insightstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Store struct {
	filestore filestore.Store
}

func NewInsightFileStore(fs filestore.Store) Store {
	return Store{filestore: fs}
}

// GetReports returns data as report.
func (s *Store) GetReports(
	ctx context.Context,
	projectID string,
	appID string,
	kind model.InsightMetricsKind,
	step model.InsightStep,
	from time.Time,
	dataPointCount int) ([]Report, error) {
	from = formatFrom(from, step)

	paths := determineFilePaths(projectID, appID, kind, step, from, dataPointCount)
	var reports []Report
	for _, p := range paths {
		r, err := s.getReport(ctx, p, kind)
		if err != nil {
			return nil, err
		}

		reports = append(reports, r)

	}

	return reports, nil
}

// List returns data as insight data point.
func (s *Store) List(
	ctx context.Context,
	projectID string,
	appID string,
	kind model.InsightMetricsKind,
	step model.InsightStep,
	from time.Time,
	dataPointCount int) ([]*model.InsightDataPoint, error) {
	from = formatFrom(from, step)

	reports, err := s.GetReports(ctx, projectID, appID, kind, step, from, dataPointCount)
	if err != nil {
		return nil, err
	}
	var idps []*model.InsightDataPoint
	for _, r := range reports {
		idp, err := convertToInsightDataPoints(r, from, dataPointCount, step)
		if err != nil {
			return nil, err
		}

		idps = append(idps, idp...)

		dataPointCount = dataPointCount - r.DataCount(step)

		nextMonth := time.Date(from.Year(), from.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		from = formatFrom(nextMonth, step)
		if step == model.InsightStep_WEEKLY && from.Month() != nextMonth.Month() {
			from = from.AddDate(0, 0, 7)
		}
	}

	return idps, nil
}

// Put create of update report.
func (s *Store) Put(ctx context.Context, report Report) error {
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}
	path := report.GetFilePath()
	if path == "" {
		return fmt.Errorf("filepath not found on report struct")
	}
	return s.filestore.PutObject(ctx, path, data)
}

func (s *Store) getReport(ctx context.Context, path string, kind model.InsightMetricsKind) (Report, error) {
	obj, err := s.filestore.GetObject(ctx, path)
	if err != nil {
		return nil, err
	}

	var report Report
	switch kind {
	case model.InsightMetricsKind_DEPLOYMENT_FREQUENCY:
		var df DeployFrequencyReport
		err := json.Unmarshal(obj.Content, &df)
		if err != nil {
			return nil, err
		}
		report, err = toReport(&df)
		if err != nil {
			return nil, err
		}
	case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
		var cfr ChangeFailureRateReport
		err := json.Unmarshal(obj.Content, &cfr)
		if err != nil {
			return nil, err
		}
		report, err = toReport(&cfr)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unimpremented insight kind: %s", kind)
	}

	report.PutFilePath(path)
	return report, nil
}

func formatFrom(from time.Time, step model.InsightStep) time.Time {
	var formattedTime time.Time
	switch step {
	case model.InsightStep_DAILY:
		formattedTime = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	case model.InsightStep_WEEKLY:
		// Sunday in the week of rangeFrom
		sunday := from.AddDate(0, 0, -int(from.Weekday()))
		formattedTime = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 0, 0, 0, 0, time.UTC)
	case model.InsightStep_MONTHLY:
		formattedTime = time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	case model.InsightStep_YEARLY:
		formattedTime = time.Date(from.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return formattedTime
}
