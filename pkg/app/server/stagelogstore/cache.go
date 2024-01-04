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

package stagelogstore

import (
	"encoding/json"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/cache"
)

type stageLogCache struct {
	cache cache.Cache
}

func (c *stageLogCache) Get(deploymentID, stageID string, retriedCount int32) (logFragment, error) {
	key := cacheKey(deploymentID, stageID, retriedCount)
	var lf logFragment
	item, err := c.cache.Get(key)
	if err != nil {
		return lf, err
	}
	return lf, json.Unmarshal(item.([]byte), &lf)
}

func (c *stageLogCache) Put(deploymentID, stageID string, retriedCount int32, lf *logFragment) error {
	key := cacheKey(deploymentID, stageID, retriedCount)
	data, err := json.Marshal(lf)
	if err != nil {
		return err
	}
	return c.cache.Put(key, data)
}

func (c *stageLogCache) Delete(deploymentID, stageID string, retriedCount int32) error {
	key := cacheKey(deploymentID, stageID, retriedCount)
	return c.cache.Delete(key)
}

func cacheKey(deploymentID, stageID string, retriedCount int32) string {
	return fmt.Sprintf("piped-stage-log:%s:%s:%d", deploymentID, stageID, retriedCount)
}
