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

package livestate

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/provider"
)

var (
	_ sdk.LivestatePlugin[sdk.ConfigNone, config.DeployTargetConfig, config.ApplicationConfigSpec] = (*Plugin)(nil)
)

type Plugin struct {
}

// GetLivestate implements sdk.LivestatePlugin.
func (p *Plugin) GetLivestate(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[config.DeployTargetConfig], input *sdk.GetLivestateInput[config.ApplicationConfigSpec]) (*sdk.GetLivestateResponse, error) {
	if len(dts) != 1 {
		return nil, fmt.Errorf("only 1 deploy target is allowed but got %d", len(dts))
	}
	dt := dts[0]

	cmd, err := provider.NewTerraformCommand(ctx, input.Client, input.Request.DeploymentSource, dt)
	if err != nil {
		input.Logger.Error("Failed to initialize Terraform command", zap.Error(err))
		return nil, err
	}

	buf := &bytes.Buffer{}
	planResult, err := cmd.Plan(ctx, buf)
	if err != nil {
		input.Logger.Error("Failed to execute plan", zap.Error(err))
		return nil, err
	}

	syncState, err := makeSyncState(planResult, input.Request.DeploymentSource.CommitHash)
	if err != nil {
		input.Logger.Error("Failed to make sync state", zap.Error(err))
		return nil, err
	}

	return &sdk.GetLivestateResponse{
		// Currently, LiveState is not supported in this plugin.
		SyncState: syncState,
	}, nil
}

func makeSyncState(r provider.PlanResult, commit string) (sdk.ApplicationSyncState, error) {
	if r.NoChanges() {
		return sdk.ApplicationSyncState{
			Status:      sdk.ApplicationSyncStateSynced,
			ShortReason: "",
			Reason:      "",
		}, nil
	}

	total := r.Imports + r.Adds + r.Destroys + r.Changes
	shortReason := fmt.Sprintf("There are %d manifests that are not synced (%d imports, %d adds, %d deletes, %d changes)", total, r.Imports, r.Adds, r.Destroys, r.Changes)
	if len(commit) >= 7 {
		commit = commit[:7]
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Diff between the defined state in Git at commit %s and actual live state:\n\n", commit))
	b.WriteString("--- Actual   (LiveState)\n+++ Expected (Git)\n\n")

	details, err := r.Render()
	if err != nil {
		return sdk.ApplicationSyncState{}, err
	}
	b.WriteString(details)

	return sdk.ApplicationSyncState{
		Status:      sdk.ApplicationSyncStateOutOfSync,
		ShortReason: shortReason,
		Reason:      b.String(),
	}, nil
}
