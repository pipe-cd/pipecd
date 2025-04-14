// Copyright 2025 The PipeCD Authors.
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

package sdk

import (
	"slices"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

// Spec[T] represents both of follows
// - the type is pointer type of T
// - the type has Validate method
type Spec[T any] = config.Spec[T]

// LoadConfigSpec loads the config spec from the application config.
// The type T must be a pointer to a struct that has Validate method.
func LoadConfigSpec[T Spec[RT], RT any](ds DeploymentSource) (T, error) {
	cfg, err := config.DecodeYAML[T](ds.ApplicationConfig)
	if err != nil {
		return nil, err
	}
	return cfg.Spec, nil
}

// LoadStages loads the pipeline stages from the application config.
func LoadStages(ds DeploymentSource) (Stages, error) {
	cfg, err := config.DecodeYAML[*config.GenericApplicationSpec](ds.ApplicationConfig)
	if err != nil {
		return Stages{}, err
	}

	// If the pipeline is not defined, return an empty list.
	if cfg == nil || cfg.Spec == nil || cfg.Spec.Pipeline == nil || len(cfg.Spec.Pipeline.Stages) == 0 {
		return Stages{}, nil
	}

	// Convert the pipeline stages to a list of stage names.
	stages := make([]string, 0, len(cfg.Spec.Pipeline.Stages))
	for _, s := range cfg.Spec.Pipeline.Stages {
		stages = append(stages, string(s.Name))
	}
	return Stages{stages: stages}, nil
}

// Stages is a list of pipeline stages.
// It's intended to use for checking if a stage is included in the pipeline.
type Stages struct {
	stages []string
}

// Has checks if the given stage is included in the list.
func (s Stages) Has(stage string) bool {
	return slices.Contains(s.stages, stage)
}
