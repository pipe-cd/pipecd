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

package commandoutputstore

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/filestore"
)

var (
	ErrNotFound = errors.New("not found")
)

type Store interface {
	Get(ctx context.Context, commandID string) ([]byte, error)
	Put(ctx context.Context, commandID string, data []byte) error
}

type store struct {
	backend filestore.Store
	logger  *zap.Logger
}

func NewStore(fs filestore.Store, logger *zap.Logger) Store {
	return &store{
		backend: fs,
		logger:  logger.Named("command-output-store"),
	}
}

func (s *store) Get(ctx context.Context, commandID string) ([]byte, error) {
	path := dataPath(commandID)
	content, err := s.backend.Get(ctx, path)
	if err != nil {
		if err == filestore.ErrNotFound {
			return nil, ErrNotFound
		}
		s.logger.Error("failed to get command output from filestore",
			zap.String("command", commandID),
			zap.Error(err),
		)
		return nil, err
	}
	return content, nil
}

func (s *store) Put(ctx context.Context, commandID string, data []byte) error {
	path := dataPath(commandID)
	return s.backend.Put(ctx, path, data)
}

func dataPath(commandID string) string {
	return fmt.Sprintf("command-output/%s.json", commandID)
}
