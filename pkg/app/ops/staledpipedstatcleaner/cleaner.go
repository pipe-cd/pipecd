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
package staledpipedstatcleaner

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/model"
)

var (
	pipedStatStaledTimeout = 24 * time.Hour
	interval               = 24 * time.Hour
)

type StaledPipedStatCleaner struct {
	backend cache.Cache
	logger  *zap.Logger
}

func NewStaledPipedStatCleaner(c cache.Cache, logger *zap.Logger) *StaledPipedStatCleaner {
	return &StaledPipedStatCleaner{
		backend: c,
		logger:  logger.Named("staled-piped-stat-cleaner"),
	}
}

func (s *StaledPipedStatCleaner) Run(ctx context.Context) error {
	s.logger.Info("start running StaledPipedStatCleaner")

	t := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("staledPipedStatCleaner has been stopped")
			return nil

		case <-t.C:
			start := time.Now()
			if err := s.flushStaledPipedsStat(ctx); err == nil {
				s.logger.Info("successfully cleaned staled pipeds stat", zap.Duration("duration", time.Since(start)))
			}
		}
	}
}

func (s *StaledPipedStatCleaner) flushStaledPipedsStat(ctx context.Context) error {
	res, err := s.backend.GetAll()
	if err != nil {
		// Ignore cache not found error since there are no stats found in cache
		// means no need to flush anything.
		if !errors.Is(err, cache.ErrNotFound) {
			s.logger.Error("failed to fetch piped stats from cache", zap.Error(err))
			return err
		}
		return nil
	}

	staled := make([]string, 0)
	for k, v := range res {
		value, okValue := v.([]byte)
		if !okValue {
			err = errors.New("error value not a bulk of string value")
			s.logger.Error("failed to unmarshal piped stat data", zap.Error(err))
			return err
		}
		ps := model.PipedStat{}
		if err = json.Unmarshal(value, &ps); err != nil {
			s.logger.Error("failed to unmarshal piped stat data", zap.Error(err))
			return err
		}
		if time.Since(time.Unix(ps.Timestamp, 0)) > pipedStatStaledTimeout {
			staled = append(staled, k)
		}
	}

	// No staled pipeds' stat found.
	if len(staled) == 0 {
		return nil
	}

	for _, id := range staled {
		if err = s.backend.Delete(id); err != nil {
			s.logger.Error("failed to remove staled piped stat data", zap.String("pipedID", id), zap.Error(err))
			return err
		}
	}

	return nil
}
