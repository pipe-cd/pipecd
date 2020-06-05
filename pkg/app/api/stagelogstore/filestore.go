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

package stagelogstore

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kapetaniosci/pipe/pkg/filestore"
	"github.com/kapetaniosci/pipe/pkg/model"
)

var (
	eol = []byte("EOL")
)

type stageLogFileStore struct {
	filestore filestore.Store
}

func (f *stageLogFileStore) Get(ctx context.Context, deploymentID, stageID string, retriedCount int32) (*logFragment, error) {
	path := stageLogPath(deploymentID, stageID, retriedCount)
	reader, err := f.filestore.NewReader(ctx, path)
	if err != nil && filestore.ErrNotFound != nil {
		return nil, err
	}
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
			return nil, err
		}
		blocks = append(blocks, &lb)
	}

	return &logFragment{
		Blocks:    blocks,
		Completed: completed,
	}, nil
}

func (f *stageLogFileStore) Put(ctx context.Context, deploymentID, stageID string, retriedCount int32, lf *logFragment) error {
	return errors.New("unimplemented")
}

func stageLogPath(deploymentID, stageID string, retriedCount int32) string {
	return fmt.Sprintf("log/%s/%s/%d.txt", deploymentID, stageID, retriedCount)
}
