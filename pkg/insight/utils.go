// Copyright 2022 The PipeCD Authors.
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
	"sort"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

// NormalizeUnixTime ignores hour, minute, second and nanosecond
func NormalizeUnixTime(t int64, loc *time.Location) int64 {
	tt := time.Unix(t, 0).In(loc)
	return time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, loc).Unix()
}

func GroupDeploymentsByDaily(deployments []*model.InsightDeployment, loc *time.Location) [][]*model.InsightDeployment {
	dailyDeployments := make(map[int64][]*model.InsightDeployment)

	for _, d := range deployments {
		t := NormalizeUnixTime(d.CompletedAt, loc)
		dailyDeployments[t] = append(dailyDeployments[t], d)
	}

	keys := make([]int64, len(dailyDeployments))
	idx := 0
	for key := range dailyDeployments {
		keys[idx] = key
		idx++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	result := make([][]*model.InsightDeployment, len(dailyDeployments))
	for idx, key := range keys {
		result[idx] = dailyDeployments[key]
	}

	return result
}
