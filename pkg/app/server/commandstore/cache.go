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
	return &s, nil
}

func (c *commandCache) Put(commandID string, command *model.Command) error {
	key := cacheKey(commandID)
	data, err := json.Marshal(command)
	if err != nil {
		return err
	}
	return c.backend.Put(key, data)
}

func cacheKey(commandID string) string {
	return fmt.Sprintf("command:%s", commandID)
}
