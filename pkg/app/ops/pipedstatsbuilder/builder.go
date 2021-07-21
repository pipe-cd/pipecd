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
	"encoding/json"
	"errors"
	"io"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	liveStateExceedDuration = 2 * time.Minute
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
		b.logger.Error("failed to fetch piped stats from cache", zap.Error(err))
		return nil, err
	}
	data := make([][]byte, 0, len(res))
	for _, v := range res {
		value, okValue := v.([]byte)
		if !okValue {
			err = errors.New("error value not a bulk of string value")
			b.logger.Error("failed to unmarshal piped stat data", zap.Error(err))
			return nil, err
		}
		ps := model.PipedStat{}
		if err = json.Unmarshal(value, &ps); err != nil {
			b.logger.Error("failed to unmarshal piped stat data", zap.Error(err))
			return nil, err
		}
		// Ignore piped stat metrics if passed time from its last committed
		// timestamp longer than limit live state check.
		if time.Since(time.Unix(ps.Timestamp, 0)) > liveStateExceedDuration {
			continue
		}
		data = append(data, ps.Metrics)
	}
	return bytes.NewReader(bytes.Join(data, []byte("\n"))), nil
}
