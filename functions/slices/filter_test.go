package slices

import (
	"slices"
	"testing"

	islices "github.com/PsionicAlch/byteforge/internal/functions/slices"
)

func TestFilter(t *testing.T) {
	const max = 1000000
	largeArr := islices.IRange(1, max)

	largeExpected := make([]int, 0, max)
	for _, num := range largeArr {
		if num%2 == 0 {
			largeExpected = append(largeExpected, num)
		}
	}

	t.Run("Filter small slice of int", func(t *testing.T) {
		result := Filter(islices.IRange(1, 10), func(num int) bool {
			return num%2 == 0
		})
		expected := []int{2, 4, 6, 8, 10}

		if !slices.Equal(result, expected) {
			t.Errorf("Expected result to be %#v. Got %#v", expected, result)
		}
	})

	t.Run("Filter large slice of int", func(t *testing.T) {
		result := Filter(largeArr, func(num int) bool {
			return num%2 == 0
		})

		if !slices.Equal(result, largeExpected) {
			t.Errorf("Expected result to be %#v. Got %#v", largeExpected, result)
		}
	})
}

func TestParallelFilter(t *testing.T) {
	const max = 1000000
	largeArr := islices.IRange(1, max)

	largeExpected := make([]int, 0, max)
	for _, num := range largeArr {
		if num%2 == 0 {
			largeExpected = append(largeExpected, num)
		}
	}

	t.Run("Parallel filter small slice of int", func(t *testing.T) {
		result := ParallelFilter(islices.IRange(1, 10), func(num int) bool {
			return num%2 == 0
		})
		expected := []int{2, 4, 6, 8, 10}

		if !slices.Equal(result, expected) {
			t.Errorf("Expected result to be %#v. Got %#v", expected, result)
		}
	})

	t.Run("Parallel filter large slice of int", func(t *testing.T) {
		result := ParallelFilter(largeArr, func(num int) bool {
			return num%2 == 0
		})

		if !slices.Equal(result, largeExpected) {
			t.Errorf("Expected result to be %#v. Got %#v", largeExpected, result)
		}
	})

	t.Run("Filter with empty slice", func(t *testing.T) {
		result := ParallelFilter([]int{}, func(_ int) bool {
			return true
		})
		expected := []int{}

		if !slices.Equal(result, expected) {
			t.Errorf("Expected result to be %#v. Got %#v", expected, result)
		}
	})

	t.Run("Parallel filter with worker pool", func(t *testing.T) {
		result := ParallelFilter(largeArr, func(num int) bool {
			return num%2 == 0
		}, 100)

		if !slices.Equal(result, largeExpected) {
			t.Errorf("Expected result to be %#v. Got %#v", largeExpected, result)
		}
	})
}
