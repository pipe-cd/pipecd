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

package config

import (
	"bytes"
	"fmt"
	"text/template"
)

// SealedSecretSpec holds the data of a sealed secret.
type SealedSecretSpec struct {
	// The template used to restore the original content.
	// Empty means the original content is same with the decrypted data of the first encrypted item.
	Template string
	// A string that represents the encrypted data of the original file.
	EncryptedData map[string]string
}

func (s *SealedSecretSpec) Validate() error {
	if len(s.EncryptedData) == 0 {
		return fmt.Errorf("encryptedData must contain at least one item")
	}
	return nil
}

func (s *SealedSecretSpec) RenderOriginalContent(secrets map[string]string) ([]byte, error) {
	if len(secrets) == 0 {
		return nil, fmt.Errorf("require at least one secret")
	}

	// If the template was not configured, the first secret will be used as the original content.
	if s.Template == "" {
		for _, v := range secrets {
			return []byte(v), nil
		}
	}

	tmpl, err := template.New("sealedsecret").Option("missingkey=error").Parse(s.Template)
	if err != nil {
		return nil, fmt.Errorf("unable to parse secret template (%w)", err)
	}

	var out bytes.Buffer
	data := map[string]interface{}{
		"encryptedData": secrets,
	}
	if err := tmpl.Execute(&out, data); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
