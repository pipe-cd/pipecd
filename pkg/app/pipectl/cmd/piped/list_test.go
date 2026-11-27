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

package piped

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/client"
)

func TestNewListCommand(t *testing.T) {
	c := &command{
		clientOptions: &client.Options{},
	}
	cmd := newListCommand(c)

	assert.Equal(t, "list", cmd.Use)
	assert.NotEmpty(t, cmd.Short)

	flag := cmd.Flags().Lookup("with-status")
	assert.NotNil(t, flag)
	assert.Equal(t, "false", flag.DefValue)
}
