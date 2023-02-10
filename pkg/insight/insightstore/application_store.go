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

package insightstore

import (
	"context"
	"encoding/json"

	"github.com/pipe-cd/pipecd/pkg/insight"
)

func (s *store) GetApplications(ctx context.Context, projectID string) (*insight.ProjectApplicationData, error) {
	content, err := s.fileStore.Get(ctx, makeApplicationsFilePath(projectID))
	if err != nil {
		return nil, err
	}

	data := &insight.ProjectApplicationData{}
	if err := json.Unmarshal(content, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *store) PutApplications(ctx context.Context, projectID string, as *insight.ProjectApplicationData) error {
	data, err := json.Marshal(as)
	if err != nil {
		return err
	}

	return s.fileStore.Put(ctx, makeApplicationsFilePath(projectID), data)
}
