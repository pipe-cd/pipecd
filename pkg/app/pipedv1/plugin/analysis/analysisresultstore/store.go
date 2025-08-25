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

package analysisresultstore

import (
	"context"
	"encoding/json"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	key = "latest-analysis-result"
)

var (
	ErrNotFound = errors.New("not found")
)

type apiClient interface {
	GetApplicationSharedObject(ctx context.Context, key string) ([]byte, error)
	PutApplicationSharedObject(ctx context.Context, key string, object []byte) error
}

type Store interface {
	GetLatestAnalysisResult(ctx context.Context, applicationID string) (*AnalysisResult, error)
	PutLatestAnalysisResult(ctx context.Context, applicationID string, analysisResult *AnalysisResult) error
}

type store struct {
	apiClient apiClient
	logger    *zap.Logger
}

func NewStore(apiClient apiClient, logger *zap.Logger) Store {
	return &store{
		apiClient: apiClient,
		logger:    logger.Named("analysis-result-store"),
	}
}

func (s *store) GetLatestAnalysisResult(ctx context.Context, appID string) (*AnalysisResult, error) {
	resp, err := s.apiClient.GetApplicationSharedObject(ctx, key)
	if status.Code(err) == codes.NotFound {
		s.logger.Info("analysis result is not found")
		return nil, ErrNotFound
	}
	if err != nil {
		s.logger.Error("failed to get the most recent analysis result",
			zap.String("application-id", appID),
			zap.Error(err),
		)
		return nil, err
	}

	result := &AnalysisResult{}
	if err = json.Unmarshal(resp, result); err != nil {
		s.logger.Error("failed to unmarshal the analysis result",
			zap.String("application-id", appID),
			zap.Error(err),
		)
		return nil, err
	}
	return result, nil
}

func (s *store) PutLatestAnalysisResult(ctx context.Context, appID string, analysisResult *AnalysisResult) error {
	json, err := json.Marshal(analysisResult)
	if err != nil {
		s.logger.Error("failed to marshal the analysis result",
			zap.String("application-id", appID),
			zap.Error(err),
		)
		return err
	}

	if err = s.apiClient.PutApplicationSharedObject(ctx, key, json); err != nil {
		s.logger.Error("failed to put the most recent analysis result",
			zap.String("application-id", appID),
			zap.Error(err),
		)
		return err
	}
	return nil
}
