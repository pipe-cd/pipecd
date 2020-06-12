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

package kubernetes

import (
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	eventCacheSize        = 900
	eventCacheMaxSize     = 1000
	eventCacheCleanOffset = 50
)

type store struct {
	apps map[string]*appNodes
	// The map with the key is "resource's uid" and the value is "appResource".
	// Because the depended resource does not include the appID in its annotations
	// so this is used to determine the application of a depended resource.
	resources map[string]appResource
	mu        sync.RWMutex

	events         []model.KubernetesResourceStateEvent
	iterators      map[int]int
	nextIteratorID int
	eventMu        sync.Mutex
}

type appResource struct {
	appID    string
	owners   []metav1.OwnerReference
	resource *unstructured.Unstructured
}

func (s *store) initialize() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	// Try to determine the application ID of all resources.
	for uid, an := range s.resources {
		// Resource has already assigned into an application.
		if an.appID != "" {
			continue
		}
		appID := s.findAppIDByOwners(an.owners)
		if appID == "" {
			continue
		}

		// Add the missing resource into the dependedResources of the app.
		key := provider.MakeResourceKey(an.resource)
		s.apps[appID].addDependedResource(uid, key, an.resource, now)

		an.appID = appID
		s.resources[uid] = an
	}

	// Remove all resources which do not have appID.
	for uid, an := range s.resources {
		if an.appID == "" {
			delete(s.resources, uid)
		}
	}

	// Clean all initial events.
	s.events = nil
}

func (s *store) onAddResource(obj *unstructured.Unstructured) {
	var (
		uid    = string(obj.GetUID())
		appID  = obj.GetAnnotations()[provider.LabelApplication]
		key    = provider.MakeResourceKey(obj)
		owners = obj.GetOwnerReferences()
		now    = time.Now()
	)

	// If this is a resource managed by PipeCD
	// it must contain appID in its annotations and has no owners.
	if appID != "" && len(owners) == 0 {
		// When this obj is for a new application
		// we register a new application to the apps.
		s.mu.Lock()
		app, ok := s.apps[appID]
		if !ok {
			app = &appNodes{
				appID:         appID,
				managingNodes: make(map[string]node),
				dependedNodes: make(map[string]node),
				version: model.ApplicationLiveStateVersion{
					Timestamp: now.Unix(),
				},
			}
			s.apps[appID] = app
		}
		s.mu.Unlock()

		// Append the resource to the application's managingNodes.
		if event, ok := app.addManagingResource(uid, key, obj, now); ok {
			s.addEvent(event)
		}

		// And update the resources.
		s.mu.Lock()
		s.resources[uid] = appResource{appID: appID, owners: owners, resource: obj}
		s.mu.Unlock()
		return
	}

	// Try to determine the application ID by traveling its owners.
	s.mu.RLock()
	appID = s.findAppIDByOwners(owners)
	s.mu.RUnlock()

	// Append the resource to the application's dependedNodes.
	if appID != "" {
		s.mu.RLock()
		app, ok := s.apps[appID]
		s.mu.RUnlock()
		if ok {
			if event, ok := app.addDependedResource(uid, key, obj, now); ok {
				s.addEvent(event)
			}
		}
	}

	// And update the resources.
	s.mu.Lock()
	s.resources[uid] = appResource{appID: appID, owners: owners, resource: obj}
	s.mu.Unlock()
}

func (s *store) onUpdateResource(oldObj, obj *unstructured.Unstructured) {
	s.onAddResource(obj)
}

func (s *store) onDeleteResource(obj *unstructured.Unstructured) {
	var (
		uid    = string(obj.GetUID())
		appID  = obj.GetAnnotations()[provider.LabelApplication]
		key    = provider.MakeResourceKey(obj)
		owners = obj.GetOwnerReferences()
		now    = time.Now()
	)

	// If this is a resource managed by PipeCD
	// it must contain appID in its annotations and has no owners.
	if appID != "" && len(owners) == 0 {
		s.mu.Lock()
		delete(s.resources, uid)
		s.mu.Unlock()

		s.mu.RLock()
		app, ok := s.apps[appID]
		s.mu.RUnlock()
		if ok {
			if event, ok := app.deleteManagingResource(uid, key, now); ok {
				s.addEvent(event)
			}
		}
		return
	}

	// Try to determine the application ID by traveling its owners.
	s.mu.RLock()
	appID = s.findAppIDByOwners(owners)
	s.mu.RUnlock()

	// Delete the resource to the application's dependedNodes.
	s.mu.RLock()
	app, ok := s.apps[appID]
	s.mu.RUnlock()
	if ok {
		if event, ok := app.deleteDependedResource(uid, key, now); ok {
			s.addEvent(event)
		}
	}

	s.mu.Lock()
	delete(s.resources, uid)
	s.mu.Unlock()
}

func (s *store) getAppManagingNodes(appID string) map[string]node {
	s.mu.RLock()
	app, ok := s.apps[appID]
	s.mu.RUnlock()

	if !ok {
		return nil
	}
	return app.getManagingNodes()
}

func (s *store) getAppNodes(appID string) map[string]node {
	s.mu.RLock()
	app, ok := s.apps[appID]
	s.mu.RUnlock()

	if !ok {
		return nil
	}
	return app.getNodes()
}

func (s *store) findAppIDByOwners(owners []metav1.OwnerReference) string {
	for _, ref := range owners {
		owner, ok := s.resources[string(ref.UID)]
		// Owner does not present in the resources.
		if !ok {
			continue
		}
		// The owner is containing the appID.
		if owner.appID != "" {
			return owner.appID
		}
		// Try with the owners of the owner.
		if appID := s.findAppIDByOwners(owner.owners); appID != "" {
			return appID
		}
	}
	return ""
}

func (s *store) getAppLiveState(appID string) AppState {
	nodes := s.getAppNodes(appID)
	resources := make([]*model.KubernetesResourceState, 0, len(nodes))
	for _, n := range nodes {
		resources = append(resources, &n.state)
	}
	state := AppState{
		Resources: resources,
	}
	return state
}

func (s *store) GetAppLiveManifests(appID string) []provider.Manifest {
	s.mu.RLock()
	app, ok := s.apps[appID]
	s.mu.RUnlock()

	if !ok {
		return nil
	}
	nodes := app.getManagingNodes()
	manifests := make([]provider.Manifest, 0, len(nodes))
	for i := range nodes {
		manifests = append(manifests, nodes[i].Manifest())
	}
	return manifests
}

func (s *store) addEvent(event model.KubernetesResourceStateEvent) {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()

	s.events = append(s.events, event)
	if len(s.events) < eventCacheMaxSize {
		return
	}

	num := len(s.events) - eventCacheSize
	s.removeOldEvents(num)
}

func (s *store) nextEvents(iteratorID, maxNum int) []model.KubernetesResourceStateEvent {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()

	var (
		from   = s.iterators[iteratorID]
		to     = len(s.events)
		length = to - from
	)
	if length <= 0 {
		return nil
	}
	if length > maxNum {
		to = from + maxNum - 1
	}

	events := s.events[from:to]
	s.iterators[iteratorID] = to

	s.cleanStaleEvents()
	return events
}

func (s *store) cleanStaleEvents() {
	var min int
	for _, v := range s.iterators {
		if v < min {
			min = v
		}
	}
	if min < eventCacheCleanOffset {
		return
	}
	s.removeOldEvents(min)
}

func (s *store) removeOldEvents(num int) {
	if len(s.events) < num {
		return
	}
	s.events = s.events[num-1:]
	for k := range s.iterators {
		newIndex := s.iterators[k] - num
		if newIndex < 0 {
			newIndex = 0
		}
		s.iterators[k] = newIndex
	}
}

func (s *store) newEventIterator() EventIterator {
	s.eventMu.Lock()
	id := s.nextIteratorID
	s.nextIteratorID++
	s.eventMu.Unlock()

	return EventIterator{
		id:    id,
		store: s,
	}
}
