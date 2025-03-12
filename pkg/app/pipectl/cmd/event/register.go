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

package event

import (
	"context"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
)

var eventKeyFormatRegex = regexp.MustCompile(`^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$`)

type register struct {
	root *command

	name     string
	data     string
	labels   map[string]string
	contexts map[string]string

	// information of commit that triggers the event
	commitHash      string
	commitTitle     string
	commitMessage   string
	commitURL       string
	commitAuthor    string
	commitTimestamp int64
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
	cmd.Flags().StringToStringVar(&r.contexts, "contexts", r.contexts, "The list of the values for the event context. Format: key=value,key2=value2. The Key Format is [a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$")

	cmd.Flags().StringVar(&r.commitHash, "commit-hash", r.commitHash, "The commit hash that triggers the event.")
	cmd.Flags().StringVar(&r.commitTitle, "commit-title", r.commitTitle, "The title of commit that triggers the event.")
	cmd.Flags().StringVar(&r.commitMessage, "commit-message", r.commitMessage, "The message of commit that triggers the event.")
	cmd.Flags().StringVar(&r.commitURL, "commit-url", r.commitURL, "The URL of commit that triggers the event.")
	cmd.Flags().StringVar(&r.commitAuthor, "commit-author", r.commitAuthor, "The author of commit that triggers the event.")
	cmd.Flags().Int64Var(&r.commitTimestamp, "commit-timestamp", r.commitTimestamp, "The timestamp of commit that triggers the event.")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("data")

	return cmd
}

func (r *register) run(ctx context.Context, input cli.Input) error {
	if err := r.validateEventContexts(); err != nil {
		return fmt.Errorf("failed to validate event context: %w", err)
	}

	cli, err := r.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	req := &apiservice.RegisterEventRequest{
		Name:            r.name,
		Data:            r.data,
		Labels:          r.labels,
		Contexts:        r.contexts,
		CommitHash:      r.commitHash,
		CommitTitle:     r.commitTitle,
		CommitMessage:   r.commitMessage,
		CommitUrl:       r.commitURL,
		CommitAuthor:    r.commitAuthor,
		CommitTimestamp: r.commitTimestamp,
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

func (r *register) validateEventContexts() error {
	for key := range r.contexts {
		if !eventKeyFormatRegex.MatchString(key) {
			return fmt.Errorf("invalid format key '%s', should be ^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$", key)
		}
	}

	return nil
}
