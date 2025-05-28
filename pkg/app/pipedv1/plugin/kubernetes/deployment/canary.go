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
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/yamlprocessor"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func (p *Plugin) executeK8sCanaryRolloutStage(_ context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], _ []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	input.Client.LogPersister().Error("Canary rollout is not yet implemented")
	return sdk.StageStatusFailure
}

func (p *Plugin) executeK8sCanaryCleanStage(_ context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], _ []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	input.Client.LogPersister().Error("Canary clean is not yet implemented")
	return sdk.StageStatusFailure
}

func findConfigMapManifests(manifests []provider.Manifest) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		if !m.IsConfigMap() {
			continue
		}
		out = append(out, m)
	}
	return out
}

func findSecretManifests(manifests []provider.Manifest) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		if !m.IsSecret() {
			continue
		}
		out = append(out, m)
	}
	return out
}

type patcher func(m provider.Manifest, cfg config.K8sResourcePatch) (*provider.Manifest, error)

func patchManifests(manifests []provider.Manifest, patches []config.K8sResourcePatch, patcher patcher) ([]provider.Manifest, error) {
	if len(patches) == 0 {
		return manifests, nil
	}

	out := make([]provider.Manifest, len(manifests))
	copy(out, manifests)

	for _, p := range patches {
		target := -1
		for i, m := range out {
			if m.Key().Kind() != p.Target.Kind {
				continue
			}
			if m.Key().Name() != p.Target.Name {
				continue
			}
			target = i
			break
		}
		if target < 0 {
			return nil, fmt.Errorf("no manifest matches the given patch: kind=%s, name=%s", p.Target.Kind, p.Target.Name)
		}
		patched, err := patcher(out[target], p)
		if err != nil {
			return nil, fmt.Errorf("failed to patch manifest: %s, error: %w", out[target].Key(), err)
		}
		out[target] = *patched
	}

	return out, nil
}

func patchManifest(m provider.Manifest, patch config.K8sResourcePatch) (*provider.Manifest, error) {
	if len(patch.Ops) == 0 {
		return &m, nil
	}

	fullBytes, err := m.YamlBytes()
	if err != nil {
		return nil, err
	}

	process := func(bytes []byte) ([]byte, error) {
		p, err := yamlprocessor.NewProcessor(bytes)
		if err != nil {
			return nil, err
		}

		for _, o := range patch.Ops {
			switch o.Op {
			case config.K8sResourcePatchOpYAMLReplace:
				if err := p.ReplaceString(o.Path, o.Value); err != nil {
					return nil, fmt.Errorf("failed to replace value at path: %s, error: %w", o.Path, err)
				}
			default:
				// TODO: Support more patch operation for K8sCanaryRolloutStageOptions.
				return nil, fmt.Errorf("%s operation is not supported currently", o.Op)
			}
		}

		return p.Bytes(), nil
	}

	buildManifest := func(bytes []byte) (*provider.Manifest, error) {
		manifests, err := provider.ParseManifests(string(bytes))
		if err != nil {
			return nil, err
		}
		if len(manifests) != 1 {
			return nil, fmt.Errorf("unexpected number of manifests, expected 1, got %d", len(manifests))
		}
		return &manifests[0], nil
	}

	// When the target is the whole manifest,
	// just pass full bytes to process and build a new manifest based on the returned data.
	root := patch.Target.DocumentRoot
	if root == "" {
		out, err := process(fullBytes)
		if err != nil {
			return nil, err
		}
		return buildManifest(out)
	}

	// When the target is a manifest field specified by documentRoot,
	// we have to extract that field value as a string.
	p, err := yamlprocessor.NewProcessor(fullBytes)
	if err != nil {
		return nil, err
	}

	v, err := p.GetValue(root)
	if err != nil {
		return nil, err
	}
	sv, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("the value for the specified root %s must be a string", root)
	}

	// And process that field data.
	out, err := process([]byte(sv))
	if err != nil {
		return nil, err
	}

	// Then rewrite the new data into the specified root.
	if err := p.ReplaceString(root, string(out)); err != nil {
		return nil, err
	}

	return buildManifest(p.Bytes())
}
