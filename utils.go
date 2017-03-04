package gotiny

import (
	"encoding"
	"encoding/gob"
	"reflect"
	"unsafe"
)

const ptrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

type refVal struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
	uintptr
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

//int -5 -4 -3 -2 -1 0 1 2 3 4 5  6
//uint 9  7  5  3  1 0 2 4 6 8 10 12
func intToUint(v int64) uint64 {
	return uint64((v << 1) ^ (v >> 63))
}

//uint 9  7  5  3  1 0 2 4 6 8 10 12
//int -5 -4 -3 -2 -1 0 1 2 3 4 5  6
func uintToInt(u uint64) int64 {
	v := int64(u)
	return (-(v & 1)) ^ (v>>1)&0x7FFFFFFFFFFFFFFF
}

//func implementsInterface(rt reflect.Type) (reflect.Value, reflect.Value, bool) {
//	encm, has := rt.MethodByName("GobEncode")
//	decm, has2 := reflect.PtrTo(rt).MethodByName("GobDecode")
//	if has && has2 {
//		return encm.Func ,decm.Func,true
//	}
//	encm, has = rt.MethodByName("MarshalBinary")
//	decm, has2 = reflect.PtrTo(rt).MethodByName("UnmarshalBinary")
//	return encm.Func ,decm.Func, has&&has2
//}

func implementsGob(rt reflect.Type) (func(gob.GobEncoder) ([]byte, error), func(gob.GobDecoder, []byte) error, bool) {
	_, has := rt.MethodByName("GobEncode")
	_, has2 := reflect.PtrTo(rt).MethodByName("GobDecode")
	return gob.GobEncoder.GobEncode, gob.GobDecoder.GobDecode, has && has2
}
func implementsBin(rt reflect.Type) (func(encoding.BinaryMarshaler) ([]byte, error), func(encoding.BinaryUnmarshaler, []byte) error, bool) {
	_, has := rt.MethodByName("MarshalBinary")
	_, has2 := reflect.PtrTo(rt).MethodByName("UnmarshalBinary")
	return encoding.BinaryMarshaler.MarshalBinary, encoding.BinaryUnmarshaler.UnmarshalBinary, has && has2
}

func implementsGotiny(rt reflect.Type) (func(GoTinySerializer, []byte) []byte, func(GoTinySerializer, []byte) int, bool) {
	_, has := rt.MethodByName("GotinyEncode")
	_, has2 := rt.MethodByName("GotinyDecode")
	return GoTinySerializer.GotinyEncode, GoTinySerializer.GotinyDecode, has && has2
}

func isNil(p unsafe.Pointer) bool {
	return *(*unsafe.Pointer)(p) == nil
}

func elem(p unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(p)
}

//只应该有指针来实现该接口
type GoTinySerializer interface {
	//编码方法，将对象的序列化结果append到入参数并返回，方法不应该操纵入参数值原有的值
	GotinyEncode([]byte) []byte
	//解码方法，将入参解码到对象里并返回使用的长度。方法从入参的第0个字节开始使用，并且不应该修改入参中的任何数据
	GotinyDecode([]byte) int
}
