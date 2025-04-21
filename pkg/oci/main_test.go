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

package oci

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	env        = "OCI_REGISTRY_HOST"
	port       = "5000"
	repository = "registry"
	tag        = "3.0.0@sha256:1fc7de654f2ac1247f0b67e8a459e273b0993be7d2beda1f3f56fbf1001ed3e7"
)

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

	log.Printf("Waiting for registry to be ready: %s", host)
	time.Sleep(1 * time.Second)

	code := m.Run()

	if err := res.Close(); err != nil {
		log.Fatalf("Failed to purge resource: %s", err)
	}

	os.Exit(code)
}
