package migrate

import (
	"context"
	"os"
	"path/filepath"
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
