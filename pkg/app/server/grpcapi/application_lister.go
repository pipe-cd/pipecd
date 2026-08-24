// Copyright 2026 The PipeCD Authors.
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

package grpcapi

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// listApplicationsPageSize is the number of applications fetched per datastore
// query while aggregating the full application list.
const listApplicationsPageSize = 100

// applicationLister is the subset of the application stores needed to page
// through applications.
type applicationLister interface {
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Application, string, error)
}

// listAllApplications pages through the datastore so that fetching every
// application matching opts cannot be done by one unbounded query. Cursor
// paging requires a stable order, so callers must set Orders on opts; the
// returned slice aggregates all pages because neither RPC response carries a
// cursor field.
func listAllApplications(ctx context.Context, store applicationLister, opts datastore.ListOptions) ([]*model.Application, error) {
	apps := make([]*model.Application, 0, listApplicationsPageSize)
	for {
		page, cursor, err := store.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		apps = append(apps, page...)
		if cursor == "" {
			return apps, nil
		}
		opts.Cursor = cursor
	}
}
