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

package toolregistry

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
)

const (
	defaultKubectlVersion   = "1.18.2"
	defaultKustomizeVersion = "3.8.1"
	defaultHelmVersion      = "3.8.2"
	defaultTerraformVersion = "0.13.0"
)

var (
	kubectlInstallScriptTmpl   = template.Must(template.New("kubectl").Parse(kubectlInstallScript))
	kustomizeInstallScriptTmpl = template.Must(template.New("kustomize").Parse(kustomizeInstallScript))
	helmInstallScriptTmpl      = template.Must(template.New("helm").Parse(helmInstallScript))
	terraformInstallScriptTmpl = template.Must(template.New("terraform").Parse(terraformInstallScript))
)

func (r *registry) installKubectl(ctx context.Context, version string) error {
	workingDir, err := os.MkdirTemp("", "kubectl-install")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workingDir)

	asDefault := version == ""
	if asDefault {
		version = defaultKubectlVersion
	}

	var (
		buf  bytes.Buffer
		data = map[string]interface{}{
			"WorkingDir": workingDir,
			"Version":    version,
			"BinDir":     r.binDir,
			"AsDefault":  asDefault,
		}
	)
	if err := kubectlInstallScriptTmpl.Execute(&buf, data); err != nil {
		r.logger.Error("failed to render kubectl install script",
			zap.String("version", version),
			zap.Error(err),
		)
		return fmt.Errorf("failed to install kubectl %s (%v)", version, err)
	}

	var (
		script = buf.String()
		cmd    = exec.CommandContext(ctx, "/bin/sh", "-c", script)
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		r.logger.Error("failed to install kubectl",
			zap.String("version", version),
			zap.String("script", script),
			zap.String("out", string(out)),
			zap.Error(err),
		)
		return fmt.Errorf("failed to install kubectl %s (%v)", version, err)
	}

	r.logger.Info("just installed kubectl", zap.String("version", version))
	return nil
}

func (r *registry) installKustomize(ctx context.Context, version string) error {
	workingDir, err := os.MkdirTemp("", "kustomize-install")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workingDir)

	asDefault := version == ""
	if asDefault {
		version = defaultKustomizeVersion
	}

	var (
		buf  bytes.Buffer
		data = map[string]interface{}{
			"WorkingDir": workingDir,
			"Version":    version,
			"BinDir":     r.binDir,
			"AsDefault":  asDefault,
		}
	)
	if err := kustomizeInstallScriptTmpl.Execute(&buf, data); err != nil {
		r.logger.Error("failed to render kustomize install script",
			zap.String("version", version),
			zap.Error(err),
		)
		return fmt.Errorf("failed to install kustomize %s (%v)", version, err)
	}

	var (
		script = buf.String()
		cmd    = exec.CommandContext(ctx, "/bin/sh", "-c", script)
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		r.logger.Error("failed to install kustomize",
			zap.String("version", version),
			zap.String("script", script),
			zap.String("out", string(out)),
			zap.Error(err),
		)
		return fmt.Errorf("failed to install kustomize %s (%v)", version, err)
	}

	r.logger.Info("just installed kustomize", zap.String("version", version))
	return nil
}

func (r *registry) installHelm(ctx context.Context, version string) error {
	workingDir, err := os.MkdirTemp("", "helm-install")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workingDir)

	asDefault := version == ""
	if asDefault {
		version = defaultHelmVersion
	}

	var (
		buf  bytes.Buffer
		data = map[string]interface{}{
			"WorkingDir": workingDir,
			"Version":    version,
			"BinDir":     r.binDir,
			"AsDefault":  asDefault,
		}
	)
	if err := helmInstallScriptTmpl.Execute(&buf, data); err != nil {
		r.logger.Error("failed to render helm install script",
			zap.String("version", version),
			zap.Error(err),
		)
		return fmt.Errorf("failed to install helm %s (%v)", version, err)
	}

	var (
		script = buf.String()
		cmd    = exec.CommandContext(ctx, "/bin/sh", "-c", script)
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		r.logger.Error("failed to install helm",
			zap.String("version", version),
			zap.String("script", script),
			zap.String("out", string(out)),
			zap.Error(err),
		)
		return fmt.Errorf("failed to install helm %s (%v)", version, err)
	}

	r.logger.Info("just installed helm", zap.String("version", version))
	return nil
}

func (r *registry) installTerraform(ctx context.Context, version string) error {
	workingDir, err := os.MkdirTemp("", "terraform-install")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workingDir)

	asDefault := version == ""
	if asDefault {
		version = defaultTerraformVersion
	}

	var (
		buf  bytes.Buffer
		data = map[string]interface{}{
			"WorkingDir": workingDir,
			"Version":    version,
			"BinDir":     r.binDir,
			"AsDefault":  asDefault,
		}
	)
	if err := terraformInstallScriptTmpl.Execute(&buf, data); err != nil {
		r.logger.Error("failed to render terraform install script",
			zap.String("version", version),
			zap.Error(err),
		)
		return fmt.Errorf("failed to install terraform %s (%w)", version, err)
	}

	var (
		script = buf.String()
		cmd    = exec.CommandContext(ctx, "/bin/sh", "-c", script)
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		r.logger.Error("failed to install terraform",
			zap.String("version", version),
			zap.String("script", script),
			zap.String("out", string(out)),
			zap.Error(err),
		)
		return fmt.Errorf("failed to install terraform %s, %s (%w)", version, string(out), err)
	}

	r.logger.Info("just installed terraform", zap.String("version", version))
	return nil
}

func (r *registry) installExternalTool(ctx context.Context, config config.ExternalTool) error {
	script := fmt.Sprintf("asdf install %s %s", config.Package, config.Version)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctxWithTimeout, "/bin/sh", "-c", script)
	if out, err := cmd.CombinedOutput(); err != nil {
		r.logger.Error("failed to install %s %s",
			zap.String("package", config.Package),
			zap.String("version", config.Version),
			zap.String("out", string(out)),
			zap.Error(err),
		)
		if errors.Is(ctxWithTimeout.Err(), context.DeadlineExceeded) {
			return errors.Errorf("failed to install %s %s (%v) because of timeout", config.Package, config.Version, err)
		}
		return errors.Errorf("failed to install %s %s (%v)", config.Package, config.Version, err)
	}

	r.logger.Info("just installed external tool",
		zap.String("package", config.Package),
		zap.String("version", config.Version),
	)
	return nil
}
