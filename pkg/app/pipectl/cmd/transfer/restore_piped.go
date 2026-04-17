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

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type restorePiped struct {
	root *command

	inputFile  string
	outputFile string
}

func newRestorePipedCommand(root *command) *cobra.Command {
	c := &restorePiped{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "piped",
		Short: "Register pipeds from a backup file on the target control plane.",
		Long: `Register pipeds registers each piped from the backup file on the target control plane
and writes an old-to-new piped ID mapping to --output-file (or stdout).

Because the target control plane assigns new IDs and API keys, you must update each
piped's configuration file with the values from the mapping output before restarting
the piped agents.

After all piped agents have connected to the target control plane and registered their
repository configurations, run 'pipectl transfer restore application' to complete the restore.`,
		RunE: cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.inputFile, "input-file", c.inputFile, "The path of the backup JSON file produced by 'pipectl transfer backup'.")
	cmd.Flags().StringVar(&c.outputFile, "output-file", c.outputFile, "Path to write the piped ID mapping JSON. Defaults to stdout when not set.")
	cmd.MarkFlagRequired("input-file")

	return cmd
}

func (c *restorePiped) run(ctx context.Context, input cli.Input) error {
	input.Logger.Info("Restoring pipeds...", zap.String("input-file", c.inputFile))

	data, err := readBackupFile(c.inputFile)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}
	input.Logger.Info(fmt.Sprintf("Found %d piped(s) in backup (created at %s)", len(data.Pipeds), data.CreatedAt))

	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	mappings, _, err := registerPipeds(ctx, cli, data.Pipeds, input.Logger)
	if err != nil {
		return fmt.Errorf("failed to register pipeds: %w", err)
	}

	result := &RestoreResult{
		PipedMappings: mappings,
	}
	if err := writeMappingOutput(c.outputFile, result, input.Logger); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	input.Logger.Info("Piped restore completed. Update each piped config with the new ID and key, then restart the piped agents before running 'restore application'.",
		zap.Int("pipeds-registered", len(mappings)),
	)
	return nil
}

// readBackupFile reads and deserialises a backup JSON file.
func readBackupFile(path string) (*BackupData, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var data BackupData
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup data: %w", err)
	}
	return &data, nil
}

// registerPipeds calls RegisterPiped for each piped in the backup and returns the ID mapping.
func registerPipeds(ctx context.Context, cli apiservice.Client, pipeds []*model.Piped, logger *zap.Logger) ([]PipedMapping, map[string]string, error) {
	var mappings []PipedMapping
	pipedIDMap := make(map[string]string, len(pipeds))

	for _, p := range pipeds {
		resp, err := cli.RegisterPiped(ctx, &apiservice.RegisterPipedRequest{
			Name: p.Name,
			Desc: p.Desc,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to register piped %q (%s): %w", p.Name, p.Id, err)
		}

		mappings = append(mappings, PipedMapping{
			OldPipedID: p.Id,
			NewPipedID: resp.Id,
			NewKey:     resp.Key,
			PipedName:  p.Name,
		})
		pipedIDMap[p.Id] = resp.Id

		logger.Info("Registered piped",
			zap.String("name", p.Name),
			zap.String("old-id", p.Id),
			zap.String("new-id", resp.Id),
		)
	}
	return mappings, pipedIDMap, nil
}

// writeMappingOutput serialises result to JSON and writes it to outputPath (or stdout if empty).
func writeMappingOutput(outputPath string, result *RestoreResult, logger *zap.Logger) error {
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal mapping result: %w", err)
	}

	if outputPath == "" {
		logger.Info("Piped ID mapping (update your piped config files with the new IDs and keys):")
		fmt.Println(string(b))
		return nil
	}

	if err := os.WriteFile(outputPath, b, 0o600); err != nil {
		return fmt.Errorf("failed to write mapping file: %w", err)
	}
	logger.Info("Piped ID mapping written", zap.String("path", outputPath))
	return nil
}
