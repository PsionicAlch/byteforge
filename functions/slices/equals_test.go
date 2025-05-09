package slices

import (
	"testing"
)

func TestShallowEquals(t *testing.T) {
	tests := []struct {
		name     string
		s1       any
		s2       any
		expected bool
	}{
		{
			name:     "Empty slices",
			s1:       []int{},
			s2:       []int{},
			expected: true,
		},
		{
			name:     "Same elements same order",
			s1:       []int{1, 2, 3, 4, 5},
			s2:       []int{1, 2, 3, 4, 5},
			expected: true,
		},
		{
			name:     "Same elements different order",
			s1:       []int{1, 2, 3, 4, 5},
			s2:       []int{5, 4, 3, 2, 1},
			expected: true,
		},
		{
			name:     "Different lengths",
			s1:       []int{1, 2, 3},
			s2:       []int{1, 2, 3, 4},
			expected: false,
		},
		{
			name:     "Different elements",
			s1:       []int{1, 2, 3, 4, 5},
			s2:       []int{1, 2, 3, 4, 6},
			expected: false,
		},
		{
			name:     "Different frequency of elements",
			s1:       []int{1, 2, 3, 4, 5},
			s2:       []int{1, 2, 3, 5, 5},
			expected: false,
		},
		{
			name:     "With duplicate elements",
			s1:       []int{1, 2, 2, 3, 3},
			s2:       []int{3, 3, 2, 2, 1},
			expected: true,
		},
		{
			name:     "String slices equal",
			s1:       []string{"apple", "banana", "cherry"},
			s2:       []string{"cherry", "banana", "apple"},
			expected: true,
		},
		{
			name:     "String slices not equal",
			s1:       []string{"apple", "banana", "cherry"},
			s2:       []string{"apple", "banana", "orange"},
			expected: false,
		},
		{
			name:     "With negative numbers",
			s1:       []int{-1, -2, -3, 0, 1},
			s2:       []int{1, 0, -3, -2, -1},
			expected: true,
		},
		{
			name:     "All same element",
			s1:       []int{5, 5, 5, 5},
			s2:       []int{5, 5, 5, 5},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch s1 := tt.s1.(type) {
			case []int:
				s2 := tt.s2.([]int)
				if got := ShallowEquals(s1, s2); got != tt.expected {
					t.Errorf("ShallowEquals() = %v, want %v for %v and %v", got, tt.expected, s1, s2)
				}
			case []string:
				s2 := tt.s2.([]string)
				if got := ShallowEquals(s1, s2); got != tt.expected {
					t.Errorf("ShallowEquals() = %v, want %v for %v and %v", got, tt.expected, s1, s2)
				}
			}
		})
	}
}

func TestDeepEquals(t *testing.T) {
	tests := []struct {
		name     string
		s1       any
		s2       any
		expected bool
	}{
		{
			name:     "Empty slices",
			s1:       []int{},
			s2:       []int{},
			expected: true,
		},
		{
			name:     "Same elements same order",
			s1:       []int{1, 2, 3, 4, 5},
			s2:       []int{1, 2, 3, 4, 5},
			expected: true,
		},
		{
			name:     "Same elements different order",
			s1:       []int{1, 2, 3, 4, 5},
			s2:       []int{5, 4, 3, 2, 1},
			expected: false, // Order matters for DeepEquals
		},
		{
			name:     "Different lengths",
			s1:       []int{1, 2, 3},
			s2:       []int{1, 2, 3, 4},
			expected: false,
		},
		{
			name:     "Different elements",
			s1:       []int{1, 2, 3, 4, 5},
			s2:       []int{1, 2, 3, 4, 6},
			expected: false,
		},
		{
			name:     "String slices equal",
			s1:       []string{"apple", "banana", "cherry"},
			s2:       []string{"apple", "banana", "cherry"},
			expected: true,
		},
		{
			name:     "String slices equal order different",
			s1:       []string{"apple", "banana", "cherry"},
			s2:       []string{"cherry", "banana", "apple"},
			expected: false, // Order matters for DeepEquals
		},
		{
			name:     "With negative numbers same order",
			s1:       []int{-1, -2, -3, 0, 1},
			s2:       []int{-1, -2, -3, 0, 1},
			expected: true,
		},
		{
			name:     "Empty vs non-empty",
			s1:       []int{},
			s2:       []int{1},
			expected: false,
		},
		{
			name:     "Single element",
			s1:       []int{42},
			s2:       []int{42},
			expected: true,
		},
		{
			name:     "Duplicate elements in same positions",
			s1:       []int{1, 2, 2, 3, 3},
			s2:       []int{1, 2, 2, 3, 3},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch s1 := tt.s1.(type) {
			case []int:
				s2 := tt.s2.([]int)
				if got := DeepEquals(s1, s2); got != tt.expected {
					t.Errorf("DeepEquals() = %v, want %v for %v and %v", got, tt.expected, s1, s2)
				}
			case []string:
				s2 := tt.s2.([]string)
				if got := DeepEquals(s1, s2); got != tt.expected {
					t.Errorf("DeepEquals() = %v, want %v for %v and %v", got, tt.expected, s1, s2)
				}
			}
		})
	}
}
