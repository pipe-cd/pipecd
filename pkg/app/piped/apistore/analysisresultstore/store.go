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

package analysisresultstore

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/model"
)

type apiClient interface {
	GetMostRecentSuccessfulAnalysisMetadata(ctx context.Context, req *pipedservice.GetMostRecentSuccessfulAnalysisMetadataRequest, opts ...grpc.CallOption) (*pipedservice.GetMostRecentSuccessfulAnalysisMetadataResponse, error)
	PutMostRecentSuccessfulAnalysisMetadata(ctx context.Context, req *pipedservice.PutMostRecentSuccessfulAnalysisMetadataRequest, opts ...grpc.CallOption) (*pipedservice.PutMostRecentSuccessfulAnalysisMetadataResponse, error)
}

type Store interface {
	GetMostRecentSuccessfulAnalysisMetadata(ctx context.Context, applicationID string) (*model.AnalysisMetadata, error)
	PutMostRecentSuccessfulAnalysisMetadata(ctx context.Context, applicationID string, analysisMetadata *model.AnalysisMetadata) error
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

func (s *store) GetMostRecentSuccessfulAnalysisMetadata(ctx context.Context, applicationID string) (*model.AnalysisMetadata, error) {
	resp, err := s.apiClient.GetMostRecentSuccessfulAnalysisMetadata(ctx, &pipedservice.GetMostRecentSuccessfulAnalysisMetadataRequest{ApplicationId: applicationID})
	if err != nil {
		s.logger.Error("failed to get the most recent analysis metadata", zap.Error(err))
		return nil, err
	}
	return resp.AnalysisMetadata, nil
}

func (s *store) PutMostRecentSuccessfulAnalysisMetadata(ctx context.Context, applicationID string, analysisMetadata *model.AnalysisMetadata) error {
	_, err := s.apiClient.PutMostRecentSuccessfulAnalysisMetadata(ctx, &pipedservice.PutMostRecentSuccessfulAnalysisMetadataRequest{
		ApplicationId:    applicationID,
		AnalysisMetadata: analysisMetadata,
	})
	if err != nil {
		s.logger.Error("failed to put the most recent analysis metadata", zap.Error(err))
		return err
	}
	return nil
}
