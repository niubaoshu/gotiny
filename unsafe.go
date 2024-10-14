package gotiny

import (
	"reflect"
	"unsafe"
)

type refVal struct {
	_    unsafe.Pointer
	ptr  unsafe.Pointer
	flag uintptr
}

const flagIndir uintptr = 1 << 7

func getUnsafePointer(rv reflect.Value) unsafe.Pointer {
	vv := (*refVal)(unsafe.Pointer(&rv))
	if vv.flag&flagIndir == 0 {
		return unsafe.Pointer(&vv.ptr)
	} else {
		return vv.ptr
	}
}

type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}
