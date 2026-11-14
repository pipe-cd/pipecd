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

package commandstore

import (
	"encoding/json"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type commandCache struct {
	backend cache.Cache
}

func (c *commandCache) Get(commandID string) (*model.Command, error) {
	key := cacheKey(commandID)
	item, err := c.backend.Get(key)
	if err != nil {
		return nil, err
	}
	var s model.Command
	if err := json.Unmarshal(item.([]byte), &s); err != nil {
		return nil, err
	}
	// A cached entry that doesn't satisfy the same validation required to
	// store a command (e.g. a zero-value Command decoded from a JSON null
	// written before Put started rejecting nil commands) must not be
	// treated as a valid cache hit. Evict it on a best-effort basis and
	// report a miss so the caller falls back to the datastore.
	if err := s.Validate(); err != nil {
		_ = c.backend.Delete(key)
		return nil, cache.ErrNotFound
	}
	return &s, nil
}

func (c *commandCache) Put(commandID string, command *model.Command) error {
	if command == nil {
		return fmt.Errorf("command must not be nil")
	}
	if err := command.Validate(); err != nil {
		return fmt.Errorf("command is invalid: %w", err)
	}
	key := cacheKey(commandID)
	data, err := json.Marshal(command)
	if err != nil {
		return err
	}
	return c.backend.Put(key, data)
}

func (c *commandCache) Delete(commandID string) error {
	return c.backend.Delete(cacheKey(commandID))
}

func cacheKey(commandID string) string {
	return fmt.Sprintf("command:%s", commandID)
}
