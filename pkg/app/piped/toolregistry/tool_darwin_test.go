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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDarwinInstallScriptsUseHostArch(t *testing.T) {
	data := map[string]interface{}{
		"WorkingDir": "/tmp/work",
		"Version":    "1.2.3",
		"BinDir":     "/tmp/bin",
		"AsDefault":  false,
		"Arch":       "arm64",
	}

	var buf bytes.Buffer

	buf.Reset()
	assert.NoError(t, kubectlInstallScriptTmpl.Execute(&buf, data))
	assert.Contains(t, buf.String(), "bin/darwin/arm64/kubectl")
	assert.NotContains(t, buf.String(), "bin/darwin/amd64/kubectl")

	buf.Reset()
	assert.NoError(t, kustomizeInstallScriptTmpl.Execute(&buf, data))
	assert.Contains(t, buf.String(), "kustomize_v1.2.3_darwin_arm64.tar.gz")
	assert.NotContains(t, buf.String(), "kustomize_v1.2.3_darwin_amd64.tar.gz")

	buf.Reset()
	assert.NoError(t, helmInstallScriptTmpl.Execute(&buf, data))
	assert.Contains(t, buf.String(), "helm-v1.2.3-darwin-arm64.tar.gz")
	assert.NotContains(t, buf.String(), "helm-v1.2.3-darwin-amd64.tar.gz")

	buf.Reset()
	assert.NoError(t, terraformInstallScriptTmpl.Execute(&buf, data))
	assert.Contains(t, buf.String(), "terraform_1.2.3_darwin_arm64.zip")
	assert.NotContains(t, buf.String(), "terraform_1.2.3_darwin_amd64.zip")
}
