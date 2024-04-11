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
	"github.com/pipe-cd/pipecd/pkg/config"
)

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

func PrepareSecretData(enc config.SecretEncryption, dcr secretDecrypter) (map[string](map[string]string), error) {
	if len(enc.DecryptionTargets) == 0 {
		return make(map[string]map[string]string), nil
	}
	if len(enc.EncryptedSecrets) == 0 {
		return make(map[string]map[string]string), nil
	}

	secrets := make(map[string]string, len(enc.EncryptedSecrets))
	for k, v := range enc.EncryptedSecrets {
		ds, err := dcr.Decrypt(v)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt %s secret (%w)", k, err)
		}
		secrets[k] = ds
	}
	data := map[string](map[string]string){
		"encryptedSecrets": secrets,
	}

	return data, nil
}

func EmbedSecret(appDir string, targets []string, data map[string](map[string]string)) error {
	for _, t := range targets {
		targetPath := filepath.Join(appDir, t)
		fileName := filepath.Base(targetPath)
		tmpl := template.
			New(fileName).
			Funcs(sprig.TxtFuncMap()).
			Option("missingkey=error")
		tmpl, err := tmpl.ParseFiles(targetPath)
		if err != nil {
			return fmt.Errorf("failed to parse decryption target %s (%w)", t, err)
		}

		f, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open decryption target %s (%w)", t, err)
		}

		if err := tmpl.Execute(f, data); err != nil {
			f.Close()
			return fmt.Errorf("failed to render decryption target %s (%w)", t, err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close decryption target %s (%w)", t, err)
		}
	}

	return nil
}

func DecryptSecrets(appDir string, enc config.SecretEncryption, dcr secretDecrypter) error {
	data, err := PrepareSecretData(enc, dcr)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	if err := EmbedSecret(appDir, enc.DecryptionTargets, data); err != nil {
		return err
	}

	return nil
}
