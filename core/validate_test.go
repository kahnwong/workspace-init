package core

import (
	"reflect"
	"testing"
)

func TestSubtractArrays(t *testing.T) {
	tests := []struct {
		name     string
		array1   []string
		array2   []string
		expected []string
	}{
		{
			name:     "subtract some elements",
			array1:   []string{"a", "b", "c", "d"},
			array2:   []string{"b", "d"},
			expected: []string{"a", "c"},
		},
		{
			name:     "subtract all elements",
			array1:   []string{"a", "b", "c"},
			array2:   []string{"a", "b", "c"},
			expected: []string{},
		},
		{
			name:     "subtract no elements",
			array1:   []string{"a", "b", "c"},
			array2:   []string{"x", "y", "z"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty first array",
			array1:   []string{},
			array2:   []string{"a", "b"},
			expected: []string{},
		},
		{
			name:     "empty second array",
			array1:   []string{"a", "b", "c"},
			array2:   []string{},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "both arrays empty",
			array1:   []string{},
			array2:   []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := subtractArrays(tt.array1, tt.array2)
			if result == nil {
				result = []string{}
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("subtractArrays(%v, %v) = %v, want %v", tt.array1, tt.array2, result, tt.expected)
			}
		})
	}
}
