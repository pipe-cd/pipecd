// Copyright 2020 The PipeCD Authors.
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

package diff

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func TestDiff(t *testing.T) {
	testcases := []struct {
		name       string
		yamlFile   string
		options    []Option
		diffNum    int
		diffString string
	}{
		{
			name:     "no diff",
			yamlFile: "testdata/no_diff.yaml",
			diffNum:  0,
		},
		{
			name:     "no diff by ignoring all adding map keys",
			yamlFile: "testdata/ignore_adding_map_keys.yaml",
			options: []Option{
				WithIgnoreAddingMapKeys(),
			},
			diffNum: 0,
		},
		{
			name:     "has diff",
			yamlFile: "testdata/has_diff.yaml",
			diffNum:  6,
			diffString: `  #spec.replicas
  spec:
-   replicas: 2
+   replicas: 3
    #spec.template.metadata.labels.app
    template:
      metadata:
        labels:
-         app: simple
+         app: simple2
          #spec.template.metadata.labels.component
-         component: foo
      #spec.template.spec.containers.0.args.1
      spec:
        containers:
          - args:
-             - hello
          #spec.template.spec.containers.1
+         - image: new-image
+           livenessProbe:
+             exec:
+               command:
+                 - cat
+                 - /tmp/healthy
+             initialDelaySeconds: 5
+           name: foo
        #spec.template.spec.strategy
+       strategy:
+         rollingUpdate:
+           maxSurge: 25%
+           maxUnavailable: 25%
+         type: RollingUpdate
`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			objs, err := loadUnstructureds(tc.yamlFile)
			require.NoError(t, err)
			require.Equal(t, 2, len(objs))

			result, err := DiffUnstructureds(objs[0], objs[1], tc.options...)
			require.NoError(t, err)
			assert.Equal(t, tc.diffNum, result.Num())

			result.SetLeftPadding(1)
			ds := result.DiffString()
			assert.Equal(t, tc.diffString, ds)
		})
	}
}

func loadUnstructureds(path string) ([]unstructured.Unstructured, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	const separator = "\n---"
	parts := strings.Split(string(data), separator)
	out := make([]unstructured.Unstructured, 0, len(parts))

	for _, part := range parts {
		//	Ignore all the cases where no content between separator.
		part = strings.TrimSpace(part)
		if len(part) == 0 {
			continue
		}
		var obj unstructured.Unstructured
		if err := yaml.Unmarshal([]byte(part), &obj); err != nil {
			return nil, err
		}
		out = append(out, obj)
	}
	return out, nil
}
