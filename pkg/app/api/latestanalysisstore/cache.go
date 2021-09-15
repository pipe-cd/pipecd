// Copyright 2021 The PipeCD Authors.
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

package latestanalysisstore

import (
	"encoding/json"
	"fmt"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/model"
)

type analysisCache struct {
	backend cache.Cache
}

func (c *analysisCache) Get(applicationID string) (*model.AnalysisMetadata, error) {
	key := cacheKey(applicationID)
	item, err := c.backend.Get(key)
	if err != nil {
		return nil, err
	}
	var a model.AnalysisMetadata
	if err := json.Unmarshal(item.([]byte), &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (c *analysisCache) Put(applicationID string, analysisMetadata *model.AnalysisMetadata) error {
	key := cacheKey(applicationID)
	data, err := json.Marshal(analysisMetadata)
	if err != nil {
		return err
	}
	return c.backend.Put(key, data)
}

func cacheKey(applicationID string) string {
	return fmt.Sprintf("most-recent-successful-analysis:%s", applicationID)
}
