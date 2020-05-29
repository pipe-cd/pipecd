// Copyright 2020 The PipeCD Authors.
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
	"io/ioutil"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/app/piped/analysisprovider/log/stackdriver"
	"github.com/kapetaniosci/pipe/pkg/config"
)

type Factory struct {
	logger *zap.Logger
}

func NewFactory(logger *zap.Logger) *Factory {
	return &Factory{logger: logger}
}

// NewProvider generates an appropriate provider according to analysis provider config.
func (f *Factory) NewProvider(providerCfg *config.AnalysisProvider) (provider Provider, err error) {
	switch {
	case providerCfg.Stackdriver != nil:
		cfg := providerCfg.Stackdriver
		sa, err := ioutil.ReadFile(cfg.ServiceAccountFile)
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
