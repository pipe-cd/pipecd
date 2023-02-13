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

package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Redis interface {
	Get() redis.Conn
	Close() error
}

type options struct {
	maxIdle            int
	maxActive          int
	idleTimeout        time.Duration
	connCheckInterval  time.Duration
	dialConnectTimeout time.Duration
}

type Option func(*options)

// WithMaxIdle specifies the maximum number of idle connections in the pool.
func WithMaxIdle(num int) Option {
	return func(opts *options) {
		opts.maxIdle = num
	}
}

// WithMaxActive specifies the maximum number of connections allocated by the pool at a given time.
// When zero, there is no limit on the number of connections in the pool.
func WithMaxActive(num int) Option {
	return func(opts *options) {
		opts.maxActive = num
	}
}

// WithIdleTimeout specifies the timeout for idle connection.
func WithIdleTimeout(duration time.Duration) Option {
	return func(opts *options) {
		opts.idleTimeout = duration
	}
}

// WithConnCheckInterval specifies the minimum interval between two times of checking idle connection health.
func WithConnCheckInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.connCheckInterval = interval
	}
}

// WithDialConnectTimeout specifies the timeout for connecting to the Redis server.
func WithDialConnectTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.dialConnectTimeout = timeout
	}
}

func NewRedis(address, password string, opts ...Option) Redis {
	opt := &options{
		maxIdle:            10,
		maxActive:          0,
		idleTimeout:        4 * time.Minute,
		connCheckInterval:  time.Minute,
		dialConnectTimeout: 30 * time.Second,
	}
	for _, o := range opts {
		o(opt)
	}
	return createPool(address, password, opt)
}

func createPool(address, password string, opt *options) *redis.Pool {
	dialOptions := make([]redis.DialOption, 0, 2)
	if password != "" {
		dialOptions = append(dialOptions, redis.DialPassword(password))
	}
	if opt.dialConnectTimeout > 0 {
		dialOptions = append(dialOptions, redis.DialConnectTimeout(opt.dialConnectTimeout))
	}
	return &redis.Pool{
		MaxIdle:     opt.maxIdle,
		MaxActive:   opt.maxActive,
		IdleTimeout: opt.idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address, dialOptions...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < opt.connCheckInterval {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
