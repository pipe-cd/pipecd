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

package model

import (
	"crypto/subtle"
	"fmt"
	"net/url"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubScopes = []string{"read:org"}

	builtinAdminRBACRole = &ProjectRBACRole{
		Name:      BuiltinRBACRoleAdmin.String(),
		Policies:  builtinAdminRBACPolicies,
		IsBuiltin: true,
	}
	builtinAdminRBACPolicies = []*ProjectRBACPolicy{
		{
			Resources: []*ProjectRBACResource{
				{
					Type: ProjectRBACResource_ALL,
				},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_ALL,
			},
		},
	}
	builtinEditorRBACRole = &ProjectRBACRole{
		Name:      BuiltinRBACRoleEditor.String(),
		Policies:  builtinEditorRBACPolicies,
		IsBuiltin: true,
	}
	builtinEditorRBACPolicies = []*ProjectRBACPolicy{
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_APPLICATION},
				{Type: ProjectRBACResource_DEPLOYMENT},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_ALL,
			},
		},
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_PIPED},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_GET,
				ProjectRBACPolicy_LIST,
			},
		},
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_PROJECT},
				{Type: ProjectRBACResource_INSIGHT},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_GET,
			},
		},
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_EVENT},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_LIST,
			},
		},
	}
	builtinViewerRBACRole = &ProjectRBACRole{
		Name:      BuiltinRBACRoleViewer.String(),
		Policies:  builtinViewerRBACPolicies,
		IsBuiltin: true,
	}
	builtinViewerRBACPolicies = []*ProjectRBACPolicy{
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_APPLICATION},
				{Type: ProjectRBACResource_DEPLOYMENT},
				{Type: ProjectRBACResource_PIPED},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_GET,
				ProjectRBACPolicy_LIST,
			},
		},
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_PROJECT},
				{Type: ProjectRBACResource_INSIGHT},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_GET,
			},
		},
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_EVENT},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_LIST,
			},
		},
	}
)

type BuiltinRBACRole string

const (
	BuiltinRBACRoleAdmin  BuiltinRBACRole = "Admin"
	BuiltinRBACRoleEditor BuiltinRBACRole = "Editor"
	BuiltinRBACRoleViewer BuiltinRBACRole = "Viewer"
)

func (b BuiltinRBACRole) String() string {
	return string(b)
}

type encrypter interface {
	Encrypt(text string) (string, error)
}

type decrypter interface {
	Decrypt(encryptedText string) (string, error)
}

// SetStaticAdmin sets admin data.
func (p *Project) SetStaticAdmin(username, password string) error {
	encoded, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.StaticAdmin = &ProjectStaticUser{
		Username:     username,
		PasswordHash: string(encoded),
	}
	return nil
}

// RedactSensitiveData redacts sensitive data.
func (p *Project) RedactSensitiveData() {
	if p.StaticAdmin != nil {
		p.StaticAdmin.RedactSensitiveData()
	}
	if p.Sso != nil {
		p.Sso.RedactSensitiveData()
	}
}

func (p *Project) SetUpdatedAt(t int64) {
	p.UpdatedAt = t
}

// RedactSensitiveData redacts sensitive data.
func (p *ProjectStaticUser) RedactSensitiveData() {
	p.PasswordHash = redactedMessage
}

// Update updates ProjectStaticUser with given data.
func (p *ProjectStaticUser) Update(username, password string) error {
	if username != "" {
		p.Username = username
	}
	if password != "" {
		encoded, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		p.PasswordHash = string(encoded)
	}
	return nil
}

// Auth confirms username and password.
func (p *ProjectStaticUser) Auth(username, password string) error {
	if username == "" {
		return fmt.Errorf("username is empty")
	}
	if subtle.ConstantTimeCompare([]byte(p.Username), []byte(username)) != 1 {
		return fmt.Errorf("wrong username %q", username)
	}
	if password == "" {
		return fmt.Errorf("password is empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(password)); err != nil {
		return fmt.Errorf("wrong password for username %q: %v", username, err)
	}
	return nil
}

// RedactSensitiveData redacts sensitive data.
func (p *ProjectSSOConfig) RedactSensitiveData() {
	if p.Github == nil {
		return
	}
	p.Github.RedactSensitiveData()
}

// Update updates ProjectSSOConfig with given data.
func (p *ProjectSSOConfig) Update(sso *ProjectSSOConfig) error {
	p.Provider = sso.Provider
	if sso.Github != nil {
		if p.Github == nil {
			p.Github = &ProjectSSOConfig_GitHub{}
		}
		if err := p.Github.Update(sso.Github); err != nil {
			return err
		}
	}
	if sso.Google != nil {
	}
	return nil
}

// Encrypt encrypts sensitive data in ProjectSSOConfig.
func (p *ProjectSSOConfig) Encrypt(encrypter encrypter) error {
	if p.Github == nil {
		return nil
	}
	if err := p.Github.Encrypt(encrypter); err != nil {
		return err
	}
	return nil
}

// Decrypt decrypts encrypted data in ProjectSSOConfig.
func (p *ProjectSSOConfig) Decrypt(decrypter decrypter) error {
	if p.Github == nil {
		return nil
	}
	if err := p.Github.Decrypt(decrypter); err != nil {
		return err
	}
	return nil
}

// GenerateAuthCodeURL generates an auth URL for the specified configuration.
func (p *ProjectSSOConfig) GenerateAuthCodeURL(project, callbackURL, state string) (string, error) {
	switch p.Provider {
	case ProjectSSOConfig_GITHUB:
		if p.Github == nil {
			return "", fmt.Errorf("missing GitHub oauth in the SSO configuration")
		}
		return p.Github.GenerateAuthCodeURL(project, callbackURL, state)

	default:
		return "", fmt.Errorf("not implemented")
	}
}

// RedactSensitiveData redacts sensitive data.
func (p *ProjectSSOConfig_GitHub) RedactSensitiveData() {
	p.ClientId = redactedMessage
	p.ClientSecret = redactedMessage
}

// Update updates ProjectSSOConfig with given data.
func (p *ProjectSSOConfig_GitHub) Update(input *ProjectSSOConfig_GitHub) error {
	if input.ClientId != "" {
		p.ClientId = input.ClientId
	}
	if input.ClientSecret != "" {
		p.ClientSecret = input.ClientSecret
	}
	if input.BaseUrl != "" {
		p.BaseUrl = input.BaseUrl
	}
	if input.UploadUrl != "" {
		p.UploadUrl = input.UploadUrl
	}
	return nil
}

// Encrypt encrypts sensitive data in ProjectSSOConfig.
func (p *ProjectSSOConfig_GitHub) Encrypt(encrypter encrypter) error {
	if p.ClientId != "" {
		encrypedClientID, err := encrypter.Encrypt(p.ClientId)
		if err != nil {
			return err
		}
		p.ClientId = encrypedClientID
	}
	if p.ClientSecret != "" {
		encryptedClientSecret, err := encrypter.Encrypt(p.ClientSecret)
		if err != nil {
			return err
		}
		p.ClientSecret = encryptedClientSecret
	}
	return nil
}

// Decrypt decrypts ProjectSSOConfig.
func (p *ProjectSSOConfig_GitHub) Decrypt(decrypter decrypter) error {
	if p.ClientId != "" {
		decrypedClientID, err := decrypter.Decrypt(p.ClientId)
		if err != nil {
			return err
		}
		p.ClientId = decrypedClientID
	}
	if p.ClientSecret != "" {
		decryptedClientSecret, err := decrypter.Decrypt(p.ClientSecret)
		if err != nil {
			return err
		}
		p.ClientSecret = decryptedClientSecret
	}
	return nil
}

// GenerateAuthCodeURL generates an auth URL for the specified configuration.
func (p *ProjectSSOConfig_GitHub) GenerateAuthCodeURL(project, callbackURL, state string) (string, error) {
	cfg := oauth2.Config{
		ClientID: p.ClientId,
		Endpoint: github.Endpoint,
	}
	if p.BaseUrl != "" {
		u, err := url.Parse(p.BaseUrl)
		if err != nil {
			return "", err
		}
		cfg.Endpoint.AuthURL = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, "/login/oauth/authorize")
	}
	cfg.Scopes = githubScopes
	cfg.RedirectURL = fmt.Sprintf("%s?project=%s", callbackURL, project)
	authURL := cfg.AuthCodeURL(state, oauth2.ApprovalForce, oauth2.AccessTypeOnline)

	return authURL, nil
}

// HasRBACRole checks whether the RBAC role is exists.
func (p *Project) HasRBACRole(name string) bool {
	for _, v := range p.RbacRoles {
		if v.Name == name {
			return true
		}
	}
	return false
}

// HasUserGroup checks whether the user group is exists.
func (p *Project) HasUserGroup(sso string) bool {
	for _, v := range p.UserGroups {
		if v.SsoGroup == sso {
			return true
		}
	}
	if r := p.Rbac; r != nil {
		return r.Admin == sso || r.Editor == sso || r.Viewer == sso
	}
	return false
}

// SetLegacyUserGroups sets the legacy RBAC config as user groups if exists.
// If the same team exists in the legacy RBAC config, this method just only sets the user group that has the highest authority level.
func (p *Project) SetLegacyUserGroups() {
	rbac := p.Rbac
	if rbac == nil {
		return
	}

	// The full list also contains 3 legacy user groups.
	all := make([]*ProjectUserGroup, 0, len(p.UserGroups)+3)
	if rbac.Admin != "" {
		all = append(all, &ProjectUserGroup{
			SsoGroup: rbac.Admin,
			Role:     BuiltinRBACRoleAdmin.String(),
		})
	}
	if rbac.Editor != "" && rbac.Editor != rbac.Admin {
		all = append(all, &ProjectUserGroup{
			SsoGroup: rbac.Editor,
			Role:     BuiltinRBACRoleEditor.String(),
		})
	}
	if rbac.Viewer != "" && rbac.Viewer != rbac.Admin && rbac.Viewer != rbac.Editor {
		all = append(all, &ProjectUserGroup{
			SsoGroup: rbac.Viewer,
			Role:     BuiltinRBACRoleViewer.String(),
		})
	}
	all = append(all, p.UserGroups...)
	p.UserGroups = all
}

// AddUserGroup adds a user group.
func (p *Project) AddUserGroup(sso, role string) error {
	if p.HasUserGroup(sso) {
		return fmt.Errorf("%s is already being used. The SSO group must be unique", sso)
	}
	if !p.HasRBACRole(role) && !isBuiltinRBACRole(role) {
		return fmt.Errorf("%s role does not exist", role)
	}
	p.UserGroups = append(p.UserGroups, &ProjectUserGroup{
		SsoGroup: sso,
		Role:     role,
	})
	return nil
}

// DeleteUserGroup deletes a user group.
func (p *Project) DeleteUserGroup(sso string) error {
	for i, v := range p.UserGroups {
		if v.SsoGroup == sso {
			c := copy(p.UserGroups[i:], p.UserGroups[i+1:])
			p.UserGroups = p.UserGroups[:i+c]
			return nil
		}
	}
	deleted := false
	if p.Rbac != nil {
		if p.Rbac.Admin == sso {
			p.Rbac.Admin, deleted = "", true
		}
		if p.Rbac.Editor == sso {
			p.Rbac.Editor, deleted = "", true
		}
		if p.Rbac.Viewer == sso {
			p.Rbac.Viewer, deleted = "", true
		}
	}
	if deleted {
		return nil
	}
	return fmt.Errorf("%s user group does not exist", sso)
}

// SetBuiltinRBACRoles sets built-in roles.
func (p *Project) SetBuiltinRBACRoles() {
	builtin := []*ProjectRBACRole{
		builtinAdminRBACRole,
		builtinEditorRBACRole,
		builtinViewerRBACRole,
	}
	all := make([]*ProjectRBACRole, 0, len(p.RbacRoles)+len(builtin))
	// Set built-in rbac role.
	all = append(all, builtin...)
	// Set custom rbac role.
	all = append(all, p.RbacRoles...)
	p.RbacRoles = all
}

// isBuiltinRBACRole checks whether the name is the name of built-in role.
func isBuiltinRBACRole(name string) bool {
	return name == BuiltinRBACRoleAdmin.String() ||
		name == BuiltinRBACRoleEditor.String() ||
		name == BuiltinRBACRoleViewer.String()
}

// AddRBACRole adds a custom RBAC role.
func (p *Project) AddRBACRole(name string, policies []*ProjectRBACPolicy) error {
	if p.HasRBACRole(name) {
		return fmt.Errorf("the name of %s is already used", name)
	}
	if isBuiltinRBACRole(name) {
		return fmt.Errorf("the name of built-in role cannot be used")
	}
	p.RbacRoles = append(p.RbacRoles, &ProjectRBACRole{
		Name:     name,
		Policies: policies,
	})
	return nil
}

// UpdateRBACRole updates a custom RBAC role.
// Built-in role cannot be updated.
func (p *Project) UpdateRBACRole(name string, policies []*ProjectRBACPolicy) error {
	for _, v := range p.RbacRoles {
		if v.Name == name {
			v.Policies = policies
			return nil
		}
	}
	if isBuiltinRBACRole(name) {
		return fmt.Errorf("built-in role cannot be updated")
	}
	return fmt.Errorf("%s role does not exist", name)
}

// DeleteRBACRole deletes a custom RBAC role.
// Built-in role cannot be deleted.
func (p *Project) DeleteRBACRole(name string) error {
	if isBuiltinRBACRole(name) {
		return fmt.Errorf("built-in role cannot be deleted")
	}
	for i, v := range p.RbacRoles {
		if v.Name == name {
			c := copy(p.RbacRoles[i:], p.RbacRoles[i+1:])
			p.RbacRoles = p.RbacRoles[:i+c]
			return nil
		}
	}
	return fmt.Errorf("%s role does nott exist", name)
}

func (p *ProjectRBACRole) HasPermission(typ ProjectRBACResource_ResourceType, action ProjectRBACPolicy_Action) bool {
	for _, v := range p.Policies {
		if v.HasPermission(typ, action) {
			return true
		}
	}
	return false
}

func (p *ProjectRBACPolicy) HasPermission(typ ProjectRBACResource_ResourceType, action ProjectRBACPolicy_Action) bool {
	var hasResource bool
	for _, r := range p.Resources {
		if r.Type == typ || r.Type == ProjectRBACResource_ALL {
			hasResource = true
			break
		}
	}

	if !hasResource {
		return false
	}

	for _, a := range p.Actions {
		if a == action || a == ProjectRBACPolicy_ALL {
			return true
		}
	}
	return false
}
