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

package launcher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGracefulStopCommand(t *testing.T) {
	testcases := []struct {
		name      string
		stopAfter time.Duration
	}{
		{
			name:      "graceful stop after very short time",
			stopAfter: time.Nanosecond,
		},
		{
			name:      "graceful stop after second",
			stopAfter: time.Second,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := runBinary("sh", []string{"sleep", "1m"})
			require.NoError(t, err)
			require.NotNil(t, cmd)

			assert.True(t, cmd.IsRunning())
			cmd.GracefulStop(tc.stopAfter)
			assert.False(t, cmd.IsRunning())
		})
	}
}
