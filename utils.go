package gotiny

import (
	"encoding"
	"encoding/gob"
	"reflect"
	"unsafe"
)

func floatToUint(f float64) uint64 {
	return reverseByte(*((*uint64)(unsafe.Pointer(&f))))
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
