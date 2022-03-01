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

package filedb

import (
	"github.com/pipe-cd/pipecd/pkg/datastore"
)

// TODO: Implement filterable interface for each collection.
type filterable interface {
	Match(e interface{}, filters []datastore.ListFilter) (bool, error)
}

func filter(col datastore.Collection, e interface{}, filters []datastore.ListFilter) (bool, error) {
	fcol, ok := col.(filterable)
	if !ok {
		return false, datastore.ErrUnsupported
	}

	return fcol.Match(e, filters)
}
