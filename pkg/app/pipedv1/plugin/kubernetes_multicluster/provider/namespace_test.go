// Copyright 2026 The PipeCD Authors.
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

package provider

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveDefaultNamespace(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name       string
		kubeconfig string
		want       string
	}{
		{
			name:       "context without namespace",
			kubeconfig: "apiVersion: v1\nkind: Config\nclusters:\n- name: c\n  cluster:\n    server: https://127.0.0.1:6443\nusers:\n- name: u\n  user: {}\ncontexts:\n- name: ctx\n  context:\n    cluster: c\n    user: u\ncurrent-context: ctx\n",
			want:       "default",
		},
		{
			name:       "context with namespace",
			kubeconfig: "apiVersion: v1\nkind: Config\nclusters:\n- name: c\n  cluster:\n    server: https://127.0.0.1:6443\nusers:\n- name: u\n  user: {}\ncontexts:\n- name: ctx\n  context:\n    cluster: c\n    user: u\n    namespace: team-a\ncurrent-context: ctx\n",
			want:       "team-a",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), "kubeconfig")
			require.NoError(t, os.WriteFile(path, []byte(tc.kubeconfig), 0o600))

			assert.Equal(t, tc.want, ResolveDefaultNamespace(path))
		})
	}

	t.Run("unreadable kubeconfig", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "default", ResolveDefaultNamespace(filepath.Join(t.TempDir(), "missing")))
	})
}
