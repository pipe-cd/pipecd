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

package applicationsharedobjectstore

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/filestore"
)

type Store interface {
	// GetObject retrieves a shared object from the store.
	GetObject(ctx context.Context, appID, pluginName, key string) ([]byte, error)
	// PutObject stores a shared object to the store.
	PutObject(ctx context.Context, appID, pluginName, key string, data []byte) error
}

type store struct {
	backend filestore.Store
	logger  *zap.Logger
}

func NewStore(fs filestore.Store, logger *zap.Logger) Store {
	return &store{
		backend: fs,
		logger:  logger.Named("application-shared-object-store"),
	}
}

func (s *store) GetObject(ctx context.Context, appID, pluginName, key string) ([]byte, error) {
	path := buildPath(appID, pluginName, key)
	content, err := s.backend.Get(ctx, path)
	if err != nil {
		if !errors.Is(err, filestore.ErrNotFound) {
			s.logger.Error("failed to get object from filestore",
				zap.String("application-id", appID),
				zap.String("plugin-name", pluginName),
				zap.String("key", key),
				zap.Error(err),
			)
		}
		return nil, err
	}
	return content, nil
}

func (s *store) PutObject(ctx context.Context, appID, pluginName, key string, data []byte) error {
	path := buildPath(appID, pluginName, key)
	if err := s.backend.Put(ctx, path, data); err != nil {
		s.logger.Error("failed to put object to filestore",
			zap.String("application-id", appID),
			zap.String("plugin-name", pluginName),
			zap.String("key", key),
			zap.Int("data-size", len(data)),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func buildPath(appID, pluginName, key string) string {
	// Although the file might not be json, here we use .json for simplicity.
	return fmt.Sprintf("application-shared-objects/%s/%s/%s.json", appID, pluginName, key)
}
