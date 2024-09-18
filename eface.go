package ebus

import "unsafe"

type eface struct {
	_type *rtype
	data  unsafe.Pointer
}

type rtype struct {
	_    uintptr // size of the type
	_    uintptr // number of bytes in the type that are pointers
	hash uint32  // hash of the type
}

func typeHash(v any) uint32 {
	ef := *(*eface)(unsafe.Pointer(&v))
	return ef._type.hash
}

func efaceData(v any) unsafe.Pointer {
	ef := *(*eface)(unsafe.Pointer(&v))
	return ef.data
}

//go:inline
func same(a, b any) bool {
	return efaceData(a) == efaceData(b)
}
