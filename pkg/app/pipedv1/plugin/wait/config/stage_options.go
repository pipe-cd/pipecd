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
	"encoding/json"

	"github.com/creasty/defaults"
	"sigs.k8s.io/yaml"
)

type StageOptions interface {
	Validate() error
}

// DecodeYAML unmarshals stageOptions YAML data to specified StageOptions type and validates the result.
func DecodeStageOptionsYAML[T StageOptions](data []byte) (*T, error) {
	js, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	var o T
	if err := json.Unmarshal(js, &o); err != nil {
		return nil, err
	}
	if err := defaults.Set(&o); err != nil {
		return nil, err
	}
	if err := o.Validate(); err != nil {
		return nil, err
	}
	return &o, nil
}
