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

var asdfInstallScript = `
git clone https://github.com/asdf-vm/asdf.git {{ .HomeDir }}/.asdf --branch v{{ .Version }}
cat <<EOT > {{ .HomeDir }}/.profile
export ASDF_DIR="$HOME/.asdf"
. {{ .HomeDir }}/.asdf/asdf.sh
EOT
`
