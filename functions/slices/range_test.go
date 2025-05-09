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
			out, err := IRange(tt.min, tt.max, tt.step...)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !slices.Equal(out, tt.expected) {
				t.Errorf("IRange(%v, %v, %v) = %v; want %v", tt.min, tt.max, tt.step, out, tt.expected)
			}
		}
	})

	t.Run("Invalid range", func(t *testing.T) {
		_, err := IRange(1, 10, -1)
		if err == nil {
			t.Error("expected error for step in wrong direction, got nil")
		}

		_, err = IRange(10, 1, 1)
		if err == nil {
			t.Error("expected error for step in wrong direction, got nil")
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
			out, err := ERange(tt.min, tt.max, tt.step...)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !slices.Equal(out, tt.expected) {
				t.Errorf("ERange(%v, %v, %v) = %v; want %v", tt.min, tt.max, tt.step, out, tt.expected)
			}
		}
	})

	t.Run("Invalid range", func(t *testing.T) {
		_, err := ERange(1, 10, -1)
		if err == nil {
			t.Error("expected error for step in wrong direction, got nil")
		}

		_, err = ERange(10, 1, 1)
		if err == nil {
			t.Error("expected error for step in wrong direction, got nil")
		}
	})
}

func TestValidateRangeParams(t *testing.T) {
	type testCase[T any] struct {
		min, max, step T
		expectedStep   T
		expectErr      bool
	}

	tests := []testCase[int]{
		{1, 5, 0, 1, false},
		{5, 1, 0, -1, false},
		{5, 1, -1, -1, false},
		{1, 5, -1, -1, true},
		{5, 1, 1, 1, true},
	}

	for _, tt := range tests {
		step, err := validateRangeParams(tt.min, tt.max, tt.step)

		if tt.expectErr {
			if err == nil {
				t.Errorf("expected error for min=%v, max=%v, step=%v", tt.min, tt.max, tt.step)
			}

			continue
		}

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if step != tt.expectedStep {
			t.Errorf("got step %v, want %v", step, tt.expectedStep)
		}
	}
}
