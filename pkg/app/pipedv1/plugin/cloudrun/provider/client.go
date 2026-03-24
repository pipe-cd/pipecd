// Copyright 2026 The PipeCD Authors.
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

package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
)

var (
	ErrServiceNotFound  = errors.New("service is not found")
	ErrRevisionNotFound = errors.New("revision is not found")
)

type ListOptions struct {
	LabelSelector string
	Limit         int64
	Cursor        string
}

type ListRevisionsOptions struct {
	LabelSelector string
	Limit         int64
	Cursor        string
}

type Client interface {
	Create(ctx context.Context, sm ServiceManifest) (*run.Service, error)
	Update(ctx context.Context, sm ServiceManifest) (*run.Service, error)
	GetService(ctx context.Context, name string) (*run.Service, error)
	List(ctx context.Context, options *ListOptions) ([]*run.Service, error)
	GetRevision(ctx context.Context, name string) (*run.Revision, error)
	ListRevisions(ctx context.Context, options *ListRevisionsOptions) ([]*run.Revision, error)
}

type client struct {
	projectID string
	region    string
	client    *run.APIService
	logger    *zap.Logger
}

func NewClient(ctx context.Context, projectID, region, credentialsFile string, logger *zap.Logger) (Client, error) {
	c := &client{
		projectID: projectID,
		region:    region,
		logger:    logger.Named("cloudrun"),
	}

	var options []option.ClientOption
	if len(credentialsFile) > 0 {
		data, err := os.ReadFile(credentialsFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read credentials file: %w", err)
		}
		options = append(options, option.WithCredentialsJSON(data))
	} else {
		credentials, err := google.FindDefaultCredentials(ctx, run.CloudPlatformScope)
		if err != nil {
			return nil, fmt.Errorf("unable to find default credentials: %w", err)
		}
		options = append(options, option.WithCredentials(credentials))
	}
	options = append(options, option.WithEndpoint(fmt.Sprintf("https://%s-run.googleapis.com/", region)))

	runClient, err := run.NewService(ctx, options...)
	if err != nil {
		return nil, err
	}
	c.client = runClient

	return c, nil
}

func (c *client) Create(ctx context.Context, sm ServiceManifest) (*run.Service, error) {
	svcCfg, err := sm.RunService()
	if err != nil {
		return nil, err
	}

	var (
		svc    = run.NewNamespacesServicesService(c.client)
		parent = makeCloudRunParent(c.projectID)
		call   = svc.Create(parent, svcCfg)
	)
	call.Context(ctx)

	service, err := call.Do()
	if err != nil {
		var apierr *googleapi.Error
		if errors.As(err, &apierr) {
			return nil, fmt.Errorf("failed to create service: code=%d, message=%s, details=%s", apierr.Code, apierr.Message, apierr.Details)
		}
		return nil, err
	}
	return service, nil
}

func (c *client) Update(ctx context.Context, sm ServiceManifest) (*run.Service, error) {
	svcCfg, err := sm.RunService()
	if err != nil {
		return nil, err
	}

	var (
		svc  = run.NewNamespacesServicesService(c.client)
		name = makeCloudRunServiceName(c.projectID, sm.Name())
		call = svc.ReplaceService(name, svcCfg)
	)
	call.Context(ctx)

	service, err := call.Do()
	if err != nil {
		var apierr *googleapi.Error
		if errors.As(err, &apierr) && apierr.Code == http.StatusNotFound {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}
	return service, nil
}

func (c *client) GetService(ctx context.Context, name string) (*run.Service, error) {
	var (
		svc  = run.NewNamespacesServicesService(c.client)
		id   = makeCloudRunServiceName(c.projectID, name)
		call = svc.Get(id)
	)
	call.Context(ctx)

	service, err := call.Do()
	if err != nil {
		var apierr *googleapi.Error
		if errors.As(err, &apierr) && apierr.Code == http.StatusNotFound {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}
	return service, nil
}

func (c *client) List(ctx context.Context, options *ListOptions) ([]*run.Service, error) {
	var (
		svc    = run.NewNamespacesServicesService(c.client)
		parent = makeCloudRunParent(c.projectID)
		call   = svc.List(parent)
	)
	call.Context(ctx)
	if options.Limit != 0 {
		call.Limit(options.Limit)
	}
	if options.LabelSelector != "" {
		call.LabelSelector(options.LabelSelector)
	}
	if options.Cursor != "" {
		call.Continue(options.Cursor)
	}

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (c *client) GetRevision(ctx context.Context, name string) (*run.Revision, error) {
	var (
		svc  = run.NewNamespacesRevisionsService(c.client)
		id   = makeCloudRunRevisionName(c.projectID, name)
		call = svc.Get(id)
	)
	call.Context(ctx)

	revision, err := call.Do()
	if err != nil {
		var apierr *googleapi.Error
		if errors.As(err, &apierr) && apierr.Code == http.StatusNotFound {
			return nil, ErrRevisionNotFound
		}
		return nil, err
	}
	return revision, nil
}

func (c *client) ListRevisions(ctx context.Context, options *ListRevisionsOptions) ([]*run.Revision, error) {
	var (
		rev    = run.NewNamespacesRevisionsService(c.client)
		parent = makeCloudRunParent(c.projectID)
		call   = rev.List(parent)
	)
	call.Context(ctx)
	if options.Limit != 0 {
		call.Limit(options.Limit)
	}
	if options.LabelSelector != "" {
		call.LabelSelector(options.LabelSelector)
	}
	if options.Cursor != "" {
		call.Continue(options.Cursor)
	}

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func makeCloudRunParent(projectID string) string {
	return fmt.Sprintf("namespaces/%s", projectID)
}

func makeCloudRunServiceName(projectID, serviceID string) string {
	return fmt.Sprintf("namespaces/%s/services/%s", projectID, serviceID)
}

func makeCloudRunRevisionName(projectID, revisionID string) string {
	return fmt.Sprintf("namespaces/%s/revisions/%s", projectID, revisionID)
}
