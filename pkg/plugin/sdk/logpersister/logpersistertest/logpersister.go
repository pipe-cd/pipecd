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

package logpersistertest

import (
	"testing"
	"time"

	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
)

// NewTestLogPersister creates a new testLogPersister for testing.
func NewTestLogPersister(t *testing.T) TestLogPersister {
	return TestLogPersister{t}
}

// TestLogPersister implements logpersister for testing.
type TestLogPersister struct {
	t *testing.T
}

func (lp TestLogPersister) StageLogPersister(deploymentID, stageID string) logpersister.StageLogPersister {
	return lp
}

func (lp TestLogPersister) Write(log []byte) (int, error) {
	// Write the log to the test logger
	lp.t.Log(string(log))
	return 0, nil
}
func (lp TestLogPersister) Info(log string) {
	lp.t.Log("INFO", log)
}
func (lp TestLogPersister) Infof(format string, a ...interface{}) {
	lp.t.Logf("INFO "+format, a...)
}
func (lp TestLogPersister) Success(log string) {
	lp.t.Log("SUCCESS", log)
}
func (lp TestLogPersister) Successf(format string, a ...interface{}) {
	lp.t.Logf("SUCCESS "+format, a...)
}
func (lp TestLogPersister) Error(log string) {
	lp.t.Log("ERROR", log)
}
func (lp TestLogPersister) Errorf(format string, a ...interface{}) {
	lp.t.Logf("ERROR "+format, a...)
}
func (lp TestLogPersister) Complete(timeout time.Duration) error {
	lp.t.Logf("Complete stage log persister with timeout: %v", timeout)
	return nil
}
