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

package quickstart

var helmInstallScript = `
curl -L https://get.helm.sh/helm-v{{ .Version }}-linux-amd64.tar.gz | tar xvz -C {{ .BinDir }}
mv {{ .BinDir }}/linux-amd64/helm {{ .BinDir }}/helm
chmod +x {{ .BinDir }}/helm
`
