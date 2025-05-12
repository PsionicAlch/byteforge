package utils

import "unsafe"

func SortByAddress[Q any, T *Q](a, b T) (T, T) {
	if uintptr(unsafe.Pointer(a)) < uintptr(unsafe.Pointer(b)) {
		return a, b
	}

	return b, a
}
