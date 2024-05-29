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

package kubernetes

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	provider "github.com/pipe-cd/pipecd/pkg/app/pipedv1/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type appNodes struct {
	appID         string
	managingNodes map[string]node
	dependedNodes map[string]node
	version       model.ApplicationLiveStateVersion
	mu            sync.RWMutex
}

type node struct {
	// The unique identifier of the resource generated by Kubernetes.
	uid          string
	appID        string
	key          provider.ResourceKey
	unstructured *unstructured.Unstructured
	state        model.KubernetesResourceState
}

func (n node) Manifest() provider.Manifest {
	return provider.MakeManifest(n.key, n.unstructured)
}

func (a *appNodes) addManagingResource(uid string, key provider.ResourceKey, obj *unstructured.Unstructured, now time.Time) (model.KubernetesResourceStateEvent, bool) {
	// Some resources in Kubernetes (e.g. Deployment) are producing multiple keys
	// for the same uid. So we use the configured original API version to ignore them.
	originalAPIVersion := obj.GetAnnotations()[provider.LabelOriginalAPIVersion]
	if originalAPIVersion != key.APIVersion {
		return model.KubernetesResourceStateEvent{}, false
	}

	n := node{
		uid:          uid,
		appID:        a.appID,
		key:          key,
		unstructured: obj,
		state:        provider.MakeKubernetesResourceState(uid, key, obj, now),
	}

	a.mu.Lock()
	oriNode, hasOriNode := a.managingNodes[uid]
	version := a.version
	a.managingNodes[uid] = n
	a.updateVersion(now)
	a.mu.Unlock()

	// No diff compared to previous state.
	if hasOriNode && !oriNode.state.HasDiff(n.state) {
		return model.KubernetesResourceStateEvent{}, false
	}

	return model.KubernetesResourceStateEvent{
		Id:              uuid.New().String(),
		ApplicationId:   a.appID,
		Type:            model.KubernetesResourceStateEvent_ADD_OR_UPDATED,
		State:           &n.state,
		SnapshotVersion: &version,
		CreatedAt:       now.Unix(),
	}, true
}

func (a *appNodes) deleteManagingResource(uid string, _ provider.ResourceKey, now time.Time) (model.KubernetesResourceStateEvent, bool) {
	a.mu.Lock()
	n, ok := a.managingNodes[uid]
	if !ok {
		a.mu.Unlock()
		return model.KubernetesResourceStateEvent{}, false
	}

	version := a.version
	delete(a.managingNodes, uid)
	a.updateVersion(now)
	a.mu.Unlock()

	return model.KubernetesResourceStateEvent{
		Id:              uuid.New().String(),
		ApplicationId:   a.appID,
		Type:            model.KubernetesResourceStateEvent_DELETED,
		State:           &n.state,
		SnapshotVersion: &version,
		CreatedAt:       now.Unix(),
	}, true
}

func (a *appNodes) addDependedResource(uid string, key provider.ResourceKey, obj *unstructured.Unstructured, now time.Time) (model.KubernetesResourceStateEvent, bool) {
	n := node{
		uid:          uid,
		appID:        a.appID,
		key:          key,
		unstructured: obj,
		state:        provider.MakeKubernetesResourceState(uid, key, obj, now),
	}

	a.mu.Lock()
	oriNode, hasOriNode := a.dependedNodes[uid]
	version := a.version
	a.dependedNodes[uid] = n
	a.updateVersion(now)
	a.mu.Unlock()

	// No diff compared to previous state.
	if hasOriNode && !oriNode.state.HasDiff(n.state) {
		return model.KubernetesResourceStateEvent{}, false
	}

	return model.KubernetesResourceStateEvent{
		Id:              uuid.New().String(),
		ApplicationId:   a.appID,
		Type:            model.KubernetesResourceStateEvent_ADD_OR_UPDATED,
		State:           &n.state,
		SnapshotVersion: &version,
		CreatedAt:       now.Unix(),
	}, true
}

func (a *appNodes) deleteDependedResource(uid string, _ provider.ResourceKey, now time.Time) (model.KubernetesResourceStateEvent, bool) {
	a.mu.Lock()
	n, ok := a.dependedNodes[uid]
	if !ok {
		a.mu.Unlock()
		return model.KubernetesResourceStateEvent{}, false
	}

	version := a.version
	delete(a.dependedNodes, uid)
	a.updateVersion(now)
	a.mu.Unlock()

	return model.KubernetesResourceStateEvent{
		Id:              uuid.New().String(),
		ApplicationId:   a.appID,
		Type:            model.KubernetesResourceStateEvent_DELETED,
		State:           &n.state,
		SnapshotVersion: &version,
		CreatedAt:       now.Unix(),
	}, true
}

func (a *appNodes) getManagingNodes() map[string]node {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.managingNodes
}

func (a *appNodes) getNodes() (map[string]node, model.ApplicationLiveStateVersion) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var (
		version = a.version
		nodes   = make(map[string]node, len(a.managingNodes)+len(a.dependedNodes))
	)
	for k, n := range a.dependedNodes {
		nodes[k] = n
	}
	for k, n := range a.managingNodes {
		nodes[k] = n
	}
	return nodes, version
}

func (a *appNodes) updateVersion(now time.Time) {
	if a.version.Timestamp == now.Unix() {
		a.version.Index++
		return
	}

	a.version.Timestamp = now.Unix()
	a.version.Index = 0
}
