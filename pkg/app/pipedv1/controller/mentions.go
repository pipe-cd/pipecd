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

package controller

import (
	"encoding/json"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/metadatastore"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// getApplicationNotificationMentions returns the list of users and groups who should be mentioned in the notification.
func getApplicationNotificationMentions(mds *metadatastore.MetadataStore, event model.NotificationEventType) ([]string, []string, error) {
	md, found := mds.SharedGet(model.MetadataKeyDeploymentNotification)
	if !found {
		return nil, nil, nil
	}

	var notification config.DeploymentNotification
	if err := json.Unmarshal([]byte(md), &notification); err != nil {
		return nil, nil, fmt.Errorf("could not extract mentions config: %w", err)
	}

	return notification.FindSlackUsers(event), notification.FindSlackGroups(event), nil
}
