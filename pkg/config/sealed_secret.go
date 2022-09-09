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

	"github.com/Masterminds/sprig/v3"
)

// SealedSecretSpec holds the data of a sealed secret.
type SealedSecretSpec struct {
	// A string that represents the encrypted data of the original file.
	// When this is configured, the template and encryptedItems fields will be ignored.
	EncryptedData string
	// The template used to restore the original content.
	Template string
	// A list of encrypted items that will be decrypted and inserted to
	// the specified template to render the original content.
	EncryptedItems map[string]string
}

func (s *SealedSecretSpec) Validate() error {
	if s.EncryptedData != "" {
		return nil
	}
	if len(s.EncryptedItems) == 0 {
		return fmt.Errorf("either encryptedData or encryptedItems must be set")
	}
	if s.Template == "" {
		return fmt.Errorf("the template must be set")
	}
	return nil
}

type sealedSecretDecrypter interface {
	Decrypt(string) (string, error)
}

func (s *SealedSecretSpec) RenderOriginalContent(dcr sealedSecretDecrypter) ([]byte, error) {
	if s.EncryptedData != "" {
		decryptedData, err := dcr.Decrypt(s.EncryptedData)
		if err != nil {
			return nil, err
		}
		return []byte(decryptedData), nil
	}

	decryptedItems := make(map[string]string, len(s.EncryptedItems))
	for k, v := range s.EncryptedItems {
		text, err := dcr.Decrypt(v)
		if err != nil {
			return nil, fmt.Errorf("unable to decrypt %s item (%w)", k, err)
		}
		decryptedItems[k] = text
	}

	tmpl, err := template.New("sealedsecret").Funcs(sprig.TxtFuncMap()).Option("missingkey=error").Parse(s.Template)
	if err != nil {
		return nil, fmt.Errorf("unable to parse secret template (%w)", err)
	}

	var out bytes.Buffer
	data := map[string]interface{}{
		"encryptedItems": decryptedItems,
	}
	if err := tmpl.Execute(&out, data); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
