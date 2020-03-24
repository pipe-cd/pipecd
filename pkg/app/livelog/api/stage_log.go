// Copyright 2020 The Pipe Authors.
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

package api

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kapetaniosci/pipe/pkg/app/livelog/service"
)

const (
	logFileName            = "stage.log"
	completionMarkFileName = "completed"
)

func (a *api) SendStageLiveLog(stream service.LiveLog_SendStageLiveLogServer) error {
	// Wait the first packet to get more information about pipeline stage.
	req, streamErr := stream.Recv()
	if streamErr != nil {
		a.logger.Error("failed to get the first packet", zap.Error(streamErr))
		return streamErr
	}
	//ctx := stream.Context()
	// if _, err := a.checkAgentKey(ctx, req.PipelineId); err != nil {
	// 	return err
	// }
	// Start appending stream message into log file.
	err := appendLog(a.dataDir, req.PipelineId, req.StageId, a.logger, func() (string, error) {
		if streamErr != nil {
			return "", streamErr
		}
		newLog := req.Log
		req, streamErr = stream.Recv()
		return newLog, nil
	})
	if err != nil {
		return err
	}
	return stream.SendAndClose(&service.SendStageLiveLogResponse{})
}

func appendLog(dataDir, pipelineID, buildStep string, logger *zap.Logger, next func() (string, error)) error {
	var (
		dir  = stageLogDir(dataDir, pipelineID, buildStep)
		path = stageLogFilePath(dataDir, pipelineID, buildStep)
	)
	logger = logger.With(zap.String("path", path))
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logger.Error("failed to create stage log directory", zap.Error(err))
		return status.Error(codes.Internal, "unabled to create stage log directory")
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Error("failed to open log file", zap.Error(err))
		return status.Error(codes.Internal, "unabled to open log file")
	}
	defer f.Close()

	logger.Info(fmt.Sprintf("start appending log into: %s", path))
	var count int
	for {
		newLog, err := next()
		if err == io.EOF {
			logger.Info("received all logs from the streaming client", zap.Int("count", count))
			if err := markLogCompletion(dataDir, pipelineID, buildStep, "completed"); err != nil {
				logger.Error("failed to create completion file", zap.Error(err))
			}
			return nil
		}
		if err != nil {
			logger.Error("failed while receiving log from the streaming client",
				zap.Int("count", count),
				zap.Error(err),
			)
			return err
		}
		if _, err := f.Write([]byte(newLog)); err != nil {
			logger.Error("failed to append new log into file", zap.Error(err))
		}
		count++
	}
}

func (a *api) GetStageLogSnapshot(ctx context.Context, req *service.GetStageLogSnapshotRequest) (*service.GetStageLogSnapshotResponse, error) {
	path := stageLogFilePath(a.dataDir, req.PipelineId, req.StageId)
	if _, err := os.Stat(path); err == os.ErrNotExist {
		return nil, status.Error(codes.NotFound, "log file was not found")
	}
	log, err := ioutil.ReadFile(path)
	if err != nil {
		a.logger.Error("failed to read log file", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	return &service.GetStageLogSnapshotResponse{
		Log: string(log),
	}, nil
}

func markLogCompletion(dataDir, pipelineID, stageID, content string) error {
	path := completionFilePath(dataDir, pipelineID, stageID)
	return ioutil.WriteFile(path, []byte(content), os.ModePerm)
}

func stageLogDir(dataDir, pipelineID, stageID string) string {
	return fmt.Sprintf("%s/%s_%s", dataDir, pipelineID, stageID)
}

func stageLogFilePath(dataDir, pipelineID, stageID string) string {
	return fmt.Sprintf("%s/%s_%s/%s", dataDir, pipelineID, stageID, logFileName)
}

func completionFilePath(dataDir, pipelineID, stageID string) string {
	return fmt.Sprintf("%s/%s_%s/%s", dataDir, pipelineID, stageID, completionMarkFileName)
}
