// Copyright 2020 The Pipe Authors.
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

package datastore

import "github.com/kapetaniosci/pipe/pkg/model"

type ListFilter struct {
	Field    string
	Operator string
	Value    interface{}
}

type ListOption struct {
	Page     int
	PageSize int
	Filters  []ListFilter
}

type ApplicationStore interface {
	AddApplication(app *model.Application) error
	DisableApplication(id string) error
	ListApplications(opts ListOption) ([]model.Application, error)
}

type ApplicationResourceStore interface {
	GetApplicationResourceTree(id string) (*model.ApplicationResourceTree, error)
}

type PipelineStore interface {
	ListPipelines(opts ListOption) ([]model.Pipeline, error)
}
