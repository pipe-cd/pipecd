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
	"fmt"

	unit "github.com/pipe-cd/piped-plugin-sdk-go/unit"
)

// AnalysisStageOptions contains all configurable values for a ANALYSIS stage.
type AnalysisStageOptions struct {
	// How long the analysis process should be executed.
	Duration unit.Duration `json:"duration,omitempty"`
	// TODO: Consider about how to handle a pod restart
	// possible count of pod restarting
	RestartThreshold int                          `json:"restartThreshold,omitempty"`
	Metrics          []TemplatableAnalysisMetrics `json:"metrics,omitempty"`
	Logs             []TemplatableAnalysisLog     `json:"logs,omitempty"`
	HTTPS            []TemplatableAnalysisHTTP    `json:"https,omitempty"`
}

func (a *AnalysisStageOptions) Validate() error {
	if a.Duration == 0 {
		return fmt.Errorf("the ANALYSIS stage requires duration field")
	}

	for _, m := range a.Metrics {
		if m.Template.Name != "" {
			if err := m.Template.Validate(); err != nil {
				return fmt.Errorf("one of metrics configurations of ANALYSIS stage is invalid: %w", err)
			}
			continue
		}
		if err := m.Validate(); err != nil {
			return fmt.Errorf("one of metrics configurations of ANALYSIS stage is invalid: %w", err)
		}
	}

	for _, l := range a.Logs {
		if l.Template.Name != "" {
			if err := l.Template.Validate(); err != nil {
				return fmt.Errorf("one of log configurations of ANALYSIS stage is invalid: %w", err)
			}
			continue
		}
		if err := l.Validate(); err != nil {
			return fmt.Errorf("one of log configurations of ANALYSIS stage is invalid: %w", err)
		}
	}
	for _, h := range a.HTTPS {
		if h.Template.Name != "" {
			if err := h.Template.Validate(); err != nil {
				return fmt.Errorf("one of http configurations of ANALYSIS stage is invalid: %w", err)
			}
			continue
		}
		if err := h.Validate(); err != nil {
			return fmt.Errorf("one of http configurations of ANALYSIS stage is invalid: %w", err)
		}
	}
	return nil
}

type AnalysisTemplateRef struct {
	Name    string            `json:"name"`
	AppArgs map[string]string `json:"appArgs"`
}

func (a *AnalysisTemplateRef) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("the reference of analysis template name is empty")
	}
	return nil
}

// TemplatableAnalysisMetrics wraps AnalysisMetrics to allow specify template to use.
type TemplatableAnalysisMetrics struct {
	AnalysisMetrics
	Template AnalysisTemplateRef `json:"template"`
}

// TemplatableAnalysisLog wraps AnalysisLog to allow specify template to use.
type TemplatableAnalysisLog struct {
	AnalysisLog
	Template AnalysisTemplateRef `json:"template"`
}

// TemplatableAnalysisHTTP wraps AnalysisHTTP to allow specify template to use.
type TemplatableAnalysisHTTP struct {
	AnalysisHTTP
	Template AnalysisTemplateRef `json:"template"`
}
