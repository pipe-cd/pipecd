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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/redis"
)

const defaultTTL = 10 * time.Second

func redisCli() redis.Redis {
	return redis.NewRedis("localhost:6379", "")
}

func TestGetAll(t *testing.T) {
	hcache := NewTTLHashCache(redisCli(), defaultTTL, "HASHKEY")
	hcache.Put("0", []byte("abc"))
	hcache.Put("1", []byte("xyz"))
	rep, err := hcache.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rep))
	assert.Equal(t, "abc", rep["0"])
	assert.Equal(t, "xyz", rep["1"])
}

func _TestGet(t *testing.T) {
	hcache := NewTTLHashCache(redisCli(), defaultTTL, "HASHKEY")
	// hcache.Put("0", "abc")
	// hcache.Put("1", "xyz")
	rep, err := hcache.Get("0")
	assert.NoError(t, err)
	assert.Equal(t, "abc", rep)
}
