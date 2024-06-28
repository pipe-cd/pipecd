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

package resource

type PodTemplateSpec struct {
	Spec PodSpec
}

type PodSpec struct {
	InitContainers []Container
	Containers     []Container
	Volumes        []Volume
}

type Container struct {
	Name         string
	Image        string
	VolumeMounts []VolumeMount
}

type Volume struct {
	Name         string
	VolumeSource `json:",inline"`
}

type VolumeSource struct {
	Secret    *SecretVolumeSource
	ConfigMap *ConfigMapVolumeSource
}

type SecretVolumeSource struct {
	SecretName string
}

type LocalObjectReference struct {
	Name string
}

type ConfigMapVolumeSource struct {
	LocalObjectReference `json:",inline"`
}

type VolumeMount struct {
	Name string
}
