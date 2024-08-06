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

package lambda

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	t.Parallel()

	old, err := loadFunctionManifest("testdata/old_function.yaml")
	require.NoError(t, err)

	new, err := loadFunctionManifest("testdata/new_function.yaml")
	require.NoError(t, err)

	// Expect to have change.
	diff, err := Diff(old, new)
	require.NoError(t, err)
	noChange := diff.NoChange()
	require.False(t, noChange)

	// Expect no change.
	diff, err = Diff(old, old)
	require.NoError(t, err)
	noChange = diff.NoChange()
	require.True(t, noChange)
}

func TestDiffResult_Render(t *testing.T) {
	old, err := loadFunctionManifest("testdata/old_function.yaml")
	require.NoError(t, err)

	new, err := loadFunctionManifest("testdata/new_function.yaml")
	require.NoError(t, err)

	result, err := Diff(old, new)
	require.NoError(t, err)

	// Not using diff command
	opt := DiffRenderOptions{}
	actual := result.Render(opt)
	expected := `  spec:
    environments:
      #spec.environments.FOO
-     FOO: bar
+     FOO: bar2

    #spec.image
-   image: ecr.ap-northeast-1.amazonaws.com/lambda-test:v0.0.1
+   image: ecr.ap-northeast-1.amazonaws.com/lambda-test:v0.0.2


`

	require.Equal(t, expected, actual)

	// Use diff command
	opt = DiffRenderOptions{UseDiffCommand: true}
	actual = result.Render(opt)
	expected = `@@ -2,9 +2,9 @@
 kind: LambdaFunction
 spec:
   environments:
-    FOO: bar
+    FOO: bar2
   handler: ""
-  image: ecr.ap-northeast-1.amazonaws.com/lambda-test:v0.0.1
+  image: ecr.ap-northeast-1.amazonaws.com/lambda-test:v0.0.2
   memory: 512
   name: TestFunction
   role: arn:aws:iam:region:account-id:role/lambda-role
`

	require.Equal(t, expected, actual)
}
