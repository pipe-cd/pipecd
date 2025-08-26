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

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
	"sigs.k8s.io/yaml"
)

type AnalysisTemplateConfig struct {
	Kind       string
	APIVersion string
	Spec       AnalysisTemplateSpec
}

type AnalysisTemplateSpec struct {
	Metrics map[string]AnalysisMetrics `json:"metrics"`
	Logs    map[string]AnalysisLog     `json:"logs"`
	HTTPS   map[string]AnalysisHTTP    `json:"https"`
}

const (
	// KindAnalysisTemplate represents shared analysis template for a repository.
	// This configuration file should be placed in .pipe directory
	// at the root of the repository.
	KindAnalysisTemplate string = "AnalysisTemplate"
)

var errNotFound = errors.New("not found")

// LoadAnalysisTemplate finds the config file for the analysis template in the .pipe
// directory first up. And returns parsed config, ErrNotFound is returned if not found.
func LoadAnalysisTemplate(sharedConfigDir string) (*AnalysisTemplateSpec, error) {
	files, err := os.ReadDir(sharedConfigDir)
	if os.IsNotExist(err) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", sharedConfigDir, err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name())
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			continue
		}
		path := filepath.Join(sharedConfigDir, f.Name())
		cfg, err := loadFromYAML(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file %s: %w", path, err)
		}
		if cfg.Kind == KindAnalysisTemplate {
			return &cfg.Spec, nil
		}
	}
	return nil, errNotFound
}

func (s *AnalysisTemplateSpec) Validate() error {
	return nil
}

// loadFromYAML reads and decodes a yaml file to construct the Config.
func loadFromYAML(file string) (*AnalysisTemplateConfig, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return decodeYAML(data)
}

// decodeYAML unmarshals config YAML data to config struct.
// It also validates the configuration after decoding.
func decodeYAML(data []byte) (*AnalysisTemplateConfig, error) {
	js, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}
	c := &AnalysisTemplateConfig{}
	if err := json.Unmarshal(js, c); err != nil {
		return nil, err
	}
	if err := defaults.Set(c); err != nil {
		return nil, err
	}
	if err := c.Spec.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}
