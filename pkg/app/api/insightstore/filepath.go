package insightstore

import (
	"fmt"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

// insight file paths according to the following format
//
// insights
//  ├─ project-id
//    ├─ deployment-frequency
//        ├─ project  # aggregated from all applications
//            ├─ years.json
//            ├─ 2020-01.json
//            ├─ 2020-02.json
//            ...
//        ├─ app-id
//            ├─ years.json
//            ├─ 2020-01.json
//            ├─ 2020-02.json
//            ...
func newYearlyFilePath(projectID string, metricsKind model.InsightMetricsKind, appID string) string {
	metricsKindKebab := getKebabCaseMetricsKind(metricsKind)
	return fmt.Sprintf("insights/%s/%s/%s/years.json", projectID, metricsKindKebab, appID)
}
func newMonthlyFilePath(projectID string, metricsKind model.InsightMetricsKind, appID string, month string) string {
	metricsKindKebab := getKebabCaseMetricsKind(metricsKind)
	return fmt.Sprintf("insights/%s/%s/%s/%s.json", projectID, metricsKindKebab, appID, month)
}

func searchFilePaths(projectID string, appID string, from time.Time, dataPointCount int, metricsKind model.InsightMetricsKind, step model.InsightStep) []string {
	if appID == "" {
		appID = "project"
	}
	switch step {
	case model.InsightStep_YEARLY:
		return []string{newYearlyFilePath(projectID, metricsKind, appID)}
	default:
		months := getPointsMonths(from, dataPointCount, step)
		var paths []string
		for _, m := range months {
			path := newMonthlyFilePath(projectID, metricsKind, appID, m)
			paths = append(paths, path)
		}
		return paths
	}
}

func getKebabCaseMetricsKind(kind model.InsightMetricsKind) string {
	var kebabKind string
	switch kind {
	case model.InsightMetricsKind_DEPLOYMENT_FREQUENCY:
		kebabKind = "deployment_frequency"
	case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
		kebabKind = "change_failure_rate"
	case model.InsightMetricsKind_MTTR:
		kebabKind = "mean_time_to_restore"
	case model.InsightMetricsKind_LEAD_TIME:
		kebabKind = "lead_time"
	}
	return kebabKind
}

// getPointsMonths return months between two dates.
func getPointsMonths(date time.Time, count int, step model.InsightStep) []string {
	var to time.Time

	switch step {
	case model.InsightStep_YEARLY:
		to = date.AddDate(count-1, 0, 0)
	case model.InsightStep_MONTHLY:
		to = date.AddDate(0, count-1, 0)
	case model.InsightStep_WEEKLY:
		to = date.AddDate(0, 0, (count-1)*7)
	case model.InsightStep_DAILY:
		to = date.AddDate(0, 0, count-1)
	}

	var months []string
	y1, m1, _ := to.Date()
	for {
		// 2015-05-05 08:05:15.828452891 +0900 UST → 2015-05
		months = append(months, date.Format("2006-01"))
		y2, m2, _ := date.Date()
		if y1 == y2 && m1 == m2 {
			return months
		}

		date = date.AddDate(0, 1, 0)
	}
}
