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

package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
)

type LauncherConfig struct {
	Kind       Kind         `json:"kind"`
	APIVersion string       `json:"apiVersion,omitempty"`
	Spec       LauncherSpec `json:"spec"`
}

func (c *LauncherConfig) Validate() error {
	if c.Kind != KindPiped {
		return fmt.Errorf("wrong configuration kind for piped: %v", c.Kind)
	}
	if c.Spec.ProjectID == "" {
		return errors.New("projectID must be set")
	}
	if c.Spec.PipedID == "" {
		return errors.New("pipedID must be set")
	}
	if c.Spec.PipedKeyData == "" && c.Spec.PipedKeyFile == "" {
		return errors.New("either pipedKeyFile or pipedKeyData must be set")
	}
	if c.Spec.PipedKeyData != "" && c.Spec.PipedKeyFile != "" {
		return errors.New("only pipedKeyFile or pipedKeyData can be set")
	}
	if c.Spec.APIAddress == "" {
		return errors.New("apiAddress must be set")
	}
	return nil
}

type LauncherSpec struct {
	// The identifier of the PipeCD project where this piped belongs to.
	ProjectID string
	// The unique identifier generated for this piped.
	PipedID string
	// The path to the file containing the generated Key string for this piped.
	PipedKeyFile string
	// Base64 encoded string of Piped key.
	PipedKeyData string
	// The address used to connect to the control-plane's API.
	APIAddress string `json:"apiAddress"`
}

func (s *LauncherSpec) LoadPipedKey() ([]byte, error) {
	if s.PipedKeyData != "" {
		return base64.StdEncoding.DecodeString(s.PipedKeyData)
	}
	if s.PipedKeyFile != "" {
		return os.ReadFile(s.PipedKeyFile)
	}
	return nil, errors.New("either pipedKeyFile or pipedKeyData must be set")
}
