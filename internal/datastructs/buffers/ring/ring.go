// Package ring provides a generic ring buffer (circular buffer) implementation.
// This internal version serves as the core data structure for FIFO queues and other public-facing abstractions.
// It supports dynamic resizing and is optimized for enqueue/dequeue performance without relying on third-party libraries.
package ring

import "slices"

// InternalRingBuffer is a generic dynamically resizable circular buffer.
// It supports enqueue and dequeue operations in constant amortized time,
// and grows or shrinks based on usage to optimize memory consumption.
//
// T represents the type of elements stored in the buffer.
type InternalRingBuffer[T any] struct {
	data       []T
	head, tail int
	size       int
	capacity   int
}

// New returns a new InternalRingBuffer with an optional initial capacity.
// If no capacity is provided or the provided value is <= 0, a default of 8 is used.
func New[T any](capacity ...int) *InternalRingBuffer[T] {
	cap := 8
	if len(capacity) > 0 && capacity[0] > 0 {
		cap = capacity[0]
	}

	return &InternalRingBuffer[T]{
		data:     make([]T, cap),
		capacity: cap,
	}
}

// FromSlice creates a new InternalRingBuffer from a given slice.
// An optional capacity may be provided. If the capacity is less than the slice length,
// the slice length is used as the minimum capacity.
func FromSlice[T any, A ~[]T](s A, capacity ...int) *InternalRingBuffer[T] {
	desiredCapacity := 8

	if len(capacity) > 0 && capacity[0] > desiredCapacity {
		desiredCapacity = capacity[0]
	} else if len(s) > 0 {
		desiredCapacity = len(s)
	}

	var data []T

	if desiredCapacity > len(s) {
		data = make([]T, desiredCapacity)
		for i := 0; i < len(s); i++ {
			data[i] = s[i]
		}
	} else {
		data = slices.Clone(s)
	}

	return &InternalRingBuffer[T]{
		data:     data,
		capacity: desiredCapacity,
		tail:     len(s),
		size:     len(s),
	}
}

// Len returns the number of elements currently stored in the buffer.
func (rb *InternalRingBuffer[T]) Len() int {
	return rb.size
}

// Cap returns the total capacity of the buffer.
func (rb *InternalRingBuffer[T]) Cap() int {
	return rb.capacity
}

// IsEmpty returns true if the buffer contains no elements.
func (rb *InternalRingBuffer[T]) IsEmpty() bool {
	return rb.size == 0
}

// Enqueue appends one or more values to the end of the buffer.
// If necessary, the buffer is resized to accommodate the new values.
func (rb *InternalRingBuffer[T]) Enqueue(values ...T) {
	required := rb.size + len(values)
	if required > rb.capacity {
		newCap := rb.capacity * 2
		for newCap < required {
			newCap *= 2
		}

		rb.resize(newCap)
	}

	for _, value := range values {
		rb.data[rb.tail] = value
		rb.tail = (rb.tail + 1) % rb.capacity
		rb.size++
	}
}

// Dequeue removes and returns the element at the front of the buffer.
// If the buffer is empty, it returns the zero value of T and false.
// The buffer may shrink if usage falls below 25% of capacity.
func (rb *InternalRingBuffer[T]) Dequeue() (T, bool) {
	var zero T
	if rb.size == 0 {
		return zero, false
	}

	val := rb.data[rb.head]
	rb.head = (rb.head + 1) % rb.capacity
	rb.size--

	if rb.capacity > 1 && rb.size <= rb.capacity/4 {
		rb.resize(rb.capacity / 2)
	}

	return val, true
}

// Peek returns the element at the front of the buffer without removing it.
// If the buffer is empty, it returns the zero value of T and false.
func (rb *InternalRingBuffer[T]) Peek() (T, bool) {
	var zero T
	if rb.size == 0 {
		return zero, false
	}

	return rb.data[rb.head], true
}

// ToSlice returns a new slice containing all elements in the buffer in their logical order.
// The returned slice is independent of the internal buffer state.
func (rb *InternalRingBuffer[T]) ToSlice() []T {
	if rb.size == 0 {
		return make([]T, 0)
	}

	result := make([]T, rb.size)
	for i := 0; i < rb.size; i++ {
		result[i] = rb.data[(rb.head+i)%rb.capacity]
	}

	return result
}

// Clone creates a deep copy of the source InternalRingBuffer.
func (rb *InternalRingBuffer[T]) Clone() *InternalRingBuffer[T] {
	newData := make([]T, rb.capacity)
	for i := 0; i < rb.size; i++ {
		newData[i] = rb.data[(rb.head+i)%rb.capacity]
	}

	return &InternalRingBuffer[T]{
		data:     newData,
		head:     0,
		tail:     rb.size,
		size:     rb.size,
		capacity: rb.capacity,
	}
}

// resize adjusts the capacity of the buffer to the specified value,
// reordering the contents so that head = 0 and tail = size.
func (rb *InternalRingBuffer[T]) resize(newCap int) {
	newData := make([]T, newCap)
	for i := 0; i < rb.size; i++ {
		newData[i] = rb.data[(rb.head+i)%rb.capacity]
	}

	rb.data = newData
	rb.head = 0
	rb.tail = rb.size
	rb.capacity = newCap
}
