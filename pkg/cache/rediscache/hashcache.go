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
	"time"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/redis"
)

type RedisHashCache struct {
	redis redis.Redis
	ttl   uint
	key   string
}

func NewHashCache(redis redis.Redis, key string) *RedisHashCache {
	return &RedisHashCache{
		redis: redis,
		key:   key,
	}
}

func NewTTLHashCache(redis redis.Redis, ttl time.Duration, key string) *RedisHashCache {
	return &RedisHashCache{
		redis: redis,
		ttl:   uint(ttl.Seconds()),
		key:   key,
	}
}

func (r *RedisHashCache) Get(k interface{}) (interface{}, error) {
	conn := r.redis.Get()
	defer conn.Close()
	reply, err := conn.Do("HGET", r.key, k)
	if err != nil {
		if err == redigo.ErrNil {
			return nil, cache.ErrNotFound
		}
		return nil, err
	}
	if reply == nil {
		return nil, cache.ErrNotFound
	}
	if err, ok := reply.(redigo.Error); ok {
		return nil, err
	}
	return reply, nil
}

func (r *RedisHashCache) Put(k interface{}, v interface{}) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", r.key, k, v)
	if r.ttl != 0 {
		_, err = conn.Do("EXPIRE", r.key, r.ttl)
	}
	return err
}

func (r *RedisHashCache) Delete(k interface{}) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", r.key, k)
	return err
}

func (r *RedisHashCache) GetAll() (map[interface{}]interface{}, error) {
	conn := r.redis.Get()
	defer conn.Close()
	reply, err := redigo.StringMap(conn.Do("HGETALL", r.key))
	if err != nil {
		if err == redigo.ErrNil {
			return nil, cache.ErrNotFound
		}
		return nil, err
	}
	if len(reply) == 0 {
		return nil, cache.ErrNotFound
	}

	out := make(map[interface{}]interface{}, len(reply))
	for k, v := range reply {
		out[k] = v
	}

	return out, nil
}
