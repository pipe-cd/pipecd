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

package rediscache

import (
	"errors"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/redis"
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

func (r *RedisHashCache) Get(k string) (interface{}, error) {
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

// Put implementation for Putter interface.
// For TTLHashCache, Put acts mostly as normal Cache except it will
// check if the TTL for hashkey (not the key k), in case the TTL for
// hashkey is not yet existed or unset, EXPIRE will be called and set
// TTL time for the whole hashkey.
//
// It is caller's responsibility to encode Go struct.
func (r *RedisHashCache) Put(k string, v interface{}) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", r.key, k, v)
	if err != nil {
		return err
	}

	// Skip set TTL if unnecessary.
	if r.ttl == 0 {
		return nil
	}

	rep, err := redigo.Int(conn.Do("TTL", r.key))
	if err != nil {
		return err
	}
	// Only set TTL for hashkey in case TTL command return key has no TTL.
	// ref: https://redis.io/commands/TTL
	if rep < 0 {
		_, err = conn.Do("EXPIRE", r.key, r.ttl)
	}
	return err
}

func (r *RedisHashCache) Delete(k string) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", r.key, k)
	return err
}

func (r *RedisHashCache) GetAll() (map[string]interface{}, error) {
	conn := r.redis.Get()
	defer conn.Close()
	reply, err := redigo.Values(conn.Do("HGETALL", r.key))
	if err != nil {
		if err == redigo.ErrNil {
			return nil, cache.ErrNotFound
		}
		return nil, err
	}
	if len(reply) == 0 {
		return nil, cache.ErrNotFound
	}
	if len(reply)%2 != 0 {
		return nil, errors.New("invalid key-value pair contained")
	}

	out := make(map[string]interface{}, len(reply)/2)
	for i := 0; i < len(reply); i += 2 {
		key, okKey := reply[i].([]byte)
		if !okKey {
			return nil, errors.New("error key not a bulk string value")
		}
		out[string(key)] = reply[i+1]
	}

	return out, nil
}
