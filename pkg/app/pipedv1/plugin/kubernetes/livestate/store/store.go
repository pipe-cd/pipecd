package store

import (
	"context"
	"sync"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

var _ resourceEventHandler = (*store)(nil)

// applicationResources is a collection of resources that belong to the same application.
// It is used to store the resources and to calculate the livestate of the application.
type applicationResources struct {
	deployTarget string
	resources    map[types.UID]provider.Manifest
	mu           sync.RWMutex
}

func newApplicationResources(deployTarget string) *applicationResources {
	return &applicationResources{
		deployTarget: deployTarget,
		resources:    make(map[types.UID]provider.Manifest),
	}
}

// addResource adds a resource to the application resources.
func (a *applicationResources) addResource(resource provider.Manifest) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.resources[resource.UID()] = resource
}

// removeResource removes a resource from the application resources.
func (a *applicationResources) removeResource(resource provider.Manifest) {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.resources, resource.UID())
}

// livestate returns the livestate of the application resources.
func (a *applicationResources) livestate() []sdk.ResourceState {
	a.mu.RLock()
	defer a.mu.RUnlock()

	resources := make([]sdk.ResourceState, 0, len(a.resources))
	for _, resource := range a.resources {
		resources = append(resources, resource.ToResourceState(a.deployTarget))
	}

	return resources
}

// store is a store for the application resources for one deploy target.
type store struct {
	// The deploy target of this store.
	deployTarget string
	// The map with the key is "application ID" and the value is "application resources".
	applications map[string]*applicationResources
	// The map with the key is "resource UID" and the value is "application ID".
	applicationIDReferences map[types.UID]string

	mu sync.RWMutex
}

// newStore creates a new store for the given deploy target.
func newStore(deployTarget string) *store {
	return &store{
		deployTarget:            deployTarget,
		applications:            make(map[string]*applicationResources),
		applicationIDReferences: make(map[types.UID]string),
	}
}

// getApplicationResources returns the application resources by the application ID.
func (s *store) getApplicationResources(appID string) *applicationResources {
	s.mu.RLock()
	defer s.mu.RUnlock()

	app, ok := s.applications[appID]
	if !ok {
		app = newApplicationResources(s.deployTarget)
		s.applications[appID] = app
	}

	return app
}

// getApplicationIDByResource returns the application ID by the resource.
func (s *store) getApplicationIDByResource(resource provider.Manifest) (string, bool) {
	appID := resource.ApplicationID()
	if appID != "" {
		return appID, true
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if appID, ok := s.applicationIDReferences[resource.UID()]; ok {
		return appID, true
	}

	ownerRefs := resource.OwnerReferences()
	for _, ref := range ownerRefs {
		if appID, ok := s.applicationIDReferences[ref]; ok {
			return appID, true
		}
	}

	return "", false
}

// addResource adds a resource to the store.
func (s *store) addResource(resource provider.Manifest) {
	appID, ok := s.getApplicationIDByResource(resource)
	if !ok {
		return
	}

	app := s.getApplicationResources(appID)

	app.addResource(resource)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.applicationIDReferences[resource.UID()] = appID
}

// removeResource removes a resource from the store.
func (s *store) removeResource(resource provider.Manifest) {
	appID, ok := s.getApplicationIDByResource(resource)
	if !ok {
		return
	}

	app := s.getApplicationResources(appID)

	app.removeResource(resource)

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.applicationIDReferences, resource.UID())
}

// Livestate returns the livestate of the application.
func (s *store) Livestate(_ context.Context, appID string) ([]sdk.ResourceState, error) {
	app := s.getApplicationResources(appID)

	return app.livestate(), nil
}

// matchResource returns true if the resource is managed by the application.
func (s *store) matchResource(resource provider.Manifest) bool {
	_, ok := s.getApplicationIDByResource(resource)
	return ok
}

// onAdd adds a resource to the store.
func (s *store) onAdd(resource provider.Manifest) {
	s.addResource(resource)
}

// onUpdate updates a resource in the store.
func (s *store) onUpdate(old, new provider.Manifest) {
	s.removeResource(old)
	s.addResource(new)
}

// onDelete removes a resource from the store.
func (s *store) onDelete(resource provider.Manifest) {
	s.removeResource(resource)
}
