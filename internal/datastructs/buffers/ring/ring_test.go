package ring

import (
	"slices"
	"testing"
)

func TestInternalRingBuffer_New(t *testing.T) {
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
			buf := New[int](scenario.capacity...)

			if buf.size != scenario.expectedSize {
				t.Errorf("Expected buffer's size to be %d. Got %d.", scenario.expectedSize, buf.size)
			}

			if buf.capacity != scenario.expectedCap {
				t.Errorf("Expected buffer's capacity to be %d. Got %d.", scenario.expectedCap, buf.capacity)
			}
		})
	}
}

func TestInternalRingBuffer_FromSlice(t *testing.T) {
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
			buf := FromSlice(scenario.slice, scenario.capacity...)

			if buf.size != scenario.expectedSize {
				t.Errorf("Expected buffer's size to be %d. Got %d.", scenario.expectedSize, buf.size)
			}

			if buf.capacity != scenario.expectedCap {
				t.Errorf("Expected buffer's capacity to be %d. Got %d.", scenario.expectedCap, buf.capacity)
			}

			for _, item := range scenario.slice {
				if !slices.Contains(buf.data, item) {
					t.Errorf("Expected buffer to contain %d", item)
				}
			}
		})
	}
}

func TestInternalRingBuffer_Len(t *testing.T) {
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

			if len(scenario.data) > 0 {
				buf.Enqueue(scenario.data...)
			}

			if buf.Len() != scenario.expectedLen {
				t.Errorf("Expected buf.Len() to be %d. Got %d", scenario.expectedLen, buf.Len())
			}
		})
	}
}

func TestInternalRingBuffer_Cap(t *testing.T) {
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

			if buf.Cap() != scenario.expectedCapacity {
				t.Errorf("Expected buf.Cap() to be %d. Got %d", scenario.expectedCapacity, buf.Cap())
			}
		})
	}
}

func TestInternalRingBuffer_IsEmpty(t *testing.T) {
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
			buf := New[int](scenario.desiredCapacity...)
			buf.Enqueue(scenario.data...)

			if buf.IsEmpty() != scenario.isEmpty {
				t.Errorf("Expected buf.IsEmpty() to be %t. Got %t", scenario.isEmpty, buf.IsEmpty())
			}
		})
	}
}

func TestInternalRingBuffer_Enqueue(t *testing.T) {
	scenarios := []struct {
		name               string
		initial            []int
		enqueue            []int
		expectedOrder      []int
		expectedSize       int
		expectedCapAtLeast int
	}{
		{
			name:               "Append within capacity",
			initial:            []int{1, 2},
			enqueue:            []int{3, 4},
			expectedOrder:      []int{1, 2, 3, 4},
			expectedSize:       4,
			expectedCapAtLeast: 4,
		},
		{
			name:               "Trigger resize on append",
			initial:            []int{1, 2, 3, 4},
			enqueue:            []int{5, 6, 7},
			expectedOrder:      []int{1, 2, 3, 4, 5, 6, 7},
			expectedSize:       7,
			expectedCapAtLeast: 7,
		},
		{
			name:               "Append nothing (no-op)",
			initial:            []int{9, 8, 7},
			enqueue:            []int{},
			expectedOrder:      []int{9, 8, 7},
			expectedSize:       3,
			expectedCapAtLeast: 3,
		},
		{
			name:               "Append to empty buffer",
			initial:            []int{},
			enqueue:            []int{10, 20, 30},
			expectedOrder:      []int{10, 20, 30},
			expectedSize:       3,
			expectedCapAtLeast: 3,
		},
		{
			name:               "Large append exceeds multiple resizes",
			initial:            []int{},
			enqueue:            makeRange(1, 50),
			expectedOrder:      makeRange(1, 50),
			expectedSize:       50,
			expectedCapAtLeast: 50,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buf := New[int](len(scenario.initial))
			buf.Enqueue(scenario.initial...)
			buf.Enqueue(scenario.enqueue...)

			if buf.Len() != scenario.expectedSize {
				t.Errorf("Expected size %d, got %d", scenario.expectedSize, buf.Len())
			}

			if buf.Cap() < scenario.expectedCapAtLeast {
				t.Errorf("Expected capacity >= %d, got %d", scenario.expectedCapAtLeast, buf.Cap())
			}

			for _, want := range scenario.expectedOrder {
				got, ok := buf.Dequeue()

				if !ok {
					t.Fatalf("Expected to dequeue %d but got nothing", want)
				}

				if got != want {
					t.Errorf("Expected dequeue value %d, got %d", want, got)
				}
			}

			if !buf.IsEmpty() {
				t.Errorf("Expected buffer to be empty after all dequeues")
			}
		})
	}
}

func TestInternalRingBuffer_Dequeue(t *testing.T) {
	scenarios := []struct {
		name              string
		initial           []int
		expectedDequeues  []int
		expectedEmpty     bool
		expectedCapAtMost int
	}{
		{
			name:              "Dequeue from empty buffer",
			initial:           []int{},
			expectedDequeues:  []int{},
			expectedEmpty:     true,
			expectedCapAtMost: 8,
		},
		{
			name:              "Dequeue single element",
			initial:           []int{42},
			expectedDequeues:  []int{42},
			expectedEmpty:     true,
			expectedCapAtMost: 8,
		},
		{
			name:              "Dequeue multiple elements",
			initial:           []int{1, 2, 3, 4},
			expectedDequeues:  []int{1, 2, 3, 4},
			expectedEmpty:     true,
			expectedCapAtMost: 8,
		},
		{
			name:              "Dequeue triggers downsize",
			initial:           makeRange(1, 16),
			expectedDequeues:  makeRange(1, 14),
			expectedEmpty:     false,
			expectedCapAtMost: 16,
		},
		{
			name:              "Dequeue wraparound case",
			initial:           []int{10, 20, 30},
			expectedDequeues:  []int{10, 20, 30},
			expectedEmpty:     true,
			expectedCapAtMost: 8,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buf := New[int]()
			buf.Enqueue(scenario.initial...)

			for i, expected := range scenario.expectedDequeues {
				val, ok := buf.Dequeue()

				if !ok {
					t.Fatalf("Expected Dequeue #%d to return value %d, got nothing", i, expected)
				}

				if val != expected {
					t.Errorf("Dequeue #%d: expected %d, got %d", i, expected, val)
				}
			}

			if buf.IsEmpty() != scenario.expectedEmpty {
				t.Errorf("Expected IsEmpty to be %v, got %v", scenario.expectedEmpty, buf.IsEmpty())
			}

			if buf.Cap() > scenario.expectedCapAtMost {
				t.Errorf("Expected capacity <= %d, got %d", scenario.expectedCapAtMost, buf.Cap())
			}
		})
	}
}

func TestInternalRingBuffer_Peek(t *testing.T) {
	scenarios := []struct {
		name         string
		initial      []int
		expectedPeek int
		expectOK     bool
		expectSize   int
	}{
		{
			name:         "Peek empty buffer",
			initial:      []int{},
			expectedPeek: 0,
			expectOK:     false,
			expectSize:   0,
		},
		{
			name:         "Peek one item",
			initial:      []int{7},
			expectedPeek: 7,
			expectOK:     true,
			expectSize:   1,
		},
		{
			name:         "Peek multiple items",
			initial:      []int{3, 4, 5},
			expectedPeek: 3,
			expectOK:     true,
			expectSize:   3,
		},
		{
			name:         "Peek after internal wraparound",
			initial:      makeRange(1, 10),
			expectedPeek: 1,
			expectOK:     true,
			expectSize:   10,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buf := New[int]()
			buf.Enqueue(scenario.initial...)

			val, ok := buf.Peek()

			if ok != scenario.expectOK {
				t.Fatalf("Expected Peek ok=%v, got %v", scenario.expectOK, ok)
			}

			if ok && val != scenario.expectedPeek {
				t.Errorf("Expected Peek value %d, got %d", scenario.expectedPeek, val)
			}

			if buf.Len() != scenario.expectSize {
				t.Errorf("Expected buffer size %d after Peek, got %d", scenario.expectSize, buf.Len())
			}
		})
	}
}

func TestInternalRingBuffer_ToSlice(t *testing.T) {
	scenarios := []struct {
		name           string
		setup          func() *InternalRingBuffer[int]
		expectedSlice  []int
		expectedSize   int
		expectedBuffer []int
	}{
		{
			name: "Empty buffer",
			setup: func() *InternalRingBuffer[int] {
				return New[int]()
			},
			expectedSlice:  []int{},
			expectedSize:   0,
			expectedBuffer: []int{},
		},
		{
			name: "Buffer with single element",
			setup: func() *InternalRingBuffer[int] {
				buf := New[int]()
				buf.Enqueue(42)
				return buf
			},
			expectedSlice:  []int{42},
			expectedSize:   1,
			expectedBuffer: []int{42},
		},
		{
			name: "Buffer with multiple elements",
			setup: func() *InternalRingBuffer[int] {
				buf := New[int]()
				buf.Enqueue(1, 2, 3, 4, 5)
				return buf
			},
			expectedSlice:  []int{1, 2, 3, 4, 5},
			expectedSize:   5,
			expectedBuffer: []int{1, 2, 3, 4, 5},
		},
		{
			name: "Buffer with wrapped elements (head > 0)",
			setup: func() *InternalRingBuffer[int] {
				buf := New[int](5)
				buf.Enqueue(1, 2, 3, 4, 5)
				buf.Dequeue()
				buf.Dequeue()
				buf.Enqueue(6, 7)
				return buf
			},
			expectedSlice:  []int{3, 4, 5, 6, 7},
			expectedSize:   5,
			expectedBuffer: []int{3, 4, 5, 6, 7},
		},
		{
			name: "Buffer after resize",
			setup: func() *InternalRingBuffer[int] {
				buf := New[int](4)
				buf.Enqueue(1, 2, 3, 4)
				buf.Enqueue(5)
				return buf
			},
			expectedSlice:  []int{1, 2, 3, 4, 5},
			expectedSize:   5,
			expectedBuffer: []int{1, 2, 3, 4, 5},
		},
		{
			name: "Buffer after many enqueue/dequeue operations",
			setup: func() *InternalRingBuffer[int] {
				buf := New[int](3)
				buf.Enqueue(1, 2, 3)
				buf.Dequeue()
				buf.Dequeue()
				buf.Enqueue(4, 5)
				buf.Dequeue()
				buf.Enqueue(6)
				return buf
			},
			expectedSlice:  []int{4, 5, 6},
			expectedSize:   3,
			expectedBuffer: []int{4, 5, 6},
		},
		{
			name: "FromSlice initialization",
			setup: func() *InternalRingBuffer[int] {
				return FromSlice([]int{10, 20, 30, 40})
			},
			expectedSlice:  []int{10, 20, 30, 40},
			expectedSize:   4,
			expectedBuffer: []int{10, 20, 30, 40},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buf := scenario.setup()

			if buf.Len() != scenario.expectedSize {
				t.Errorf("Before ToSlice: Expected size %d, got %d", scenario.expectedSize, buf.Len())
			}

			result := buf.ToSlice()

			if len(result) != len(scenario.expectedSlice) {
				t.Errorf("ToSlice result length: Expected %d, got %d",
					len(scenario.expectedSlice), len(result))
			}

			for i, expected := range scenario.expectedSlice {
				if i >= len(result) {
					t.Errorf("Missing expected element at index %d: %v", i, expected)
					continue
				}
				if result[i] != expected {
					t.Errorf("Element mismatch at index %d: Expected %v, got %v",
						i, expected, result[i])
				}
			}

			bufferSlice := []int{}
			originalSize := buf.Len()
			for i := 0; i < originalSize; i++ {
				val, ok := buf.Dequeue()
				if !ok {
					t.Fatalf("Failed to dequeue element %d", i)
				}
				bufferSlice = append(bufferSlice, val)
			}

			if len(bufferSlice) != len(scenario.expectedBuffer) {
				t.Errorf("After ToSlice - Buffer size changed: Expected %d, got %d",
					len(scenario.expectedBuffer), len(bufferSlice))
			}

			for i, expected := range scenario.expectedBuffer {
				if i >= len(bufferSlice) {
					t.Errorf("After ToSlice - Missing expected element at index %d: %v", i, expected)
					continue
				}
				if bufferSlice[i] != expected {
					t.Errorf("After ToSlice - Element mismatch at index %d: Expected %v, got %v",
						i, expected, bufferSlice[i])
				}
			}
		})
	}
}

func TestInternalRingBuffer_Clone(t *testing.T) {
	src := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	dst := src.Clone()

	if !slices.Equal(src.ToSlice(), dst.ToSlice()) {
		t.Error("Expected src and destination to have the same underlying data.")
	}

	src.Enqueue(11, 12, 13, 14, 15, 16, 17, 18, 19, 20)

	if slices.Equal(src.ToSlice(), dst.ToSlice()) {
		t.Error("Expected src and destination to have the varying underlying data.")
	}

	dstSlice := dst.ToSlice()
	for _, num := range []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20} {
		if slices.Contains(dstSlice, num) {
			t.Errorf("dst slice was not supposed to contain %d", num)
		}
	}
}

func TestInternalRingBuffer_resize(t *testing.T) {
	scenarios := []struct {
		name         string
		setup        func() *InternalRingBuffer[int]
		resizeTo     int
		expectedData []int
	}{
		{
			name: "Resize larger with contiguous data",
			setup: func() *InternalRingBuffer[int] {
				buf := New[int]()
				buf.Enqueue(1, 2, 3)
				return buf
			},
			resizeTo:     8,
			expectedData: []int{1, 2, 3},
		},
		{
			name: "Resize with wrapped data",
			setup: func() *InternalRingBuffer[int] {
				buf := New[int]()
				buf.Enqueue(1, 2, 3, 4)
				_, _ = buf.Dequeue() // Remove 1
				_, _ = buf.Dequeue() // Remove 2
				buf.Enqueue(5, 6)    // Wrap around
				return buf
			},
			resizeTo:     10,
			expectedData: []int{3, 4, 5, 6},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buf := scenario.setup()
			buf.resize(scenario.resizeTo)

			if buf.capacity != scenario.resizeTo {
				t.Errorf("Expected capacity %d, got %d", scenario.resizeTo, buf.capacity)
			}
			if buf.head != 0 {
				t.Errorf("Expected head to be reset to 0, got %d", buf.head)
			}
			if buf.tail != buf.size {
				t.Errorf("Expected tail == size (%d), got %d", buf.size, buf.tail)
			}
			if buf.size != len(scenario.expectedData) {
				t.Fatalf("Expected size %d, got %d", len(scenario.expectedData), buf.size)
			}

			for i := 0; i < buf.size; i++ {
				actual := buf.data[i]
				expected := scenario.expectedData[i]
				if actual != expected {
					t.Errorf("At index %d: expected %d, got %d", i, expected, actual)
				}
			}
		})
	}
}

func makeRange(start, end int) []int {
	out := make([]int, end-start+1)
	for i := range out {
		out[i] = start + i
	}

	return out
}
