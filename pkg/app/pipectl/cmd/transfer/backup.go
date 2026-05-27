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

package transfer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type backup struct {
	root *command

	outputFile string
}

func newBackupCommand(root *command) *cobra.Command {
	c := &backup{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup piped and application data from the source control plane to a local file.",
		Long: `Backup exports all pipeds (discovered via their applications) and all applications
from the source control plane into a single JSON file. Use the parent --address and --api-key
flags to point at the source control plane.

Note: deployment history is not included because the API does not expose a write endpoint for deployments.`,
		RunE: cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.outputFile, "output-file", c.outputFile, "The path of the output JSON file to save the backup data.")
	cmd.MarkFlagRequired("output-file")

	return cmd
}

func (c *backup) run(ctx context.Context, input cli.Input) error {
	input.Logger.Info("Starting control plane backup...")

	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	// Collect all applications via paginated ListApplications calls.
	applications, err := listAllApplications(ctx, cli, input.Logger)
	if err != nil {
		return fmt.Errorf("failed to list applications: %w", err)
	}
	input.Logger.Info(fmt.Sprintf("Found %d application(s)", len(applications)))

	// Discover unique piped IDs from the application list, then fetch piped details.
	pipeds, err := fetchPipeds(ctx, cli, applications, input.Logger)
	if err != nil {
		return fmt.Errorf("failed to fetch pipeds: %w", err)
	}
	input.Logger.Info(fmt.Sprintf("Found %d piped(s)", len(pipeds)))

	data := &BackupData{
		Version:      "1",
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		Pipeds:       pipeds,
		Applications: applications,
	}

	if err := writeBackupFile(c.outputFile, data); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	input.Logger.Info("Backup completed successfully", zap.String("output-file", c.outputFile))
	return nil
}

// listAllApplications fetches all applications (both enabled and disabled) from the control plane
// using cursor-based pagination. The ListApplications API filters by the disabled field as a strict
// equality match, so two separate paginated sweeps are required.
func listAllApplications(ctx context.Context, cli apiservice.Client, logger *zap.Logger) ([]*model.Application, error) {
	enabled, err := listApplicationsByDisabledStatus(ctx, cli, false, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to list enabled applications: %w", err)
	}
	disabled, err := listApplicationsByDisabledStatus(ctx, cli, true, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to list disabled applications: %w", err)
	}
	all := append(enabled, disabled...)
	return all, nil
}

// listApplicationsByDisabledStatus paginates through ListApplications for a given disabled status.
func listApplicationsByDisabledStatus(ctx context.Context, cli apiservice.Client, disabled bool, logger *zap.Logger) ([]*model.Application, error) {
	var (
		all    []*model.Application
		cursor string
	)
	for {
		resp, err := cli.ListApplications(ctx, &apiservice.ListApplicationsRequest{
			Disabled: disabled,
			Cursor:   cursor,
			Limit:    500,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Applications...)
		if resp.Cursor == "" || len(resp.Applications) == 0 {
			break
		}
		cursor = resp.Cursor
	}
	label := "enabled"
	if disabled {
		label = "disabled"
	}
	logger.Info(fmt.Sprintf("Fetched %d %s application(s)", len(all), label))
	return all, nil
}

// fetchPipeds collects the unique piped IDs from the applications and fetches each piped's details.
func fetchPipeds(ctx context.Context, cli apiservice.Client, applications []*model.Application, logger *zap.Logger) ([]*model.Piped, error) {
	seen := make(map[string]struct{})
	pipeds := make([]*model.Piped, 0, 10)

	for _, app := range applications {
		if _, ok := seen[app.PipedId]; ok {
			continue
		}
		seen[app.PipedId] = struct{}{}

		resp, err := cli.GetPiped(ctx, &apiservice.GetPipedRequest{PipedId: app.PipedId})
		if err != nil {
			logger.Warn("failed to fetch piped, skipping", zap.String("piped-id", app.PipedId), zap.Error(err))
			continue
		}
		pipeds = append(pipeds, resp.Piped)
	}
	if len(pipeds) == 0 {
		return nil, fmt.Errorf("no piped for backup")
	}
	return pipeds, nil
}

// writeBackupFile serialises data to JSON and writes it to the given path.
func writeBackupFile(path string, data *BackupData) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %w", err)
	}
	return os.WriteFile(path, b, 0o600)
}
