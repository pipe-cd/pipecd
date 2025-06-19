package main

import (
	"encoding/json"
	"fmt"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

type scriptRunStageOptions struct {
	Env        map[string]string `json:"env"`
	Run        string            `json:"run"`
	Timeout    config.Duration   `json:"timeout" default:"6h"`
	OnRollback string            `json:"onRollback"`
}

func (o scriptRunStageOptions) validate() error {
	if o.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}
	return nil
}

// decode decodes the raw JSON data and validates it.
func decode(data json.RawMessage) (scriptRunStageOptions, error) {
	var opts scriptRunStageOptions
	if err := json.Unmarshal(data, &opts); err != nil {
		return scriptRunStageOptions{}, fmt.Errorf("failed to unmarshal the config: %w", err)
	}
	if err := opts.validate(); err != nil {
		return scriptRunStageOptions{}, fmt.Errorf("failed to validate the config: %w", err)
	}
	return opts, nil
}
