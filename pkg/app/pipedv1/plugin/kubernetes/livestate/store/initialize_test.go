// Copyright 2026 The PipeCD Authors.
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
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

func initializeTestObject(id string) *unstructured.Unstructured {
	u := &unstructured.Unstructured{Object: map[string]any{"apiVersion": "v1", "kind": "ConfigMap"}}
	u.SetName(id)
	u.SetUID(types.UID(id))
	u.SetAnnotations(map[string]string{provider.LabelApplication: "app", provider.LabelManagedBy: provider.ManagedByPiped})
	return u
}

func initializeTestManifest(id string) provider.Manifest {
	return provider.FromUnstructured(initializeTestObject(id))
}

func TestDeployTargetResources_initializeConcurrent(t *testing.T) {
	t.Parallel()

	s := newDeployTargetResources("test")
	for i := 0; i < 3000; i++ {
		s.onAdd(initializeTestManifest(fmt.Sprint(i)))
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3000; i++ {
			updated := initializeTestObject(fmt.Sprint(i))
			updated.SetName("updated-" + fmt.Sprint(i))
			s.onUpdate(initializeTestManifest(fmt.Sprint(i)), provider.FromUnstructured(updated))
		}
	}()
	s.initialize()
	wg.Wait()
	resources, err := s.Livestate(context.Background(), "app")
	require.NoError(t, err)
	require.Len(t, resources, 3000)
	ids := make(map[types.UID]bool, len(resources))
	for _, resource := range resources {
		ids[resource.UID()] = true
		assert.Equal(t, "updated-"+string(resource.UID()), resource.Key().Name())
	}
	for i := 0; i < 3000; i++ {
		assert.True(t, ids[types.UID(fmt.Sprint(i))])
	}
}

func TestDeployTargetResources_initialize(t *testing.T) {
	t.Parallel()

	s := newDeployTargetResources("target")
	child := initializeTestObject("child")
	child.SetAnnotations(nil)
	child.SetOwnerReferences([]metav1.OwnerReference{{UID: "parent", APIVersion: "apps/v1", Kind: "Deployment", Name: "parent"}})
	s.onAdd(provider.FromUnstructured(child))
	s.onAdd(initializeTestManifest("orphan"))
	orphan := initializeTestObject("unowned")
	orphan.SetAnnotations(nil)
	s.onAdd(provider.FromUnstructured(orphan))
	s.onAdd(initializeTestManifest("parent"))
	s.onDelete(initializeTestManifest("orphan"))
	s.initialize()
	got, err := s.Livestate(context.Background(), "app")
	require.NoError(t, err)
	ids := make([]types.UID, 0, len(got))
	for _, m := range got {
		ids = append(ids, m.UID())
	}
	assert.ElementsMatch(t, []types.UID{"child", "parent"}, ids)
	s.mu.RLock()
	assert.Len(t, s.resources, 2)
	s.mu.RUnlock()
}

func TestDeployTargetResources_initializeConcurrentDelete(t *testing.T) {
	t.Parallel()

	s := newDeployTargetResources("target")
	for i := 0; i < 3000; i++ {
		s.onAdd(initializeTestManifest(fmt.Sprint(i)))
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3000; i++ {
			s.onDelete(initializeTestManifest(fmt.Sprint(i)))
		}
	}()
	s.initialize()
	wg.Wait()
	got, err := s.Livestate(context.Background(), "app")
	require.NoError(t, err)
	assert.Empty(t, got)
	s.mu.RLock()
	assert.Empty(t, s.resources)
	s.mu.RUnlock()
}

func TestDeployTargetResources_initializeConcurrentAdd(t *testing.T) {
	t.Parallel()

	s := newDeployTargetResources("target")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3000; i++ {
			s.onAdd(initializeTestManifest(fmt.Sprint(i)))
		}
	}()
	s.initialize()
	wg.Wait()
	got, err := s.Livestate(context.Background(), "app")
	require.NoError(t, err)
	require.Len(t, got, 3000)
	ids := make(map[types.UID]bool, len(got))
	for _, m := range got {
		ids[m.UID()] = true
	}
	for i := 0; i < 3000; i++ {
		assert.True(t, ids[types.UID(fmt.Sprint(i))])
	}
	s.mu.RLock()
	assert.Len(t, s.resources, 3000)
	s.mu.RUnlock()
}

func TestDeployTargetResources_initializeApplicationChange(t *testing.T) {
	t.Parallel()

	s := newDeployTargetResources("target")
	old := initializeTestObject("resource")
	s.onAdd(provider.FromUnstructured(old))
	s.initialize()
	updated := old.DeepCopy()
	annotations := updated.GetAnnotations()
	annotations[provider.LabelApplication] = "other"
	updated.SetAnnotations(annotations)
	s.onUpdate(provider.FromUnstructured(old), provider.FromUnstructured(updated))
	s.initialize()
	before, err := s.Livestate(context.Background(), "app")
	require.NoError(t, err)
	assert.Empty(t, before)
	after, err := s.Livestate(context.Background(), "other")
	require.NoError(t, err)
	require.Len(t, after, 1)
	assert.Equal(t, types.UID("resource"), after[0].UID())
	assert.Equal(t, "other", after[0].ApplicationID())
	managed := s.getApplicationResources("other").getManagedResources()
	require.Len(t, managed, 1)
	assert.Equal(t, "other", managed[0].ApplicationID())
}

func TestDeployTargetResources_initializeMixed(t *testing.T) {
	t.Parallel()

	s := newDeployTargetResources("target")
	for i := 0; i < 1000; i++ {
		s.onAdd(initializeTestManifest(fmt.Sprint(i)))
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := 0; i < 300; i++ {
			s.onDelete(initializeTestManifest(fmt.Sprint(i)))
		}
	}()
	go func() {
		defer wg.Done()
		for i := 300; i < 1000; i++ {
			u := initializeTestObject(fmt.Sprint(i))
			u.SetName("updated-" + fmt.Sprint(i))
			s.onUpdate(initializeTestManifest(fmt.Sprint(i)), provider.FromUnstructured(u))
		}
	}()
	go func() {
		defer wg.Done()
		for i := 1000; i < 1500; i++ {
			s.onAdd(initializeTestManifest(fmt.Sprint(i)))
		}
	}()
	s.initialize()
	wg.Wait()
	got, err := s.Livestate(context.Background(), "app")
	require.NoError(t, err)
	require.Len(t, got, 1200)
	ids := make(map[types.UID]bool, len(got))
	for _, m := range got {
		ids[m.UID()] = true
		if id, err := strconv.Atoi(string(m.UID())); err == nil && id >= 300 && id < 1000 {
			assert.Equal(t, "updated-"+string(m.UID()), m.Key().Name())
		}
	}
	for i := 0; i < 1500; i++ {
		assert.Equal(t, i >= 300, ids[types.UID(fmt.Sprint(i))])
	}
	s.mu.RLock()
	assert.Len(t, s.resources, 1200)
	s.mu.RUnlock()
}
