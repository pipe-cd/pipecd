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
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"
)

const emulatorHost = "localhost:8080"

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	os.Setenv("FIRESTORE_EMULATOR_HOST", emulatorHost)
	cmd := exec.CommandContext(ctx, "gcloud", "beta", "emulators", "firestore", "start", fmt.Sprintf("--host-port=%s", emulatorHost))

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 256, 256)
	// Spawn another goroutine to see the emulator output.
	go func() {
		for {
			n, err := stderr.Read(buf[:])
			if err != nil && err == io.EOF {
				break
			}
			if n > 0 {
				log.Printf("%s", string(buf[:n]))
			}
		}
	}()

	os.Exit(m.Run())
}
