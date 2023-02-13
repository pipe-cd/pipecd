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
	listIndexes(ctx context.Context) ([]index, error)
}

type indexEnsurer struct {
	firestoreClient
	collectionNamePrefix string
	logger               *zap.Logger
}

func NewIndexEnsurer(gcloudPath, projectID, serviceAcccountFile, colPrefix string, logger *zap.Logger) IndexEnsurer {
	return &indexEnsurer{
		// TODO: Use Go SDK to create Firebase indexes upon it will be started providing
		firestoreClient:      newGcloud(gcloudPath, projectID, serviceAcccountFile, logger),
		collectionNamePrefix: colPrefix,
		logger:               logger.Named("firestore-index-ensurer"),
	}
}

// CreateIndexes creates needed composite indexes for Google Firestore based on
// well-defined indexes list.
func (e *indexEnsurer) CreateIndexes(ctx context.Context) error {
	e.logger.Info("start ensuring the existence of composite indexes for Google Cloud Firestore")
	if err := e.authorize(ctx); err != nil {
		return fmt.Errorf("failed to authorize: %w", err)
	}

	indexes, err := parseIndexes()
	if err != nil {
		return err
	}

	if p := e.collectionNamePrefix; p != "" {
		prefixIndexes(indexes, p)
	}

	exists, err := e.listIndexes(ctx)
	if err != nil {
		return err
	}

	filtered := filterIndexes(indexes, exists)
	if len(filtered) == 0 {
		return nil
	}

	e.logger.Info(fmt.Sprintf("%d missing Firebase composite indexes found", len(filtered)))
	for i := 0; i < len(filtered); i++ {
		if err := e.createIndex(ctx, &filtered[i]); err != nil {
			e.logger.Error("failed to create a Firestore composite index",
				zap.Any("index", filtered[i]),
				zap.Error(err),
			)
		}
	}
	return nil
}
