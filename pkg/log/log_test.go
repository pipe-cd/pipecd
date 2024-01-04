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

package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerOK(t *testing.T) {
	validLevels := []string{
		"debug",
		"info",
		"warn",
		"error",
		"dpanic",
		"panic",
		"fatal",
	}
	validEncodings := []EncodingType{
		JSONEncoding,
		ConsoleEncoding,
		HumanizeEncoding,
	}
	for _, level := range validLevels {
		for _, encoding := range validEncodings {
			cfg := Configs{
				Level:    level,
				Encoding: encoding,
				ServiceContext: &ServiceContext{
					Service: "test-service",
					Version: "1.0.0",
				},
			}
			logger, err := NewLogger(cfg)
			des := fmt.Sprintf("level: %s, encoding: %s", level, encoding)
			assert.Nil(t, err, des)
			assert.NotNil(t, logger, des)
		}
	}
}

func TestNewLoggerFailed(t *testing.T) {
	configs := []Configs{
		Configs{
			Level:    "foo",
			Encoding: "json",
		},
		Configs{
			Level:    "info",
			Encoding: "foo",
		},
	}
	for _, cfg := range configs {
		logger, err := NewLogger(cfg)
		des := fmt.Sprintf("level: %s, encoding: %s", cfg.Level, cfg.Encoding)
		assert.NotNil(t, err, des)
		assert.Nil(t, logger, des)
	}
}
