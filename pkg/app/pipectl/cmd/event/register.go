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

package event

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
)

type register struct {
	root *command

	name   string
	data   string
	labels map[string]string
}

func newRegisterCommand(root *command) *cobra.Command {
	r := &register{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register an event.",
		RunE:  cli.WithContext(r.run),
	}

	cmd.Flags().StringVar(&r.name, "name", r.name, "The name of event.")
	cmd.Flags().StringVar(&r.data, "data", r.data, "The string value of event data.")
	cmd.Flags().StringToStringVar(&r.labels, "labels", r.labels, "The list of labels for event. Format: key=value,key2=value2")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("data")

	return cmd
}

func (r *register) run(ctx context.Context, input cli.Input) error {
	cli, err := r.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	req := &apiservice.RegisterEventRequest{
		Name:   r.name,
		Data:   r.data,
		Labels: r.labels,
	}

	res, err := cli.RegisterEvent(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to register event: %w", err)
	}

	input.Logger.Info("Successfully registered event",
		zap.String("id", res.EventId),
	)
	return nil
}
