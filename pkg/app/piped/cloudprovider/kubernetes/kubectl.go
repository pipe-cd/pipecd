// Copyright 2020 The PipeCD Authors.
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

package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"k8s.io/client-go/rest"

	"github.com/pipe-cd/pipecd/pkg/app/piped/cloudprovider/kubernetes/kubernetesmetrics"
)

var (
	errorReplaceNotFound = errors.New("specified resource is not found")
	errorNotFoundLiteral = "Error from server (NotFound)"
)

type Kubectl struct {
	version  string
	execPath string
	config   *rest.Config
}

func NewKubectl(version, path string) *Kubectl {
	return &Kubectl{
		version:  version,
		execPath: path,
	}
}

func (c *Kubectl) Apply(ctx context.Context, namespace string, manifest Manifest) (err error) {
	defer func() {
		kubernetesmetrics.IncKubectlCallsCounter(
			c.version,
			kubernetesmetrics.LabelApplyCommand,
			err == nil,
		)
	}()

	data, err := manifest.YamlBytes()
	if err != nil {
		return err
	}

	args := make([]string, 0, 5)
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "apply", "-f", "-")

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	r := bytes.NewReader(data)
	cmd.Stdin = r

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply: %s (%v)", string(out), err)
	}
	return nil
}

func (c *Kubectl) Create(ctx context.Context, namespace string, manifest Manifest) (err error) {
	defer func() {
		kubernetesmetrics.IncKubectlCallsCounter(
			c.version,
			kubernetesmetrics.LabelCreateCommand,
			err == nil,
		)
	}()

	data, err := manifest.YamlBytes()
	if err != nil {
		return err
	}

	args := make([]string, 0, 5)
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "create", "-f", "-")

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	r := bytes.NewReader(data)
	cmd.Stdin = r

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create: %s (%w)", string(out), err)
	}
	return nil
}

func (c *Kubectl) Replace(ctx context.Context, namespace string, manifest Manifest) (err error) {
	defer func() {
		kubernetesmetrics.IncKubectlCallsCounter(
			c.version,
			kubernetesmetrics.LabelReplaceCommand,
			err == nil,
		)
	}()

	data, err := manifest.YamlBytes()
	if err != nil {
		return err
	}

	args := make([]string, 0, 5)
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "replace", "-f", "-")

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	r := bytes.NewReader(data)
	cmd.Stdin = r

	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	if strings.HasPrefix(err.Error(), errorNotFoundLiteral) {
		return errorReplaceNotFound
	}

	return fmt.Errorf("failed to replace: %s (%w)", string(out), err)
}

func (c *Kubectl) Delete(ctx context.Context, namespace string, r ResourceKey) (err error) {
	defer func() {
		kubernetesmetrics.IncKubectlCallsCounter(
			c.version,
			kubernetesmetrics.LabelDeleteCommand,
			err == nil,
		)
	}()

	args := make([]string, 0, 5)
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "delete", r.Kind, r.Name)

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	out, err := cmd.CombinedOutput()

	if strings.Contains(string(out), "(NotFound)") {
		return fmt.Errorf("failed to delete: %s, (%w), %v", string(out), ErrNotFound, err)
	}
	if err != nil {
		return fmt.Errorf("failed to delete: %s, %v", string(out), err)
	}
	return nil
}
