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
	"context"
	"testing"

	"github.com/pipe-cd/pipe/pkg/datastore"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func makeTestableDynamoDBStore() *DynamoDB {
	options := []Option{
		WithCredentialsFile("default", "/Users/s12228/.aws/credentials"),
	}
	s, _ := NewDynamoDB("ap-northeast-1", "https://dynamodb.ap-northeast-1.amazonaws.com", options...)
	return s
}

func TestNewDynamoDBClient(t *testing.T) {
	s := makeTestableDynamoDBStore()
	assert.NotEqual(t, nil, s.client)

	input := &dynamodb.ListTablesInput{}
	_, err := s.client.ListTables(input)
	assert.Equal(t, nil, err)
}

type sample struct {
	ProjectId string `json:"projectId"`
	Id        string `json:"id"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func _TestPaginationOnFindQuery(t *testing.T) {
	s := makeTestableDynamoDBStore()
	filters := []datastore.ListFilter{
		{
			Field:    "projectId",
			Operator: "==",
			Value:    "project-1",
		},
	}
	iter, err := s.Find(context.TODO(), "demo-application", datastore.ListOptions{
		Filters:  filters,
		PageSize: 2,
	})
	assert.Equal(t, nil, err)
	results := make([]*sample, 0)
	for {
		var d sample
		err := iter.Next(&d)
		if err == datastore.ErrIteratorDone {
			break
		}
		results = append(results, &d)
	}
	assert.Equal(t, 2, len(results))
}

func _TestPutItem(t *testing.T) {
	s := makeTestableDynamoDBStore()
	data := &sample{
		ProjectId: "project-2",
		Id:        "app-3",
		CreatedAt: 105,
		UpdatedAt: 1005,
	}
	err := s.Put(context.TODO(), "demo-application", data.Id, data)
	assert.Equal(t, nil, err)
}

func _TestCreateItem(t *testing.T) {
	s := makeTestableDynamoDBStore()
	data := &sample{
		ProjectId: "project-2",
		Id:        "app-3",
		CreatedAt: 107,
		UpdatedAt: 107,
	}
	err := s.Create(context.TODO(), "demo-application", data.Id, data)
	assert.NotEqual(t, nil, err)
}

var sampleFactory = func() interface{} {
	return &sample{}
}

var sampleUpdater = func(s *sample) error {
	s.UpdatedAt++
	return nil
}

func _TestUpdateItem(t *testing.T) {
	s := makeTestableDynamoDBStore()
	err := s.Update(context.TODO(), "demo-application", "app-3", sampleFactory, func(e interface{}) error {
		d := e.(*sample)
		if err := sampleUpdater(d); err != nil {
			return err
		}
		return nil
	})
	assert.Equal(t, nil, err)
}

func _TestGetItem(t *testing.T) {
	s := makeTestableDynamoDBStore()
	var out sample
	err := s.Get(context.TODO(), "demo-application", "app-3", &out)
	assert.Equal(t, nil, err)
}
