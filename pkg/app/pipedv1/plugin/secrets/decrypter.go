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

package secrets

import (
	"context"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/cmd/piped/service"
)

type Decrypter struct {
	client service.PluginServiceClient
}

func (d *Decrypter) Decrypt(src string) (string, error) {
	r, err := d.client.DecryptSecret(context.TODO(), &service.DecryptSecretRequest{Secret: src})
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret: %w", err)
	}
	return r.GetDecryptedSecret(), nil
}
