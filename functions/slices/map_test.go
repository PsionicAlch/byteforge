package slices

import (
	"slices"
	"strconv"
	"testing"

	islices "github.com/PsionicAlch/byteforge/internal/functions/slices"
)

func TestMap(t *testing.T) {
	const max = 1000000
	largeArr := islices.ERange(0, max)

	largeExpected := make([]int, max)
	for i := 0; i < max; i++ {
		largeExpected[i] = i * 2
	}

	t.Run("Map from int to int", func(t *testing.T) {
		result := Map([]int{0, 1, 2, 3, 4, 5}, func(num int) int {
			return num * 2
		})
		expected := []int{0, 2, 4, 6, 8, 10}

		if !slices.Equal(result, expected) {
			t.Errorf("Expected result to be %#v. Got %#v", expected, result)
		}
	})

	t.Run("Map from int to string", func(t *testing.T) {
		result := Map([]int{0, 1, 2, 3, 4, 5}, func(num int) string {
			return strconv.Itoa(num)
		})
		expected := []string{"0", "1", "2", "3", "4", "5"}

		if !slices.Equal(result, expected) {
			t.Errorf("Expected result to be %#v. Got %#v", expected, result)
		}
	})

	t.Run("Map with empty slice", func(t *testing.T) {
		result := Map([]int{}, func(num int) int {
			return num
		})
		expected := []int{}

		if !slices.Equal(result, expected) {
			t.Errorf("Expected result to be %#v. Got %#v", expected, result)
		}
	})

	t.Run("Map with huge slice", func(t *testing.T) {
		result := Map(largeArr, func(num int) int {
			return num * 2
		})

		if !slices.Equal(result, largeExpected) {
			t.Errorf("Expected result to be %#v. Got %#v", largeExpected, result)
		}
	})
}

func TestParallelMap(t *testing.T) {
	const max = 1000000
	largeArr := islices.ERange(0, max)

	largeExpected := make([]int, max)
	for i := 0; i < max; i++ {
		largeExpected[i] = i * 2
	}

	t.Run("Map from int to int", func(t *testing.T) {
		result := ParallelMap([]int{0, 1, 2, 3, 4, 5}, func(num int) int {
			return num * 2
		})
		expected := []int{0, 2, 4, 6, 8, 10}

		if !slices.Equal(result, expected) {
			t.Errorf("Expected result to be %#v. Got %#v", expected, result)
		}
	})

	t.Run("Map from int to string", func(t *testing.T) {
		result := ParallelMap([]int{0, 1, 2, 3, 4, 5}, func(num int) string {
			return strconv.Itoa(num)
		})
		expected := []string{"0", "1", "2", "3", "4", "5"}

		if !slices.Equal(result, expected) {
			t.Errorf("Expected result to be %#v. Got %#v", expected, result)
		}
	})

	t.Run("Map with empty slice", func(t *testing.T) {
		result := ParallelMap([]int{}, func(num int) int {
			return num
		})
		expected := []int{}

		if !slices.Equal(result, expected) {
			t.Errorf("Expected result to be %#v. Got %#v", expected, result)
		}
	})

	t.Run("Map with huge slice", func(t *testing.T) {
		result := ParallelMap(largeArr, func(num int) int {
			return num * 2
		})

		if !slices.Equal(result, largeExpected) {
			t.Errorf("Expected result to be %#v. Got %#v", largeExpected, result)
		}
	})

	t.Run("Map with negative worker pool", func(t *testing.T) {
		result := ParallelMap([]int{0, 1, 2, 3, 4, 5}, func(num int) int {
			return num * 2
		}, -10)
		expected := []int{0, 2, 4, 6, 8, 10}

		if !slices.Equal(result, expected) {
			t.Errorf("Expected result to be %#v. Got %#v", expected, result)
		}
	})

	t.Run("Map with positive worker pool", func(t *testing.T) {
		result := ParallelMap(largeArr, func(num int) int {
			return num * 2
		}, 50)

		if !slices.Equal(result, largeExpected) {
			t.Errorf("Expected result to be %#v. Got %#v", largeExpected, result)
		}
	})
}
