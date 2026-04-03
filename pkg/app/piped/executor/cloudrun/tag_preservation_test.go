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

package cloudrun

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/run/v1"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/cloudrun"
)

type mockClient struct {
	service *provider.Service
	getErr  error
}

func (m *mockClient) Create(ctx context.Context, sm provider.ServiceManifest) (*provider.Service, error) {
	return nil, nil
}

func (m *mockClient) Update(ctx context.Context, sm provider.ServiceManifest) (*provider.Service, error) {
	return nil, nil
}

func (m *mockClient) Get(ctx context.Context, serviceName string) (*provider.Service, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.service, nil
}

func (m *mockClient) List(ctx context.Context, options *provider.ListOptions) ([]*provider.Service, string, error) {
	return nil, "", nil
}

func (m *mockClient) GetRevision(ctx context.Context, name string) (*provider.Revision, error) {
	return nil, nil
}

func (m *mockClient) ListRevisions(ctx context.Context, options *provider.ListRevisionsOptions) ([]*provider.Revision, string, error) {
	return nil, "", nil
}

type mockLogPersister struct {
	logs []string
}

func (m *mockLogPersister) Info(msg string) {
	m.logs = append(m.logs, "INFO: "+msg)
}

func (m *mockLogPersister) Infof(format string, args ...interface{}) {
	m.logs = append(m.logs, "INFO: "+format)
}

func (m *mockLogPersister) Success(msg string) {
	m.logs = append(m.logs, "SUCCESS: "+msg)
}

func (m *mockLogPersister) Successf(format string, args ...interface{}) {
	m.logs = append(m.logs, "SUCCESS: "+format)
}

func (m *mockLogPersister) Error(msg string) {
	m.logs = append(m.logs, "ERROR: "+msg)
}

func (m *mockLogPersister) Errorf(format string, args ...interface{}) {
	m.logs = append(m.logs, "ERROR: "+format)
}

func (m *mockLogPersister) Write(data []byte) (int, error) {
	return len(data), nil
}

func TestGetExistingRevisionTags(t *testing.T) {
	ctx := context.Background()
	client := &mockClient{
		service: &provider.Service{
			Status: &run.ServiceStatus{
				Traffic: []*run.TrafficTarget{
					{
						RevisionName: "revision-1",
						Percent:      80,
						Tag:          "stable",
					},
					{
						RevisionName: "revision-2",
						Percent:      20,
						Tag:          "canary",
					},
					{
						RevisionName: "revision-3",
						Percent:      0,
					},
				},
			},
		},
	}
	lp := &mockLogPersister{}

	tags := getExistingRevisionTags(ctx, client, "test-service", lp)

	assert.Equal(t, 2, len(tags))
	assert.Equal(t, "stable", tags["revision-1"])
	assert.Equal(t, "canary", tags["revision-2"])
	assert.NotContains(t, tags, "revision-3")
}

func TestMergeTrafficWithExistingTags(t *testing.T) {
	traffics := []provider.RevisionTraffic{
		{
			RevisionName: "revision-1",
			Percent:      50,
		},
		{
			RevisionName: "revision-2",
			Percent:      50,
			Tag:          "explicit-tag",
		},
		{
			RevisionName: "new-revision",
			Percent:      0,
		},
	}
	existingTags := map[string]string{
		"revision-1": "stable",
		"revision-2": "old-tag",
	}

	result := mergeTrafficWithExistingTags(traffics, existingTags)

	assert.Equal(t, 3, len(result))
	assert.Equal(t, "stable", result[0].Tag, "Should preserve existing tag")
	assert.Equal(t, "explicit-tag", result[1].Tag, "Should keep explicit tag, not override")
	assert.Equal(t, "", result[2].Tag, "New revision should have no tag")
}
