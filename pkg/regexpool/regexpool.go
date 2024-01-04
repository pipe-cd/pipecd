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

package regexpool

import (
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/pipe-cd/pipecd/pkg/cache/memorycache"
)

var (
	defaultSize = 1000
	defaultPool *Pool
)

// Pool is the representation of a pool of regular expression.
type Pool struct {
	cache *memorycache.LRUCache
	fails map[string]struct{}
	mu    sync.RWMutex
}

func init() {
	var err error
	if defaultPool, err = NewPool(defaultSize); err != nil {
		log.Fatal(err)
	}
}

// DefaultPool returns the default pool which was initialized at init by default.
func DefaultPool() *Pool {
	return defaultPool
}

// NewPool initializes a new pool with a specified size.
func NewPool(size int) (*Pool, error) {
	cache, err := memorycache.NewLRUCache(size)
	if err != nil {
		return nil, err
	}
	return &Pool{
		fails: make(map[string]struct{}),
		cache: cache,
	}, nil
}

// Get retrieves the Regexp object from the pool or
// creates a new one if it does not exist.
func (p *Pool) Get(expr string) (*regexp.Regexp, error) {
	regex, err := p.cache.Get(expr)
	if err == nil {
		return regex.(*regexp.Regexp), nil
	}
	// Check if the given expression was unable to compile before
	// then return the error fast.
	p.mu.RLock()
	_, ok := p.fails[expr]
	p.mu.RUnlock()
	if ok {
		return nil, fmt.Errorf("unable to compile: %s", expr)
	}
	// Compile the expression string and cache its result.
	reg, err := regexp.Compile(expr)
	if err == nil {
		p.cache.Put(expr, reg)
		return reg, nil
	}
	p.mu.Lock()
	p.fails[expr] = struct{}{}
	p.mu.Unlock()
	return nil, fmt.Errorf("unable to compile: %s", expr)
}
