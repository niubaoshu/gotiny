package gotiny

// import (
// 	"reflect"
// 	"unsafe"
// )

// func (e *Encoder) Encodes(is ...interface{}) {
// 	for i, _ := range is {
// 		e.encode(i)
// 	}
// }

// func (e *Encoder) Encode(i interface{}) {
// 	switch v := i.(type) {
// 	case bool:
// 		e.encBool(v)
// 	case uint8:
// 		e.encUint8(v)
// 	case int8:
// 		e.encInt8(v)
// 	case uint, uint16, uint32, uint64, uintptr:
// 		e.encUint(uint64(v))
// 	case int, int16, int32, int64:
// 		e.encInt(int64(v))
// 	case float32, float64:
// 		e.encFloat(v)
// 	case complex64, complex128:
// 		e.encComplex(v)
// 	case string:
// 		e.encString(v)
// 	case []byte:
// 		e.encSliceByte(v)
// 	default:
// 		e.encOther(v)
// 	}
// }

// func (e *Encoder) encBool(v bool) {
// 	if e.boolBit == 0 {
// 		e.boolBit = 1
// 		e.boolPos = len(e.buff)
// 		e.buff = append(e.buff, 0)
// 	}
// 	if v {
// 		e.buff[e.boolPos] |= e.boolBit
// 	}
// 	e.boolBit <<= 1
// }

// func (e *Encoder) encSliceByte(v []byte) { e.encUint(uint64(len(v))); e.buff = append(e.buff, v...) }
// func (e *Encoder) encUint8(v uint8)      { e.buff = append(e.buff, v) }
// func (e *Encoder) encInt8(v int8)        { e.encUint8(uint8(v)) }

// func (e *Encoder) encUint(v uint64) {
// 	for v >= 0x80 {
// 		e.buff = append(e.buff, byte(v)|0x80)
// 		v >>= 7
// 	}
// 	e.buff = append(e.buff, byte(v))
// }

// //int -5 -4 -3 -2 -1 0 1 2 3 4 5 6
// //uint 9  7  5  3  1 0 2 4 6 8 10 12
// func (e *Encoder) encInt(v int64) {
// 	x := uint64(v) << 1
// 	if v < 0 {
// 		x = ^x
// 	}
// 	e.encUint(x)
// }

// func (e *Encoder) encFloat(v float64)      { e.encUint(floatBits(v)) }
// func (e *Encoder) encComplex(v complex128) { e.encFloat(real(v)); e.encFloat(imag(v)) }
// func (e *Encoder) encString(v string) {
// 	e.encUint(uint64(len(v)))
// 	e.buff = append(e.buff, []byte(v)...)
// }

// func floatBits(f float64) uint64 {
// 	u := *((*uint64)(unsafe.Pointer(&f)))
// 	var v uint64
// 	for i := 0; i < 8; i++ {
// 		v <<= 8
// 		v |= u & 0xFF
// 		u >>= 8
// 	}
// 	return v
// }

// func (e *Encoder) encOther(v interface{}) {
// 	e.encodeValue(reflect.ValueOf(v))
// }

// func (e *Encoder) EncodeValue(v reflect.Value) {
// 	switch v.Kind() {
// 	case reflect.Array:
// 		for i := 0; i < v.Len(); i++ {
// 			e.Encode(v.Index(i).Interface())
// 		}
// 	case reflect.Map:
// 		e.EncUint(uint64(v.Len()))
// 		keys := v.MapKeys()
// 		for _, key := range keys {
// 			e.Encode(key.Interface())
// 			e.Encode(v.MapIndex(key).Interface())
// 		}
// 	case reflect.Ptr:
// 		if v.IsNil() {
// 			panic("totiny: cannot encode nil pointer of type ")
// 		}
// 		e.Encode(v.Elem().Interface())
// 	case reflect.Slice:
// 		l := v.Len()
// 		e.EncUint(uint64(l))
// 		for i := 0; i < l; i++ {
// 			e.Encode(v.Index(i).Interface())
// 		}
// 	case reflect.Struct:
// 		vt := v.Type()
// 		for i := 0; i < v.NumField(); i++ {
// 			if vt.Field(i).PkgPath == "" { // vt.Field(i).PkgPath 等于 ""代表导出字段
// 				e.Encode(v.Field(i).Interface())
// 			}
// 		}
// 	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Invalid:
// 		//panic("暂不支持这些类型")
// 	}
// }
