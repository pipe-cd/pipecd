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

package toolregistry

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	"go.uber.org/zap"
)

const (
	defaultKubectlVersion   = "1.18.2"
	defaultKustomizeVersion = "3.8.1"
	defaultHelmVersion      = "3.8.2"
	defaultTerraformVersion = "0.13.0"
)

var installScripts = map[string]map[string]string{
	"kubectl": {
		"linux": `
cd {{ .WorkingDir }}
curl -LO https://dl.k8s.io/release/v{{ .Version }}/bin/linux/amd64/kubectl
mv kubectl {{ .BinDir }}/kubectl-{{ .Version }}
chmod +x {{ .BinDir }}/kubectl-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/kubectl-{{ .Version }} {{ .BinDir }}/kubectl
{{ end }}
`,
		"darwin": `
cd {{ .WorkingDir }}
curl -LO https://dl.k8s.io/release/v{{ .Version }}/bin/darwin/amd64/kubectl
mv kubectl {{ .BinDir }}/kubectl-{{ .Version }}
chmod +x {{ .BinDir }}/kubectl-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/kubectl-{{ .Version }} {{ .BinDir }}/kubectl
{{ end }}
`,
	},
	"kustomize": {
		"linux": `
cd {{ .WorkingDir }}
curl -L https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v{{ .Version }}/kustomize_v{{ .Version }}_linux_amd64.tar.gz | tar xvz
mv kustomize {{ .BinDir }}/kustomize-{{ .Version }}
chmod +x {{ .BinDir }}/kustomize-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/kustomize-{{ .Version }} {{ .BinDir }}/kustomize
{{ end }}
`,
		"darwin": `
cd {{ .WorkingDir }}
curl -L https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v{{ .Version }}/kustomize_v{{ .Version }}_darwin_amd64.tar.gz | tar xvz
mv kustomize {{ .BinDir }}/kustomize-{{ .Version }}
chmod +x {{ .BinDir }}/kustomize-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/kustomize-{{ .Version }} {{ .BinDir }}/kustomize
{{ end }}
`,
	},
	"helm": {
		"linux": `
cd {{ .WorkingDir }}
curl -L https://get.helm.sh/helm-v{{ .Version }}-linux-amd64.tar.gz | tar xvz
mv linux-amd64/helm {{ .BinDir }}/helm-{{ .Version }}
chmod +x {{ .BinDir }}/helm-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/helm-{{ .Version }} {{ .BinDir }}/helm
{{ end }}
`,
		"darwin": `
cd {{ .WorkingDir }}
curl -L https://get.helm.sh/helm-v{{ .Version }}-darwin-amd64.tar.gz | tar xvz
mv darwin-amd64/helm {{ .BinDir }}/helm-{{ .Version }}
chmod +x {{ .BinDir }}/helm-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/helm-{{ .Version }} {{ .BinDir }}/helm
{{ end }}
`,
	},
	"terraform": {
		"linux": `
cd {{ .WorkingDir }}
curl https://releases.hashicorp.com/terraform/{{ .Version }}/terraform_{{ .Version }}_linux_amd64.zip -o terraform_{{ .Version }}_linux_amd64.zip
unzip terraform_{{ .Version }}_linux_amd64.zip
mv terraform {{ .BinDir }}/terraform-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/terraform-{{ .Version }} {{ .BinDir }}/terraform
{{ end }}
`,
		"darwin": `
cd {{ .WorkingDir }}
curl https://releases.hashicorp.com/terraform/{{ .Version }}/terraform_{{ .Version }}_darwin_amd64.zip -o terraform_{{ .Version }}_linux_amd64.zip
unzip terraform_{{ .Version }}_linux_amd64.zip
mv terraform {{ .BinDir }}/terraform-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/terraform-{{ .Version }} {{ .BinDir }}/terraform
{{ end }}
`,
	},
}

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
		data = map[string]any{
			"WorkingDir": workingDir,
			"Version":    version,
			"BinDir":     r.binDir,
			"AsDefault":  asDefault,
		}
	)
	kubectlInstallScriptTmpl := template.Must(template.New("kubectl").Parse(installScripts["kubectl"][runtime.GOOS]))
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
		data = map[string]any{
			"WorkingDir": workingDir,
			"Version":    version,
			"BinDir":     r.binDir,
			"AsDefault":  asDefault,
		}
	)
	kustomizeInstallScriptTmpl := template.Must(template.New("kustomize").Parse(installScripts["kustomize"][runtime.GOOS]))
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
		data = map[string]any{
			"WorkingDir": workingDir,
			"Version":    version,
			"BinDir":     r.binDir,
			"AsDefault":  asDefault,
		}
	)
	helmInstallScriptTmpl := template.Must(template.New("helm").Parse(installScripts["helm"][runtime.GOOS]))
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
		data = map[string]any{
			"WorkingDir": workingDir,
			"Version":    version,
			"BinDir":     r.binDir,
			"AsDefault":  asDefault,
		}
	)
	terraformInstallScriptTmpl := template.Must(template.New("terraform").Parse(installScripts["terraform"][runtime.GOOS]))
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
