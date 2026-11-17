// Copyright 2026 The PipeCD Authors.
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

package migrate

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
)

// migrateApplicationConfig reads a user-supplied YAML file into map[string]any,
// so every lookup out of it can hold something other than the expected shape.
// A malformed file has to be reported, not panic the CLI.
func TestMigrateApplicationConfigMalformed(t *testing.T) {
	cases := map[string]string{
		"spec is a string":   "apiVersion: pipecd.dev/v1beta1\nkind: KubernetesApp\nspec: oops\n",
		"spec missing":       "apiVersion: pipecd.dev/v1beta1\nkind: KubernetesApp\n",
		"kind missing":       "apiVersion: pipecd.dev/v1beta1\nspec:\n  name: x\n",
		"kind is a number":   "apiVersion: pipecd.dev/v1beta1\nkind: 42\nspec:\n  name: x\n",
		"pipeline is a list": "apiVersion: pipecd.dev/v1beta1\nkind: KubernetesApp\nspec:\n  pipeline:\n    - a\n",
		"stages is a string": "apiVersion: pipecd.dev/v1beta1\nkind: KubernetesApp\nspec:\n  pipeline:\n    stages: nope\n",
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			f := filepath.Join(dir, "app.pipecd.yaml")
			if err := os.WriteFile(f, []byte(body), 0o600); err != nil {
				t.Fatal(err)
			}
			c := &applicationConfig{}
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked on a user-supplied config: %v", r)
				}
			}()
			if err := c.migrateApplicationConfig(context.Background(), f, zap.NewNop()); err == nil {
				t.Error("expected an error for a malformed config, got nil")
			}
		})
	}
}

// The v0 loader accepts an empty `pipeline:` the same as no pipeline at all,
// so the migration has to keep those files working.
func TestMigrateApplicationConfigEmptyPipeline(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "app.pipecd.yaml")
	body := "apiVersion: pipecd.dev/v1beta1\nkind: KubernetesApp\nspec:\n  name: x\n  pipeline:\n"
	if err := os.WriteFile(f, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	c := &applicationConfig{}
	if err := c.migrateApplicationConfig(context.Background(), f, zap.NewNop()); err != nil {
		t.Fatalf("expected an empty pipeline to migrate, got %v", err)
	}
	out, err := os.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "pipeline") {
		t.Errorf("expected no pipeline in the migrated config, got:\n%s", out)
	}
}
