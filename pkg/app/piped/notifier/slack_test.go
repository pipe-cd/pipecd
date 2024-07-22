// Copyright 2024 The PipeCD Authors.
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

package notifier

import "testing"

func Test_getAccountsAsString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		accounts []string
		want     string
	}{
		{
			name:     "empty",
			accounts: []string{},
			want:     "",
		},
		{
			name:     "single",
			accounts: []string{"foo"},
			want:     "<@foo>",
		},
		{
			name:     "multiple",
			accounts: []string{"foo", "bar"},
			want:     "<@foo> <@bar>",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getAccountsAsString(tt.accounts)
			if got != tt.want {
				t.Errorf("getAccountsAsString(): got %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_getGroupsAsString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		groups []string
		want   string
	}{
		{
			name:   "empty",
			groups: []string{},
			want:   "",
		},
		{
			name:   "single",
			groups: []string{"foo"},
			want:   "<!subteam^foo>",
		},
		{
			name:   "with correct format <!subteam^foo>",
			groups: []string{"<!subteam^foo>"},
			want:   "<!subteam^foo>",
		},
		{
			name:   "multiple",
			groups: []string{"foo", "bar"},
			want:   "<!subteam^foo> <!subteam^bar>",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getGroupsAsString(tt.groups)
			if got != tt.want {
				t.Errorf("getGroupsAsString(): got %s, want %s", got, tt.want)
			}
		})
	}
}
