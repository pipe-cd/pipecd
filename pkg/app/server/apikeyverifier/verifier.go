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

package apikeyverifier

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/cache/memorycache"
	"github.com/pipe-cd/pipe/pkg/model"
)

type apiKeyGetter interface {
	GetAPIKey(ctx context.Context, id string) (*model.APIKey, error)
}

type Verifier struct {
	apiKeyCache cache.Cache
	apiKeyStore apiKeyGetter
	logger      *zap.Logger
}

func NewVerifier(ctx context.Context, getter apiKeyGetter, logger *zap.Logger) *Verifier {
	return &Verifier{
		apiKeyCache: memorycache.NewTTLCache(ctx, 5*time.Minute, time.Minute),
		apiKeyStore: getter,
		logger:      logger,
	}
}

func (v *Verifier) Verify(ctx context.Context, key string) (*model.APIKey, error) {
	keyID, err := model.ExtractAPIKeyID(key)
	if err != nil {
		return nil, err
	}

	var apiKey *model.APIKey
	item, err := v.apiKeyCache.Get(keyID)
	if err == nil {
		apiKey = item.(*model.APIKey)
		if err := checkAPIKey(apiKey, keyID, key); err != nil {
			return nil, err
		}
		return apiKey, nil
	}

	// If the cache data was not found,
	// we have to retrieve from datastore and save it to the cache.
	apiKey, err = v.apiKeyStore.GetAPIKey(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("unable to find API key %s from datastore, %w", keyID, err)
	}

	if err := v.apiKeyCache.Put(keyID, apiKey); err != nil {
		v.logger.Warn("unable to store API key in memory cache", zap.Error(err))
	}
	if err := checkAPIKey(apiKey, keyID, key); err != nil {
		return nil, err
	}

	return apiKey, nil
}

func checkAPIKey(apiKey *model.APIKey, id, key string) error {
	if apiKey.Disabled {
		return fmt.Errorf("the api key %s was already disabled", id)
	}

	if err := apiKey.CompareKey(key); err != nil {
		return fmt.Errorf("invalid api key %s: %w", id, err)
	}

	return nil
}
