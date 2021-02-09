// Copyright 2021 The PipeCD Authors.
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
	"fmt"

	"github.com/pipe-cd/pipe/pkg/insight"
)

// LoadApplicationCount loads insight.ApplicationCount.
func (s *Store) LoadApplicationCount(ctx context.Context, projectID string) (*insight.ApplicationCount, error) {
	a := &insight.ApplicationCount{}
	obj, err := s.filestore.GetObject(ctx, determineFilePath(projectID))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(obj.Content, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// PutApplicationCount creates or updates insight.ApplicationCount.
func (s *Store) PutApplicationCount(ctx context.Context, ac *insight.ApplicationCount, projectID string) error {
	data, err := json.Marshal(ac)
	if err != nil {
		return err
	}
	return s.filestore.PutObject(ctx, determineFilePath(projectID), data)
}

// File paths according to the following format.
//
// insights
//  ├─ projects  # aggregated application counts in all projects
//    ├─ applications-count
//       ├─ applications-count.json
//  ├─ project-id
//    ├─ applications-count
//       ├─ applications-count.json
func determineFilePath(projectID string) string {
	const applicationsCountFilePathFormat = "insights/%s/applications-count/applications-count.json"

	if projectID == "" {
		projectID = "projects"
	}
	return fmt.Sprintf(applicationsCountFilePathFormat, projectID)
}
