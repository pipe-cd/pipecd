package promptreader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockReadString(t *testing.T) {
	testcases := []struct {
		name          string
		inputs        []string
		expected      string
		expectedArray []string
		expectedErr   bool
	}{
		{
			name:          "valid input",
			inputs:        []string{"abc"},
			expected:      "abc",
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "empty input",
			inputs:        []string{""},
			expected:      "",
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "too many arguments",
			inputs:        []string{"abc def"},
			expected:      "",
			expectedArray: []string{},
			expectedErr:   true,
		},
		{
			name:          "do not read next element",
			inputs:        []string{"abc", "def"},
			expected:      "abc",
			expectedArray: []string{"def"},
			expectedErr:   false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := MockPromptReader{inputs: tc.inputs}
			actual, e := r.readString("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedArray, r.inputs)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}

}

func TestMockReadStrings(t *testing.T) {
	testcases := []struct {
		name          string
		inputs        []string
		expected      []string
		expectedArray []string
		expectedErr   bool
	}{
		{
			name:          "valid input",
			inputs:        []string{"abc def"},
			expected:      []string{"abc", "def"},
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "empty input",
			inputs:        []string{""},
			expected:      []string{},
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "do not read next element",
			inputs:        []string{"abc def", "ghi jkl"},
			expected:      []string{"abc", "def"},
			expectedArray: []string{"ghi jkl"},
			expectedErr:   false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := MockPromptReader{inputs: tc.inputs}
			actual, e := r.readStrings("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedArray, r.inputs)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestMockReadInt(t *testing.T) {
	testcases := []struct {
		name          string
		inputs        []string
		expected      int
		expectedArray []string
		expectedErr   bool
	}{
		{
			name:          "valid input",
			inputs:        []string{"123"},
			expected:      123,
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "empty input",
			inputs:        []string{""},
			expected:      0,
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "string input",
			inputs:        []string{"abc"},
			expected:      0,
			expectedArray: []string{},
			expectedErr:   true,
		},
		{
			name:          "do not read next element",
			inputs:        []string{"1", "23"},
			expected:      1,
			expectedArray: []string{"23"},
			expectedErr:   false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := MockPromptReader{inputs: tc.inputs}
			actual, e := r.readInt("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedArray, r.inputs)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestMockReadStringRequired(t *testing.T) {
	testcases := []struct {
		name          string
		inputs        []string
		expected      string
		expectedArray []string
		expectedErr   bool
	}{
		{
			name:          "valid input",
			inputs:        []string{"abc"},
			expected:      "abc",
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "empty input",
			inputs:        []string{""},
			expected:      "",
			expectedArray: []string{},
			expectedErr:   true,
		},
		{
			name:          "too many arguments",
			inputs:        []string{"abc def"},
			expected:      "",
			expectedArray: []string{},
			expectedErr:   true,
		},
		{
			name:          "do not read next element",
			inputs:        []string{"abc", "def"},
			expected:      "abc",
			expectedArray: []string{"def"},
			expectedErr:   false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := MockPromptReader{inputs: tc.inputs}
			actual, e := r.readStringRequired("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedArray, r.inputs)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}
