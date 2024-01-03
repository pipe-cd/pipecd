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

package pipedstatsbuilder

import (
	"bytes"
	"errors"
	"io"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/model"
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

func (b *PipedStatsBuilder) Build() (io.Reader, error) {
	res, err := b.backend.GetAll()
	if err != nil {
		// Only show error in case it's not cache not found error.
		if !errors.Is(err, cache.ErrNotFound) {
			b.logger.Error("failed to fetch piped stats from cache", zap.Error(err))
			return nil, err
		}
		return bytes.NewReader([]byte("")), nil
	}
	data := make([][]byte, 0, len(res))
	for _, v := range res {
		ps := model.PipedStat{}
		if err = model.UnmarshalPipedStat(v, &ps); err != nil {
			b.logger.Error("failed to unmarshal piped stat data", zap.Error(err))
			return nil, err
		}

		// Ignore piped stat metrics if passed time from its last committed
		// timestamp longer than limit live state check.
		if ps.IsStaled(model.PipedStatsRetention) {
			continue
		}
		data = append(data, ps.Metrics)
	}
	return bytes.NewReader(bytes.Join(data, []byte("\n"))), nil
}
