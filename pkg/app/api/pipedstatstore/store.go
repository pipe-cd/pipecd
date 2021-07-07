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

package pipedstatstore

import (
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/cache"
)

type Store interface {
	PutPipedStat(pipedID string, stat []byte) error
	GetCurrentStats() ([]byte, error)
}

type store struct {
	backend cache.Cache
	logger  *zap.Logger
}

func NewStore(c cache.Cache, logger *zap.Logger) Store {
	return &store{
		backend: c,
		logger:  logger,
	}
}

func (s *store) PutPipedStat(pipedID string, stat []byte) error {
	if err := s.backend.PutHash(pipedID, stat); err != nil {
		s.logger.Error("failed to store piped's stat", zap.Error(err))
		return err
	}
	return nil
}

// TODO: Implement GetCurrentStats so that ops can use it to expose pipeds' stat to prometheus.
func (c *store) GetCurrentStats() ([]byte, error) {
	return []byte(""), nil
}
