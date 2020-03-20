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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendLog(t *testing.T) {
	dir, err := ioutil.TempDir("", "data")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	var (
		count      = 0
		maxLine    = 3
		pipelineID = "test-pipeline"
		stageID    = "test-stage"
		path       = stageLogFilePath(dir, pipelineID, stageID)
		logger     = zap.NewNop()
		eof        = false
		next       = func() (string, error) {
			count++
			if count < maxLine {
				return fmt.Sprintf("line-%d\n", count), nil
			}
			if eof {
				return "", io.EOF
			}
			return "", fmt.Errorf("")
		}
	)

	err = appendLog(dir, pipelineID, stageID, logger, next)
	require.Error(t, err)

	data, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, `line-1
line-2
`, string(data))

	count = 2
	maxLine = 5
	eof = true
	err = appendLog(dir, pipelineID, stageID, logger, next)
	require.NoError(t, err)

	data, err = ioutil.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, `line-1
line-2
line-3
line-4
`,
		string(data))

	data, err = ioutil.ReadFile(completionFilePath(dir, pipelineID, stageID))
	require.NoError(t, err)
	assert.Equal(t, "completed", string(data))
}
