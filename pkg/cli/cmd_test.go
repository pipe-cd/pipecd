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

package cli

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunWithContext(t *testing.T) {
	calls := 0
	piped := func(ctx context.Context, t Input) error {
		calls++
		timeout := time.NewTimer(time.Second)
		select {
		case <-ctx.Done():
			return nil
		case <-timeout.C:
			return fmt.Errorf("timed out")
		}
	}
	ch := make(chan os.Signal, 1)
	ch <- syscall.SIGINT
	app := NewApp("app", "test")
	err := runWithContext(app.rootCmd, ch, piped)
	assert.NoError(t, err)
	assert.Equal(t, 1, calls)
}

func TestNewLogger(t *testing.T) {
	logger, err := newLogger("service", "1.0.0", "debug", "json")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}
