// Package ring provides a generic ring buffer (circular buffer) implementation with thread-safety.
package ring

import (
	"sync"

	"github.com/PsionicAlch/byteforge/internal/datastructs/buffers/ring"
)

// SyncRingBuffer is a generic dynamically resizable circular buffer
// with thread-safety. It supports enqueue and dequeue operations in
// constant amortized time, and grows or shrinks based on usage to
// optimize memory consumption.
//
// T represents the type of elements stored in the buffer.
type SyncRingBuffer[T any] struct {
	buffer *ring.InternalRingBuffer[T]
	mu     sync.RWMutex
}

// SyncNew returns a new SyncRingBuffer with an optional initial capacity.
// If no capacity is provided or the provided value is <= 0, a default of 8 is used.
func NewSync[T any](capacity ...int) *SyncRingBuffer[T] {
	return &SyncRingBuffer[T]{
		buffer: ring.New[T](capacity...),
	}
}

// SyncFromSlice creates a new SyncRingBuffer from a given slice.
// An optional capacity may be provided. If the capacity is less than the slice length,
// the slice length is used as the minimum capacity.
func SyncFromSlice[T any, A ~[]T](s A, capacity ...int) *SyncRingBuffer[T] {
	return &SyncRingBuffer[T]{
		buffer: ring.FromSlice(s, capacity...),
	}
}

// SyncFromRingBuffer creates a new SyncRingBuffer from a given RingBuffer.
// This results in a deep copy so the underlying buffer won't be connected
// to the original RingBuffer.
func SyncFromRingBuffer[T any](src *RingBuffer[T]) *SyncRingBuffer[T] {
	return &SyncRingBuffer[T]{
		buffer: src.buffer.Clone(),
	}
}

// Len returns the number of elements currently stored in the buffer.
func (rb *SyncRingBuffer[T]) Len() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.buffer.Len()
}

// Cap returns the total capacity of the buffer.
func (rb *SyncRingBuffer[T]) Cap() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.buffer.Cap()
}

// IsEmpty returns true if the buffer contains no elements.
func (rb *SyncRingBuffer[T]) IsEmpty() bool {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.buffer.IsEmpty()
}

// Enqueue appends one or more values to the end of the buffer.
// If necessary, the buffer is resized to accommodate the new values.
func (rb *SyncRingBuffer[T]) Enqueue(values ...T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.buffer.Enqueue(values...)
}

// Dequeue removes and returns the element at the front of the buffer.
// If the buffer is empty, it returns the zero value of T and false.
// The buffer may shrink if usage falls below 25% of capacity.
func (rb *SyncRingBuffer[T]) Dequeue() (T, bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	return rb.buffer.Dequeue()
}

// Peek returns the element at the front of the buffer without removing it.
// If the buffer is empty, it returns the zero value of T and false.
func (rb *SyncRingBuffer[T]) Peek() (T, bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	return rb.buffer.Peek()
}

// ToSlice returns a new slice containing all elements in the buffer in their logical order.
// The returned slice is independent of the internal buffer state.
func (rb *SyncRingBuffer[T]) ToSlice() []T {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	return rb.buffer.ToSlice()
}

// Clone creates a deep copy of the source SyncRingBuffer.
func (rb *SyncRingBuffer[T]) Clone() *SyncRingBuffer[T] {
	return &SyncRingBuffer[T]{
		buffer: rb.buffer.Clone(),
	}
}
