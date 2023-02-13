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

package datastore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestAddEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := model.Event{
		Id:        "id",
		Name:      "name",
		Data:      "data",
		ProjectId: "project",
		EventKey:  "82a3537ff0dbce7eec35d69edc3a189ee6f17d82f353a553f9aa96cb0be3ce89",
		CreatedAt: 12345,
		UpdatedAt: 12345,
	}

	testcases := []struct {
		name    string
		event   model.Event
		ds      DataStore
		wantErr bool
	}{
		{
			name:  "Invalid event",
			event: model.Event{},
			ds: func() DataStore {
				return NewMockDataStore(ctrl)
			}(),
			wantErr: true,
		},
		{
			name:  "OK",
			event: event,
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Create(gomock.Any(), gomock.Any(), event.Id, &event).
					Return(nil)
				return ds
			}(),
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewEventStore(tc.ds, TestCommander)
			err := s.Add(context.Background(), tc.event)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
