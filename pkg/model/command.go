// Copyright 2023 The PipeCD Authors.
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

package model

import "context"

const (
	MetadataKeyTriggeredDeploymentID = "TriggeredDeploymentID"
)

type ReportableCommand struct {
	*Command
	Report func(ctx context.Context, status CommandStatus, metadata map[string]string, output []byte) error
}

func (c *Command) IsHandled() bool {
	return c.Status != CommandStatus_COMMAND_NOT_HANDLED_YET
}

func (c *Command) IsSyncApplicationCmd() bool {
	return c.GetSyncApplication() != nil
}

func (c *Command) IsChainSyncApplicationCmd() bool {
	return c.GetChainSyncApplication() != nil
}

func (c *Command) IsRestartPipedCmd() bool {
	return c.GetRestartPiped() != nil
}

func (c *Command) SetUpdatedAt(t int64) {
	c.UpdatedAt = t
}
