package ebus

import "unsafe"

//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

//go:inline
func noescapeVal[T any](p *T) *T {
	return (*T)(noescape(unsafe.Pointer(p)))
}
