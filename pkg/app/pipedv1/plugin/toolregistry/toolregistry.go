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

package toolregistry

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/cmd/piped/service"
)

type ToolRegistry struct {
	client service.PluginServiceClient
}

func (r *ToolRegistry) InstallTool(ctx context.Context, name, version, script string) (path string, err error) {
	res, err := r.client.InstallTool(ctx, &service.InstallToolRequest{
		Name:          name,
		Version:       version,
		InstallScript: script,
	})

	if err != nil {
		return "", err
	}

	return res.GetInstalledPath(), nil
}
