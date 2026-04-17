// Copyright 2026 The PipeCD Authors.
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

package transfer

import (
	"github.com/pipe-cd/pipecd/pkg/model"
)

// BackupData holds the exported piped and application data from a source control plane.
type BackupData struct {
	// Version of the backup format.
	Version string `json:"version"`
	// CreatedAt is the RFC3339 timestamp when the backup was created.
	CreatedAt string `json:"created_at"`
	// Pipeds is the list of piped records discovered from the source control plane.
	Pipeds []*model.Piped `json:"pipeds"`
	// Applications is the full list of applications from the source control plane.
	Applications []*model.Application `json:"applications"`
}

// PipedMapping records the mapping of a piped from the source to the target control plane.
type PipedMapping struct {
	// OldPipedID is the piped ID in the source control plane.
	OldPipedID string `json:"old_piped_id"`
	// NewPipedID is the newly registered piped ID in the target control plane.
	NewPipedID string `json:"new_piped_id"`
	// NewKey is the newly issued piped API key. Update the piped config with this value.
	NewKey string `json:"new_key"`
	// PipedName is the human-readable name of the piped for reference.
	PipedName string `json:"piped_name"`
}

// RestoreResult summarises what was created on the target control plane.
type RestoreResult struct {
	// PipedMappings maps each source piped to its new registration on the target.
	PipedMappings []PipedMapping `json:"piped_mappings"`
	// RestoredApplicationCount is the number of applications successfully created.
	RestoredApplicationCount int `json:"restored_application_count"`
	// FailedApplicationCount is the number of applications that failed to be created.
	FailedApplicationCount int `json:"failed_application_count"`
}
