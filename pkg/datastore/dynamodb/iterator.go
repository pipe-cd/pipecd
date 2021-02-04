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

package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

// Iterator is wrapper of queried data pool
type Iterator struct {
	datapool []map[string]*dynamodb.AttributeValue
}

// Next implementation of DynamoDB Iterator
func (it *Iterator) Next(dst interface{}) error {
	// If the data pool is empty, it means all items is popped
	if len(it.datapool) == 0 {
		return datastore.ErrIteratorDone
	}
	// Pop the first item from data pool
	item := it.datapool[0]
	it.datapool = it.datapool[1:]

	return dynamodbattribute.UnmarshalMap(item, dst)
}

// Cursor implementation of DynamoDB Iterator
func (it *Iterator) Cursor() (string, error) {
	// Note: Perhaps, not required.
	return "", datastore.ErrUnimplemented
}
