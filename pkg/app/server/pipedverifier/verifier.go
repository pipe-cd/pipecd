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

package pipedverifier

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/memorycache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type projectGetter interface {
	Get(ctx context.Context, id string) (*model.Project, error)
}

type pipedGetter interface {
	Get(ctx context.Context, id string) (*model.Piped, error)
}

type Verifier struct {
	config          *config.ControlPlaneSpec
	projectCache    cache.Cache
	projectStore    projectGetter
	pipedCache      cache.Cache
	pipedStore      pipedGetter
	invalidKeyCache cache.Cache
	logger          *zap.Logger
}

func NewVerifier(
	ctx context.Context,
	cfg *config.ControlPlaneSpec,
	projectGetter projectGetter,
	pipedGetter pipedGetter,
	logger *zap.Logger,
) *Verifier {
	return &Verifier{
		config:          cfg,
		projectCache:    memorycache.NewTTLCache(ctx, 12*time.Hour, time.Hour),
		projectStore:    projectGetter,
		pipedCache:      memorycache.NewTTLCache(ctx, 30*time.Minute, 5*time.Minute),
		pipedStore:      pipedGetter,
		invalidKeyCache: memorycache.NewTTLCache(ctx, 30*time.Minute, 5*time.Minute),
		logger:          logger,
	}
}

func (v *Verifier) Verify(ctx context.Context, projectID, pipedID, pipedKey string) error {
	// Check the project information.
	if err := v.verifyProject(ctx, projectID, pipedID); err != nil {
		return err
	}

	// Fail-fast for the invalid keys.
	keyID := fmt.Sprintf("%s#%s#%s", projectID, pipedID, pipedKey)
	if item, err := v.invalidKeyCache.Get(keyID); err == nil {
		return item.(error)
	}

	// Check the piped information.
	var piped *model.Piped
	item, err := v.pipedCache.Get(pipedID)
	if err == nil {
		piped = item.(*model.Piped)
		// When an error is returned, it is probably because
		// the requested key has just been added but not updated in the cache.
		// So in that case, we should refresh the cache data.
		if _, err := checkPiped(piped, projectID, pipedID, pipedKey); err == nil {
			return nil
		}
	}

	// If the cache data was not found or stale,
	// we have to retrieve from datastore and save it to the cache.
	piped, err = v.pipedStore.Get(ctx, pipedID)
	if err != nil {
		return fmt.Errorf("unable to find piped %s from datastore, %w", pipedID, err)
	}
	if err := v.pipedCache.Put(pipedID, piped); err != nil {
		v.logger.Warn("unable to store piped in memory cache", zap.Error(err))
	}

	keyNotMatch, err := checkPiped(piped, projectID, pipedID, pipedKey)
	if err != nil {
		v.logger.Info("detected an invalid piped key",
			zap.String("project", projectID),
			zap.String("piped-id", pipedID),
			zap.String("piped-key", pipedKey),
		)
		if keyNotMatch {
			v.invalidKeyCache.Put(keyID, err)
		}
		return err
	}

	return nil
}

func (v *Verifier) verifyProject(ctx context.Context, projectID, pipedID string) error {
	// Firstly, we check from the list specified in the Control Plane configuration.
	if _, ok := v.config.FindProject(projectID); ok {
		return nil
	}

	// If not found, we check from the list inside datastore.
	if _, err := v.projectCache.Get(projectID); err == nil {
		return nil
	}

	if _, err := v.projectStore.Get(ctx, projectID); err != nil {
		return fmt.Errorf("project %s for piped %s was not found", projectID, pipedID)
	}

	if err := v.projectCache.Put(projectID, true); err != nil {
		v.logger.Warn("unable to store project in memory cache", zap.Error(err))
	}

	return nil
}

func checkPiped(piped *model.Piped, projectID, pipedID, pipedKey string) (keyNotMatch bool, err error) {
	if piped.ProjectId != projectID {
		return false, fmt.Errorf("the project of piped %s is not matched, expected=%s, got=%s", pipedID, projectID, piped.ProjectId)
	}
	if piped.Disabled {
		return false, fmt.Errorf("piped %s was already disabled", pipedID)
	}
	if err := piped.CheckKey(pipedKey); err != nil {
		return true, fmt.Errorf("the key of piped %s is not matched, %v", pipedID, err)
	}
	return false, nil
}
