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

// Package firestoreindexensurer automatically creates/deletes the needed composite indexes
// for Google Firestore, based on well-defined JSON format indexes list.
package firestoreindexensurer

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type IndexEnsurer interface {
	CreateIndexes(ctx context.Context) error
}

type firestoreClient interface {
	authorize(ctx context.Context) error
	createIndex(ctx context.Context, idx *index) error
}

type indexEnsurer struct {
	firestoreClient
	logger *zap.Logger
}

func NewIndexEnsurer(gcloudPath, projectID, serviceAcccountFile string, logger *zap.Logger) IndexEnsurer {
	return &indexEnsurer{
		// TODO: Use Go SDK to create Firebase indexes upon it will be started providing
		firestoreClient: newGcloud(gcloudPath, projectID, serviceAcccountFile, logger),
		logger:          logger.Named("firestore-index-ensurer"),
	}
}

// CreateIndexes creates needed composite indexes for Google Firestore based on
// well-defined indexes list.
func (c *indexEnsurer) CreateIndexes(ctx context.Context) error {
	if err := c.authorize(ctx); err != nil {
		return fmt.Errorf("failed to authorize: %w", err)
	}

	idxs, err := parseIndexes()
	if err != nil {
		return err
	}
	for i := 0; i < len(idxs.Indexes); i++ {
		// TODO: Check if the index is already added and if so skip creating
		if err := c.createIndex(ctx, &idxs.Indexes[i]); err != nil {
			c.logger.Error("failed to create a Firestore composite index",
				zap.Any("index", idxs.Indexes[i]),
				zap.Error(err),
			)
		}
	}
	return nil
}
