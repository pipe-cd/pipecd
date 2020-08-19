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

package cloudrun

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
	"sigs.k8s.io/yaml"
)

type client struct {
	projectID string
	region    string
	client    *run.APIService
	logger    *zap.Logger
}

func newClient(ctx context.Context, projectID, region, credentialsFile string, logger *zap.Logger) (*client, error) {
	c := &client{
		projectID: projectID,
		region:    region,
		logger:    logger.Named("cloudrun"),
	}

	var options []option.ClientOption
	if len(credentialsFile) > 0 {
		data, err := ioutil.ReadFile(credentialsFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read credentials file (%w)", err)
		}
		options = append(options, option.WithCredentialsJSON(data))
	}
	options = append(options,
		option.WithEndpoint(fmt.Sprintf("https://%s-run.googleapis.com/", region)),
	)

	runClient, err := run.NewService(ctx, options...)
	if err != nil {
		return nil, err
	}
	c.client = runClient

	return c, nil
}

func (c *client) Apply(ctx context.Context, sm ServiceManifest) (*Service, error) {
	service, err := manifestToRunService(sm)
	if err != nil {
		return nil, err
	}

	var (
		svc  = run.NewNamespacesServicesService(c.client)
		name = makeCloudRunServiceName(c.projectID, sm.Name)
		call = svc.ReplaceService(name, service)
	)
	call.Context(ctx)
	updatedService, err := call.Do()
	if err != nil {
		if e, ok := err.(*googleapi.Error); ok && e.Code == http.StatusNotFound {
			return nil, fmt.Errorf("service %s was not found (%w)", name, ErrServiceNotFound)
		}
		return nil, err
	}

	return (*Service)(updatedService), nil
}

func (c *client) List(ctx context.Context) error {
	var (
		svc    = run.NewNamespacesServicesService(c.client)
		parent = makeCloudRunParent(c.projectID)
		call   = svc.List(parent)
	)

	call.Context(ctx)
	resp, err := call.Do()
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}

func makeCloudRunParent(projectID string) string {
	return fmt.Sprintf("namespaces/%s", projectID)
}

func makeCloudRunServiceName(projectID, serviceID string) string {
	return fmt.Sprintf("namespaces/%s/services/%s", projectID, serviceID)
}

func manifestToRunService(sm ServiceManifest) (*run.Service, error) {
	data, err := sm.YamlBytes()
	if err != nil {
		return nil, err
	}

	var s run.Service
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}
