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

package datastore

import (
	"context"

	"github.com/kapetaniosci/pipe/pkg/model"
)

type ListFilter struct {
	Field    string
	Operator string
	Value    interface{}
}

type ListOptions struct {
	Page     int
	PageSize int
	Filters  []ListFilter
}

type ProjectStore interface {
	AddProject(ctx context.Context, proj *model.Project) error
	ListProjects(ctx context.Context, opts ListOptions) ([]model.Project, error)
}

type EnvironmentStore interface {
	AddEnvironment(ctx context.Context, proj *model.Environment) error
	ListEnvironments(ctx context.Context, opts ListOptions) ([]model.Environment, error)
}

type PipedStore interface {
	AddPiped(ctx context.Context, proj *model.Piped) error
	ListPipeds(ctx context.Context, opts ListOptions) ([]model.Piped, error)
}

type ApplicationStore interface {
	AddApplication(ctx context.Context, app *model.Application) error
	DisableApplication(ctx context.Context, id string) error
	ListApplications(ctx context.Context, opts ListOptions) ([]model.Application, error)
}

type DeploymentStore interface {
	ListDeployments(ctx context.Context, opts ListOptions) ([]model.Deployment, error)
}

type CommandStore interface {
	AddCommand(ctx context.Context, proj *model.Command) error
	ListCommands(ctx context.Context, opts ListOptions) ([]model.Command, error)
}

type PipedStatsStore interface {
	AddPipedStats(ctx context.Context, proj *model.PipedStats) error
	ListPipedStatss(ctx context.Context, opts ListOptions) ([]model.PipedStats, error)
}
