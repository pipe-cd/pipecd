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

package cloudrun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeCloudRunParent(t *testing.T) {
	t.Parallel()

	const projectID = "projectID"
	got := makeCloudRunParent(projectID)
	want := "namespaces/projectID"
	assert.Equal(t, want, got)
}

func TestMakeCloudRunServiceName(t *testing.T) {
	t.Parallel()

	const (
		projectID = "projectID"
		serviceID = "serviceID"
	)
	got := makeCloudRunServiceName(projectID, serviceID)
	want := "namespaces/projectID/services/serviceID"
	assert.Equal(t, want, got)
}

func TestMakeCloudRunRevisionName(t *testing.T) {
	t.Parallel()

	const (
		projectID  = "projectID"
		revisionID = "revisionID"
	)
	got := makeCloudRunRevisionName(projectID, revisionID)
	want := "namespaces/projectID/revisions/revisionID"
	assert.Equal(t, want, got)
}
