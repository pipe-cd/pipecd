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

package webservice

import (
	"github.com/pipe-cd/pipe/pkg/model"
)

func MakePiped(input *model.Piped) *Piped {
	if input == nil {
		return nil
	}
	return &Piped{
		Id: input.Id,
		Desc: input.Desc,
		ProjectId: input.ProjectId,
		Version: input.Version,
		StartedAt: input.StartedAt,
		CloudProviders: input.CloudProviders,
		RepositoryIds: input.RepositoryIds,
		Disabled: input.Disabled,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
}
