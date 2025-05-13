package utils

import (
	"testing"
	"unsafe"
)

func TestSortByAddress(t *testing.T) {
	type sample struct {
		value int
	}

	x := &sample{value: 1}
	y := &sample{value: 2}

	xAddr := uintptr(unsafe.Pointer(x))
	yAddr := uintptr(unsafe.Pointer(y))

	if xAddr == yAddr {
		t.Fatalf("unexpected: x and y have the same address (%p)", x)
	}

	a1, b1 := SortByAddress(x, y)
	if uintptr(unsafe.Pointer(a1)) > uintptr(unsafe.Pointer(b1)) {
		t.Errorf("SortByAddress(x, y) returned addresses out of order: %p > %p", a1, b1)
	}

	a2, b2 := SortByAddress(y, x)
	if uintptr(unsafe.Pointer(a2)) > uintptr(unsafe.Pointer(b2)) {
		t.Errorf("SortByAddress(y, x) returned addresses out of order: %p > %p", a2, b2)
	}

	if a1 != a2 || b1 != b2 {
		t.Errorf("SortByAddress(x, y) and SortByAddress(y, x) gave inconsistent results:\n  (a1, b1) = (%p, %p)\n  (a2, b2) = (%p, %p)", a1, b1, a2, b2)
	}
}
