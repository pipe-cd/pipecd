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

	provider "github.com/kapetaniosci/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type store struct {
	apps map[string]*appNodes
	// The map with the key is "resource's uid" and the value is "appResource".
	// Because the depended resource does not include the appID in its annotations
	// so this is used to determine the application of a depended resource.
	resources map[string]appResource
	mu        sync.RWMutex

	events         []model.KubernetesResourceEvent
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
				updatedAt:     now,
			}
			s.apps[appID] = app
		}
		s.mu.Unlock()

		// Append the resource to the application's managingNodes.
		app.addManagingResource(uid, key, obj, now)

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
			app.addDependedResource(uid, key, obj, now)
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
			app.deleteManagingResource(uid, key, now)
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
		app.deleteDependedResource(uid, key, now)
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

func (s *store) nextEvents(iteratorID, maxNum int) []model.KubernetesResourceEvent {
	return nil
}

func (s *store) newEventIterator() EventIterator {
	s.eventMu.Lock()
	id := s.nextIteratorID
	s.nextIteratorID++
	s.eventMu.Unlock()

	return EventIterator{id: id}
}
