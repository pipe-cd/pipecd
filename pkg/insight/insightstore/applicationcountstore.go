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

// LoadApplicationCounts loads ApplicationCounts data for a specific project from file store.
func (s *store) LoadApplicationCounts(ctx context.Context, projectID string) (*insight.ApplicationCounts, error) {
	content, err := s.filestore.Get(ctx, determineFilePath(projectID))
	if err != nil {
		return nil, err
	}

	counts := &insight.ApplicationCounts{}
	if err := json.Unmarshal(content, counts); err != nil {
		return nil, err
	}

	return counts, nil
}

// PutApplicationCounts saves the ApplicationCounts data for a specific project into file store.
func (s *store) PutApplicationCounts(ctx context.Context, projectID string, counts insight.ApplicationCounts) error {
	data, err := json.Marshal(counts)
	if err != nil {
		return err
	}
	return s.filestore.Put(ctx, determineFilePath(projectID), data)
}

// File paths will be decided as below:
//
// insights
//  ├─ project-id
//    ├─ applications-count
//       ├─ applications-counts.json
func determineFilePath(projectID string) string {
	const path = "insights/%s/applications-count/applications-counts.json"
	return fmt.Sprintf(path, projectID)
}
