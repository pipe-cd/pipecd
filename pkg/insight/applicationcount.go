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
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type ApplicationCounts struct {
	Counts    []model.InsightApplicationCount `json:"counts"`
	UpdatedAt int64                           `json:"updated_at"`
}

func MakeApplicationCounts(apps []*model.Application, now time.Time) ApplicationCounts {
	if len(apps) == 0 {
		return ApplicationCounts{
			UpdatedAt: now.Unix(),
		}
	}

	type key struct {
		status string
		kind   string
	}
	m := make(map[key]int)
	for _, app := range apps {
		k := key{
			status: model.ApplicationActiveStatus_ENABLED.String(),
			kind:   app.Kind.String(),
		}
		if app.Disabled {
			k.status = model.ApplicationActiveStatus_DISABLED.String()
		}
		m[k] = m[k] + 1
	}

	counts := make([]model.InsightApplicationCount, 0, len(m))
	for k, c := range m {
		counts = append(counts, model.InsightApplicationCount{
			Labels: map[string]string{
				model.InsightApplicationCountLabelKey_KIND.String():          k.kind,
				model.InsightApplicationCountLabelKey_ACTIVE_STATUS.String(): k.status,
			},
			Count: int32(c),
		})
	}

	return ApplicationCounts{
		Counts:    counts,
		UpdatedAt: now.Unix(),
	}
}

func determineApplicationStatus(app *model.Application) model.ApplicationActiveStatus {
	if app.Deleted {
		return model.ApplicationActiveStatus_DELETED
	}
	if app.Disabled {
		return model.ApplicationActiveStatus_DISABLED
	}
	return model.ApplicationActiveStatus_ENABLED
}
