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

func (v ApplicationLiveStateVersion) IsBefore(a ApplicationLiveStateVersion) bool {
	if v.Timestamp < a.Timestamp {
		return true
	}
	if v.Timestamp > a.Timestamp {
		return false
	}
	return v.Index < a.Index
}

func (s KubernetesResourceState) HasDiff(a KubernetesResourceState) bool {
	if s.ApiVersion != a.ApiVersion {
		return true
	}
	if s.Namespace != a.Namespace {
		return true
	}
	if s.HealthStatus != a.HealthStatus {
		return true
	}
	if s.HealthDescription != a.HealthDescription {
		return true
	}
	if len(s.OwnerIds) != len(a.OwnerIds) {
		return true
	}
	if len(s.ParentIds) != len(a.ParentIds) {
		return true
	}

	for i := range s.OwnerIds {
		if s.OwnerIds[i] != a.OwnerIds[i] {
			return false
		}
	}

	for i := range s.ParentIds {
		if s.ParentIds[i] != a.ParentIds[i] {
			return false
		}
	}

	return false
}

// DetermineAppHealthStatus updates its own health status, which is determined based on its resources status.
// TODO: Determine health state of other than k8s and cloud run app
func (s *ApplicationLiveStateSnapshot) DetermineAppHealthStatus() {
	switch s.Kind {
	case ApplicationKind_KUBERNETES:
		s.determineKubernetesAppHealthStatus()
	case ApplicationKind_CLOUDRUN:
		s.determineCloudRunAppHealthStatus()
	}
	return
}

func (s *ApplicationLiveStateSnapshot) determineKubernetesAppHealthStatus() {
	app := s.Kubernetes
	if app == nil {
		return
	}
	status := ApplicationLiveStateSnapshot_HEALTHY
	for _, r := range app.Resources {
		if r.HealthStatus == KubernetesResourceState_OTHER {
			status = ApplicationLiveStateSnapshot_OTHER
			break
		}
	}
	s.HealthStatus = status
}

func (s *ApplicationLiveStateSnapshot) determineCloudRunAppHealthStatus() {
	app := s.Cloudrun
	if app == nil {
		return
	}
	for _, r := range app.Resources {
		if r.HealthStatus == CloudRunResourceState_OTHER {
			s.HealthStatus = ApplicationLiveStateSnapshot_OTHER
			return
		}

		if r.HealthStatus == CloudRunResourceState_UNKNOWN {
			s.HealthStatus = ApplicationLiveStateSnapshot_UNKNOWN
			return
		}
	}
	s.HealthStatus = ApplicationLiveStateSnapshot_HEALTHY
}
