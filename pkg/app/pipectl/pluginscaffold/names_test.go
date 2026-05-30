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

package pluginscaffold

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStageNames(t *testing.T) {
	t.Parallel()

	assert.NoError(t, ValidateStageNames([]string{"DEMO_SYNC"}))
	assert.Error(t, ValidateStageNames([]string{"demo_sync"}))
	assert.Error(t, ValidateStageNames(nil))
}

func TestTypePrefix(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "MyPlugin", TypePrefix("my-plugin"))
	assert.Equal(t, "Demo", TypePrefix("demo"))
}

func TestStageFuncName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "executeDemoSync", StageFuncName("DEMO_SYNC"))
}

func TestStageConstSuffix(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "DemoSync", StageConstSuffix("DEMO_SYNC"))
}

func TestFindRollbackStage(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "DEMO_ROLLBACK", FindRollbackStage([]string{"DEMO_SYNC", "DEMO_ROLLBACK"}))
	assert.Empty(t, FindRollbackStage([]string{"DEMO_SYNC"}))
}
