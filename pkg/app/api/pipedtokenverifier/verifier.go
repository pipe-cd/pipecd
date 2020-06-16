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

package pipedtokenverifier

import (
	"context"
	"fmt"
	"time"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/cache/memorycache"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Verifier struct {
	config       *config.ControlPlaneSpec
	projectCache cache.Cache
	projectStore datastore.ProjectStore
	pipedCache   cache.Cache
	pipedStore   datastore.PipedStore
}

func NewVerifier(ctx context.Context, cfg *config.ControlPlaneSpec, ds datastore.DataStore) *Verifier {
	return &Verifier{
		config:       cfg,
		projectCache: memorycache.NewTTLCache(ctx, 12*time.Hour, time.Hour),
		projectStore: datastore.NewProjectStore(ds),
		pipedCache:   memorycache.NewTTLCache(ctx, 30*time.Minute, 5*time.Minute),
		pipedStore:   datastore.NewPipedStore(ds),
	}
}

func (v *Verifier) Verify(ctx context.Context, projectID, pipedID, pipedKey string) error {
	// Check the project information.
	// Firstly, we check from the list in Control Plane configuration.
	_, ok := v.config.GetProject(projectID)
	if !ok {
		// Check the projects inside datastore.
		if _, err := v.projectCache.Get(projectID); err != nil {
			if _, err := v.projectStore.GetProject(ctx, projectID); err != nil {
				return fmt.Errorf("project %s for piped %s was not found", projectID, pipedID)
			}
			v.projectCache.Put(projectID, true)
		}
	}

	// Check the piped information.
	var piped *model.Piped
	item, err := v.pipedCache.Get(pipedID)
	if err == nil {
		piped = item.(*model.Piped)
		return checkPiped(piped, projectID, pipedID, pipedKey)
	}

	// If memory cache data is not found,
	// we have to find from datastore and save to the cache.
	piped, err = v.pipedStore.GetPiped(ctx, pipedID)
	if err != nil {
		return fmt.Errorf("unabled to find piped %s from datastore, %w", pipedID, err)
	}
	v.pipedCache.Put(pipedID, piped)

	return checkPiped(piped, projectID, pipedID, pipedKey)
}

func checkPiped(piped *model.Piped, projectID, pipedID, pipedKey string) error {
	if piped.ProjectId != projectID {
		return fmt.Errorf("the project of piped %s is not matched, expected = %s, got = %s", pipedID, projectID, piped.ProjectId)
	}
	if piped.Disabled {
		return fmt.Errorf("piped %s was already disabled", pipedID)
	}
	if err := piped.CompareKey(pipedKey); err != nil {
		return fmt.Errorf("the key of piped %s is not matched, %v", pipedID, err)
	}
	return nil
}
