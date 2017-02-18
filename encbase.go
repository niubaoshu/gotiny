package gotiny

import (
	//"fmt"
	"reflect"
)

func (e *Encoder) encBool(v bool) {
	if e.boolBit == 0 {
		e.boolBit = 1
		e.boolPos = e.index
		e.index++
	}
	if v {
		e.buf[e.boolPos] |= e.boolBit
	}
	e.boolBit <<= 1
}

func (e *Encoder) varuint(v uint64) (i int) {
	buf, i := e.buf[e.index:], 0
	for v >= 0x80 {
		buf[i] = uint8(v) | 0x80
		i++
		v >>= 7
	}
	buf[i] = uint8(v)
	i++
	e.index += i
	return i
	//i = varuint(, v)
	//e.index += i
	//return
}
func (e *Encoder) encUint8(v uint64)            { e.buf[e.index] = uint8(v); e.index++ }
func (e *Encoder) encUint16(v uint64)           { e.reqLen -= byte2 - e.varuint(v)<<3 }
func (e *Encoder) encUint32(v uint64)           { e.reqLen -= byte4 - e.varuint(v)<<3 }
func (e *Encoder) encUint64(v uint64)           { e.reqLen -= byte8 - e.varuint(v)<<3 }
func (e *Encoder) encUint(v uint64)             { e.reqLen -= word1 - e.varuint(v)<<3 }
func (e *Encoder) encUintptr(v uint64)          { e.encUint(v) }
func (e *Encoder) encInt8(v int64)              { e.encUint8(uint64(v)) }
func (e *Encoder) encInt16(v int64)             { e.encUint16(intToUint(v)) }
func (e *Encoder) encInt32(v int64)             { e.encUint32(intToUint(v)) }
func (e *Encoder) encInt64(v int64)             { e.encUint64(intToUint(v)) }
func (e *Encoder) encInt(v int64)               { e.encUint(intToUint(v)) }
func (e *Encoder) encFloat32(v float64)         { e.encUint32(floatToUint(v)) }
func (e *Encoder) encFloat64(v float64)         { e.encUint64(floatToUint(v)) }
func (e *Encoder) encComplex64(v complex128)    { e.encFloat32(real(v)); e.encFloat32(imag(v)) }
func (e *Encoder) encComplex128(v complex128)   { e.encFloat64(real(v)); e.encFloat64(imag(v)) }
func encignore(e *Encoder, v reflect.Value)     {}
func encBool(e *Encoder, v reflect.Value)       { e.encBool(v.Bool()) }
func encUint8(e *Encoder, v reflect.Value)      { e.encUint8(v.Uint()) }
func encUint16(e *Encoder, v reflect.Value)     { e.encUint16(v.Uint()) }
func encUint32(e *Encoder, v reflect.Value)     { e.encUint32(v.Uint()) }
func encUint64(e *Encoder, v reflect.Value)     { e.encUint64(v.Uint()) }
func encUint(e *Encoder, v reflect.Value)       { e.encUint(v.Uint()) }
func encUintptr(e *Encoder, v reflect.Value)    { e.encUintptr(v.Uint()) }
func encInt8(e *Encoder, v reflect.Value)       { e.encInt8(v.Int()) }
func encInt16(e *Encoder, v reflect.Value)      { e.encInt16(v.Int()) }
func encInt32(e *Encoder, v reflect.Value)      { e.encInt32(v.Int()) }
func encInt64(e *Encoder, v reflect.Value)      { e.encInt64(v.Int()) }
func encInt(e *Encoder, v reflect.Value)        { e.encInt(v.Int()) }
func encFloat32(e *Encoder, v reflect.Value)    { e.encFloat32(v.Float()) }
func encFloat64(e *Encoder, v reflect.Value)    { e.encFloat64(v.Float()) }
func encComplex64(e *Encoder, v reflect.Value)  { e.encComplex64(v.Complex()) }
func encComplex128(e *Encoder, v reflect.Value) { e.encComplex128(v.Complex()) }
func encString(e *Encoder, v reflect.Value)     { e.encString(v.String()) }

func (e *Encoder) encString(v string) {
	l := len(v)
	e.encUint(uint64(l))
	l = l << 3
	e.reqLen += l
	e.append(l)
	e.index += copy(e.buf[e.index:], v)
}
