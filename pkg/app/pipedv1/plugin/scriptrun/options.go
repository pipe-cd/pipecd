package main

import (
	"encoding/json"
	"fmt"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"time"
)

type scriptRunStageOptions struct {
	// user provided env variables to run the script with
	Env map[string]string `json:"env,omitempty"`
	// the command(s) to run
	Run string `json:"run,omitempty"`
	// timeout limit for this stage
	Timeout config.Duration `json:"timeout,omitempty" default:"6h"`
	// the rollback command(s) to run if deployment fails
	OnRollback string `json:"onRollback,omitempty"`
}

func (o scriptRunStageOptions) validate() error {
	if o.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}
	return nil
}

// decode decodes the raw JSON data and validates it.
func decode(data json.RawMessage) (scriptRunStageOptions, error) {
	opts := scriptRunStageOptions{
		Timeout: config.Duration(6 * time.Hour),
	}
	if err := json.Unmarshal(data, &opts); err != nil {
		return scriptRunStageOptions{}, fmt.Errorf("failed to unmarshal the config: %w", err)
	}
	if err := opts.validate(); err != nil {
		return scriptRunStageOptions{}, fmt.Errorf("failed to validate the config: %w", err)
	}
	return opts, nil
}
