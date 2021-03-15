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
	"errors"

	"github.com/pipe-cd/pipe/pkg/model"
)

type ApplicationStatus string

var (
	ApplicationStatusUnknown ApplicationStatus = "unknown"
	ApplicationStatusEnable  ApplicationStatus = "enable"
	ApplicationStatusDisable ApplicationStatus = "disable"
	ApplicationStatusDeleted ApplicationStatus = "deleted"
)

var statuses = []ApplicationStatus{ApplicationStatusEnable, ApplicationStatusDisable, ApplicationStatusDeleted}

type ApplicationCount struct {
	Counts          []ApplicationCountByLabelSet `json:"counts"`
	AccumulatedFrom int64                        `json:"accumulated_from"`
	AccumulatedTo   int64                        `json:"accumulated_to"`
}

type ApplicationCountByLabelSet struct {
	LabelSet ApplicationCountLabelSet `json:"label_set"`
	Count    int                      `json:"count"`
}

type ApplicationCountLabelSet struct {
	// KUBERNETES, TERRAFORM, CLOUDRUN...
	Kind model.ApplicationKind `json:"kind"`
	// enable, disable or deleted
	Status ApplicationStatus `json:"status"`
}

func NewApplicationCount() *ApplicationCount {
	counts := make([]ApplicationCountByLabelSet, len(model.ApplicationKind_name)*len(statuses))
	for _, k := range model.ApplicationKind_value {
		for j, s := range statuses {
			counts[int(k)*len(statuses)+j] = ApplicationCountByLabelSet{
				LabelSet: ApplicationCountLabelSet{
					Kind:   model.ApplicationKind(k),
					Status: s,
				},
			}
		}
	}
	return &ApplicationCount{
		Counts: counts,
	}
}

// MigrateApplicationCount add new labelset on count.
func (a *ApplicationCount) MigrateApplicationCount() {
	new := NewApplicationCount()
	for _, c := range new.Counts {
		if _, err := a.Find(c.LabelSet); err != nil {
			if err == ErrCountNotFound {
				a.Counts = append(a.Counts, c)
			}
		}
	}
}

var ErrCountNotFound = errors.New("error application count by label set not found")

// Find finds the count by labelset
func (a *ApplicationCount) Find(labelSet ApplicationCountLabelSet) (ApplicationCountByLabelSet, error) {
	for _, c := range a.Counts {
		if c.LabelSet == labelSet {
			return c, nil
		}
	}
	return ApplicationCountByLabelSet{}, ErrCountNotFound
}

// UpdateCount update the count
func (a *ApplicationCount) UpdateCount(apps []*model.Application) {
	// init appmac
	appmap := map[ApplicationStatus]map[model.ApplicationKind]int{}
	for _, s := range statuses {
		appmap[s] = map[model.ApplicationKind]int{}
	}
	// classify and aggregate applications
	for _, app := range apps {
		s := determineApplicationStatus(app)
		for _, k := range model.ApplicationKind_value {
			kind := model.ApplicationKind(k)
			if kind == app.Kind {
				appmap[s][kind]++
			}
		}
	}

	for i := 0; i < len(a.Counts); i++ {
		c := &a.Counts[i]
		c.Count = appmap[c.LabelSet.Status][c.LabelSet.Kind]
	}
}

// determineApplicationStatus uniquely determine the application status
func determineApplicationStatus(app *model.Application) ApplicationStatus {
	if app.Deleted {
		return ApplicationStatusDeleted
	}
	if app.Disabled {
		return ApplicationStatusDisable
	}
	if app.Deploying {
		return ApplicationStatusEnable
	}
	return ApplicationStatusUnknown
}
