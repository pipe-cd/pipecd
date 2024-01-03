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

package cache

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrUnimplemented = errors.New("unimplemented")
)

// Getter wraps a method to read from cache.
type Getter interface {
	Get(key string) (interface{}, error)
	GetAll() (map[string]interface{}, error)
}

// Putter wraps a method to write to cache.
type Putter interface {
	Put(key string, value interface{}) error
}

// Deleter wraps a method to delete from cache.
type Deleter interface {
	Delete(key string) error
}

// Cache groups Getter, Putter and Deleter.
type Cache interface {
	Getter
	Putter
	Deleter
}

type multiGetter struct {
	getters []Getter
}

// MultiGetter combines a list of getters into a single getter.
func MultiGetter(getters ...Getter) Getter {
	all := make([]Getter, 0, len(getters))
	for _, r := range getters {
		if mg, ok := r.(*multiGetter); ok {
			all = append(all, mg.getters...)
		} else {
			all = append(all, r)
		}
	}
	return &multiGetter{
		getters: all,
	}
}

func (mg *multiGetter) Get(key string) (interface{}, error) {
	if len(mg.getters) == 0 {
		return nil, ErrNotFound
	}
	if len(mg.getters) == 1 {
		return mg.getters[0].Get(key)
	}
	var firstErr error
	for i := range mg.getters {
		e, err := mg.getters[i].Get(key)
		if firstErr == nil && err != nil {
			firstErr = err
		}
		if err == nil {
			return e, nil
		}
	}
	return nil, firstErr
}

func (mg *multiGetter) GetAll() (map[string]interface{}, error) {
	return nil, ErrUnimplemented
}
