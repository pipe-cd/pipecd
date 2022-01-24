// Copyright 2022 The PipeCD Authors.
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

	"github.com/stretchr/testify/require"
)

func TestMakeCloudRunParent(t *testing.T) {
	const projectID = "projectID"
	got := makeCloudRunParent(projectID)
	want := "namespaces/projectID"
	require.Equal(t, want, got)
}

func TestMakeCloudRunServiceName(t *testing.T) {
	const (
		projectID = "projectID"
		serviceID = "serviceID"
	)
	got := makeCloudRunServiceName(projectID, serviceID)
	want := "namespaces/projectID/services/serviceID"
	require.Equal(t, want, got)
}

func TestMakeCloudRunRevisionName(t *testing.T) {
	const (
		projectID  = "projectID"
		revisionID = "revisionID"
	)
	got := makeCloudRunRevisionName(projectID, revisionID)
	want := "namespaces/projectID/revisions/revisionID"
	require.Equal(t, want, got)
}

func TestManifestToRunService(t *testing.T) {
	sm, err := ParseServiceManifest([]byte(manifest))
	require.NoError(t, err)
	require.NotEmpty(t, sm)

	got, err := manifestToRunService(sm)
	require.NoError(t, err)
	require.NotEmpty(t, got)
}

func TestService(t *testing.T) {
	sm, err := ParseServiceManifest([]byte(manifest))
	require.NoError(t, err)

	svc, err := manifestToRunService(sm)
	require.NoError(t, err)

	s := (*Service)(svc)
	got, err := s.ServiceManifest()
	require.NoError(t, err)
	require.Equal(t, sm, got)
}
