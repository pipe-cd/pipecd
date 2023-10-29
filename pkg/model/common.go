// Copyright 2023 The PipeCD Authors.
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

package model

// ApplicationKindStrings returns a list of available deployment kinds in string.
func ApplicationKindStrings() []string {
	out := make([]string, 0, len(ApplicationKind_value))
	for k := range ApplicationKind_value {
		out = append(out, k)
	}
	return out
}

func (ak ApplicationKind) ToRollbackKind() RollbackKind {
	switch ak {
	case ApplicationKind_KUBERNETES:
		return RollbackKind_Rollback_KUBERNETES
	case ApplicationKind_TERRAFORM:
		return RollbackKind_Rollback_TERRAFORM
	case ApplicationKind_LAMBDA:
		return RollbackKind_Rollback_LAMBDA
	case ApplicationKind_CLOUDRUN:
		return RollbackKind_Rollback_CLOUDRUN
	case ApplicationKind_ECS:
		return RollbackKind_Rollback_ECS
	default:
		return RollbackKind_Rollback_KUBERNETES
	}
}

// ContainLabels checks if it has all the given labels.
func (a *ApplicationInfo) ContainLabels(labels map[string]string) bool {
	if len(a.Labels) < len(labels) {
		return false
	}

	for k, v := range labels {
		value, ok := a.Labels[k]
		if !ok {
			return false
		}
		if value != v {
			return false
		}
	}
	return true
}
