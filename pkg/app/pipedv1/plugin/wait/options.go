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

package main

import (
	"encoding/json"
	"fmt"

	"github.com/creasty/defaults"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

// WaitStageOptions contains configurable values for a WAIT stage.
type WaitStageOptions struct {
	Duration config.Duration `json:"duration,omitempty" default:"1m"`
	// TODO: Handle SkipOn options.
	// SkipOn   config.SkipOptions `json:"skipOn,omitempty"`
}

func (o WaitStageOptions) validate() error {
	if o.Duration <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}
	return nil
}

// decode decodes the raw JSON data and validates it.
func decode(data json.RawMessage) (WaitStageOptions, error) {
	var opts WaitStageOptions
	if err := json.Unmarshal(data, &opts); err != nil {
		return WaitStageOptions{}, fmt.Errorf("failed to unmarshal the config: %w", err)
	}
	if err := defaults.Set(&opts); err != nil {
		return WaitStageOptions{}, fmt.Errorf("failed to set default values: %w", err)
	}
	if err := opts.validate(); err != nil {
		return WaitStageOptions{}, fmt.Errorf("failed to validate the config: %w", err)
	}
	return opts, nil
}
