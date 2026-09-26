package oci

import (
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

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Printf("Docker daemon is not available: %s. Skipping pkg/oci tests.", err)
		os.Exit(0)
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
	res, err := pool.RunWithOptions(opts, hcOpts)
	if err != nil {
		// Check if the error indicates Docker daemon is down or unreachable
		errMsg := err.Error()
		if strings.Contains(errMsg, "docker.sock") || strings.Contains(errMsg, "Connection refused") || strings.Contains(errMsg, "is the daemon running") {
			log.Printf("Docker daemon is unavailable: %s. Skipping pkg/oci tests.", err)
			os.Exit(0)
		}
		// If it's some other unexpected error (e.g., config error), fail loudly as expected
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