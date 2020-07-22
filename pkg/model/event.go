// Copyright 2020 The PipeCD Authors.
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

type Event struct {
	Type     EventType
	Metadata interface{}
}

func (e Event) Group() EventGroup {
	switch {
	case e.Type < 100:
		return EventGroup_EVENT_DEPLOYMENT
	case e.Type < 200:
		return EventGroup_EVENT_APPLICATION_SYNC
	case e.Type < 300:
		return EventGroup_EVENT_APPLICATION_HEALTH
	case e.Type < 400:
		return EventGroup_EVENT_PIPED
	default:
		return EventGroup_EVENT_NONE
	}
}

func (e *EventDeploymentTriggered) AppName() string {
	return e.Deployment.ApplicationName
}

func (e *EventDeploymentPlanned) AppName() string {
	return e.Deployment.ApplicationName
}

func (e *EventDeploymentApproved) AppName() string {
	return e.Deployment.ApplicationName
}

func (e *EventDeploymentRollingBack) AppName() string {
	return e.Deployment.ApplicationName
}

func (e *EventDeploymentSucceeded) AppName() string {
	return e.Deployment.ApplicationName
}

func (e *EventDeploymentFailed) AppName() string {
	return e.Deployment.ApplicationName
}

func (e *EventApplicationSynced) AppName() string {
	return e.Application.Id
}

func (e *EventApplicationOutOfSync) AppName() string {
	return e.Application.Id
}

func (e *EventDeploymentTriggered) EnvID() string {
	return e.Deployment.EnvId
}

func (e *EventDeploymentPlanned) EnvID() string {
	return e.Deployment.EnvId
}

func (e *EventDeploymentApproved) EnvID() string {
	return e.Deployment.EnvId
}

func (e *EventDeploymentRollingBack) EnvID() string {
	return e.Deployment.EnvId
}

func (e *EventDeploymentSucceeded) EnvID() string {
	return e.Deployment.EnvId
}

func (e *EventDeploymentFailed) EnvID() string {
	return e.Deployment.EnvId
}

func (e *EventApplicationSynced) EnvID() string {
	return e.Application.EnvId
}

func (e *EventApplicationOutOfSync) EnvID() string {
	return e.Application.EnvId
}
