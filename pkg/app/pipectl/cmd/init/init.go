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

package init

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/cli"
)

type command struct {
	someTextOption string
}

func NewCommand() *cobra.Command {
	c := &command{
		someTextOption: "default-value",
	}
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Create a app.pipecd.yaml easily (interactively)",
		Example: `  pipectl init`,
		Long:    "Create a app.pipecd.yaml easily, interactively selecting options.",
		RunE:    cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.someTextOption, "some-text-option", c.someTextOption, "Some text option")

	return cmd
}

func (c *command) run(ctx context.Context, input cli.Input) error {
	input.Logger.Warn("not implemented yet")
	return nil
}
