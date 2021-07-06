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

package rediscache

import (
	redigo "github.com/gomodule/redigo/redis"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/redis"
)

type RedisHashCache struct {
	RedisCache
	key string
}

func NewRedisHashCache(redis redis.Redis, key string) *RedisHashCache {
	return &RedisHashCache{
		RedisCache: RedisCache{
			redis: redis,
		},
		key: key,
	}
}

func (r *RedisHashCache) Put(k interface{}, v interface{}) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", r.key, k, v)
	return err
}

func (r *RedisHashCache) GetAll() ([]interface{}, error) {
	conn := r.redis.Get()
	defer conn.Close()
	reply, err := redigo.StringMap(conn.Do("HGETALL", r.key))
	if err != nil {
		if err == redigo.ErrNil {
			return nil, cache.ErrNotFound
		}
		return nil, err
	}
	if reply == nil {
		return nil, cache.ErrNotFound
	}

	out := make([]interface{}, len(reply))
	for _, v := range reply {
		out = append(out, v)
	}
	return out, nil
}
