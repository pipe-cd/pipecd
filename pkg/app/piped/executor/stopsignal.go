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

package executor

import (
	"context"

	"go.uber.org/atomic"
)

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
	Terminated() bool
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

func (s *stopSignal) Terminated() bool {
	value := s.signal.Load()
	return StopSignalType(value) == StopSignalTerminate
}
