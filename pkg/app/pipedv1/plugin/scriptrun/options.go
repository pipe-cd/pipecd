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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/creasty/defaults"
	unit "github.com/pipe-cd/piped-plugin-sdk-go/unit"
)

type scriptRunStageOptions struct {
	// user provided env variables to run the script with
	Env map[string]string `json:"env,omitempty"`
	// the command(s) to run
	Run string `json:"run,omitempty"`
	// timeout limit for this stage
	Timeout unit.Duration `json:"timeout,omitempty" default:"6h"`
	// the rollback command(s) to run if deployment fails
	OnRollback string `json:"onRollback,omitempty"`
}

func (o scriptRunStageOptions) validate() error {
	if o.Run == "" {
		return fmt.Errorf("SCRIPT_RUN stage requires run field")
	}
	if o.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}
	return nil
}

// decode decodes the raw JSON data and validates it.
func decode(data json.RawMessage) (scriptRunStageOptions, error) {
	opts := scriptRunStageOptions{}
	if err := json.Unmarshal(data, &opts); err != nil {
		return scriptRunStageOptions{}, fmt.Errorf("failed to unmarshal the stage config: %w", err)
	}
	if err := defaults.Set(&opts); err != nil {
		return scriptRunStageOptions{}, fmt.Errorf("failed to set default values for stage config: %w", err)
	}
	if err := opts.validate(); err != nil {
		return scriptRunStageOptions{}, fmt.Errorf("failed to validate the stage config: %w", err)
	}
	return opts, nil
}
