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

import "fmt"

// SealedSecretSpec holds the data of a sealed secret.
type SealedSecretSpec struct {
	// A string that represents the encrypted data of the original file.
	EncryptedData string
}

func (s *SealedSecretSpec) Validate() error {
	if s.EncryptedData == "" {
		return fmt.Errorf("encryptedData must be set")
	}
	return nil
}
