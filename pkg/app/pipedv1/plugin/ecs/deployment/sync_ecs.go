// Copyright 2026 The PipeCD Authors.
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

package deployment

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

// syncECSController rolls out the task definition to a service managed by the ECS deployment controller.
// The rollout itself is performed by ECS: updating the service to the new task definition starts a
// service deployment, which this function waits for.
func syncECSController(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	taskDef types.TaskDefinition,
	serviceDef types.Service,
) error {
	// The commit hash tag must reflect what is actually deployed,
	// so it is only written after the service deployment succeeds.
	commitTag, otherTags := splitCommitHashTag(serviceDef.Tags)
	serviceDef.Tags = otherTags

	lp.Info("Start applying the ECS task definition")
	td, err := applyTaskDefinition(ctx, client, taskDef)
	if err != nil {
		return fmt.Errorf("failed to apply task definition: %w", err)
	}
	serviceDef.TaskDefinition = td.TaskDefinitionArn

	prev, err := client.GetLatestServiceDeployment(ctx, *serviceDef.ClusterArn, *serviceDef.ServiceName)
	if err != nil {
		return fmt.Errorf("failed to get the latest service deployment: %w", err)
	}
	var prevArn *string
	if prev != nil {
		prevArn = prev.ServiceDeploymentArn
	}

	lp.Info("Start applying the ECS service definition")
	service, err := applyServiceDefinition(ctx, lp, client, serviceDef)
	if err != nil {
		return fmt.Errorf("failed to apply service definition: %w", err)
	}

	lp.Infof("Wait for the service deployment of service %s to complete", *service.ServiceName)
	if err := client.WaitServiceDeploymentComplete(ctx, *service.ClusterArn, *service.ServiceName, prevArn); err != nil {
		return fmt.Errorf("service deployment did not complete: %w", err)
	}

	if commitTag != nil {
		if err := client.TagResource(ctx, *service.ServiceArn, []types.Tag{*commitTag}); err != nil {
			return fmt.Errorf("failed to tag service %s with the deployed commit: %w", *service.ServiceName, err)
		}
	}
	return nil
}

func splitCommitHashTag(tags []types.Tag) (*types.Tag, []types.Tag) {
	var commitTag *types.Tag
	others := make([]types.Tag, 0, len(tags))
	for _, t := range tags {
		if t.Key != nil && *t.Key == provider.LabelCommitHash {
			t := t
			commitTag = &t
			continue
		}
		others = append(others, t)
	}
	return commitTag, others
}
