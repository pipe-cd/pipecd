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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	env        = "OCI_REGISTRY_HOST"
	port       = "5000"
	repository = "registry"
	tag        = "3.0.0"
)

var (
	registryHost     string
	registryResource *dockertest.Resource
	registrySetupErr error
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		registrySetupErr = err
		code := m.Run()
		os.Exit(code)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %s", err)
	}

	opts := &dockertest.RunOptions{
		Repository: repository,
		Tag:        tag,
		Env: []string{
			"REGISTRY_AUTH=htpasswd",
			"REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm",
			"REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd",
		},
		Mounts: []string{
			filepath.Join(wd, "testdata", "auth") + ":/auth",
		},
	}
	hcOpts := func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	}
	if err := pool.Client.Ping(); err != nil {
		registrySetupErr = err
		code := m.Run()
		os.Exit(code)
	}

	registryResource, err = pool.RunWithOptions(opts, hcOpts)
	if err != nil {
		registrySetupErr = err
		code := m.Run()
		os.Exit(code)
	}

	portID := fmt.Sprintf("%s/tcp", port)
	registryHost = fmt.Sprintf("localhost:%s", registryResource.GetPort(portID))
	os.Setenv(env, registryHost)

	log.Printf("Waiting for registry to be ready: %s", registryHost)
	time.Sleep(1 * time.Second)

	code := m.Run()

	if err := registryResource.Close(); err != nil {
		log.Fatalf("Failed to purge resource: %s", err)
	}

	os.Exit(code)
}

func isDockerUnavailable(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, os.ErrNotExist) {
		return true
	}

	msg := strings.ToLower(err.Error())
	for _, marker := range []string{
		"cannot connect to the docker daemon",
		"dial unix",
		"is the docker daemon running",
		"no such file or directory",
	} {
		if strings.Contains(msg, marker) {
			return true
		}
	}

	return false
}

func requireOCIRegistry(t *testing.T, repo string) string {
	t.Helper()

	if registrySetupErr != nil {
		if isDockerUnavailable(registrySetupErr) {
			t.Skipf("Skipping Docker-dependent OCI test because Docker is unavailable: %v", registrySetupErr)
		}
		t.Fatalf("could not set up OCI registry: %v", registrySetupErr)
	}

	return fmt.Sprintf("oci://%s/%s", registryHost, repo)
}
