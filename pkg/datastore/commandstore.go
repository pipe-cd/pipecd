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

const commandModelKind = "Command"

var commandFactory = func() interface{} {
	return &model.Command{}
}

type CommandStore interface {
	AddCommand(ctx context.Context, cmd *model.Command) error
	ListCommands(ctx context.Context, opts ListOptions) ([]model.Command, error)
}

type commandStore struct {
	backend
	nowFunc func() time.Time
}

func NewCommandStore(ds DataStore) CommandStore {
	return &commandStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *commandStore) AddCommand(ctx context.Context, cmd *model.Command) error {
	now := s.nowFunc().Unix()
	if cmd.CreatedAt == 0 {
		cmd.CreatedAt = now
	}
	if cmd.UpdatedAt == 0 {
		cmd.UpdatedAt = now
	}
	if err := cmd.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, commandModelKind, cmd.Id, cmd)
}

func (s *commandStore) ListCommands(ctx context.Context, opts ListOptions) ([]model.Command, error) {
	it, err := s.ds.Find(ctx, commandModelKind, opts)
	if err != nil {
		return nil, err
	}
	cmds := make([]model.Command, 0)
	for {
		var cmd model.Command
		err := it.Next(&cmd)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}
