// Copyright 2021 The PipeCD Authors.
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

	"github.com/pipe-cd/pipe/pkg/model"
)

type ApplicationCounts struct {
	Counts        []model.InsightApplicationCount `json:"counts"`
	AccumulatedTo int64                           `json:"accumulated_to"`
}

func MakeApplicationCounts(apps []*model.Application, now time.Time) (*ApplicationCounts, error) {
	return nil, nil
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
