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

package applicationsharedobjectstore

import (
	"context"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/filestore"
)

type appObjectFileStore struct {
	backend filestore.Store
}

func (f *appObjectFileStore) Get(ctx context.Context, path string) ([]byte, error) {
	return f.backend.Get(ctx, path)
}

func (f *appObjectFileStore) Put(ctx context.Context, path string, data []byte) error {
	return f.backend.Put(ctx, path, data)
}

func buildPath(appID, pluginName, key string) string {
	// Although the file might not be json, here we use .json for simplicity.
	return fmt.Sprintf("application-shared-objects/%s/%s/%s.json", appID, pluginName, key)
}
