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

package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationCollection struct {
	requestedBy Commander
}

func (a *applicationCollection) Kind() string {
	return "Application"
}

func (a *applicationCollection) Factory() Factory {
	return func() interface{} {
		return &model.Application{}
	}
}

func (a *applicationCollection) ListInUsedShards() []Shard {
	return []Shard{
		ClientShard,
		AgentShard,
	}
}

func (a *applicationCollection) GetUpdatableShard() (Shard, error) {
	switch a.requestedBy {
	case WebCommander, PipectlCommander:
		return ClientShard, nil
	case PipedCommander:
		return AgentShard, nil
	default:
		return "", ErrUnsupported
	}
}

func (a *applicationCollection) Decode(e interface{}, parts map[Shard][]byte) error {
	errFmt := "failed while decode Application object: %s"

	if len(parts) != len(a.ListInUsedShards()) {
		return fmt.Errorf(errFmt, "shards count not matched")
	}

	app, ok := e.(*model.Application)
	if !ok {
		return fmt.Errorf(errFmt, "type not matched")
	}

	var (
		kind             model.ApplicationKind
		name             string
		pipedID          string
		configFilename   string
		platformProvider string
		updatedAt        int64
	)
	for shard, p := range parts {
		if err := json.Unmarshal(p, &app); err != nil {
			return err
		}
		if updatedAt < app.UpdatedAt {
			updatedAt = app.UpdatedAt
		}
		// Values of below fields from ClientShard have a higher priority:
		// - PipedId
		// - GitPath.ConfigFilename
		// - PlatformProvider
		if shard == ClientShard {
			pipedID = app.PipedId
			configFilename = app.GitPath.ConfigFilename
			platformProvider = app.PlatformProvider
		}
		// Values of below fields from AgentShard have a higher priority:
		// - Kind
		// - Name
		if shard == AgentShard {
			kind = app.Kind
			name = app.Name
		}
	}

	app.Kind = kind
	app.Name = name
	app.PipedId = pipedID
	app.PlatformProvider = platformProvider
	app.GitPath.ConfigFilename = configFilename
	app.UpdatedAt = updatedAt
	return nil
}

func (a *applicationCollection) Encode(e interface{}) (map[Shard][]byte, error) {
	const errFmt = "failed while encode Application object: %s"

	me, ok := e.(*model.Application)
	if !ok {
		return nil, fmt.Errorf(errFmt, "type not matched")
	}

	// TODO: Find a way to generate function to build this kind of object by specifying a field tag in the proto file.
	// For example:
	// ```proto
	// message Application {
	//   // The generated unique identifier.
	//   string id = 1 [(validate.rules).string.min_len = 1, shard=client];
	// ```
	clientShardStruct := model.Application{
		// Fields which required in all shard for validation on update.
		Id:        me.Id,
		ProjectId: me.ProjectId,
		CreatedAt: me.CreatedAt,
		UpdatedAt: me.UpdatedAt,
		// Fields which exist in both AgentShard and ClientShard but AgentShard has
		// a higher priority since those fields can only be updated by PipedCommander.
		Kind: me.Kind,
		Name: me.Name,
		// Fields which exist in both AgentShard and ClientShard but ClientShard has
		// a higher priority since those fields can only be updated by WebCommander.
		PipedId:          me.PipedId,
		PlatformProvider: me.PlatformProvider,
		// Note: Only GitPath.ConfigFilename is changeable.
		GitPath: me.GitPath,
		// Fields which exist only in ClientShard.
		Disabled:  me.Disabled,
		Deleted:   me.Deleted,
		DeletedAt: me.DeletedAt,
	}
	cdata, err := json.Marshal(&clientShardStruct)
	if err != nil {
		return nil, fmt.Errorf(errFmt, "unable to marshal entity data")
	}

	agentShardStruct := model.Application{
		// Fields which required in all shard for validation on update.
		Id:        me.Id,
		ProjectId: me.ProjectId,
		CreatedAt: me.CreatedAt,
		UpdatedAt: me.UpdatedAt,
		// Fields which exist in both AgentShard and ClientShard but AgentShard has
		// a higher priority since those fields can only be updated by PipedCommander.
		Kind: me.Kind,
		Name: me.Name,
		// Fields which exist in both AgentShard and ClientShard but ClientShard has
		// a higher priority since those fields can only be updated by WebCommander.
		PipedId:          me.PipedId,
		GitPath:          me.GitPath,
		PlatformProvider: me.PlatformProvider,
		// Fields which exist only in AgentShard.
		Description:                      me.Description,
		Labels:                           me.Labels,
		SyncState:                        me.SyncState,
		Deploying:                        me.Deploying,
		MostRecentlySuccessfulDeployment: me.MostRecentlySuccessfulDeployment,
		MostRecentlyTriggeredDeployment:  me.MostRecentlyTriggeredDeployment,
	}
	adata, err := json.Marshal(&agentShardStruct)
	if err != nil {
		return nil, fmt.Errorf(errFmt, "unable to marshal entity data")
	}

	return map[Shard][]byte{
		ClientShard: cdata,
		AgentShard:  adata,
	}, nil
}

type ApplicationStore interface {
	Add(ctx context.Context, app *model.Application) error
	Get(ctx context.Context, id string) (*model.Application, error)
	List(ctx context.Context, opts ListOptions) ([]*model.Application, string, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, app *model.Application) error
	Enable(ctx context.Context, id string) error
	Disable(ctx context.Context, id string) error
	UpdateSyncState(ctx context.Context, id string, syncState *model.ApplicationSyncState) error
	UpdateMostRecentDeployment(ctx context.Context, id string, status model.DeploymentStatus, deployment *model.ApplicationDeploymentReference) error
	UpdateConfigFilename(ctx context.Context, id, configFilename string) error
	UpdateDeployingStatus(ctx context.Context, id string, deploying bool) error
	UpdateBasicInfo(ctx context.Context, id, name, description string, labels map[string]string) error
	UpdateConfiguration(ctx context.Context, id, pipedID, platformProvider, configFilename string) error
	UpdatePlatformProvider(ctx context.Context, id string, provider string) error
}

type applicationStore struct {
	backend
	commander Commander
	nowFunc   func() time.Time
}

func NewApplicationStore(ds DataStore, c Commander) ApplicationStore {
	return &applicationStore{
		backend: backend{
			ds:  ds,
			col: &applicationCollection{requestedBy: c},
		},
		commander: c,
		nowFunc:   time.Now,
	}
}

func (s *applicationStore) Add(ctx context.Context, app *model.Application) error {
	now := s.nowFunc().Unix()
	if app.CreatedAt == 0 {
		app.CreatedAt = now
	}
	if app.UpdatedAt == 0 {
		app.UpdatedAt = now
	}
	if err := app.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, s.col, app.Id, app)
}

func (s *applicationStore) Get(ctx context.Context, id string) (*model.Application, error) {
	var entity model.Application
	if err := s.ds.Get(ctx, s.col, id, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *applicationStore) List(ctx context.Context, opts ListOptions) ([]*model.Application, string, error) {
	it, err := s.ds.Find(ctx, s.col, opts)
	if err != nil {
		return nil, "", err
	}
	apps := make([]*model.Application, 0)
	for {
		var app model.Application
		err := it.Next(&app)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, "", err
		}
		apps = append(apps, &app)
	}

	// In case there is no more elements found, cursor should be set to empty too.
	if len(apps) == 0 {
		return apps, "", nil
	}
	cursor, err := it.Cursor()
	if err != nil {
		return nil, "", err
	}
	return apps, cursor, nil
}

func (s *applicationStore) Update(ctx context.Context, app *model.Application) error {
	// UpdateXXX does not fail, so Update will not result in an error.
	s.UpdateBasicInfo(ctx, app.Id, app.Name, app.Description, app.Labels)
	s.UpdateConfiguration(ctx, app.Id, app.PipedId, app.PlatformProvider, app.GitPath.ConfigFilename)
	s.UpdateSyncState(ctx, app.Id, app.SyncState)
	return nil
}

func (s *applicationStore) Delete(ctx context.Context, id string) error {
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		now := s.nowFunc().Unix()
		app := e.(*model.Application)
		app.Deleted = true
		app.Disabled = true
		app.DeletedAt = now
		app.UpdatedAt = now
		return nil
	})
}

func (s *applicationStore) Enable(ctx context.Context, id string) error {
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		app := e.(*model.Application)
		if app.Deleted {
			return fmt.Errorf("cannot enable a deleted application: %w", ErrInvalidArgument)
		}
		app.Disabled = false
		app.UpdatedAt = s.nowFunc().Unix()
		return nil
	})
}

func (s *applicationStore) Disable(ctx context.Context, id string) error {
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		app := e.(*model.Application)
		if app.Deleted {
			return fmt.Errorf("cannot disable a deleted application: %w", ErrInvalidArgument)
		}
		app.Disabled = true
		app.UpdatedAt = s.nowFunc().Unix()
		return nil
	})
}

func (s *applicationStore) update(ctx context.Context, id string, updater func(*model.Application) error) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		a := e.(*model.Application)
		if a.Deleted {
			return fmt.Errorf("cannot update a deleted application: %w", ErrInvalidArgument)
		}
		if err := updater(a); err != nil {
			return err
		}
		a.UpdatedAt = now
		return a.Validate()
	})
}

func (s *applicationStore) UpdateSyncState(ctx context.Context, id string, syncState *model.ApplicationSyncState) error {
	return s.update(ctx, id, func(a *model.Application) error {
		a.SyncState = syncState
		return nil
	})
}

func (s *applicationStore) UpdateMostRecentDeployment(ctx context.Context, id string, status model.DeploymentStatus, deployment *model.ApplicationDeploymentReference) error {
	return s.update(ctx, id, func(a *model.Application) error {
		switch status {
		case model.DeploymentStatus_DEPLOYMENT_SUCCESS:
			a.MostRecentlySuccessfulDeployment = deployment
		case model.DeploymentStatus_DEPLOYMENT_PENDING:
			a.MostRecentlyTriggeredDeployment = deployment
		}
		return nil
	})
}

func (s *applicationStore) UpdateConfigFilename(ctx context.Context, id, configFilename string) error {
	return s.update(ctx, id, func(app *model.Application) error {
		app.GitPath.ConfigFilename = configFilename
		return nil
	})
}

func (s *applicationStore) UpdateDeployingStatus(ctx context.Context, id string, deploying bool) error {
	return s.update(ctx, id, func(app *model.Application) error {
		app.Deploying = deploying
		return nil
	})
}

func (s *applicationStore) UpdateBasicInfo(ctx context.Context, id, name, description string, labels map[string]string) error {
	return s.update(ctx, id, func(app *model.Application) error {
		app.Name = name
		app.Description = description
		app.Labels = labels
		return nil
	})
}

func (s *applicationStore) UpdateConfiguration(ctx context.Context, id, pipedID, platformProvider, configFilename string) error {
	return s.update(ctx, id, func(app *model.Application) error {
		app.PipedId = pipedID
		app.PlatformProvider = platformProvider
		app.CloudProvider = platformProvider
		app.GitPath.ConfigFilename = configFilename
		return nil
	})
}

func (s *applicationStore) UpdatePlatformProvider(ctx context.Context, id string, provider string) error {
	return s.update(ctx, id, func(app *model.Application) error {
		app.PlatformProvider = provider
		return nil
	})
}
