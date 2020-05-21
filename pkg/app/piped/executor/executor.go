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

package executor

import (
	"context"

	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Executor interface {
	// Execute starts running executor until completion
	// or the StopSignal has emitted.
	Execute(sig StopSignal) model.StageStatus
}

type Factory func(in Input) Executor

type LogPersister interface {
	Append(log string, s model.LogSeverity)
	AppendInfo(log string)
	AppendSuccess(log string)
	AppendError(log string)
}

type MetadataStore interface {
	Get(key string) (string, bool)
	Set(ctx context.Context, key, value string) error

	GetStageMetadata(stageID string, metadata interface{}) error
	SetStageMetadata(ctx context.Context, stageID string, metadata interface{}) error
}

type CommandLister interface {
	ListCommands() []model.ReportableCommand
}

type Input struct {
	Stage *model.PipelineStage
	// Readonly deployment model.
	Deployment       *model.Deployment
	DeploymentConfig *config.Config
	PipedConfig      *config.PipedSpec
	WorkingDir       string
	RepoDir          string
	StageWorkingDir  string
	CommandLister    CommandLister
	LogPersister     LogPersister
	MetadataStore    MetadataStore
	Logger           *zap.Logger
}

type StopSignalType string

const (
	// StopSignalTerminate means the executor should stop its execution
	// because the program was asked to terminate.
	StopSignalTerminate StopSignalType = "terminate"
	// StopSignalCancel means the executor should stop its execution
	// because the deployment was cancelled.
	StopSignalCancel StopSignalType = "cancel"
	// StopSignalTimeout means the executor should stop its execution
	// because of timeout.
	StopSignalTimeout StopSignalType = "timeout"
	// StopSignalNone means the excutor can be continuously executed.
	StopSignalNone StopSignalType = "none"
)

type StopSignal interface {
	Context() context.Context
	Ch() <-chan StopSignalType
	Signal() StopSignalType
	Stopped() bool
}

type StopSignalHandler interface {
	Cancel()
	Timeout()
	Terminate()
}

type stopSignal struct {
	ctx    context.Context
	cancel func()
	ch     chan StopSignalType
	signal *atomic.String
}

func NewStopSignal() (StopSignal, StopSignalHandler) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &stopSignal{
		ctx:    ctx,
		cancel: cancel,
		ch:     make(chan StopSignalType, 1),
		signal: atomic.NewString(string(StopSignalNone)),
	}
	return s, s
}

func (s *stopSignal) Cancel() {
	s.signal.Store(string(StopSignalCancel))
	s.cancel()
	s.ch <- StopSignalCancel
	close(s.ch)
}

func (s *stopSignal) Timeout() {
	s.signal.Store(string(StopSignalTimeout))
	s.cancel()
	s.ch <- StopSignalTimeout
	close(s.ch)
}

func (s *stopSignal) Terminate() {
	s.signal.Store(string(StopSignalTerminate))
	s.cancel()
	s.ch <- StopSignalTerminate
	close(s.ch)
}

func (s *stopSignal) Context() context.Context {
	return s.ctx
}

func (s *stopSignal) Ch() <-chan StopSignalType {
	return s.ch
}

func (s *stopSignal) Signal() StopSignalType {
	value := s.signal.Load()
	return StopSignalType(value)
}

func (s *stopSignal) Stopped() bool {
	value := s.signal.Load()
	return StopSignalType(value) != StopSignalNone
}

func DetermineStageStatus(sig StopSignalType, ori, got model.StageStatus) model.StageStatus {
	switch sig {
	case StopSignalNone:
		return got
	case StopSignalTerminate:
		return ori
	case StopSignalCancel:
		return model.StageStatus_STAGE_CANCELLED
	case StopSignalTimeout:
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_FAILURE
}
