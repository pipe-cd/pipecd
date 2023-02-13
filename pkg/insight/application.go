// Copyright 2023 The PipeCD Authors.
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
	"github.com/pipe-cd/pipecd/pkg/model"
)

type ApplicationData struct {
	Id     string            `json:"id"`
	Labels map[string]string `json:"labels"`
	Kind   string            `json:"kind"`
	Status string            `json:"status"`
}

type ProjectApplicationData struct {
	Applications []*ApplicationData `json:"applications"`
	UpdatedAt    int64              `json:"updated_at"`
}

func BuildApplicationData(a *model.Application) ApplicationData {
	status := determineApplicationStatus(a)

	return ApplicationData{
		Id:     a.Id,
		Labels: a.Labels,
		Kind:   a.Kind.String(),
		Status: status.String(),
	}
}

func BuildProjectApplicationData(apps []*ApplicationData, updatedAt int64) ProjectApplicationData {
	return ProjectApplicationData{
		Applications: apps,
		UpdatedAt:    updatedAt,
	}
}

type ApplicationCounts struct {
	Counts    []model.InsightApplicationCount
	UpdatedAt int64
}
