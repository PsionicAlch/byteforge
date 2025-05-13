package ring

import (
	"slices"
	"sync"
	"testing"
)

func TestSyncRingBuffer_New(t *testing.T) {
	scenarios := []struct {
		name         string
		capacity     []int
		expectedSize int
		expectedCap  int
	}{
		{"Empty capacity", []int{}, 0, 8},
		{"Non-empty capacity", []int{5}, 0, 5},
		{"Negative capacity", []int{-10}, 0, 8},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buf := NewSync[int](scenario.capacity...)

			if buf == nil {
				t.Fatal("Expected buf to not be nil")
			}

			if buf.buffer == nil {
				t.Fatal("Expected buf.buffer to not be nil")
			}

			if buf.Len() != scenario.expectedSize {
				t.Errorf("Expected buffer's size to be %d. Got %d.", scenario.expectedSize, buf.Len())
			}

			if buf.Cap() != scenario.expectedCap {
				t.Errorf("Expected buffer's capacity to be %d. Got %d.", scenario.expectedCap, buf.Cap())
			}
		})
	}
}

func TestSyncRingBuffer_FromSlice(t *testing.T) {
	scenarios := []struct {
		name         string
		slice        []int
		capacity     []int
		expectedSize int
		expectedCap  int
	}{
		{"Empty slice and empty capacity", []int{}, []int{}, 0, 8},
		{"Non-empty slice and empty capacity", []int{1, 2, 3}, []int{}, 3, 3},
		{"Non-empty slice and non-empty capacity", []int{1, 2, 3}, []int{10}, 3, 10},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buf := SyncFromSlice(scenario.slice, scenario.capacity...)

			if buf == nil {
				t.Fatal("Expected buf to not be nil")
			}

			if buf.buffer == nil {
				t.Fatal("Expected buf.buffer to not be nil")
			}

			if buf.Len() != scenario.expectedSize {
				t.Errorf("Expected buffer's size to be %d. Got %d.", scenario.expectedSize, buf.Len())
			}

			if buf.Cap() != scenario.expectedCap {
				t.Errorf("Expected buffer's capacity to be %d. Got %d.", scenario.expectedCap, buf.Cap())
			}

			for _, item := range scenario.slice {
				if !slices.Contains(buf.buffer.ToSlice(), item) {
					t.Errorf("Expected buffer to contain %d", item)
				}
			}
		})
	}
}

func TestSyncRingBuffer_SyncFromRingBuffer(t *testing.T) {
	src := FromSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	dst := SyncFromRingBuffer(src)

	if dst == nil {
		t.Fatal("Expected dst to not be nil")
	}

	if dst.buffer == nil {
		t.Fatal("Expected dst.buffer to not be nil")
	}

	if dst.Len() != src.Len() {
		t.Error("Expected dst.Len() to be equal to src.Len()")
	}

	if dst.Cap() != src.Cap() {
		t.Error("Expected dst.Cap() to be equal to src.Cap()")
	}

	if !slices.Equal(src.ToSlice(), dst.ToSlice()) {
		t.Error("Expected src.ToSlice to be equal to dst.ToSlice.")
	}

	src.Enqueue(10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)

	if dst.Len() == src.Len() {
		t.Error("Did not expect dst.Len() to be equal to src.Len()")
	}

	if dst.Cap() == src.Cap() {
		t.Error("Did not expect dst.Cap() to be equal to src.Cap()")
	}

	if slices.Equal(src.ToSlice(), dst.ToSlice()) {
		t.Error("Did not expect src.ToSlice to be equal to dst.ToSlice.")
	}
}

func TestSyncRingBuffer_Len(t *testing.T) {
	scenarios := []struct {
		name        string
		data        []int
		expectedLen int
	}{
		{"Lenght with 0 items", nil, 0},
		{"Lenght with 3 items", []int{1, 2, 3}, 3},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buf := New[int]()
			buf.Enqueue(scenario.data...)

			var wg sync.WaitGroup

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					if buf.Len() != scenario.expectedLen {
						t.Errorf("Expected buf.Len() to be %d. Got %d", scenario.expectedLen, buf.Len())
					}
				}()
			}

			wg.Wait()
		})
	}
}

func TestSyncRingBuffer_Cap(t *testing.T) {
	scenarios := []struct {
		name             string
		data             []int
		desiredCapacity  []int
		expectedCapacity int
	}{
		{"Capacity with 0 items and no desired capacity", nil, nil, 8},
		{"Capacity with 0 items and desired capacity of 1", nil, []int{1}, 1},
		{"Capacity with 0 items and desired capacity of -1", nil, []int{-1}, 8},
		{"Capacity with 9 items and no desired capacity", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil, 16},
		{"Capacity with 5 items and desired capacity of 3", []int{1, 2, 3, 4, 5}, []int{3}, 6},
		{"Capacity with 9 items and desired capacity of -9", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{-9}, 16},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buf := New[int](scenario.desiredCapacity...)
			buf.Enqueue(scenario.data...)

			var wg sync.WaitGroup

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					if buf.Cap() != scenario.expectedCapacity {
						t.Errorf("Expected buf.Cap() to be %d. Got %d", scenario.expectedCapacity, buf.Cap())
					}
				}()
			}

			wg.Wait()
		})
	}
}

func TestSyncRingBuffer_IsEmpty(t *testing.T) {
	scenarios := []struct {
		name            string
		data            []int
		desiredCapacity []int
		isEmpty         bool
	}{
		{"IsEmpty with 0 items and no desired capacity", nil, nil, true},
		{"IsEmpty with 0 items and desired capacity of 1", nil, []int{1}, true},
		{"IsEmpty with 0 items and desired capacity of -1", nil, []int{-1}, true},
		{"IsEmpty with 9 items and no desired capacity", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil, false},
		{"IsEmpty with 5 items and desired capacity of 3", []int{1, 2, 3, 4, 5}, []int{3}, false},
		{"IsEmpty with 9 items and desired capacity of -9", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{-9}, false},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buf := NewSync[int](scenario.desiredCapacity...)
			buf.Enqueue(scenario.data...)

			var wg sync.WaitGroup

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					if buf.IsEmpty() != scenario.isEmpty {
						t.Errorf("Expected buf.IsEmpty() to be %t. Got %t", scenario.isEmpty, buf.IsEmpty())
					}
				}()
			}

			wg.Wait()
		})
	}
}

func TestSyncRingBuffer_Enqueue(t *testing.T) {
	const max = 1000

	buf := NewSync[int]()

	var wg sync.WaitGroup

	for i := 0; i < max; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			buf.Enqueue(i)
		}()
	}

	wg.Wait()

	if buf.Len() != max {
		t.Errorf("Expected buf.Len() to be %d. Got %d", max, buf.Len())
	}
}

func TestSyncRingBuffer_Dequeue(t *testing.T) {
	const max = 1000

	buf := NewSync[int]()

	var wg sync.WaitGroup

	for i := 1; i <= max; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			buf.Enqueue(i)
		}()
	}

	wg.Wait()

	for i := 1; i <= max; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			value, found := buf.Dequeue()

			if !found {
				t.Error("Expected to find value from buf.Dequeue")
			}

			if value < 1 || value > max {
				t.Errorf("Expected value to be in range [0, %d]. Got %d", max, value)
			}
		}()
	}
}

func TestSyncRingBuffer_Peek(t *testing.T) {
	buf := SyncFromSlice([]int{1, 2, 3, 4, 5})
	expectedValue := 1

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			value, found := buf.Peek()

			if !found {
				t.Error("Expected buf.Peek() to return value")
			}

			if value != expectedValue {
				t.Errorf("Expected to find %d. Got %d", expectedValue, value)
			}
		}()
	}

	wg.Wait()
}

func TestSyncRingBuffer_ToSlice(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	buf := SyncFromSlice(data)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if !slices.Equal(data, buf.ToSlice()) {
				t.Error("Expected buf.ToSlice() to be equal to data.")
			}
		}()
	}

	wg.Wait()
}

func TestSyncRingBuffer_Clone(t *testing.T) {
	src := SyncFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			dst := src.Clone()

			if !slices.Equal(src.ToSlice(), dst.ToSlice()) {
				t.Error("Expected src and destination to have the same underlying data.")
			}

			dst.Enqueue(11, 12, 13, 14, 15, 16, 17, 18, 19, 20)

			if slices.Equal(src.ToSlice(), dst.ToSlice()) {
				t.Error("Expected src and destination to have the varying underlying data.")
			}

			srcSlice := src.ToSlice()
			for _, num := range []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20} {
				if slices.Contains(srcSlice, num) {
					t.Errorf("srcSlice was not supposed to contain %d", num)
				}
			}
		}()
	}

	wg.Wait()
}
