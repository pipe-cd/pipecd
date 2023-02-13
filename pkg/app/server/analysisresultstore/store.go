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

package analysisresultstore

import (
	"context"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type Store interface {
	// GetLatestAnalysisResult gives back the most recent successful analysis result of the specified application.
	GetLatestAnalysisResult(ctx context.Context, applicationID string) (*model.AnalysisResult, error)
	// PutStateSnapshot updates the most recent successful analysis result of the specified application.
	PutLatestAnalysisResult(ctx context.Context, applicationID string, snapshot *model.AnalysisResult) error
}

type store struct {
	backend *analysisFileStore
	logger  *zap.Logger
}

func NewStore(fs filestore.Store, logger *zap.Logger) Store {
	return &store{
		backend: &analysisFileStore{
			backend: fs,
		},
		logger: logger.Named("latest-analysis-store"),
	}
}

func (s *store) GetLatestAnalysisResult(ctx context.Context, applicationID string) (*model.AnalysisResult, error) {
	resp, err := s.backend.Get(ctx, applicationID)
	if err != nil {
		s.logger.Error("failed to get the most recent successful analysis result from filestore", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *store) PutLatestAnalysisResult(ctx context.Context, applicationID string, snapshot *model.AnalysisResult) error {
	if err := s.backend.Put(ctx, applicationID, snapshot); err != nil {
		s.logger.Error("failed to put the most recent successful analysis result to filestore", zap.Error(err))
		return err
	}
	return nil
}
