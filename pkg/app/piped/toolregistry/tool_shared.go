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

var asdfInstallScriptBash = `
git clone https://github.com/asdf-vm/asdf.git {{ .HomeDir }}/.asdf --branch v0.11.3
echo -e "\n. {{ .HomeDir }}/.asdf/asdf.sh" >> {{ .HomeDir }}/.bashrc
echo -e "\n. {{ .HomeDir }}/.asdf/completions/asdf.bash" >> {{ .HomeDir }}/.bashrc
`

var asdfInstallScriptFish = `
git clone https://github.com/asdf-vm/asdf.git {{ .HomeDir }}/.asdf --branch v0.11.3
echo -e "\nsource {{ .HomeDir }}/.asdf/asdf.fish" >> {{ .HomeDir }}/.config/fish/config.fish
mkdir -p {{ .HomeDir }}/.config/fish/completions; and ln -s {{ .HomeDir }}/.asdf/completions/asdf.fish {{ .HomeDir }}/.config/fish/completions
`

var asdfInstallScriptElvish = `
git clone https://github.com/asdf-vm/asdf.git {{ .HomeDir }}/.asdf --branch v0.11.3
mkdir -p {{ .HomeDir }}/.config/elvish/lib; ln -s {{ .HomeDir }}/.asdf/asdf.elv {{ .HomeDir }}/.config/elvish/lib/asdf.elv
echo "\n"'use asdf _asdf; var asdf~ = $_asdf:asdf~' >> {{ .HomeDir }}/.config/elvish/rc.elv
echo "\n"'set edit:completion:arg-completer[asdf] = $_asdf:arg-completer~' >> {{ .HomeDir }}/.config/elvish/rc.elv
`

var asdfInstallScriptZsh = `
git clone https://github.com/asdf-vm/asdf.git {{ .HomeDir }}/.asdf --branch v0.11.3
echo -e "\n. {{ .HomeDir }}/.asdf/asdf.sh\nfpath=(${ASDF_DIR}/completions $fpath)\nautoload -Uz compinit && compinit" >> {{ .HomeDir }}/.config/fish/config.fish
`
