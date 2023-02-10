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

package analysisresultstore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type analysisFileStore struct {
	backend filestore.Store
}

func (f *analysisFileStore) Get(ctx context.Context, applicationID string) (*model.AnalysisResult, error) {
	path := buildPath(applicationID)

	content, err := f.backend.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	var a model.AnalysisResult
	if err := json.Unmarshal(content, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (f *analysisFileStore) Put(ctx context.Context, applicationID string, analysisResult *model.AnalysisResult) error {
	path := buildPath(applicationID)
	data, err := json.Marshal(analysisResult)
	if err != nil {
		return err
	}
	return f.backend.Put(ctx, path, data)
}

func buildPath(applicationID string) string {
	return fmt.Sprintf("latest-analysis-result/%s.json", applicationID)
}
