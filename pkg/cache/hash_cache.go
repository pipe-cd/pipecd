// Copyright 2022 The PipeCD Authors.
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

// HGetter wraps a method to read from hash cache.
type HGetter interface {
	HGet(hash, key string) (interface{}, error)
	HGetAll(hash string) (map[string]interface{}, error)
}

// HPutter wraps a method to write to hash cache.
type HPutter interface {
	HPut(hash, key string, value interface{}) error
}

// HDeleter wraps a method to delete from hash cache.
type HDeleter interface {
	HDelete(hahs, key string) error
}

// HCache groups HGetter, HPutter and HDeleter.
type HCache interface {
	HGetter
	HPutter
	HDeleter
}
