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

var kubectlInstallScript = `
cd {{ .WorkingDir }}
curl -LO https://storage.googleapis.com/kubernetes-release/release/v{{ .Version }}/bin/linux/amd64/kubectl
mv kubectl {{ .BinDir }}/kubectl-{{ .Version }}
chmod +x {{ .BinDir }}/kubectl-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/kubectl-{{ .Version }} {{ .BinDir }}/kubectl
{{ end }}
`

var kustomizeInstallScript = `
cd {{ .WorkingDir }}
curl -L https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v{{ .Version }}/kustomize_v{{ .Version }}_linux_amd64.tar.gz | tar xvz
mv kustomize {{ .BinDir }}/kustomize-{{ .Version }}
chmod +x {{ .BinDir }}/kustomize-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/kustomize-{{ .Version }} {{ .BinDir }}/kustomize
{{ end }}
`

var helmInstallScript = `
cd {{ .WorkingDir }}
curl -L https://get.helm.sh/helm-v{{ .Version }}-linux-amd64.tar.gz | tar xvz
mv linux-amd64/helm {{ .BinDir }}/helm-{{ .Version }}
chmod +x {{ .BinDir }}/helm-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/helm-{{ .Version }} {{ .BinDir }}/helm
{{ end }}
`

var terraformInstallScript = `
cd {{ .WorkingDir }}
curl https://releases.hashicorp.com/terraform/{{ .Version }}/terraform_{{ .Version }}_linux_amd64.zip -o terraform_{{ .Version }}_linux_amd64.zip
unzip terraform_{{ .Version }}_linux_amd64.zip
mv terraform {{ .BinDir }}/terraform-{{ .Version }}
{{ if .AsDefault }}
cp -f {{ .BinDir }}/terraform-{{ .Version }} {{ .BinDir }}/terraform
{{ end }}
`
