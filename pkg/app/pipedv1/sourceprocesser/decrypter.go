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

package sourceprocesser

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type secretDecrypterProcessor struct {
	enc *config.SecretEncryption
	dcr secretDecrypter
}

func NewSecretDecrypterProcessor(enc *config.SecretEncryption, dcr secretDecrypter) *secretDecrypterProcessor {
	return &secretDecrypterProcessor{
		enc: enc,
		dcr: dcr,
	}
}

func (s *secretDecrypterProcessor) BuildTemplateData(appDir string) (map[string]string, error) {
	if len(s.enc.EncryptedSecrets) == 0 {
		// Skip building no error.
		return nil, nil
	}

	secrets := make(map[string]string, len(s.enc.EncryptedSecrets))
	for k, v := range s.enc.EncryptedSecrets {
		ds, err := s.dcr.Decrypt(v)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt %s secret (%w)", k, err)
		}
		secrets[k] = ds
	}

	return secrets, nil
}

func (s *secretDecrypterProcessor) TemplateKey() string {
	return "encryptedSecrets"
}

func (s *secretDecrypterProcessor) TemplateSource(appDir string, data map[string]map[string]string) error {
	for _, t := range s.enc.DecryptionTargets {
		targetPath := filepath.Join(appDir, t)
		fileName := filepath.Base(targetPath)
		tmpl := template.New(fileName).Funcs(sprig.TxtFuncMap()).Option("missingkey=error")
		tmpl, err := tmpl.ParseFiles(targetPath)
		if err != nil {
			return fmt.Errorf("failed to parse target file %s (%w)", t, err)
		}

		f, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open target file %s (%w)", t, err)
		}

		if err := tmpl.Execute(f, data); err != nil {
			f.Close()
			return fmt.Errorf("failed to render target file %s (%w)", t, err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close target file %s (%w)", t, err)
		}
	}
	return nil
}
