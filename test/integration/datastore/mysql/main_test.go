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

package mysql

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore/mysql"
	"github.com/pipe-cd/pipecd/pkg/datastore/mysql/ensurer"
)

const (
	db      = "mysql"
	version = "8.0"
)

var client *mysql.MySQL

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Failed to connect to docker: %s", err)
	}
	opts := &dockertest.RunOptions{
		Repository: db,
		Tag:        version,
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
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
		log.Fatalf("Failed to start resource: %s", err)
	}

	url := fmt.Sprintf("root:secret@tcp(localhost:%s)", res.GetPort("3306/tcp"))
	var ens ensurer.SQLEnsurer
	if err := pool.Retry(func() error {
		ens, err = ensurer.NewMySQLEnsurer(url, db, "", "", zap.NewNop())
		if err != nil {
			return err
		}
		return ens.Ping()
	}); err != nil {
		log.Fatalf("Failed to connect to docker by ensurer: %s", err)
	}

	ctx := context.Background()
	if err := ens.EnsureSchema(ctx); err != nil {
		log.Fatalf("Failed to prepare sql database: %s", err)
	}
	if err := ens.EnsureIndexes(ctx); err != nil {
		log.Fatalf("Failed to create required indexes on sql database: %s", err)
	}

	if err := pool.Retry(func() error {
		client, err = mysql.NewMySQL(url, db)
		if err != nil {
			return err
		}
		return client.Ping()
	}); err != nil {
		log.Fatalf("Failed to connect to docker by client: %s", err)
	}

	code := m.Run()

	if err := res.Close(); err != nil {
		log.Fatalf("Failed to purge resource: %s", err)
	}

	if err := ens.Close(); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}
