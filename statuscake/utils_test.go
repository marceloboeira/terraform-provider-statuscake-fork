package statuscake

import (
	"reflect"
	"testing"
)

func TestToStringSlice(t *testing.T) {
	var tests = []struct {
		name     string
		input    []interface{}
		expected []string
	}{
		{
			name:     "returns a string slice from a slice of strings",
			input:    []string{"1", "2", "3"},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "returns a string slice from a slice of integers",
			input:    []int{1, 2, 3},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "returns a string slice from a slice of booleans",
			input:    []bool{true, false, true},
			expected: []string{"true", "false", "true"},
		},
		{
			name:     "returns a string slice from a slice of mixed types",
			input:    []interface{}{1, "2", true},
			expected: []string{"1", "2", "true"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toStringSlice(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected: %+v, ggt: %+v", tt.expected, got)
			}
		})
	}
}
