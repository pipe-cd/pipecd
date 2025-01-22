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

package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStopSignal(t *testing.T) {
	signal, handler := NewStopSignal()
	assert.NotNil(t, signal)
	assert.NotNil(t, signal.Context())
	assert.NotNil(t, signal.Ch())
	assert.NotNil(t, handler)
	assert.Equal(t, StopSignalNone, signal.Signal())
}

func TestCancel(t *testing.T) {
	signal, handler := NewStopSignal()
	handler.Cancel()
	assert.Equal(t, StopSignalCancel, signal.Signal())
	assert.Equal(t, StopSignalCancel, <-signal.Ch())
}

func TestTimeout(t *testing.T) {
	signal, handler := NewStopSignal()
	handler.Timeout()
	assert.Equal(t, StopSignalTimeout, signal.Signal())
	assert.Equal(t, StopSignalTimeout, <-signal.Ch())
}

func TestTerminate(t *testing.T) {
	signal, handler := NewStopSignal()
	handler.Terminate()
	assert.Equal(t, StopSignalTerminate, signal.Signal())
	assert.Equal(t, StopSignalTerminate, <-signal.Ch())
}

func TestTerminated(t *testing.T) {
	signal, handler := NewStopSignal()
	assert.False(t, signal.Terminated())
	handler.Terminate()
	assert.True(t, signal.Terminated())
}
