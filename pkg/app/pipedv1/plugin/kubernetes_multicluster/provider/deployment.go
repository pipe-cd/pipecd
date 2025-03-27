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

package provider

import (
	"slices"
)

// FindContainerImages finds all container images that are referenced by the given manifest.
//
// It looks for container images in the following fields:
// - spec.template.spec.containers.image
//
// TODO: we should consider other fields like spec.template.spec.initContainers.image, spec.jobTempate.spec.template.spec.containers.image
func FindContainerImages(m Manifest) []string {
	var images []string

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "containers", "image"); len(n) > 0 {
		images = append(images, n...)
	}

	slices.Sort(images)
	return slices.Compact(images)
}

// FindReferencingConfigMaps finds all configmaps that are referenced by the given manifest.
//
// It looks for configmaps in the following fields:
// - spec.template.spec.volumes.configMap.name
// - spec.template.spec.initContainers.env.valueFrom.configMapKeyRef.name
// - spec.template.spec.initContainers.envFrom.configMapRef.name
// - spec.template.spec.containers.env.valueFrom.configMapKeyRef.name
// - spec.template.spec.containers.envFrom.configMapRef.name
func FindReferencingConfigMaps(m Manifest) []string {
	var configMaps []string

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "volumes", "configMap", "name"); len(n) > 0 {
		configMaps = append(configMaps, n...)
	}

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "initContainers", "env", "valueFrom", "configMapKeyRef", "name"); len(n) > 0 {
		configMaps = append(configMaps, n...)
	}

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "initContainers", "envFrom", "configMapRef", "name"); len(n) > 0 {
		configMaps = append(configMaps, n...)
	}

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "containers", "env", "valueFrom", "configMapKeyRef", "name"); len(n) > 0 {
		configMaps = append(configMaps, n...)
	}

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "containers", "envFrom", "configMapRef", "name"); len(n) > 0 {
		configMaps = append(configMaps, n...)
	}

	slices.Sort(configMaps)
	return slices.Compact(configMaps)
}

// FindReferencingSecrets finds all secrets that are referenced by the given manifest.
//
// It looks for secrets in the following fields:
// - spec.template.spec.volumes.secret.secretName
// - spec.template.spec.initContainers.env.valueFrom.secretKeyRef.name
// - spec.template.spec.initContainers.envFrom.secretRef.name
// - spec.template.spec.containers.env.valueFrom.secretKeyRef.name
// - spec.template.spec.containers.envFrom.secretRef.name
func FindReferencingSecrets(m Manifest) []string {
	var secrets []string

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "volumes", "secret", "secretName"); len(n) > 0 {
		secrets = append(secrets, n...)
	}

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "initContainers", "env", "valueFrom", "secretKeyRef", "name"); len(n) > 0 {
		secrets = append(secrets, n...)
	}

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "initContainers", "envFrom", "secretRef", "name"); len(n) > 0 {
		secrets = append(secrets, n...)
	}

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "containers", "env", "valueFrom", "secretKeyRef", "name"); len(n) > 0 {
		secrets = append(secrets, n...)
	}

	if n := nestedStringSlice(m.body.Object, "spec", "template", "spec", "containers", "envFrom", "secretRef", "name"); len(n) > 0 {
		secrets = append(secrets, n...)
	}

	slices.Sort(secrets)
	return slices.Compact(secrets)
}

// nestedStringSlice extracts a string slice from the given object by following the fields.
// It returns the extracted string slice.
// If there is []map[string]any in the middle of the fields, it will be flattened.
func nestedStringSlice(obj any, fields ...string) []string {
	// No field to extract, return the original object.
	if len(fields) == 0 {
		switch obj := obj.(type) {
		case []string:
			return obj
		case []any:
			var result []string
			for _, item := range obj {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		case string:
			return []string{obj}
		default:
			return nil
		}
	}

	switch v := obj.(type) {
	case map[string]any:
		return nestedStringSlice(v[fields[0]], fields[1:]...)
	case []any:
		var result []string
		for _, item := range v {
			result = append(result, nestedStringSlice(item, fields...)...)
		}
		return result
	default:
		return nil
	}
}
