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

package kubernetes

import (
	"log"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
)

func TestMain(m *testing.M) {
	binDir := "/tmp/piped-bin"
	if err := toolregistry.InitDefaultRegistry(binDir, zap.NewNop()); err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
