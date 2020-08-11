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

// MakePiped makes piped message without sensitive data.
func MakePiped(input *model.Piped) *Piped {
	if input == nil {
		return nil
	}
	return &Piped{
		Id:             input.Id,
		Name:           input.Name,
		Desc:           input.Desc,
		ProjectId:      input.ProjectId,
		Version:        input.Version,
		StartedAt:      input.StartedAt,
		CloudProviders: input.CloudProviders,
		RepositoryIds:  input.RepositoryIds,
		Disabled:       input.Disabled,
		CreatedAt:      input.CreatedAt,
		UpdatedAt:      input.UpdatedAt,
	}
}

// HasLabel checks if DeploymentConfigTemplate has the given label.
func (t *DeploymentConfigTemplate) HasLabel(label DeploymentConfigTemplateLabel) bool {
	for _, l := range t.Labels {
		if l == label {
			return true
		}
	}
	return false
}

// MakeProject makes project message without sensitive data.
func MakeProject(input *model.Project) *Project {
	if input == nil {
		return nil
	}
	var sso *ProjectSingleSignOn
	if input.Sso != nil {
		switch input.Sso.Provider {
		case model.ProjectSingleSignOnProvider_GITHUB:
			sso.Provider = ProjectSingleSignOnProvider_GITHUB
		case model.ProjectSingleSignOnProvider_GOOGLE:
			sso.Provider = ProjectSingleSignOnProvider_GOOGLE
		}
		if input.Sso.Github != nil {
			sso.Github = &ProjectSingleSignOn_GitHub{
				BaseUrl:    input.Sso.Github.BaseUrl,
				UploadUrl:  input.Sso.Github.UploadUrl,
				Org:        input.Sso.Github.Org,
				AdminTeam:  input.Sso.Github.AdminTeam,
				EditorTeam: input.Sso.Github.EditorTeam,
				ViewerTeam: input.Sso.Github.ViewerTeam,
			}
		}
		if input.Sso.Google != nil {
			sso.Google = &ProjectSingleSignOn_Google{}
		}
	}
	return &Project{
		Id:   input.Id,
		Desc: input.Desc,
		StaticAdmin: &ProjectStaticUser{
			Username: input.StaticAdmin.Username,
		},
		StaticAdminDisabled: input.StaticAdminDisabled,
		Sso:                 sso,
		CreatedAt:           input.CreatedAt,
		UpdatedAt:           input.UpdatedAt,
	}
}
