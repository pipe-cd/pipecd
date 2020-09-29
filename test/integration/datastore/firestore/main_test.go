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

	"go.uber.org/atomic"
)

const emulatorHost = "localhost:8080"

// buffer wraps `bytes.Buffer` to prevent Buffers from a data race.
type buffer struct {
	b *bytes.Buffer
	s *atomic.String
}

func newBuffer() *buffer {
	return &buffer{
		b: new(bytes.Buffer),
		s: atomic.NewString(""),
	}
}

func (b *buffer) Write(p []byte) (n int, err error) {
	b.s.Store(b.s.String() + string(p))
	return b.b.Write(p)
}

func (b *buffer) String() string {
	return b.s.Load()
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	os.Setenv("FIRESTORE_EMULATOR_HOST", emulatorHost)
	cmd := exec.CommandContext(ctx, "gcloud", "beta", "emulators", "firestore", "start", fmt.Sprintf("--host-port=%s", emulatorHost))

	b := newBuffer()
	cmd.Stdout = b
	cmd.Stderr = b

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	log.Printf("=== Firestore Emulator Output ===\n%s\n=== Firestore Emulator Output End ===\n", b.String())

	os.Exit(code)
}
