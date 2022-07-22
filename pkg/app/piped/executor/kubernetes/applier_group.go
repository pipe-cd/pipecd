// Copyright 2022 The PipeCD Authors.
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
	"fmt"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applierGetter interface {
	Get(k provider.ResourceKey) (provider.Applier, error)
}

type applierGroup struct {
	resourceRoutes []config.KubernetesResourceRoute
	appliers       map[string]provider.Applier
	defaultApplier provider.Applier
}

func newApplierGroup(defaultProvider string, appCfg config.KubernetesApplicationSpec, pipedCfg *config.PipedSpec, logger *zap.Logger) (*applierGroup, error) {
	cp, ok := pipedCfg.FindCloudProvider(defaultProvider, model.ApplicationKind_KUBERNETES)
	if !ok {
		return nil, fmt.Errorf("provider %s was not found", defaultProvider)
	}

	defaultApplier := provider.NewApplier(
		appCfg.Input,
		*cp.KubernetesConfig,
		logger,
	)
	d := &applierGroup{
		resourceRoutes: appCfg.ResourceRoutes,
		appliers:       map[string]provider.Applier{defaultProvider: defaultApplier},
		defaultApplier: defaultApplier,
	}

	for _, r := range appCfg.ResourceRoutes {
		if _, ok := d.appliers[r.Provider]; ok {
			continue
		}

		cp, ok := pipedCfg.FindCloudProvider(r.Provider, model.ApplicationKind_KUBERNETES)
		if !ok {
			return nil, fmt.Errorf("provider %s specified in resourceRoutes was not found", r.Provider)
		}

		d.appliers[r.Provider] = provider.NewApplier(appCfg.Input, *cp.KubernetesConfig, logger)
	}

	return d, nil
}

// TODO: Add test for this applierGroup function.
func (d applierGroup) Get(rk provider.ResourceKey) (provider.Applier, error) {
	for _, r := range d.resourceRoutes {
		if r.Match == nil {
			if a, ok := d.appliers[r.Provider]; ok {
				return a, nil
			}
			return nil, fmt.Errorf("provider %s specified in resourceRoutes was not found", r.Provider)
		}
		if k := r.Match.Kind; k != "" && k != rk.Kind {
			continue
		}
		if n := r.Match.Name; n != "" && n != rk.Name {
			continue
		}
		if a, ok := d.appliers[r.Provider]; ok {
			return a, nil
		}
		return nil, fmt.Errorf("provider %s specified in resourceRoutes was not found", r.Provider)
	}

	return d.defaultApplier, nil
}
