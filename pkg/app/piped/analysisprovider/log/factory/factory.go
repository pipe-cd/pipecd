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

package factory

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/analysisprovider/log"
	"github.com/pipe-cd/pipecd/pkg/app/piped/analysisprovider/log/stackdriver"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// NewProvider generates an appropriate provider according to analysis provider config.
func NewProvider(providerCfg *config.PipedAnalysisProvider, logger *zap.Logger) (provider log.Provider, err error) {
	switch providerCfg.Type {
	case model.AnalysisProviderStackdriver:
		cfg := providerCfg.StackdriverConfig
		sa, err := os.ReadFile(cfg.ServiceAccountFile)
		if err != nil {
			return nil, err
		}
		provider, err = stackdriver.NewProvider(sa)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("any of providers config not found")
	}
	return provider, nil
}
