package yamlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	testcases := []struct {
		name    string
		yml     string
		path    string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "empty yaml given",
			yml:     "",
			path:    "$.foo",
			wantErr: true,
		},
		{
			name:    "empty path given",
			yml:     "foo: bar",
			path:    "",
			wantErr: true,
		},
		{
			name:    "wrong path given",
			yml:     "foo: bar",
			path:    "wrong",
			wantErr: true,
		},
		{
			name:    "wrong yaml given",
			yml:     "::",
			path:    "$.foo",
			wantErr: true,
		},
		{
			name:    "lack of root element",
			yml:     "foo: bar",
			path:    "foo",
			wantErr: true,
		},
		{
			name:    "given a string path",
			yml:     "foo: bar",
			path:    "$.foo",
			want:    "bar",
			wantErr: false,
		},
		{
			name:    "given a bool path",
			yml:     "foo: true",
			path:    "$.foo",
			want:    true,
			wantErr: false,
		},
		{
			name:    "given a uint64 path",
			yml:     "foo: 1",
			path:    "$.foo",
			want:    uint64(1),
			wantErr: false,
		},
		{
			name:    "given a float64 path",
			yml:     "foo: 1.5",
			path:    "$.foo",
			want:    1.5,
			wantErr: false,
		},
		{
			name: "given a array path",
			yml: `
foo:
- bar: 1`,
			path:    "$.foo[0].bar",
			want:    uint64(1),
			wantErr: false,
		},
		{
			name: "given a array path with wildcard",
			yml: `
foo:
- bar: 1
- baz: 2`,
			path:    "$.foo[*]",
			want:    []interface{}{map[string]interface{}{"bar": uint64(1)}, map[string]interface{}{"baz": uint64(2)}},
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetValue([]byte(tc.yml), tc.path)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestReplaceValue(t *testing.T) {
	testcases := []struct {
		name    string
		yml     string
		path    string
		value   interface{}
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty yaml given",
			yml:     "",
			path:    "$.foo",
			value:   1,
			wantErr: true,
		},
		{
			name:    "empty path given",
			yml:     "foo: bar",
			path:    "",
			value:   1,
			wantErr: true,
		},
		{
			name:    "string value given",
			yml:     "foo: bar",
			path:    "$.foo",
			value:   "new-text",
			want:    []byte("foo: new-text"),
			wantErr: false,
		},
		{
			name:    "bool value given",
			yml:     "foo: bar",
			path:    "$.foo",
			value:   "true",
			want:    []byte("foo: true"),
			wantErr: false,
		},
		{
			name:    "float64 value given",
			yml:     "foo: bar",
			path:    "$.foo",
			value:   "1.5",
			want:    []byte("foo: 1.5"),
			wantErr: false,
		},
		{
			name:    "int value given",
			yml:     "foo: bar",
			path:    "$.foo",
			value:   "1",
			want:    []byte("foo: 1"),
			wantErr: false,
		},
		{
			name:    "nil given",
			yml:     "foo: bar",
			path:    "$.foo",
			value:   nil,
			want:    []byte("foo: null"),
			wantErr: false,
		},
		{
			name:    "unsupported type given",
			yml:     "foo: bar",
			path:    "$.foo",
			value:   []string{"bar"},
			wantErr: true,
		},
		{
			name: "there is an useless blank line",
			yml: `
foo:
  - bar
  - baz`,
			path:  "$.foo[0]",
			value: "new-text",
			want: []byte(`foo:
  - new-text
  - baz`),
			wantErr: false,
		},
		{
			name:    "array in flow style",
			yml:     `foo: [bar, baz]`,
			path:    "$.foo[0]",
			value:   "new-text",
			want:    []byte(`foo: [new-text, baz]`),
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ReplaceValue([]byte(tc.yml), tc.path, tc.value)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
