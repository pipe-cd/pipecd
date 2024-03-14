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
	"fmt"
	"sort"
	"strings"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applierGetter interface {
	Get(k provider.ResourceKey) (provider.Applier, error)
}

type applierGroup struct {
	resourceRoutes   []config.KubernetesResourceRoute
	appliers         map[string]provider.Applier
	labeledProviders map[string][]string
	defaultApplier   provider.Applier
}

func newApplierGroup(defaultProvider string, appCfg config.KubernetesApplicationSpec, pipedCfg *config.PipedSpec, logger *zap.Logger) (*applierGroup, error) {
	cp, ok := pipedCfg.FindPlatformProvider(defaultProvider, model.ApplicationKind_KUBERNETES)
	if !ok {
		return nil, fmt.Errorf("provider %s was not found", defaultProvider)
	}

	defaultApplier := provider.NewApplier(
		appCfg.Input,
		*cp.KubernetesConfig,
		logger,
	)
	d := &applierGroup{
		resourceRoutes:   appCfg.ResourceRoutes,
		appliers:         map[string]provider.Applier{defaultProvider: defaultApplier},
		labeledProviders: make(map[string][]string, 0),
		defaultApplier:   defaultApplier,
	}

	for _, r := range appCfg.ResourceRoutes {
		if name := r.Provider.Name; name != "" {
			if _, ok := d.appliers[name]; ok {
				continue
			}
			cp, found := pipedCfg.FindPlatformProvider(name, model.ApplicationKind_KUBERNETES)
			if !found {
				return nil, fmt.Errorf("provider %s specified in resourceRoutes was not found", name)
			}
			d.appliers[name] = provider.NewApplier(appCfg.Input, *cp.KubernetesConfig, logger)
			continue
		}
		if labels := r.Provider.Labels; len(labels) > 0 {
			cps := pipedCfg.FindPlatformProvidersByLabels(labels, model.ApplicationKind_KUBERNETES)
			if len(cps) == 0 {
				return nil, fmt.Errorf("there is no provider that matches the specified labels (%v)", labels)
			}
			names := make([]string, 0, len(cps))
			for _, cp := range cps {
				if _, ok := d.appliers[cp.Name]; !ok {
					d.appliers[cp.Name] = provider.NewApplier(appCfg.Input, *cp.KubernetesConfig, logger)
				}
				names = append(names, cp.Name)
			}
			// Save names of the labeled providers for search later.
			key := makeKeyFromProviderLabels(labels)
			d.labeledProviders[key] = names
		}
	}

	return d, nil
}

func makeKeyFromProviderLabels(labels map[string]string) string {
	labelList := make([]string, 0, len(labels))
	for k, v := range labels {
		if v != "" {
			labelList = append(labelList, fmt.Sprintf("%s:%s", k, v))
		}
	}
	sort.Strings(labelList)
	return strings.Join(labelList, ",")
}

// TODO: Add test for this applierGroup function.
func (d applierGroup) Get(rk provider.ResourceKey) (provider.Applier, error) {
	resourceMatch := func(matcher *config.KubernetesResourceRouteMatcher) bool {
		// Match any resource when the matcher was not specified.
		if matcher == nil {
			return true
		}
		if matcher.Kind != "" && matcher.Kind != rk.Kind {
			return false
		}
		if matcher.Name != "" && matcher.Name != rk.Name {
			return false
		}
		return true
	}

	for _, r := range d.resourceRoutes {
		if !resourceMatch(r.Match) {
			continue
		}
		if name := r.Provider.Name; name != "" {
			if a, ok := d.appliers[name]; ok {
				return a, nil
			}
			return nil, fmt.Errorf("provider %s specified in resourceRoutes was not found", name)
		}
		if labels := r.Provider.Labels; len(labels) > 0 {
			key := makeKeyFromProviderLabels(labels)
			cps := d.labeledProviders[key]
			if len(cps) == 0 {
				return nil, fmt.Errorf("there are no provider that matches the specified labels (%v)", labels)
			}
			as := make([]provider.Applier, 0, len(cps))
			for _, cp := range cps {
				if a, ok := d.appliers[cp]; ok {
					as = append(as, a)
					continue
				}
				return nil, fmt.Errorf("provider %s specified in resourceRoutes was not found", cp)
			}
			applier := provider.NewMultiApplier(as...)
			return applier, nil
		}
	}

	return d.defaultApplier, nil
}
