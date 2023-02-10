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

package ensurer

import (
	"context"
)

type IndexEnsurer interface {
	// EnsureIndexes loads indexes defined sql file and applies it to the database.
	// In case of indexes already existed, no errors will be returned.
	EnsureIndexes(ctx context.Context) error
}

type SchemaEnsurer interface {
	// EnsureSchema loads schema defined sql file and applies it to the database.
	EnsureSchema(ctx context.Context) error
}

type SQLEnsurer interface {
	IndexEnsurer
	SchemaEnsurer
	// Close closes database connection held by client.
	Close() error
	Ping() error
}
