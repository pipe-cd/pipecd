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

package kubernetes

import (
	"fmt"

	"github.com/kapetaniosci/pipe/pkg/cache"
)

type appManifestsCache struct {
	cache cache.Cache
}

func (c *appManifestsCache) Get(appID, commit string) ([]Manifest, error) {
	key := appManifestsCacheKey(appID, commit)
	item, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	return item.([]Manifest), nil
}

func (c *appManifestsCache) Put(appID, commit string, manifests []Manifest) error {
	key := appManifestsCacheKey(appID, commit)
	return c.cache.Put(key, manifests)
}

func appManifestsCacheKey(appID, commit string) string {
	return fmt.Sprintf("%s/%s", appID, commit)
}
