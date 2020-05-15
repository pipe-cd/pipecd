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
	"time"

	"github.com/kapetaniosci/pipe/pkg/model"
)

const applicationModelKind = "Application"

var applicationFactory = func() interface{} {
	return &model.Application{}
}

type ApplicationStore interface {
	AddApplication(ctx context.Context, app *model.Application) error
	DisableApplication(ctx context.Context, id string) error
	ListApplications(ctx context.Context, opts ListOptions) ([]model.Application, error)
}

type applicationStore struct {
	backend
	nowFunc func() time.Time
}

func NewApplicationStore(ds DataStore) ApplicationStore {
	return &applicationStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *applicationStore) AddApplication(ctx context.Context, app *model.Application) error {
	now := s.nowFunc().Unix()
	if app.CreatedAt == 0 {
		app.CreatedAt = now
	}
	if app.UpdatedAt == 0 {
		app.UpdatedAt = now
	}
	if err := app.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, applicationModelKind, app.Id, app)
}

func (s *applicationStore) DisableApplication(ctx context.Context, id string) error {
	return s.ds.Update(ctx, applicationModelKind, id, applicationFactory, func(e interface{}) error {
		app := e.(*model.Application)
		app.Disabled = true
		app.UpdatedAt = time.Now().Unix()
		return nil
	})
}

func (s *applicationStore) ListApplications(ctx context.Context, opts ListOptions) ([]model.Application, error) {
	it, err := s.ds.Find(ctx, applicationModelKind, opts)
	if err != nil {
		return nil, err
	}
	apps := make([]model.Application, 0)
	for {
		var app model.Application
		err := it.Next(&app)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}
