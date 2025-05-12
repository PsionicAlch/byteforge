// Package ring provides a generic ring buffer (circular buffer) implementation.
package ring

import "github.com/PsionicAlch/byteforge/internal/datastructs/buffers/ring"

// RingBuffer is a generic dynamically resizable circular buffer.
// It supports enqueue and dequeue operations in constant amortized time,
// and grows or shrinks based on usage to optimize memory consumption.
//
// T represents the type of elements stored in the buffer.
type RingBuffer[T any] struct {
	buffer *ring.InternalRingBuffer[T]
}

// New returns a new RingBuffer with an optional initial capacity.
// If no capacity is provided or the provided value is <= 0, a default of 8 is used.
func New[T any](capacity ...int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: ring.New[T](capacity...),
	}
}

// FromSlice creates a new RingBuffer from a given slice.
// An optional capacity may be provided. If the capacity is less than the slice length,
// the slice length is used as the minimum capacity.
func FromSlice[T any, A ~[]T](s A, capacity ...int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: ring.FromSlice(s, capacity...),
	}
}

// FromSyncRingBuffer creates a new RingBuffer from a given SyncRingBuffer.
// This results in a deep copy so the underlying buffer won't be connected
// to the original SyncRingBuffer.
func FromSyncRingBuffer[T any](src *SyncRingBuffer[T]) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: src.buffer.Clone(),
	}
}

// Len returns the number of elements currently stored in the buffer.
func (rb *RingBuffer[T]) Len() int {
	return rb.buffer.Len()
}

// Cap returns the total capacity of the buffer.
func (rb *RingBuffer[T]) Cap() int {
	return rb.buffer.Cap()
}

// IsEmpty returns true if the buffer contains no elements.
func (rb *RingBuffer[T]) IsEmpty() bool {
	return rb.buffer.IsEmpty()
}

// Enqueue appends one or more values to the end of the buffer.
// If necessary, the buffer is resized to accommodate the new values.
func (rb *RingBuffer[T]) Enqueue(values ...T) {
	rb.buffer.Enqueue(values...)
}

// Dequeue removes and returns the element at the front of the buffer.
// If the buffer is empty, it returns the zero value of T and false.
// The buffer may shrink if usage falls below 25% of capacity.
func (rb *RingBuffer[T]) Dequeue() (T, bool) {
	return rb.buffer.Dequeue()
}

// Peek returns the element at the front of the buffer without removing it.
// If the buffer is empty, it returns the zero value of T and false.
func (rb *RingBuffer[T]) Peek() (T, bool) {
	return rb.buffer.Peek()
}

// ToSlice returns a new slice containing all elements in the buffer in their logical order.
// The returned slice is independent of the internal buffer state.
func (rb *RingBuffer[T]) ToSlice() []T {
	return rb.buffer.ToSlice()
}

// Clone creates a deep copy of the source RingBuffer.
func (rb *RingBuffer[T]) Clone() *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: rb.buffer.Clone(),
	}
}
