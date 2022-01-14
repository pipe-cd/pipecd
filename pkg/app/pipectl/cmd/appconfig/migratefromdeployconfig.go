// Copyright 2022 The PipeCD Authors.
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

package appconfig

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type migrateFromDeployConfig struct {
	root *command

	repoRootPath string
	repoID       string
	envName      string
	stdout       io.Writer
}

func newListCommand(root *command) *cobra.Command {
	m := &migrateFromDeployConfig{
		root:   root,
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "migrate-from-deploy-config",
		Short: "Migrate local deployment configuration files to application configuration files. Basically this finds all deployment configuration files to add name field and env label to them.",
		RunE:  cli.WithContext(m.run),
	}

	cmd.Flags().StringVar(&m.repoRootPath, "repo-root-path", m.repoRootPath, "The absolute path to the root directory of the Git repository.")
	cmd.Flags().StringVar(&m.repoID, "repo-id", m.repoID, "The repository ID that is being registered in Piped config.")
	cmd.Flags().StringVar(&m.envName, "env-name", m.envName, "The environment name.")

	cmd.MarkFlagRequired("repo-root-path")
	cmd.MarkFlagRequired("repo-id")
	cmd.MarkFlagRequired("env-name")

	return cmd
}

func (m *migrateFromDeployConfig) run(ctx context.Context, _ cli.Input) error {
	cli, err := m.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	fmt.Fprintln(m.stdout, "Start finding and migrating deployment configuration files to application configuration files...")

	var cursor string
	var count int
	for {
		req := &apiservice.ListApplicationsRequest{
			EnvName: m.envName,
			Cursor:  cursor,
		}
		resp, err := cli.ListApplications(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list application: %w", err)
		}

		for _, app := range resp.Applications {
			if app.GitPath.Repo.Id != m.repoID {
				continue
			}
			if err := m.migrate(ctx, app); err != nil {
				return err
			}
			count++
		}
		if resp.Cursor == "" {
			break
		}
		cursor = resp.Cursor
	}

	fmt.Fprintf(m.stdout, "Successfully migrated %d applications\n", count)
	return nil
}

func (m *migrateFromDeployConfig) migrate(ctx context.Context, app *model.Application) error {
	configFilePath := filepath.Join(m.repoRootPath, app.GitPath.GetApplicationConfigFilePath())
	oriData, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	newData, err := convert(oriData, app.Name, m.envName, app.Description)
	if err != nil {
		return err
	}

	info, err := os.Stat(configFilePath)
	if err != nil {
		return err
	}
	return os.WriteFile(configFilePath, newData, info.Mode())
}

func convert(data []byte, name, env, description string) ([]byte, error) {
	var out strings.Builder
	var shouldWrite, done bool

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if shouldWrite {
			indent, ok := determineIndent(line)
			if !ok {
				fmt.Fprintf(&out, "%s\n", line)
				continue
			}

			writeNewFields(&out, name, env, description, indent)
			shouldWrite = false
			done = true
		}

		fmt.Fprintf(&out, "%s\n", line)
		if done {
			continue
		}
		if strings.HasPrefix(line, "spec:") {
			shouldWrite = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return []byte(out.String()), nil
}

func determineIndent(line string) (string, bool) {
	noIndentLine := strings.TrimLeft(line, " \t")
	// In case of just an empty line.
	if noIndentLine == "" {
		return "", false
	}
	// In case of just a comment.
	if strings.HasPrefix(noIndentLine, "#") {
		return "", false
	}
	return line[:len(line)-len(noIndentLine)], true
}

func writeNewFields(out io.Writer, name, env, description, indent string) {
	doubleIndent := strings.Repeat(indent, 2)

	fmt.Fprintf(out, "%sname: %s\n", indent, name)
	fmt.Fprintf(out, "%slabels:\n", indent)
	fmt.Fprintf(out, "%senv: %s\n", doubleIndent, env)
	if description != "" {
		fmt.Fprintf(out, "%sdescription: |\n", indent)
		for _, s := range strings.Split(description, "\n") {
			fmt.Fprintf(out, "%s%s\n", doubleIndent, s)
		}
	}
}
