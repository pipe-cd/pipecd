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

package stagelogstore

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

var (
	eol = []byte("EOL")
)

type stageLogFileStore struct {
	filestore filestore.Store
}

func (f *stageLogFileStore) Get(ctx context.Context, deploymentID, stageID string, retriedCount int32) (logFragment, error) {
	path := stageLogPath(deploymentID, stageID, retriedCount)
	lf := logFragment{}
	reader, err := f.filestore.GetReader(ctx, path)
	if err != nil {
		return lf, err
	}
	defer reader.Close()

	blocks := make([]*model.LogBlock, 0)
	scanner := bufio.NewScanner(reader)

	completed := false
	for scanner.Scan() {
		data := scanner.Bytes()
		if len(data) == 0 {
			continue
		}
		if bytes.Equal(data, eol) {
			completed = true
			break
		}
		var lb model.LogBlock
		if err := json.Unmarshal(data, &lb); err != nil {
			return lf, err
		}
		blocks = append(blocks, &lb)
	}

	lf.Blocks = blocks
	lf.Completed = completed
	return lf, nil
}

func (f *stageLogFileStore) Put(ctx context.Context, deploymentID, stageID string, retriedCount int32, lf *logFragment) error {
	path := stageLogPath(deploymentID, stageID, retriedCount)
	var buf bytes.Buffer
	for _, lb := range lf.Blocks {
		// TODO: Reduce the number of marshaling log blocks for improving performance
		raw, err := json.Marshal(lb)
		if err != nil {
			return err
		}
		buf.Write(raw)
		buf.WriteString("\n")
	}

	if lf.Completed {
		buf.Write(eol)
	}
	return f.filestore.Put(ctx, path, buf.Bytes())
}

func stageLogPath(deploymentID, stageID string, retriedCount int32) string {
	return fmt.Sprintf("log/%s/%s/%d.txt", deploymentID, stageID, retriedCount)
}
