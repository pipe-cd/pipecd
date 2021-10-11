// Copyright 2021 The PipeCD Authors.
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

package sourcedecrypter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pipe-cd/pipe/pkg/config"
)

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

func DecryptSecrets(appDir string, enc config.SecretEncryption, dcr secretDecrypter) error {
	if len(enc.DecryptionTargets) == 0 {
		return nil
	}
	if len(enc.EncryptedSecrets) == 0 {
		return fmt.Errorf("no encrypted secret was specified to decrypt (%q)", enc.DecryptionTargets)
	}

	secrets := make(map[string]string, len(enc.EncryptedSecrets))
	for k, v := range enc.EncryptedSecrets {
		ds, err := dcr.Decrypt(v)
		if err != nil {
			return fmt.Errorf("failed to decrypt %s secret (%w)", k, err)
		}
		secrets[k] = ds
	}
	data := map[string](map[string]string){
		"encryptedSecrets": secrets,
	}

	for _, t := range enc.DecryptionTargets {
		targetPath := filepath.Join(appDir, t)
		tmpl, err := template.ParseFiles(targetPath)
		if err != nil {
			return fmt.Errorf("failed to parse decryption target %s (%w)", t, err)
		}

		// Return an error immediately if the target is using a nonexistent secret.
		tmpl = tmpl.Option("missingkey=error")

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

func DecryptSealedSecrets(appDir string, secrets []config.SealedSecretMapping, dcr secretDecrypter) error {
	for _, s := range secrets {
		secretPath := filepath.Join(appDir, s.Path)
		cfg, err := config.LoadFromYAML(secretPath)
		if err != nil {
			return fmt.Errorf("unable to read sealed secret file %s (%w)", s.Path, err)
		}
		if cfg.Kind != config.KindSealedSecret {
			return fmt.Errorf("unexpected kind in sealed secret file %s, want %q but got %q", s.Path, config.KindSealedSecret, cfg.Kind)
		}

		content, err := cfg.SealedSecretSpec.RenderOriginalContent(dcr)
		if err != nil {
			return fmt.Errorf("unable to render the original content of the sealed secret file %s (%w)", s.Path, err)
		}

		outDir, outFile := filepath.Split(s.Path)
		if s.OutFilename != "" {
			outFile = s.OutFilename
		}
		if s.OutDir != "" {
			outDir = s.OutDir
		}
		// TODO: Ensure that the output directory must be inside the application directory.
		if outDir != "" {
			if err := os.MkdirAll(filepath.Join(appDir, outDir), 0700); err != nil {
				return fmt.Errorf("unable to write decrypted content of sealed secret file %s to directory %s (%w)", s.Path, outDir, err)
			}
		}
		outPath := filepath.Join(appDir, outDir, outFile)

		if err := ioutil.WriteFile(outPath, content, 0644); err != nil {
			return fmt.Errorf("unable to write decrypted content of sealed secret file %s (%w)", s.Path, err)
		}
	}
	return nil
}
