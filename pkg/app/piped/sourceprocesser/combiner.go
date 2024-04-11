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
	"maps"

	"github.com/pipe-cd/pipecd/pkg/config"
)

func EmbedCombination(appDir string, enc config.SecretEncryption, dcr secretDecrypter, atc config.Attachment) error {
	secretData, err := PrepareSecretData(enc, dcr)
	if err != nil {
		return err
	}

	attachData, err := PrepareAttachmentData(appDir, atc)
	if err != nil {
		return err
	}

	data := make(map[string](map[string]string))
	maps.Copy(data, secretData)
	maps.Copy(data, attachData)

	if len(data) == 0 {
		return nil
	}

	if err := EmbedSecret(appDir, enc.DecryptionTargets, data); err != nil {
		return err
	}

	if err := EmbedAttach(appDir, atc.Targets, data); err != nil {
		return err
	}

	return nil
}
