package promptreader

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadString(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expected    string
		expectedErr bool
	}{
		{
			name:        "valid input",
			input:       "foo\n",
			expected:    "foo",
			expectedErr: false,
		},
		{
			name:        "empty input",
			input:       "\n",
			expected:    "",
			expectedErr: false,
		},
		{
			name:        "two words",
			input:       "foo bar\n",
			expected:    "",
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.input)
			r := DefaultPromptReader{in: in}
			actual, e := r.readString("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestReadStrings(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expected    []string
		expectedErr bool
	}{
		{
			name:        "valid input",
			input:       "foo\n",
			expected:    []string{"foo"},
			expectedErr: false,
		},
		{
			name:        "empty input",
			input:       "\n",
			expected:    []string{},
			expectedErr: false,
		},
		{
			name:        "two words",
			input:       "foo bar\n",
			expected:    []string{"foo", "bar"},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.input)
			r := DefaultPromptReader{in: in}
			actual, e := r.readStrings("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestReadInt(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expected    int
		expectedErr bool
	}{
		{
			name:        "valid input",
			input:       "123\n",
			expected:    123,
			expectedErr: false,
		},
		{
			name:        "empty input",
			input:       "\n",
			expected:    0,
			expectedErr: false,
		},
		{
			name:        "string input",
			input:       "abc\n",
			expected:    0,
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.input)
			r := DefaultPromptReader{in: in}
			actual, e := r.readInt("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestReadStringRequired(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expected    string
		expectedErr bool
	}{
		{
			name:        "valid input",
			input:       "foo\n",
			expected:    "foo",
			expectedErr: false,
		},
		{
			name:        "empty input",
			input:       "\n",
			expected:    "",
			expectedErr: true,
		},
		{
			name:        "two words",
			input:       "foo bar\n",
			expected:    "",
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.input)
			r := DefaultPromptReader{in: in}
			actual, e := r.readStringRequired("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}
