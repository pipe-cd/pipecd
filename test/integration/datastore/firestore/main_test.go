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

package firestore

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
)

const emulatorHost = "localhost:8080"

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())

	os.Setenv("FIRESTORE_EMULATOR_HOST", emulatorHost)
	cmd := exec.CommandContext(ctx, "gcloud", "beta", "emulators", "firestore", "start", fmt.Sprintf("--host-port=%s", emulatorHost))

	b := new(bytes.Buffer)
	cmd.Stdout = b
	cmd.Stderr = b
	defer func() {
		cancel()
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
		log.Printf("=== Firestore Emulator Output ===\n%s\n=== Firestore Emulator Output End ===\n", b.String())
	}()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	os.Exit(code)
}
