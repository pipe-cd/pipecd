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

// Package firestoreindexcreator automatically creates the needed composite indexes for Google Firestore,
// based on well-defined JSON format indexes list with the name indexes.json.
package firestoreindexcreator

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type IndexCreator interface {
	CreateIndexes(ctx context.Context) error
}

type firestoreClient interface {
	authorize(ctx context.Context) error
	createIndex(ctx context.Context, idx *index)
}

type indexCreator struct {
	firestoreClient
	logger *zap.Logger
}

func NewIndexCreator(gcloudPath, projectID, serviceAcccountFile string, logger *zap.Logger) IndexCreator {
	return &indexCreator{
		// TODO: Use Go SDK to create Firebase indexes upon it will be started providing
		firestoreClient: newGcloud(gcloudPath, projectID, serviceAcccountFile, logger),
		logger:          logger.Named("firestore-index-creator"),
	}
}

// CreateIndexes creates needed composite indexes for Google Firestore based on
// well-defined indexes list.
func (c *indexCreator) CreateIndexes(ctx context.Context) error {
	if err := c.authorize(ctx); err != nil {
		return fmt.Errorf("failed to authorize: %w", err)
	}

	idxs, err := parseIndexes()
	if err != nil {
		return err
	}
	for i := 0; i < len(idxs.Indexes); i++ {
		// TODO: Check if the index is already added and if so skip creating
		c.createIndex(ctx, &idxs.Indexes[i])
	}
	return nil
}
