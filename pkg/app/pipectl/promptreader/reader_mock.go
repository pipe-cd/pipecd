package promptreader

import (
	"fmt"
	"strconv"
	"strings"
)

// mockPromptReader is a mock implementation of promptReader for unit tests.
// It reads input from the given slice of strings.
type mockPromptReader struct {
	inputs []string
}

func (r *mockPromptReader) readString(message string) (string, error) {
	if len(r.inputs) == 0 {
		return "", fmt.Errorf("no more input")
	}
	input := r.inputs[0]
	r.inputs = r.inputs[1:]
	f := strings.Fields(input)
	if len(f) > 1 {
		return "", fmt.Errorf("too many arguments")
	}

	if len(f) == 0 {
		return "", nil
	}
	return f[0], nil
}

func (r *mockPromptReader) readStrings(message string) ([]string, error) {
	if len(r.inputs) == 0 {
		return nil, fmt.Errorf("no more input")
	}
	input := r.inputs[0]
	r.inputs = r.inputs[1:]
	return strings.Fields(input), nil
}

func (r *mockPromptReader) readInt(message string) (int, error) {
	if len(r.inputs) == 0 {
		return 0, fmt.Errorf("no more input")
	}
	input := r.inputs[0]
	r.inputs = r.inputs[1:]
	if len(input) == 0 {
		return 0, nil
	}
	return strconv.Atoi(input)
}

func (r *mockPromptReader) readStringRequired(message string) (string, error) {
	s, e := r.readString(message)
	if e != nil {
		return "", e
	}
	if len(s) == 0 {
		return "", fmt.Errorf("empty input")
	}
	return s, e
}
