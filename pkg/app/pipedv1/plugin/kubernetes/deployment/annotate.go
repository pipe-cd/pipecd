// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deployment

import (
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

// annotateConfigHash appends a hash annotation into the workload manifests.
// The hash value is calculated by hashing the content of all configmaps/secrets
// that are referenced by the workload.
// This appending ensures that the workload should be restarted when
// one of its configurations changed.
func annotateConfigHash(manifests []provider.Manifest) error {
	if len(manifests) == 0 {
		return nil
	}

	configMaps := make(map[string]provider.Manifest)
	secrets := make(map[string]provider.Manifest)
	for _, m := range manifests {
		if m.Key.IsConfigMap() {
			configMaps[m.Key.Name] = m
			continue
		}
		if m.Key.IsSecret() {
			secrets[m.Key.Name] = m
		}
	}

	// This application is not containing any config manifests
	// so nothing to do.
	if len(configMaps)+len(secrets) == 0 {
		return nil
	}

	for _, m := range manifests {
		if m.Key.IsDeployment() {
			if err := annotateConfigHashToWorkload(m, configMaps, secrets); err != nil {
				return err
			}

			// TODO: Add support for other workload types, such as StatefulSet, DaemonSet, etc.
		}
	}

	return nil
}

func annotateConfigHashToWorkload(m provider.Manifest, managedConfigMaps, managedSecrets map[string]provider.Manifest) error {
	configMaps := provider.FindReferencingConfigMaps(m.Body)
	secrets := provider.FindReferencingSecrets(m.Body)

	// The deployment is not referencing any config resources.
	if len(configMaps)+len(secrets) == 0 {
		return nil
	}

	cfgs := make([]provider.Manifest, 0, len(configMaps)+len(secrets))
	for _, cm := range configMaps {
		m, ok := managedConfigMaps[cm]
		if !ok {
			// We do not return error here because the deployment may use
			// a config resource that is not managed by PipeCD.
			continue
		}
		cfgs = append(cfgs, m)
	}
	for _, s := range secrets {
		m, ok := managedSecrets[s]
		if !ok {
			// We do not return error here because the deployment may use
			// a config resource that is not managed by PipeCD.
			continue
		}
		cfgs = append(cfgs, m)
	}

	if len(cfgs) == 0 {
		return nil
	}

	hash, err := provider.HashManifests(cfgs)
	if err != nil {
		return err
	}

	m.AddStringMapValues(
		map[string]string{
			provider.AnnotationConfigHash: hash,
		},
		"spec",
		"template",
		"metadata",
		"annotations",
	)
	return nil
}
