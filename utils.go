package gotiny

import (
	"encoding"
	"encoding/gob"
	"reflect"
	"unsafe"
)

const ptr1Size = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

type refVal struct {
	typ  unsafe.Pointer
	ptr  unsafe.Pointer
	flag flag
}

type flag uintptr

//go:linkname flagIndir reflect.flagIndir
const flagIndir flag = 1 << 7

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}

// stringHeader is a safe version of StringHeader used within this package.
type stringHeader struct {
	data unsafe.Pointer
	len  int
}

func floatToUint(v float64) uint64 {
	return reverseByte(*(*uint64)(unsafe.Pointer(&v)))
}

func uintToFloat(u uint64) float64 {
	u = reverseByte(u)
	return *((*float64)(unsafe.Pointer(&u)))
}

func reverseByte(u uint64) uint64 {
	u = (u << 32) | (u >> 32)
	u = ((u << 16) & 0xFFFF0000FFFF0000) | ((u >> 16) & 0xFFFF0000FFFF)
	u = ((u << 8) & 0xFF00FF00FF00FF00) | ((u >> 8) & 0xFF00FF00FF00FF)
	return u
}

// int -5 -4 -3 -2 -1 0 1 2 3 4 5  6
// uint 9  7  5  3  1 0 2 4 6 8 10 12
func intToUint(v int64) uint64 {
	return uint64((v << 1) ^ (v >> 63))
}

// uint 9  7  5  3  1 0 2 4 6 8 10 12
// int -5 -4 -3 -2 -1 0 1 2 3 4 5  6
func uintToInt(u uint64) int64 {
	v := int64(u)
	return (-(v & 1)) ^ (v>>1)&0x7FFFFFFFFFFFFFFF
}

var (
	gobType  = [2]reflect.Type{reflect.TypeOf((*gob.GobEncoder)(nil)).Elem(), reflect.TypeOf((*gob.GobDecoder)(nil)).Elem()}
	binType  = [2]reflect.Type{reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem(), reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()}
	tinyType = reflect.TypeOf((*GoTinySerializer)(nil)).Elem()
)

func implementsInterface(rt reflect.Type, typ [2]reflect.Type) bool {
	return rt.Implements(typ[0]) && reflect.PtrTo(rt).Implements(typ[1])
}

func isNil(p unsafe.Pointer) bool {
	return *(*unsafe.Pointer)(p) == nil
}

// 只应该由指针来实现该接口
type GoTinySerializer interface {
	// 编码方法，将对象的序列化结果append到入参数并返回，方法不应该修改入参数值原有的值
	GotinyEncode([]byte) []byte
	// 解码方法，将入参解码到对象里并返回使用的长度。方法从入参的第0个字节开始使用，并且不应该修改入参中的任何数据
	GotinyDecode([]byte) int
}
