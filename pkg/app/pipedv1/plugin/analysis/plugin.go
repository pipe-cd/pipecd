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

package main

import (
	"time"

)

const (
	stageAnalysis = "ANALYSIS"
)

type analysisPlugin struct{}

// AnalysisStageOptions represents the configuration for ANALYSIS stage
type AnalysisStageOptions struct {
	Timeout  int                      `json:"timeout"`
	Duration time.Duration            `json:"duration"`
	Metrics  []map[string]interface{} `json:"metrics,omitempty"`
	HTTP     []map[string]interface{} `json:"http,omitempty"`
	Logs     []map[string]interface{} `json:"logs,omitempty"`
}

// AnalysisState represents the current state of the analysis
type AnalysisState struct {
	StartTime     time.Time `json:"startTime"`
	LastRunTime   time.Time `json:"lastRunTime"`
	Completed     bool      `json:"completed"`
	Success       bool      `json:"success"`
	FailureReason string    `json:"failureReason,omitempty"`
}

