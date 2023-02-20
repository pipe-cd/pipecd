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

package applicationlivestatestore

import (
	"encoding/json"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLiveStateCache struct {
	backend cache.Cache
}

func (c *applicationLiveStateCache) Get(applicationID string) (*model.ApplicationLiveStateSnapshot, error) {
	key := cacheKey(applicationID)
	item, err := c.backend.Get(key)
	if err != nil {
		return nil, err
	}
	var s model.ApplicationLiveStateSnapshot
	if err := json.Unmarshal(item.([]byte), &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *applicationLiveStateCache) Put(applicationID string, als *model.ApplicationLiveStateSnapshot) error {
	key := cacheKey(applicationID)
	data, err := json.Marshal(als)
	if err != nil {
		return err
	}
	return c.backend.Put(key, data)
}

func (c *applicationLiveStateCache) Delete(applicationID string) error {
	key := cacheKey(applicationID)
	return c.backend.Delete(key)
}

func cacheKey(applicationID string) string {
	return fmt.Sprintf("application-live-state:%s", applicationID)
}
