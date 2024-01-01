// Copyright 2024 The PipeCD Authors.
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
	"encoding/json"
	"fmt"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type commandCollection struct {
	requestedBy Commander
}

func (c *commandCollection) Kind() string {
	return "Command"
}

func (c *commandCollection) Factory() Factory {
	return func() interface{} {
		return &model.Command{}
	}
}

func (c *commandCollection) ListInUsedShards() []Shard {
	return []Shard{
		AgentShard,
		OpsShard,
	}
}

func (c *commandCollection) GetUpdatableShard() (Shard, error) {
	switch c.requestedBy {
	case PipedCommander:
		return AgentShard, nil
	case OpsCommander:
		return OpsShard, nil
	default:
		return "", ErrUnsupported
	}
}

func (c *commandCollection) Decode(e interface{}, parts map[Shard][]byte) error {
	if len(parts) != len(c.ListInUsedShards()) {
		return fmt.Errorf("failed while decode Command object: shards count not matched")
	}

	cmd, ok := e.(*model.Command)
	if !ok {
		return fmt.Errorf("failed while decode Command object: type not matched")
	}

	var (
		status    model.CommandStatus
		updatedAt int64
	)
	for _, p := range parts {
		if err := json.Unmarshal(p, &cmd); err != nil {
			return err
		}
		if updatedAt < cmd.UpdatedAt {
			updatedAt = cmd.UpdatedAt
		}
		if status < cmd.Status {
			status = cmd.Status
		}
	}

	cmd.Status = status
	cmd.UpdatedAt = updatedAt
	return nil
}

func (c *commandCollection) Encode(e interface{}) (map[Shard][]byte, error) {
	const errFmt = "failed while encode Command object: %s"

	me, ok := e.(*model.Command)
	if !ok {
		return nil, fmt.Errorf(errFmt, "type not matched")
	}

	agentShardStruct := me
	adata, err := json.Marshal(&agentShardStruct)
	if err != nil {
		return nil, fmt.Errorf(errFmt, "unable to marshal entity data")
	}

	opsShardStruct := model.Command{
		// Required fields to pass validation on update.
		Id:        me.Id,
		PipedId:   me.PipedId,
		CreatedAt: me.CreatedAt,
		UpdatedAt: me.UpdatedAt,
		// Fields which in both shards, but the higher value has
		// a higher priority on merge.
		Status: me.Status,
	}
	odata, err := json.Marshal(&opsShardStruct)
	if err != nil {
		return nil, fmt.Errorf(errFmt, "unable to marshal entity data")
	}

	return map[Shard][]byte{
		AgentShard: adata,
		OpsShard:   odata,
	}, nil
}

type CommandStore interface {
	Add(ctx context.Context, cmd *model.Command) error
	Get(ctx context.Context, id string) (*model.Command, error)
	List(ctx context.Context, opts ListOptions) ([]*model.Command, error)
	UpdateStatus(ctx context.Context, id string, status model.CommandStatus, metadata map[string]string, handledAt int64) error
}

type commandStore struct {
	backend
	commander Commander
	nowFunc   func() time.Time
}

func NewCommandStore(ds DataStore, c Commander) CommandStore {
	return &commandStore{
		backend: backend{
			ds:  ds,
			col: &commandCollection{requestedBy: c},
		},
		commander: c,
		nowFunc:   time.Now,
	}
}

func (s *commandStore) Add(ctx context.Context, cmd *model.Command) error {
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
	return s.ds.Create(ctx, s.col, cmd.Id, cmd)
}

func (s *commandStore) Get(ctx context.Context, id string) (*model.Command, error) {
	var entity model.Command
	if err := s.ds.Get(ctx, s.col, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *commandStore) List(ctx context.Context, opts ListOptions) ([]*model.Command, error) {
	it, err := s.ds.Find(ctx, s.col, opts)
	if err != nil {
		return nil, err
	}
	cmds := make([]*model.Command, 0)
	for {
		var cmd model.Command
		err := it.Next(&cmd)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, &cmd)
	}
	return cmds, nil
}

func (s *commandStore) update(ctx context.Context, id string, updater func(piped *model.Command) error) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		p := e.(*model.Command)
		if err := updater(p); err != nil {
			return err
		}
		p.UpdatedAt = now
		return p.Validate()
	})
}

func (s *commandStore) UpdateStatus(ctx context.Context, id string, status model.CommandStatus, metadata map[string]string, handledAt int64) error {
	return s.update(ctx, id, func(c *model.Command) error {
		c.Status = status
		c.Metadata = metadata
		c.HandledAt = handledAt
		return nil
	})
}
