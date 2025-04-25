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

package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/oci"
)

type push struct {
	root *command

	files      []string
	tag        string
	insecure   bool
	registry   string
	repository string
}

func newPushCommand(root *command) *cobra.Command {
	p := &push{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push a plugin to the server.",
		RunE:  cli.WithContext(p.run),
	}

	cmd.Flags().StringSliceVar(&p.files, "files", p.files, "The list of files to push. For example, 'linux/amd64=file1,linux/arm64=file2'")
	cmd.Flags().StringVar(&p.tag, "tag", p.tag, "The tag of the plugin.")
	cmd.Flags().BoolVar(&p.insecure, "insecure", p.insecure, "If true, the plugin will be pushed to the server without TLS verification.")
	cmd.Flags().StringVar(&p.registry, "registry", p.registry, "The registry of the plugin.")
	cmd.Flags().StringVar(&p.repository, "repository", p.repository, "The repository of the plugin.")

	cmd.MarkFlagRequired("files")
	cmd.MarkFlagRequired("tag")
	cmd.MarkFlagRequired("registry")
	cmd.MarkFlagRequired("repository")

	return cmd
}

func (p *push) run(ctx context.Context, input cli.Input) error {
	workdir, err := os.MkdirTemp("", "pipectl-plugin-push")
	if err != nil {
		input.Logger.Error("failed to create temp directory", zap.Error(err))
		return err
	}
	defer os.RemoveAll(workdir)

	targetURL := fmt.Sprintf("%s/%s:%s", p.registry, p.repository, p.tag)

	opts := make([]oci.PushOption, 0, 1)
	if p.insecure {
		opts = append(opts, oci.WithInsecure())
	}

	files, err := p.parseFilePaths(p.files)
	if err != nil {
		input.Logger.Error("failed to parse file paths", zap.Error(err))
		return err
	}

	input.Logger.Info("pushing plugin to the server", zap.String("target", targetURL), zap.Any("files", files))

	artifact := &oci.Artifact{
		MediaType:    oci.MediaTypePipedPlugin,
		ArtifactType: oci.ArtifactTypePipedPlugin,
		FilePaths:    files,
	}

	if err := oci.PushFilesToRegistry(ctx, workdir, artifact, targetURL, opts...); err != nil {
		input.Logger.Error("failed to push plugin to the server", zap.Error(err))
		return err
	}

	input.Logger.Info("successfully pushed plugin to the server", zap.String("target", targetURL))

	return nil
}

func (p *push) parseFilePaths(fps []string) (map[oci.Platform]string, error) {
	files := make(map[oci.Platform]string, len(fps))
	for _, fp := range fps {
		platform, path, ok := strings.Cut(fp, "=")
		if !ok {
			return nil, fmt.Errorf("invalid file format: %s", fp)
		}

		p, err := p.parsePlatform(platform)
		if err != nil {
			return nil, fmt.Errorf("invalid platform format: %s", platform)
		}

		files[p] = path
	}

	return files, nil
}

func (p *push) parsePlatform(platform string) (oci.Platform, error) {
	parts := strings.Split(platform, "/")
	if len(parts) < 2 {
		return oci.Platform{}, fmt.Errorf("invalid platform format: %s", platform)
	}
	if len(parts) > 2 {
		return oci.Platform{}, fmt.Errorf("current implementation only supports OS/Arch format and does not support variant: %s", platform)
	}

	return oci.Platform{
		OS:   parts[0],
		Arch: parts[1],
	}, nil
}
