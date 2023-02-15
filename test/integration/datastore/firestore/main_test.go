// Copyright 2023 The PipeCD Authors.
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
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/pipe-cd/pipecd/pkg/datastore/firestore"
)

const (
	env        = "FIRESTORE_EMULATOR_HOST"
	port       = "8080"
	repository = "ghcr.io/pipe-cd/firestore-emulator"
	tag        = "v0.34.0-3-gf22c209"
	project    = "pipecd-test"
)

var store *firestore.FireStore

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Failed to connect to docker: %s", err)
	}
	opts := &dockertest.RunOptions{
		Repository: repository,
		Tag:        tag,
	}
	hcOpts := func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	}
	res, err := pool.RunWithOptions(opts, hcOpts)
	if err != nil {
		log.Fatalf("Failed to start resource: %s", err)
	}

	portID := fmt.Sprintf("%s/tcp", port)
	host := fmt.Sprintf("localhost:%s", res.GetPort(portID))
	os.Setenv(env, host)

	ctx := context.Background()
	store, err = firestore.NewFireStore(ctx, project, "namespace", "environment")
	if err != nil {
		log.Fatalf("Failed to connect to docker: %s", err)
	}

	code := m.Run()

	if err := res.Close(); err != nil {
		log.Fatalf("Failed to purge resource: %s", err)
	}

	if err := store.Close(); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}
