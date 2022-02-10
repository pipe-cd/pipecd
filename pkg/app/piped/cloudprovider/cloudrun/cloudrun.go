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
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/api/run/v1"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	DefaultServiceManifestFilename = "service.yaml"
)

var (
	ErrServiceNotFound  = errors.New("not found")
	ErrRevisionNotFound = errors.New("not found")
)

type (
	Service  run.Service
	Revision run.Revision
)

const (
	LabelManagedBy   = "pipecd-dev-managed-by"  // Always be piped.
	LabelPiped       = "pipecd-dev-piped"       // The id of piped handling this application.
	LabelApplication = "pipecd-dev-application" // The application this resource belongs to.
	LabelCommitHash  = "pipecd-dev-commit-hash" // Hash value of the deployed commit.
	ManagedByPiped   = "piped"
)

type Client interface {
	Create(ctx context.Context, sm ServiceManifest) (*Service, error)
	Update(ctx context.Context, sm ServiceManifest) (*Service, error)
	List(ctx context.Context, options *ListOptions) ([]*Service, string, error)
	GetRevision(ctx context.Context, name string) (*Revision, error)
}

type ListOptions struct {
	Limit         int64
	LabelSelector string
	Cursor        string
}

type Registry interface {
	Client(ctx context.Context, name string, cfg *config.CloudProviderCloudRunConfig, logger *zap.Logger) (Client, error)
}

func LoadServiceManifest(appDir, serviceFilename string) (ServiceManifest, error) {
	if serviceFilename == "" {
		serviceFilename = DefaultServiceManifestFilename
	}
	path := filepath.Join(appDir, serviceFilename)
	return loadServiceManifest(path)
}

var defaultRegistry = &registry{
	clients:  make(map[string]Client),
	newGroup: &singleflight.Group{},
}

func DefaultRegistry() Registry {
	return defaultRegistry
}

type registry struct {
	clients  map[string]Client
	mu       sync.RWMutex
	newGroup *singleflight.Group
}

func (r *registry) Client(ctx context.Context, name string, cfg *config.CloudProviderCloudRunConfig, logger *zap.Logger) (Client, error) {
	r.mu.RLock()
	client, ok := r.clients[name]
	r.mu.RUnlock()
	if ok {
		return client, nil
	}

	c, err, _ := r.newGroup.Do(name, func() (interface{}, error) {
		return newClient(ctx, cfg.Project, cfg.Region, cfg.CredentialsFile, logger)
	})
	if err != nil {
		return nil, err
	}

	client = c.(Client)
	r.mu.Lock()
	r.clients[name] = client
	r.mu.Unlock()

	return client, nil
}

func MakeManagedByPipedLabel() string {
	return fmt.Sprintf("%s=%s", LabelManagedBy, ManagedByPiped)
}

func (s *Service) ServiceManifest() (ServiceManifest, error) {
	r := (*run.Service)(s)
	data, err := r.MarshalJSON()
	if err != nil {
		return ServiceManifest{}, err
	}
	return ParseServiceManifest(data)
}

func (s *Service) UID() (string, bool) {
	if s.Metadata == nil || s.Metadata.Uid == "" {
		return "", false
	}
	return s.Metadata.Uid, true
}

// ActiveRevisionNames returns all its active revisions which may handle the traffic.
func (s *Service) ActiveRevisionNames() []string {
	if s.Status == nil {
		return nil
	}
	tf := s.Status.Traffic
	ret := make([]string, len(tf))
	for i := range tf {
		ret[i] = tf[i].RevisionName
	}
	return ret
}

func (s *Service) HealthStatus() (model.CloudRunResourceState_HealthStatus, string) {
	var (
		status     = s.Status
		errMessage = "Unexpected error while calculating: %s"
	)
	if status == nil {
		return model.CloudRunResourceState_UNKNOWN, fmt.Sprintf(errMessage, "unable to find status")
	}

	conds := status.Conditions
	for _, c := range conds {
		isHealthy, err := strconv.ParseBool(c.Status)
		if err != nil {
			return model.CloudRunResourceState_UNKNOWN, fmt.Sprintf(errMessage, "unable to parse status: %s", c.Status)
		}
		if !isHealthy {
			return model.CloudRunResourceState_OTHER, c.Message
		}
	}
	return model.CloudRunResourceState_HEALTHY, ""
}

func (r *Revision) RevisionManifest() (RevisionManifest, error) {
	rev := (*run.Revision)(r)
	data, err := rev.MarshalJSON()
	if err != nil {
		return RevisionManifest{}, err
	}
	return ParseRevisionManifest(data)
}

func (r *Revision) HealthStatus() (model.CloudRunResourceState_HealthStatus, string) {
	var (
		status     = r.Status
		errMessage = "Unexpected error while calculating: %s"
	)
	if status == nil {
		return model.CloudRunResourceState_UNKNOWN, fmt.Sprintf(errMessage, "unable to find status")
	}

	conds := status.Conditions
	for _, c := range conds {
		isHealthy, err := strconv.ParseBool(c.Status)
		if err != nil {
			return model.CloudRunResourceState_UNKNOWN, fmt.Sprintf(errMessage, "unable to parse status: %s", c.Status)
		}
		if !isHealthy {
			return model.CloudRunResourceState_OTHER, c.Message
		}
	}
	return model.CloudRunResourceState_HEALTHY, ""
}
