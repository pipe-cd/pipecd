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
)

// WaitApprovalStageOptions contains configurable values for a WAIT_APPROVAL stage.
type waitApprovalStageOptions struct {
	Approvers      []string `json:"approvers,omitempty"`
	MinApproverNum int      `json:"minApproverNum,omitempty"`
}

func (o waitApprovalStageOptions) validate() error {
	if len(o.Approvers) == 0 {
		return fmt.Errorf("approvers must be set")
	}
	if o.MinApproverNum < 1 {
		return fmt.Errorf("minApproverNum %d should be greater than 0", o.MinApproverNum)
	}
	if o.MinApproverNum > len(o.Approvers) {
		return fmt.Errorf("minApproverNum must be less than or equal to the number of approvers")
	}
	return nil
}

// decode decodes the raw JSON data and validates it.
func decode(data json.RawMessage) (waitApprovalStageOptions, error) {
	var opts waitApprovalStageOptions
	if err := json.Unmarshal(data, &opts); err != nil {
		return waitApprovalStageOptions{}, fmt.Errorf("failed to unmarshal the config: %w", err)
	}
	if err := opts.validate(); err != nil {
		return waitApprovalStageOptions{}, fmt.Errorf("failed to validate the config: %w", err)
	}
	return opts, nil
}
