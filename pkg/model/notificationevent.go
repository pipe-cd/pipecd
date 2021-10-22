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

type NotificationEvent struct {
	Type     NotificationEventType
	Metadata interface{}
}

func (e NotificationEvent) Group() NotificationEventGroup {
	switch {
	case e.Type < 100:
		return NotificationEventGroup_EVENT_DEPLOYMENT
	case e.Type < 200:
		return NotificationEventGroup_EVENT_APPLICATION_SYNC
	case e.Type < 300:
		return NotificationEventGroup_EVENT_APPLICATION_HEALTH
	case e.Type < 400:
		return NotificationEventGroup_EVENT_PIPED
	default:
		return NotificationEventGroup_EVENT_NONE
	}
}

func (e *NotificationEventDeploymentTriggered) GetAppName() string {
	return e.Deployment.ApplicationName
}

func (e *NotificationEventDeploymentPlanned) GetAppName() string {
	return e.Deployment.ApplicationName
}

func (e *NotificationEventDeploymentApproved) GetAppName() string {
	return e.Deployment.ApplicationName
}

func (e *NotificationEventDeploymentRollingBack) GetAppName() string {
	return e.Deployment.ApplicationName
}

func (e *NotificationEventDeploymentSucceeded) GetAppName() string {
	return e.Deployment.ApplicationName
}

func (e *NotificationEventDeploymentFailed) GetAppName() string {
	return e.Deployment.ApplicationName
}

func (e *NotificationEventApplicationSynced) GetAppName() string {
	return e.Application.Name
}

func (e *NotificationEventApplicationOutOfSync) GetAppName() string {
	return e.Application.Name
}
