// Copyright 2025 The PipeCD Authors.
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

package oci

import (
	"errors"
	"os"
	"testing"
)

func TestIsDockerUnavailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "missing socket path",
			err:  errors.New("dial unix /var/run/docker.sock: connect: no such file or directory"),
			want: true,
		},
		{
			name: "daemon not running",
			err:  errors.New("Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"),
			want: true,
		},
		{
			name: "wrapped os err not exist",
			err:  errors.Join(errors.New("connect docker"), os.ErrNotExist),
			want: true,
		},
		{
			name: "registry startup failure",
			err:  errors.New("failed to pull image"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isDockerUnavailable(tt.err); got != tt.want {
				t.Fatalf("isDockerUnavailable(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
