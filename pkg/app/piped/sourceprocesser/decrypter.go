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

	"github.com/pipe-cd/pipecd/pkg/config"
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

func (s *secretDecrypterProcessor) TargetFilePaths() []string {
	return s.enc.DecryptionTargets
}

func (s *secretDecrypterProcessor) TemplateKey() string {
	return "encryptedSecrets"
}
