// Copyright 2025 The PipeCD Authors.
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

package sdktest

import (
	"os"
	"testing"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

// LoadApplicationConfig loads the application config from the given filename.
// When the error occurs, it will call t.Fatal/t.Fatalf and stop the test.
func LoadApplicationConfig[Spec any](t *testing.T, filename string) *sdk.ApplicationConfig[Spec] {
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read application config: %s", err)
	}
	cfg, err := config.DecodeYAML[*sdk.ApplicationConfig[Spec]](data)
	if err != nil {
		t.Fatalf("failed to decode application config: %s", err)
	}
	if cfg.Spec == nil {
		t.Fatal("application config is not set")
	}
	return cfg.Spec
}
