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

package diff

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func TestDiff(t *testing.T) {
	testcases := []struct {
		name        string
		yamlFile    string
		resourceKey string
		options     []Option
		diffNum     int
		diffString  string
	}{
		{
			name:     "no diff",
			yamlFile: "testdata/no_diff.yaml",
			options: []Option{
				WithEquateEmpty(),
				WithIgnoreAddingMapKeys(),
				WithCompareNumberAndNumericString(),
			},
			diffNum: 0,
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
			name:        "diff by ignoring specified field with correct key",
			yamlFile:    "testdata/has_diff.yaml",
			resourceKey: "deployment-key",
			options: []Option{
				WithIgnoreConfig(
					map[string][]string{
						"deployment-key": {
							"spec.replicas",
							"spec.template.spec.containers.0.args.1",
							"spec.template.spec.strategy.rollingUpdate.maxSurge",
							"spec.template.spec.containers.3.livenessProbe.initialDelaySeconds",
						},
					},
				),
			},
			diffNum: 6,
			diffString: `  spec:
    template:
      metadata:
        labels:
          #spec.template.metadata.labels.app
-         app: simple
+         app: simple2

          #spec.template.metadata.labels.component
-         component: foo

      spec:
        containers:
          -
            #spec.template.spec.containers.1.image
-           image: gcr.io/pipecd/helloworld:v2.0.0
+           image: gcr.io/pipecd/helloworld:v2.1.0

          -
            #spec.template.spec.containers.2.image
-           image: 

          #spec.template.spec.containers.3
+         - image: new-image
+           livenessProbe:
+             exec:
+               command:
+                 - cat
+                 - /tmp/healthy
+           name: foo

        #spec.template.spec.strategy
+       strategy:
+         rollingUpdate:
+           maxUnavailable: 25%
+         type: RollingUpdate

`,
		},
		{
			name:        "diff by ignoring specified field with wrong resource key",
			yamlFile:    "testdata/has_diff.yaml",
			resourceKey: "deployment-key",
			options: []Option{
				WithIgnoreConfig(
					map[string][]string{
						"crd-key": {
							"spec.replicas",
							"spec.template.spec.containers.0.args.1",
							"spec.template.spec.strategy.rollingUpdate.maxSurge",
							"spec.template.spec.containers.3.livenessProbe.initialDelaySeconds",
						},
					},
				),
			},
			diffNum: 8,
			diffString: `  spec:
    #spec.replicas
-   replicas: 2
+   replicas: 3

    template:
      metadata:
        labels:
          #spec.template.metadata.labels.app
-         app: simple
+         app: simple2

          #spec.template.metadata.labels.component
-         component: foo

      spec:
        containers:
          - args:
              #spec.template.spec.containers.0.args.1
-             - hello

          -
            #spec.template.spec.containers.1.image
-           image: gcr.io/pipecd/helloworld:v2.0.0
+           image: gcr.io/pipecd/helloworld:v2.1.0

          -
            #spec.template.spec.containers.2.image
-           image: 

          #spec.template.spec.containers.3
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
		{
			name:     "has diff",
			yamlFile: "testdata/has_diff.yaml",
			diffNum:  8,
			diffString: `  spec:
    #spec.replicas
-   replicas: 2
+   replicas: 3

    template:
      metadata:
        labels:
          #spec.template.metadata.labels.app
-         app: simple
+         app: simple2

          #spec.template.metadata.labels.component
-         component: foo

      spec:
        containers:
          - args:
              #spec.template.spec.containers.0.args.1
-             - hello

          -
            #spec.template.spec.containers.1.image
-           image: gcr.io/pipecd/helloworld:v2.0.0
+           image: gcr.io/pipecd/helloworld:v2.1.0

          -
            #spec.template.spec.containers.2.image
-           image: 

          #spec.template.spec.containers.3
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

			result, err := DiffUnstructureds(objs[0], objs[1], tc.resourceKey, tc.options...)
			require.NoError(t, err)
			assert.Equal(t, tc.diffNum, result.NumNodes())

			renderer := NewRenderer(WithLeftPadding(1))
			ds := renderer.Render(result.Nodes())

			assert.Equal(t, tc.diffString, ds)
		})
	}
}

func loadUnstructureds(path string) ([]unstructured.Unstructured, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	const separator = "\n---"
	parts := strings.Split(string(data), separator)
	out := make([]unstructured.Unstructured, 0, len(parts))

	for _, part := range parts {
		// Ignore all the cases where no content between separator.
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

func TestDiffStructs(t *testing.T) {
	type SubStruct struct {
		FieldSub string
	}
	type TestStruct struct {
		Field1 string
		Field2 int
		Field3 bool
		Field4 []string
		Field5 map[string]string
		Field6 SubStruct
	}

	testcases := []struct {
		name    string
		old     TestStruct
		new     TestStruct
		key     string
		diffNum int
	}{
		{
			name: "no diff",
			old: TestStruct{
				Field1: "value1",
				Field2: 1,
				Field3: true,
				Field4: []string{"a", "b"},
				Field5: map[string]string{
					"key1": "value1",
				},
				Field6: SubStruct{
					FieldSub: "valueSub",
				},
			},
			new: TestStruct{
				Field1: "value1",
				Field2: 1,
				Field3: true,
				Field4: []string{"a", "b"},
				Field5: map[string]string{
					"key1": "value1",
				},
				Field6: SubStruct{
					FieldSub: "valueSub",
				},
			},
			key:     "value1",
			diffNum: 0,
		},
		{
			name: "has diff",
			old: TestStruct{
				Field1: "value1",
				Field2: 1,
				Field3: true,
				Field4: []string{"a", "b"},
				Field5: map[string]string{
					"key1": "value1",
				},
				Field6: SubStruct{
					FieldSub: "valueSub",
				},
			},
			new: TestStruct{
				Field1: "value2",
				Field2: 10,
				Field3: false,
				Field4: []string{"a", "b", "c"},
				Field5: map[string]string{
					"key1": "valueXXX",
				},
				Field6: SubStruct{
					FieldSub: "xxx",
				},
			},
			key:     "value1",
			diffNum: 6,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			diff, err := DiffStructs(tc.old, tc.new)
			require.NoError(t, err)
			assert.Equal(t, tc.diffNum, diff.NumNodes())
		})
	}
}

func TestIsEmptyInterface(t *testing.T) {
	testcases := []struct {
		name     string
		v        interface{}
		expected bool
	}{
		{
			name:     "nil",
			v:        nil,
			expected: true,
		},
		{
			name:     "nil map",
			v:        map[string]int(nil),
			expected: true,
		},
		{
			name:     "empty map",
			v:        map[string]int{},
			expected: true,
		},
		{
			name:     "nil slice",
			v:        []int(nil),
			expected: true,
		},
		{
			name:     "empty slice",
			v:        []int{},
			expected: true,
		},
		{
			name:     "number",
			v:        1,
			expected: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := []interface{}{tc.v}
			v := reflect.ValueOf(s)

			got := isEmptyInterface(v.Index(0))
			assert.Equal(t, tc.expected, got)
		})
	}
}
