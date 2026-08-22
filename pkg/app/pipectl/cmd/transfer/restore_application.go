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
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type restoreApplication struct {
	root *command

	inputFile        string
	pipedMappingFile string
	labels           []string
}

func newRestoreApplicationCommand(root *command) *cobra.Command {
	c := &restoreApplication{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "application",
		Short: "Restore applications from a backup file to the target control plane.",
		Long: `Restore application re-creates all applications from the backup file on the target control
plane. It requires a piped ID mapping file produced by 'pipectl transfer restore piped'.

Run this command only after the newly registered piped agents have connected to the target
control plane at least once. The control plane validates each application's Git repository
against the piped's registered repos, which are populated on first piped connection.

Disabled applications from the source are restored and immediately re-disabled on the target
to preserve their original status.`,
		RunE: cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.inputFile, "input-file", c.inputFile, "The path of the backup JSON file produced by 'pipectl transfer backup'.")
	cmd.Flags().StringVar(&c.pipedMappingFile, "piped-id-mapping-file", c.pipedMappingFile, "Path to the piped ID mapping JSON produced by 'pipectl transfer restore piped'.")
	cmd.Flags().StringSliceVar(&c.labels, "labels", c.labels, "Filter applications to restore by labels. Expect input as comma-separated KEY:VALUE pairs (e.g., --labels=env:prod,team:backend).")

	cmd.MarkFlagRequired("input-file")
	cmd.MarkFlagRequired("piped-id-mapping-file")

	return cmd
}

func (c *restoreApplication) run(ctx context.Context, input cli.Input) error {
	input.Logger.Info("Restoring applications...",
		zap.String("input-file", c.inputFile),
		zap.String("piped-id-mapping-file", c.pipedMappingFile),
	)

	labels := map[string]string{}
	for _, label := range c.labels {
		sp := strings.SplitN(label, ":", 2)
		if len(sp) == 2 {
			labels[sp[0]] = sp[1]
		}
	}

	data, err := readBackupFile(c.inputFile)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}
	input.Logger.Info(fmt.Sprintf("Found %d application(s) in backup (created at %s)", len(data.Applications), data.CreatedAt))

	// Filter applications by labels if specified.
	applications := filterApplicationsByLabels(data.Applications, labels)
	if len(labels) > 0 {
		input.Logger.Info(fmt.Sprintf("Filtered to %d application(s) matching labels", len(applications)))
	}

	_, pipedIDMap, err := loadPipedMapping(c.pipedMappingFile)
	if err != nil {
		return fmt.Errorf("failed to load piped ID mapping: %w", err)
	}
	input.Logger.Info(fmt.Sprintf("Loaded mapping for %d piped(s)", len(pipedIDMap)))

	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	restored, failed := restoreApplications(ctx, cli, applications, pipedIDMap, input.Logger)

	input.Logger.Info("Application restore completed",
		zap.Int("restored", restored),
		zap.Int("failed", failed),
	)
	if failed > 0 {
		return fmt.Errorf("%d application(s) failed to restore", failed)
	}
	return nil
}

// loadPipedMapping reads a RestoreResult JSON and returns the piped ID mapping.
func loadPipedMapping(path string) ([]PipedMapping, map[string]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read mapping file: %w", err)
	}
	var result RestoreResult
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal mapping file: %w", err)
	}
	pipedIDMap := make(map[string]string, len(result.PipedMappings))
	for _, m := range result.PipedMappings {
		pipedIDMap[m.OldPipedID] = m.NewPipedID
	}
	return result.PipedMappings, pipedIDMap, nil
}

// restoreApplications creates each application on the target control plane using the piped ID mapping.
// Disabled applications are re-disabled after creation to preserve their original status.
// Returns counts of successfully created and failed applications.
func restoreApplications(ctx context.Context, cli apiservice.Client, applications []*model.Application, pipedIDMap map[string]string, logger *zap.Logger) (restored, failed int) {
	for _, app := range applications {
		newPipedID, ok := pipedIDMap[app.PipedId]
		if !ok {
			logger.Warn("No piped mapping found for application, skipping",
				zap.String("app-name", app.Name),
				zap.String("app-id", app.Id),
				zap.String("piped-id", app.PipedId),
			)
			failed++
			continue
		}

		req := &apiservice.AddApplicationRequest{
			Name:             app.Name,
			PipedId:          newPipedID,
			GitPath:          app.GitPath,
			Kind:             app.Kind,             //nolint:staticcheck // deprecated but still used in AddApplicationRequest
			PlatformProvider: app.PlatformProvider, //nolint:staticcheck
			Description:      app.Description,
		}

		resp, err := cli.AddApplication(ctx, req)
		if err != nil {
			logger.Warn("Failed to restore application",
				zap.String("app-name", app.Name),
				zap.String("app-id", app.Id),
				zap.Error(err),
			)
			failed++
			continue
		}

		// Preserve disabled status — AddApplication always creates apps as enabled.
		if app.Disabled {
			if _, err := cli.DisableApplication(ctx, &apiservice.DisableApplicationRequest{
				ApplicationId: resp.ApplicationId,
			}); err != nil {
				logger.Warn("Application restored but failed to disable it",
					zap.String("app-name", app.Name),
					zap.String("new-id", resp.ApplicationId),
					zap.Error(err),
				)
			} else {
				logger.Info("Restored application (disabled)",
					zap.String("name", app.Name),
					zap.String("old-id", app.Id),
					zap.String("new-id", resp.ApplicationId),
				)
			}
		} else {
			logger.Info("Restored application",
				zap.String("name", app.Name),
				zap.String("old-id", app.Id),
				zap.String("new-id", resp.ApplicationId),
			)
		}
		restored++
	}
	return restored, failed
}

// filterApplicationsByLabels returns applications that match all specified labels.
// If labels is empty, all applications are returned.
func filterApplicationsByLabels(applications []*model.Application, labels map[string]string) []*model.Application {
	if len(labels) == 0 {
		return applications
	}
	filtered := make([]*model.Application, 0, len(applications))
	for _, app := range applications {
		if matchLabels(app.Labels, labels) {
			filtered = append(filtered, app)
		}
	}
	return filtered
}

// matchLabels returns true if appLabels contains all key-value pairs in filterLabels.
func matchLabels(appLabels, filterLabels map[string]string) bool {
	for k, v := range filterLabels {
		if appLabels[k] != v {
			return false
		}
	}
	return true
}
