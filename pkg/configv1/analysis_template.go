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

package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Mirror of types moved under plugin analysis config to maintain configv1 API.
// Keep these in sync with pkg/app/pipedv1/plugin/analysis/config.

const (
	AnalysisStrategyThreshold      = "THRESHOLD"
	AnalysisStrategyPrevious       = "PREVIOUS"
	AnalysisStrategyCanaryBaseline = "CANARY_BASELINE"
	AnalysisStrategyCanaryPrimary  = "CANARY_PRIMARY"

	AnalysisDeviationEither = "EITHER"
	AnalysisDeviationHigh   = "HIGH"
	AnalysisDeviationLow    = "LOW"
)

type AnalysisMetrics struct {
	Strategy     string            `json:"strategy" default:"THRESHOLD"`
	Provider     string            `json:"provider"`
	Query        string            `json:"query"`
	Expected     AnalysisExpected  `json:"expected"`
	Interval     Duration          `json:"interval"`
	FailureLimit int               `json:"failureLimit"`
	SkipOnNoData bool              `json:"skipOnNoData"`
	Timeout      Duration          `json:"timeout" default:"30s"`
	Deviation    string            `json:"deviation" default:"EITHER"`
	CanaryArgs   map[string]string `json:"canaryArgs"`
	BaselineArgs map[string]string `json:"baselineArgs"`
	PrimaryArgs  map[string]string `json:"primaryArgs"`
}

func (m *AnalysisMetrics) Validate() error {
	if m.Provider == "" {
		return fmt.Errorf("missing \"provider\" field")
	}
	if m.Query == "" {
		return fmt.Errorf("missing \"query\" field")
	}
	if m.Interval == 0 {
		return fmt.Errorf("missing \"interval\" field")
	}
	if m.Deviation != AnalysisDeviationEither && m.Deviation != AnalysisDeviationHigh && m.Deviation != AnalysisDeviationLow {
		return fmt.Errorf("\"deviation\" have to be one of %s, %s or %s", AnalysisDeviationEither, AnalysisDeviationHigh, AnalysisDeviationLow)
	}
	return nil
}

type AnalysisExpected struct {
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`
}

func (e *AnalysisExpected) Validate() error {
	if e.Min == nil && e.Max == nil {
		return fmt.Errorf("expected range is undefined")
	}
	return nil
}

type AnalysisLog struct {
	Query        string   `json:"query"`
	Interval     Duration `json:"interval"`
	FailureLimit int      `json:"failureLimit"`
	SkipOnNoData bool     `json:"skipOnNoData"`
	Timeout      Duration `json:"timeout"`
	Provider     string   `json:"provider"`
}

func (a *AnalysisLog) Validate() error { return nil }

type AnalysisHTTP struct {
	URL              string               `json:"url"`
	Method           string               `json:"method"`
	Headers          []AnalysisHTTPHeader `json:"headers"`
	ExpectedCode     int                  `json:"expectedCode"`
	ExpectedResponse string               `json:"expectedResponse"`
	Interval         Duration             `json:"interval"`
	FailureLimit     int                  `json:"failureLimit"`
	SkipOnNoData     bool                 `json:"skipOnNoData"`
	Timeout          Duration             `json:"timeout"`
}

func (a *AnalysisHTTP) Validate() error { return nil }

type AnalysisHTTPHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

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
		cfg, err := LoadFromYAML[*AnalysisTemplateSpec](path)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file %s: %w", path, err)
		}
		if cfg.Kind == KindAnalysisTemplate {
			return cfg.Spec, nil
		}
	}
	return nil, ErrNotFound
}

func (s *AnalysisTemplateSpec) Validate() error {
	return nil
}
