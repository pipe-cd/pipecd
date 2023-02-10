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

package rediscache

import (
	"time"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachemetrics"
	"github.com/pipe-cd/pipecd/pkg/redis"
)

type RedisCache struct {
	redis redis.Redis
	ttl   uint
}

func NewCache(redis redis.Redis) *RedisCache {
	return &RedisCache{
		redis: redis,
	}
}

func NewTTLCache(redis redis.Redis, ttl time.Duration) *RedisCache {
	return &RedisCache{
		redis: redis,
		ttl:   uint(ttl.Seconds()),
	}
}

func (c *RedisCache) Get(k string) (interface{}, error) {
	conn := c.redis.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", k)
	if err != nil {
		if err == redigo.ErrNil {
			cachemetrics.IncGetOperationCounter(
				cachemetrics.LabelSourceRedis,
				cachemetrics.LabelStatusMiss,
			)
			return nil, cache.ErrNotFound
		}
		return nil, err
	}
	if reply == nil {
		cachemetrics.IncGetOperationCounter(
			cachemetrics.LabelSourceRedis,
			cachemetrics.LabelStatusMiss,
		)
		return nil, cache.ErrNotFound
	}
	if err, ok := reply.(redigo.Error); ok {
		cachemetrics.IncGetOperationCounter(
			cachemetrics.LabelSourceRedis,
			cachemetrics.LabelStatusMiss,
		)
		return nil, err
	}
	cachemetrics.IncGetOperationCounter(
		cachemetrics.LabelSourceRedis,
		cachemetrics.LabelStatusHit,
	)
	return reply, nil
}

// It is caller's responsibility to encode Go struct.
func (c *RedisCache) Put(k string, v interface{}) error {
	conn := c.redis.Get()
	defer conn.Close()
	var err error
	if c.ttl == 0 {
		_, err = conn.Do("SET", k, v)
	} else {
		_, err = conn.Do("SETEX", k, c.ttl, v)
	}
	return err
}

func (c *RedisCache) Delete(k string) error {
	conn := c.redis.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", k)
	return err
}

func (c *RedisCache) GetAll() (map[string]interface{}, error) {
	return nil, cache.ErrUnimplemented
}
