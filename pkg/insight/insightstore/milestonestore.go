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

package insightstore

import (
	"context"
	"encoding/json"

	"github.com/pipe-cd/pipe/pkg/insight"
)

const milestonePath = "insights/milestone.json"

func (s *Store) LoadMilestone(ctx context.Context) (*insight.Milestone, error) {
	m := &insight.Milestone{}
	obj, err := s.filestore.GetObject(ctx, milestonePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(obj.Content, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Store) PutMilestone(ctx context.Context, m *insight.Milestone) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return s.filestore.PutObject(ctx, milestonePath, data)
}
