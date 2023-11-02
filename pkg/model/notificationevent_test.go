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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		typ  NotificationEventType
		want NotificationEventGroup
	}{
		{
			name: "returns EVENT_DEPLOYMENT for type < 100",
			typ:  99,
			want: NotificationEventGroup_EVENT_DEPLOYMENT,
		},
		{
			name: "returns EVENT_APPLICATION_SYNC for type < 200",
			typ:  199,
			want: NotificationEventGroup_EVENT_APPLICATION_SYNC,
		},
		{
			name: "returns EVENT_APPLICATION_HEALTH for type < 300",
			typ:  299,
			want: NotificationEventGroup_EVENT_APPLICATION_HEALTH,
		},
		{
			name: "returns EVENT_PIPED for type < 400",
			typ:  399,
			want: NotificationEventGroup_EVENT_PIPED,
		},
		{
			name: "returns EVENT_NONE for type >= 400",
			typ:  400,
			want: NotificationEventGroup_EVENT_NONE,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := NotificationEvent{Type: tt.typ}
			got := e.Group()
			assert.Equal(t, tt.want, got)
		})
	}
}
