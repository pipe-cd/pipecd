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

package applicationlivestatestore

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type Store interface {
	// GetStateSnapshot get the specified application live state snapshot.
	GetStateSnapshot(ctx context.Context, applicationID string) (*model.ApplicationLiveStateSnapshot, error)
	// PutStateSnapshot updates completely the specified application live state snapshot with input.
	PutStateSnapshot(ctx context.Context, snapshot *model.ApplicationLiveStateSnapshot) error
	// PatchKubernetesApplicationLiveState updates the kubernetes resource state in the application live state snapshot.
	PatchKubernetesApplicationLiveState(ctx context.Context, events []*model.KubernetesResourceStateEvent)
}

type store struct {
	backend *applicationLiveStateFileStore
	cache   *applicationLiveStateCache
	logger  *zap.Logger
}

func NewStore(fs filestore.Store, c cache.Cache, logger *zap.Logger) Store {
	return &store{
		backend: &applicationLiveStateFileStore{
			backend: fs,
		},
		cache: &applicationLiveStateCache{
			backend: c,
		},
		logger: logger.Named("application-live-state-store"),
	}
}

func (s *store) GetStateSnapshot(ctx context.Context, applicationID string) (*model.ApplicationLiveStateSnapshot, error) {
	cacheResp, err := s.cache.Get(applicationID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		s.logger.Error("failed to get application live state from cache", zap.Error(err))
	}
	if cacheResp != nil {
		return cacheResp, nil
	}

	fileResp, err := s.backend.Get(ctx, applicationID)
	if err != nil {
		if !errors.Is(err, filestore.ErrNotFound) {
			s.logger.Error("failed to get application live state from filestore",
				zap.String("application-id", applicationID),
				zap.Error(err),
			)
		}
		return nil, err
	}

	if err := s.cache.Put(applicationID, fileResp); err != nil {
		s.logger.Error("failed to put application live state to cache", zap.Error(err))
	}
	return fileResp, nil
}

func (s *store) PutStateSnapshot(ctx context.Context, snapshot *model.ApplicationLiveStateSnapshot) error {
	if err := s.backend.Put(ctx, snapshot.ApplicationId, snapshot); err != nil {
		s.logger.Error("failed to put application live state snapshot to filestore", zap.Error(err))
		return err
	}

	if err := s.cache.Put(snapshot.ApplicationId, snapshot); err != nil {
		s.logger.Error("failed to put application live state snapshot to cache", zap.Error(err))
	}
	return nil
}

func (s *store) PatchKubernetesApplicationLiveState(ctx context.Context, events []*model.KubernetesResourceStateEvent) {
	snapshots := make(map[string]*model.ApplicationLiveStateSnapshot)
	for _, ev := range events {
		snapshot, ok := snapshots[ev.ApplicationId]
		if !ok {
			// Ignore error because logging is already doing in GetStateSnapshot.
			ss, _ := s.GetStateSnapshot(ctx, ev.ApplicationId)
			if ss == nil {
				s.logger.Warn("application live state snapshot was not found",
					zap.String("event-id", ev.Id),
					zap.String("application-id", ev.ApplicationId),
				)
				continue
			}
			ss.DetermineAppHealthStatus()
			snapshot = ss
			snapshots[ev.ApplicationId] = ss
		}
		if ev.SnapshotVersion.IsBefore(*snapshot.Version) {
			continue
		}
		switch ev.Type {
		case model.KubernetesResourceStateEvent_ADD_OR_UPDATED:
			snapshot.Kubernetes.Resources = mergeKubernetesResourceStatesOnAddOrUpdated(snapshot.Kubernetes.Resources, ev)
		case model.KubernetesResourceStateEvent_DELETED:
			snapshot.Kubernetes.Resources = mergeKubernetesResourceStatesOnDeleted(snapshot.Kubernetes.Resources, ev)
		}
	}

	for applicationID, snapshot := range snapshots {
		if err := s.cache.Put(applicationID, snapshot); err != nil {
			s.logger.Error("failed to put application live state snapshot to cache",
				zap.String("application-id", applicationID),
				zap.Error(err),
			)
		}
	}
}

func mergeKubernetesResourceStatesOnAddOrUpdated(prevs []*model.KubernetesResourceState, event *model.KubernetesResourceStateEvent) []*model.KubernetesResourceState {
	var found bool
	news := make([]*model.KubernetesResourceState, 0, len(prevs)+1)
	for _, rs := range prevs {
		if rs.Id != event.State.Id {
			news = append(news, rs)
			continue
		}

		news = append(news, event.State)
		found = true
	}

	if !found {
		news = append(news, event.State)
	}
	return news
}

func mergeKubernetesResourceStatesOnDeleted(prevs []*model.KubernetesResourceState, event *model.KubernetesResourceStateEvent) []*model.KubernetesResourceState {
	remains := make([]*model.KubernetesResourceState, 0, len(prevs))
	for _, state := range prevs {
		if state.Id != event.State.Id {
			remains = append(remains, state)
		}
	}
	return remains
}
