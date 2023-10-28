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

package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type AnalysisTemplateSpec struct {
	Metrics map[string]AnalysisMetrics `json:"metrics"`
	Logs    map[string]AnalysisLog     `json:"logs"`
	HTTPS   map[string]AnalysisHTTP    `json:"https"`
}

// LoadAnalysisTemplate finds the config file for the analysis template in the .pipe
// directory first up. And returns parsed config, ErrNotFound is returned if not found.
func LoadAnalysisTemplate(repoRoot string) (*AnalysisTemplateSpec, error) {
	dir := filepath.Join(repoRoot, SharedConfigurationDirName)
	files, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", dir, err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name())
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			continue
		}
		path := filepath.Join(dir, f.Name())
		cfg, err := LoadFromYAML(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file %s: %w", path, err)
		}
		if cfg.Kind == KindAnalysisTemplate {
			return cfg.AnalysisTemplateSpec, nil
		}
	}
	return nil, ErrNotFound
}

func (s *AnalysisTemplateSpec) Validate() error {
	return nil
}
