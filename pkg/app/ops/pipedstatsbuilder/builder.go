// Copyright 2021 The PipeCD Authors.
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

package pipedstatsbuilder

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/cache"
)

type PipedStatsBuilder struct {
	backend cache.Cache
	logger  *zap.Logger
}

func NewPipedStatsBuilder(c cache.Cache, logger *zap.Logger) *PipedStatsBuilder {
	return &PipedStatsBuilder{
		backend: c,
		logger:  logger.Named("piped-metrics-builder"),
	}
}

func (b *PipedStatsBuilder) Build() (io.ReadCloser, error) {
	res, err := b.backend.GetAll()
	if err != nil {
		b.logger.Error("failed to fetch piped stats from cache", zap.Error(err))
		return nil, err
	}
	data := make([][]byte, 0, len(res))
	for _, v := range res {
		value, okValue := v.([]byte)
		if !okValue {
			err = errors.New("error value not a bulk of string value")
			b.logger.Error("failed to marshal piped stat data", zap.Error(err))
			return nil, err
		}
		data = append(data, value)
	}
	return ioutil.NopCloser(bytes.NewReader(bytes.Join(data, []byte("\n")))), nil
}
