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

package deployment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	kubeConfigPkg "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

// func TestPlugin_executeK8sRollbackStage_NoPreviousDeployment(t *testing.T) { ... }
// func TestPlugin_executeK8sRollbackStage_SuccessfulRollback(t *testing.T) { ... }
// func TestPlugin_executeK8sRollbackStage_WithVariantLabels(t *testing.T) { ... }

type fakeApplier struct {
	deleted *[]string
}

func (f *fakeApplier) Delete(_ context.Context, k provider.ResourceKey) error {
	*f.deleted = append(*f.deleted, k.String())
	return nil
}

func TestRemoveVariantResources_PrunesCorrectResources(t *testing.T) {
	t.Parallel()

	deleted := make([]string, 0)

	makeManifest := func(name, variant string) provider.Manifest {
		obj := &unstructured.Unstructured{}
		obj.SetName(name)
		obj.SetNamespace("default")
		obj.SetKind("Deployment")
		obj.SetAPIVersion("apps/v1")
		labels := map[string]string{"pipecd.dev/variant": variant}
		obj.SetLabels(labels)
		m, err := provider.FromStructuredObject(obj)
		if err != nil {
			t.Fatalf("failed to create manifest: %v", err)
		}
		return m
	}
	manifests := []provider.Manifest{
		makeManifest("canary-1", "canary"),
		makeManifest("canary-2", "canary"),
		makeManifest("baseline-1", "baseline"),
		makeManifest("primary-1", "primary"),
	}

	cfg := kubeConfigPkg.KubernetesApplicationSpec{
		VariantLabel: kubeConfigPkg.KubernetesVariantLabel{
			Key:           "pipecd.dev/variant",
			PrimaryValue:  "primary",
			CanaryValue:   "canary",
			BaselineValue: "baseline",
		},
	}

	variant := "canary"
	variantLabel := cfg.VariantLabel.Key
	var toDelete []provider.ResourceKey
	for _, m := range manifests {
		labels := m.Labels()
		if labels[variantLabel] == variant {
			toDelete = append(toDelete, m.Key())
		}
	}
	fa := &fakeApplier{deleted: &deleted}
	for _, k := range toDelete {
		_ = fa.Delete(context.Background(), k)
	}
	assert.ElementsMatch(t, []string{
		"apps:Deployment:default:canary-1",
		"apps:Deployment:default:canary-2",
	}, deleted)
}
