// Copyright 2022 The PipeCD Authors.
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
	"strings"
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

var (
	TypeConditions = map[string]struct{}{
		"Active":              struct{}{},
		"Ready":               struct{}{},
		"ConfigurationsReady": struct{}{},
		"RoutesReady":         struct{}{},
		"ContainerHealthy":    struct{}{},
		"ResourcesAvailable":  struct{}{},
	}
	TypeHealthyServiceConditions = map[string]struct{}{
		"Ready":               struct{}{},
		"ConfigurationsReady": struct{}{},
		"RoutesReady":         struct{}{},
	}
	TypeHealthyRevisionConditions = map[string]struct{}{
		"Ready":              struct{}{},
		"Active":             struct{}{},
		"ContainerHealthy":   struct{}{},
		"ResourcesAvailable": struct{}{},
	}
)

// Kind represents the kind of resource.
type Kind string

const (
	KindService  Kind = "Service"
	KindRevision Kind = "Revision"
)

type (
	Service  run.Service
	Revision run.Revision

	StatusConditions struct {
		Kind      Kind
		TrueTypes map[string]struct{}

		// Eliminate duplicated messages with the same reason.
		FalseMessages   []string
		UnknownMessages []string
	}
)

const (
	LabelManagedBy    = "pipecd-dev-managed-by"    // Always be piped.
	LabelPiped        = "pipecd-dev-piped"         // The id of piped handling this application.
	LabelApplication  = "pipecd-dev-application"   // The application this resource belongs to.
	LabelCommitHash   = "pipecd-dev-commit-hash"   // Hash value of the deployed commit.
	LabelRevisionName = "pipecd-dev-revision-name" // The name of revision.
	ManagedByPiped    = "piped"
)

type Client interface {
	Create(ctx context.Context, sm ServiceManifest) (*Service, error)
	Update(ctx context.Context, sm ServiceManifest) (*Service, error)
	List(ctx context.Context, options *ListOptions) ([]*Service, string, error)
	GetRevision(ctx context.Context, name string) (*Revision, error)
	ListRevisions(ctx context.Context, options *ListRevisionsOptions) ([]*Revision, string, error)
}

type ListOptions struct {
	Limit         int64
	LabelSelector string
	Cursor        string
}

type ListRevisionsOptions struct {
	Limit         int64
	LabelSelector string
	Cursor        string
}

type Registry interface {
	Client(ctx context.Context, name string, cfg *config.PlatformProviderCloudRunConfig, logger *zap.Logger) (Client, error)
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

func (r *registry) Client(ctx context.Context, name string, cfg *config.PlatformProviderCloudRunConfig, logger *zap.Logger) (Client, error) {
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

func MakeManagedByPipedSelector() string {
	return fmt.Sprintf("%s=%s", LabelManagedBy, ManagedByPiped)
}

func MakeRevisionNamesSelector(names []string) string {
	return fmt.Sprintf("%s in (%s)", LabelRevisionName, strings.Join(names, ","))
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

func (s *Service) StatusConditions() *StatusConditions {
	var (
		trueTypes   = make(map[string]struct{}, len(TypeConditions))
		falseMsgs   = make(map[string]string, len(TypeConditions))
		unknownMsgs = make(map[string]string, len(TypeConditions))
	)

	if s.Status == nil {
		return nil
	}
	for _, cond := range s.Status.Conditions {
		if _, ok := TypeConditions[cond.Type]; !ok {
			continue
		}
		switch cond.Status {
		case "True":
			trueTypes[cond.Type] = struct{}{}
		case "False":
			falseMsgs[cond.Reason] = cond.Message
		default:
			unknownMsgs[cond.Reason] = cond.Message
		}
	}

	fMsgs := make([]string, 0, len(falseMsgs))
	for _, v := range falseMsgs {
		fMsgs = append(fMsgs, v)
	}

	uMsgs := make([]string, 0, len(unknownMsgs))
	for _, v := range unknownMsgs {
		uMsgs = append(uMsgs, v)
	}

	return &StatusConditions{
		Kind:            KindService,
		TrueTypes:       trueTypes,
		FalseMessages:   fMsgs,
		UnknownMessages: uMsgs,
	}
}

func (r *Revision) RevisionManifest() (RevisionManifest, error) {
	rev := (*run.Revision)(r)
	data, err := rev.MarshalJSON()
	if err != nil {
		return RevisionManifest{}, err
	}
	return ParseRevisionManifest(data)
}

func (r *Revision) StatusConditions() *StatusConditions {
	var (
		trueTypes   = make(map[string]struct{}, len(TypeConditions))
		falseMsgs   = make(map[string]string, len(TypeConditions))
		unknownMsgs = make(map[string]string, len(TypeConditions))
	)

	if r.Status == nil {
		return nil
	}
	for _, cond := range r.Status.Conditions {
		if _, ok := TypeConditions[cond.Type]; !ok {
			continue
		}
		switch cond.Status {
		case "True":
			trueTypes[cond.Type] = struct{}{}
		case "False":
			falseMsgs[cond.Reason] = cond.Message
		default:
			unknownMsgs[cond.Reason] = cond.Message
		}
	}

	fMsgs := make([]string, 0, len(falseMsgs))
	for _, v := range falseMsgs {
		fMsgs = append(fMsgs, v)
	}

	uMsgs := make([]string, 0, len(unknownMsgs))
	for _, v := range unknownMsgs {
		uMsgs = append(uMsgs, v)
	}

	return &StatusConditions{
		Kind:            KindRevision,
		TrueTypes:       trueTypes,
		FalseMessages:   fMsgs,
		UnknownMessages: uMsgs,
	}
}

func (s *StatusConditions) HealthStatus() (model.CloudRunResourceState_HealthStatus, string) {
	if s == nil {
		return model.CloudRunResourceState_UNKNOWN, "Unexpected error while calculating: unable to find status"
	}

	if len(s.FalseMessages) > 0 {
		return model.CloudRunResourceState_OTHER, strings.Join(s.FalseMessages, "; ")
	}

	if len(s.UnknownMessages) > 0 {
		return model.CloudRunResourceState_UNKNOWN, strings.Join(s.UnknownMessages, "; ")
	}

	mustPassConditions := TypeHealthyServiceConditions
	if s.Kind == KindRevision {
		mustPassConditions = TypeHealthyRevisionConditions
	}
	for k := range mustPassConditions {
		if _, ok := s.TrueTypes[k]; !ok {
			return model.CloudRunResourceState_UNKNOWN, fmt.Sprintf("Could not check status field %q", k)
		}
	}
	return model.CloudRunResourceState_HEALTHY, ""
}
