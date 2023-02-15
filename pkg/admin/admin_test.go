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

package admin

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleTop(t *testing.T) {
	req := httptest.NewRequest("GET", "http://admin", nil)

	testcases := []struct {
		name     string
		admin    *Admin
		expected string
	}{
		{
			name:  "no pattern",
			admin: &Admin{},
			expected: `
<!DOCTYPE html>
<html>
<body>

<h3>Admin Page</h3>
</body>
</html>
`,
		},
		{
			name: "there are some patterns",
			admin: &Admin{
				patterns: []string{"metrics", "healthz"},
			},
			expected: `
<!DOCTYPE html>
<html>
<body>

<h3>Admin Page</h3>
<p><a href="metrics">metrics</a></p>
<p><a href="healthz">healthz</a></p>
</body>
</html>
`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tc.admin.handleTop(w, req)
			body, _ := io.ReadAll(w.Body)
			assert.Equal(t, tc.expected, string(body))
		})
	}
}
