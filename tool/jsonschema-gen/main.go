// Copyright 2024 The PipeCD Authors.
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

// jsonschema-gen generates JSON Schema files from PipeCD application
// configuration structs. The generated schemas can be used by editors
// and language servers to provide validation and autocompletion for
// PipeCD YAML configuration files.
//
// Usage:
//
//	go run ./tool/jsonschema-gen -out ./docs/static/jsonschema
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/invopop/jsonschema"

	"github.com/pipe-cd/pipecd/pkg/config"
)

var outDir = flag.String("out", "docs/static/jsonschema", "output directory for generated schema files")

type schemaTarget struct {
	name     string
	instance any
}

func main() {
	flag.Parse()

	targets := []schemaTarget{
		{"ecs", &config.ECSApplicationSpec{}},
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output directory: %v\n", err)
		os.Exit(1)
	}

	r := &jsonschema.Reflector{
		KeyNamer: func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToLower(s[:1]) + s[1:]
		},
		RequiredFromJSONSchemaTags: true,
	}

	for _, t := range targets {
		schema := r.Reflect(t.instance)
		data, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error marshaling schema for %s: %v\n", t.name, err)
			os.Exit(1)
		}

		path := filepath.Join(*outDir, t.name+".json")
		if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("generated %s\n", path)
	}
}
