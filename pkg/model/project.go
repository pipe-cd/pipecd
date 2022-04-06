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
		Name:      builtinRBACRoleAdmin.String(),
		Policies:  builtinAdminRBACPolicy,
		IsBuiltin: true,
	}
	builtinAdminRBACPolicy = []*ProjectRBACPolicy{
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
		Name:      builtinRBACRoleEditor.String(),
		Policies:  builtinEditorRBACPolicy,
		IsBuiltin: true,
	}
	builtinEditorRBACPolicy = []*ProjectRBACPolicy{
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_APPLICATION},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_ALL,
			},
		},
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_DEPLOYMENT},
				{Type: ProjectRBACResource_PIPED},
				{Type: ProjectRBACResource_DEPLOYMENT_CHAIN},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_GET,
				ProjectRBACPolicy_LIST,
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
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_INSIGHT},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_GET,
			},
		},
	}
	builtinViewerRBACRole = &ProjectRBACRole{
		Name:      builtinRBACRoleViewer.String(),
		Policies:  builtinViewerRBACPolicy,
		IsBuiltin: true,
	}
	builtinViewerRBACPolicy = []*ProjectRBACPolicy{
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_APPLICATION},
				{Type: ProjectRBACResource_DEPLOYMENT},
				{Type: ProjectRBACResource_PIPED},
				{Type: ProjectRBACResource_DEPLOYMENT_CHAIN},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_GET,
				ProjectRBACPolicy_LIST,
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
		{
			Resources: []*ProjectRBACResource{
				{Type: ProjectRBACResource_INSIGHT},
			},
			Actions: []ProjectRBACPolicy_Action{
				ProjectRBACPolicy_GET,
			},
		},
	}
)

type BuiltinRBACRole string

const (
	builtinRBACRoleAdmin  BuiltinRBACRole = "Admin"
	builtinRBACRoleEditor BuiltinRBACRole = "Editor"
	builtinRBACRoleViewer BuiltinRBACRole = "Viewer"
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
	if p.Github != nil {
		p.Github.RedactSensitiveData()
	}
	if p.Google != nil {
	}
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
	if p.Github != nil {
		if err := p.Github.Encrypt(encrypter); err != nil {
			return err
		}
	}
	if p.Google != nil {
	}
	return nil
}

// Decrypt decrypts encrypted data in ProjectSSOConfig.
func (p *ProjectSSOConfig) Decrypt(decrypter decrypter) error {
	if p.Github != nil {
		if err := p.Github.Decrypt(decrypter); err != nil {
			return err
		}
	}
	if p.Google != nil {
	}
	return nil
}

// GenerateAuthCodeURL generates an auth URL for the specified configuration.
func (p *ProjectSSOConfig) GenerateAuthCodeURL(project, callbackURL, state string) (string, error) {
	switch p.Provider {
	case ProjectSSOConfig_GITHUB, ProjectSSOConfig_GITHUB_ENTERPRISE:
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

// FilterRBACRoles filter rbac roles for built-in or not.
func (p *Project) FilterRBACRoles(isBuiltin bool) []*ProjectRBACRole {
	v := make([]*ProjectRBACRole, 0, len(p.RbacRoles))
	for _, role := range p.RbacRoles {
		if role.IsBuiltin != isBuiltin {
			continue
		}
		v = append(v, role)
	}
	return v
}

// ValidateRBACRoles validate rbac roles.
func ValidateRBACRoles(roles []*ProjectRBACRole) error {
	v := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		if role.IsBuiltinName() {
			return fmt.Errorf("cannot use built-in role name")
		}
		if _, ok := v[role.Name]; ok {
			return fmt.Errorf("role name must be unique")
		}
		v[role.Name] = struct{}{}
	}
	return nil
}

func (p *ProjectRBACRole) IsBuiltinName() bool {
	return p.Name == builtinRBACRoleAdmin.String() ||
		p.Name == builtinRBACRoleEditor.String() ||
		p.Name == builtinRBACRoleViewer.String()
}

// ValidateUserGroups validate user groups.
func ValidateUserGroups(groups []*ProjectUserGroup) error {
	v := make(map[string]struct{}, len(groups))
	for _, group := range groups {
		if _, ok := v[group.SsoGroup]; ok {
			return fmt.Errorf("the SSO group must be unique")
		}
		v[group.SsoGroup] = struct{}{}
	}
	return nil
}

func (p *Project) GetAllUserGroups() []*ProjectUserGroup {
	rbac := p.Rbac
	if rbac == nil {
		return p.UserGroups
	}

	// The full list also contains 3 legacy user groups.
	all := make([]*ProjectUserGroup, 0, len(p.UserGroups) + 3)
	if rbac.Admin != "" {
		group := &ProjectUserGroup{
			SsoGroup: rbac.Admin,
			Role:     builtinRBACRoleAdmin.String(),
		}
		v = append(v, group)
	}
	if rbac.Editor != "" {
		group := &ProjectUserGroup{
			SsoGroup: rbac.Editor,
			Role:     builtinRBACRoleEditor.String(),
		}
		v = append(v, group)
	}
	if rbac.Viewer != "" {
		group := &ProjectUserGroup{
			SsoGroup: rbac.Viewer,
			Role:     builtinRBACRoleViewer.String(),
		}
		v = append(v, group)
	}
	v = append(v, p.UserGroups...)

	return v
}

func (p *Project) GetAllRBACRoles() []*ProjectRBACRole {
	builtin := []*ProjectRBACRole{
		builtinAdminRBACRole,
		builtinEditorRBACRole,
		builtinViewerRBACRole,
	}

	// Initialize the capcity since the length of RbacRoles will continue to increase.
	size := len(p.RbacRoles) + len(builtin)
	v := make([]*ProjectRBACRole, 0, size)
	// Set built-in rbac role.
	v = append(v, builtin...)
	// Set custom rbac role.
	v = append(v, p.RbacRoles...)

	return v
}
