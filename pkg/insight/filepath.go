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

package insight

import (
	"fmt"
	"strings"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

// insight file paths according to the following format.
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
func MakeYearsFilePath(projectID string, metricsKind model.InsightMetricsKind, appID string) string {
	k := strings.ToLower(metricsKind.String())
	return fmt.Sprintf("insights/%s/%s/%s/years.json", projectID, k, appID)
}

func MakeChunkFilePath(projectID string, metricsKind model.InsightMetricsKind, appID string, month string) string {
	k := strings.ToLower(metricsKind.String())
	return fmt.Sprintf("insights/%s/%s/%s/%s.json", projectID, k, appID, month)
}

func DetermineFilePaths(projectID string, appID string, kind model.InsightMetricsKind, step model.InsightStep, from time.Time, count int) []string {
	if appID == "" {
		appID = "project"
	}
	switch step {
	case model.InsightStep_YEARLY:
		return []string{MakeYearsFilePath(projectID, kind, appID)}
	default:
		keys := determineChunkKeys(step, from, count)
		var paths []string
		for _, k := range keys {
			path := MakeChunkFilePath(projectID, kind, appID, k)
			paths = append(paths, path)
		}
		return paths
	}
}

// determineChunkKeys returns a sorted list of chunk keys needed for a given time range.
func determineChunkKeys(step model.InsightStep, from time.Time, count int) []string {
	var to time.Time

	switch step {
	case model.InsightStep_YEARLY:
		to = from.AddDate(count-1, 0, 0)
	case model.InsightStep_MONTHLY:
		to = from.AddDate(0, count-1, 0)
	case model.InsightStep_WEEKLY:
		to = from.AddDate(0, 0, (count-1)*7)
	case model.InsightStep_DAILY:
		to = from.AddDate(0, 0, count-1)
	}

	from = time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	to = time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.UTC)

	var keys []string
	cur := from
	for !cur.After(to) {
		keys = append(keys, cur.Format("2006-01"))
		cur = cur.AddDate(0, 1, 0)
	}
	return keys
}
