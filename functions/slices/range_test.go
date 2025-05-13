package slices

import (
	"slices"
	"testing"
)

func TestIRange(t *testing.T) {
	t.Run("Basic use cases", func(t *testing.T) {
		tests := []struct {
			min, max int
			step     []int
			expected []int
		}{
			{1, 5, nil, []int{1, 2, 3, 4, 5}},
			{5, 1, nil, []int{5, 4, 3, 2, 1}},
			{0, 10, []int{2}, []int{0, 2, 4, 6, 8, 10}},
			{10, 0, []int{-2}, []int{10, 8, 6, 4, 2, 0}},
			{3, 3, nil, []int{3}},
		}

		for _, tt := range tests {
			out := IRange(tt.min, tt.max, tt.step...)
			if !slices.Equal(out, tt.expected) {
				t.Errorf("IRange(%v, %v, %v) = %v; want %v", tt.min, tt.max, tt.step, out, tt.expected)
			}
		}
	})

	t.Run("Invalid range", func(t *testing.T) {
		out := IRange(1, 10, -1)
		if !slices.Equal(out, []int{}) {
			t.Errorf("Expected to get empty slice. Got %#v", out)
		}

		out = IRange(10, 1, 1)
		if !slices.Equal(out, []int{}) {
			t.Errorf("Expected to get empty slice for step in wrong direction. Got %#v", out)
		}
	})
}

func TestERange(t *testing.T) {
	t.Run("Basic use cases", func(t *testing.T) {
		tests := []struct {
			min, max int
			step     []int
			expected []int
		}{
			{1, 5, nil, []int{1, 2, 3, 4}},
			{5, 1, nil, []int{5, 4, 3, 2}},
			{0, 10, []int{3}, []int{0, 3, 6, 9}},
			{10, 0, []int{-3}, []int{10, 7, 4, 1}},
			{3, 3, nil, []int{}},
		}

		for _, tt := range tests {
			out := ERange(tt.min, tt.max, tt.step...)
			if !slices.Equal(out, tt.expected) {
				t.Errorf("ERange(%v, %v, %v) = %v; want %v", tt.min, tt.max, tt.step, out, tt.expected)
			}
		}
	})

	t.Run("Invalid range", func(t *testing.T) {
		out := ERange(1, 10, -1)
		if !slices.Equal(out, []int{}) {
			t.Errorf("Expected to get empty slice for step in wrong direction. Got %#v", out)
		}

		out = ERange(10, 1, 1)
		if !slices.Equal(out, []int{}) {
			t.Errorf("Expected to get empty slice step in wrong direction. Got %#v", out)
		}
	})
}
