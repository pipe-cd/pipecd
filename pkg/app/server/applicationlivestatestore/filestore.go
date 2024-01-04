// Copyright 2024 The PipeCD Authors.
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

package applicationlivestatestore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLiveStateFileStore struct {
	backend filestore.Store
}

func (f *applicationLiveStateFileStore) Get(ctx context.Context, applicationID string) (*model.ApplicationLiveStateSnapshot, error) {
	path := applicationLiveStatePath(applicationID)

	content, err := f.backend.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	var s model.ApplicationLiveStateSnapshot
	if err := json.Unmarshal(content, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (f *applicationLiveStateFileStore) Put(ctx context.Context, applicationID string, alss *model.ApplicationLiveStateSnapshot) error {
	path := applicationLiveStatePath(applicationID)
	data, err := json.Marshal(alss)
	if err != nil {
		return err
	}
	return f.backend.Put(ctx, path, data)
}

func applicationLiveStatePath(applicationID string) string {
	return fmt.Sprintf("application-live-state/%s.json", applicationID)
}
