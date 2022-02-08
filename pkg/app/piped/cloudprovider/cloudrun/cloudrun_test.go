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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/run/v1"
	"sigs.k8s.io/yaml"
)

func TestLoadServiceManifest(t *testing.T) {
	const (
		appDir      = "testdata"
		serviceFile = "new_manifest.yaml"
	)
	// Success
	got, err := LoadServiceManifest(appDir, serviceFile)
	require.NoError(t, err)
	assert.NotEmpty(t, got)

	// Failure
	_, err = LoadServiceManifest(appDir, "")
	assert.Error(t, err)
}

func TestMakeManagedByPipedLabel(t *testing.T) {
	want := "pipecd-dev-managed-by=piped"
	got := MakeManagedByPipedLabel()
	assert.Equal(t, want, got)
}

func TestService(t *testing.T) {
	sm, err := ParseServiceManifest([]byte(serviceManifest))
	require.NoError(t, err)

	svc, err := manifestToRunService(sm)
	require.NoError(t, err)

	// ServiceManifest
	s := (*Service)(svc)
	got, err := s.ServiceManifest()
	require.NoError(t, err)
	assert.Equal(t, sm, got)

	// UID
	id := s.UID()
	assert.Equal(t, "service-uid", id)

	// RevisionNames
	names := s.RevisionNames()
	assert.Len(t, names, 1)
}

func TestRevision(t *testing.T) {
	rm, err := ParseRevisionManifest([]byte(revisionManifest))
	require.NoError(t, err)
	require.NotEmpty(t, rm)

	data, err := yaml.Marshal(rm.u)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var rev run.Revision
	err = yaml.Unmarshal(data, &rev)
	require.NoError(t, err)

	r := (*Revision)(&rev)
	got, err := r.RevisionManifest()
	require.NoError(t, err)
	assert.Equal(t, rm, got)
}
