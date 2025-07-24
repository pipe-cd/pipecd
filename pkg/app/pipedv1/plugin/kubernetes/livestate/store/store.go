// Copyright 2025 The PipeCD Authors.
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

package store

import (
	"context"
	"fmt"
	"sync"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

var _ resourceEventHandler = (*deployTargetResources)(nil)

type Store struct {
	// The map with the key is "name" of the deploy target and the value is "deploy target resources".
	// This map is immutable after the store is created.
	deployTargetResources map[string]*deployTargetResources
}

func Run(ctx context.Context, deployTargets map[string]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], logger *zap.Logger) (*Store, error) {
	deployTargetResources := make(map[string]*deployTargetResources)

	for _, deployTarget := range deployTargets {
		dtr := newDeployTargetResources(deployTarget.Name)
		kubeConfig, err := clientcmd.BuildConfigFromFlags(deployTarget.Config.MasterURL, deployTarget.Config.KubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to build kube config: %w", err)
		}

		rf := reflector{
			namespace:            deployTarget.Config.AppStateInformer.Namespace,
			targetMatcher:        newResourceMatcher(deployTarget.Config.AppStateInformer),
			resourceEventHandler: dtr,
			kubeConfig:           kubeConfig,
			logger:               logger.Named(fmt.Sprintf("livestate-store-deploy-target-%s", deployTarget.Name)),
		}

		// Start the reflector and get the list of resources that is watched by the reflector.
		watchingResourceKinds, err := rf.start(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to start reflector: %w", err)
		}
		dtr.watchingResourceKinds = watchingResourceKinds

		deployTargetResources[deployTarget.Name] = dtr
	}

	s := &Store{
		deployTargetResources: deployTargetResources,
	}

	return s, nil
}

func (s *Store) Livestate(ctx context.Context, deployTargetName string, appID string) ([]provider.Manifest, error) {
	dtr, ok := s.deployTargetResources[deployTargetName]
	if !ok {
		return nil, fmt.Errorf("deploy target %s not found", deployTargetName)
	}

	return dtr.Livestate(ctx, appID)
}

func (s *Store) ManagedResources(ctx context.Context, deployTargetName string, appID string) ([]provider.Manifest, error) {
	dtr, ok := s.deployTargetResources[deployTargetName]
	if !ok {
		return nil, fmt.Errorf("deploy target %s not found", deployTargetName)
	}

	app := dtr.getApplicationResources(appID)

	return app.getManagedResources(), nil
}

func (s *Store) WatchingResourceKinds(deployTargetName string) ([]schema.GroupVersionKind, error) {
	dtr, ok := s.deployTargetResources[deployTargetName]
	if !ok {
		return nil, fmt.Errorf("deploy target %s not found", deployTargetName)
	}

	return dtr.watchingResourceKinds, nil
}

// applicationResources is a collection of resources that belong to the same application.
// It is used to store the resources and to calculate the livestate of the application.
type applicationResources struct {
	deployTarget      string
	managedResources  map[types.UID]provider.Manifest
	dependedResources map[types.UID]provider.Manifest
	mu                sync.RWMutex
}

func newApplicationResources(deployTarget string) *applicationResources {
	return &applicationResources{
		deployTarget:      deployTarget,
		managedResources:  make(map[types.UID]provider.Manifest),
		dependedResources: make(map[types.UID]provider.Manifest),
	}
}

// addResource adds a resource to the application resources.
func (a *applicationResources) addResource(resource provider.Manifest) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if resource.IsManagedByPiped() {
		a.managedResources[resource.UID()] = resource
	} else {
		a.dependedResources[resource.UID()] = resource
	}
}

// removeResource removes a resource from the application resources.
func (a *applicationResources) removeResource(resource provider.Manifest) {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.managedResources, resource.UID())
	delete(a.dependedResources, resource.UID())
}

// livestate returns the livestate of the application resources.
func (a *applicationResources) livestate() []provider.Manifest {
	a.mu.RLock()
	defer a.mu.RUnlock()

	managedResources := make([]provider.Manifest, 0, len(a.managedResources))
	for _, resource := range a.managedResources {
		managedResources = append(managedResources, resource)
	}

	dependedResources := make([]provider.Manifest, 0, len(a.dependedResources))
	for _, resource := range a.dependedResources {
		dependedResources = append(dependedResources, resource)
	}

	return append(managedResources, dependedResources...)
}

func (a *applicationResources) getManagedResources() []provider.Manifest {
	a.mu.RLock()
	defer a.mu.RUnlock()

	resources := make([]provider.Manifest, 0, len(a.managedResources))
	for _, resource := range a.managedResources {
		resources = append(resources, resource)
	}

	return resources
}

// deployTargetResources is a deployTargetResources for the application resources for one deploy target.
type deployTargetResources struct {
	// The deploy target of this store.
	deployTarget string
	// The map with the key is "application ID" and the value is "application resources".
	applications map[string]*applicationResources
	// The map with the key is "resource UID" and the value is "application ID".
	applicationIDReferences map[types.UID]string

	// The list of resources that is watched by the reflector.
	// This is immutable after the reflector is started.
	watchingResourceKinds []schema.GroupVersionKind

	mu sync.RWMutex
}

// newDeployTargetResources creates a new deployTargetResources for the given deploy target.
func newDeployTargetResources(deployTarget string) *deployTargetResources {
	return &deployTargetResources{
		deployTarget:            deployTarget,
		applications:            make(map[string]*applicationResources),
		applicationIDReferences: make(map[types.UID]string),
	}
}

// getApplicationResources returns the application resources by the application ID.
func (s *deployTargetResources) getApplicationResources(appID string) *applicationResources {
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
func (s *deployTargetResources) getApplicationIDByResource(resource provider.Manifest) (string, bool) {
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
func (s *deployTargetResources) addResource(resource provider.Manifest) {
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
func (s *deployTargetResources) removeResource(resource provider.Manifest) {
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
func (s *deployTargetResources) Livestate(_ context.Context, appID string) ([]provider.Manifest, error) {
	app := s.getApplicationResources(appID)

	return app.livestate(), nil
}

// matchResource returns true if the resource is managed by the application.
func (s *deployTargetResources) matchResource(resource provider.Manifest) bool {
	_, ok := s.getApplicationIDByResource(resource)
	return ok
}

// onAdd adds a resource to the store.
func (s *deployTargetResources) onAdd(resource provider.Manifest) {
	s.addResource(resource)
}

// onUpdate updates a resource in the store.
func (s *deployTargetResources) onUpdate(old, new provider.Manifest) {
	s.removeResource(old)
	s.addResource(new)
}

// onDelete removes a resource from the store.
func (s *deployTargetResources) onDelete(resource provider.Manifest) {
	s.removeResource(resource)
}
