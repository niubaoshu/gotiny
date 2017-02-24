package gotiny

import (
	"reflect"
)

func (e *Encoder) encBool(v bool) {
	if e.boolBit == 0 {
		e.boolBit = 1
		e.boolPos = len(e.buf)
		e.buf = append(e.buf, 0)
	}
	if v {
		e.buf[e.boolPos] |= e.boolBit
	}
	e.boolBit <<= 1
}

func (e *Encoder) encUint(v uint64) {
	buf := e.buf
	for v >= 0x80 {
		buf = append(buf, uint8(v)|0x80)
		v >>= 7
	}
	e.buf = append(buf, uint8(v))
}

func encignore(e *Encoder, v reflect.Value) {}
func encBool(e *Encoder, v reflect.Value)   { e.encBool(v.Bool()) }
func encUint(e *Encoder, v reflect.Value)   { e.encUint(v.Uint()) }
func encUint8(e *Encoder, v reflect.Value)  { e.buf = append(e.buf, uint8(v.Uint())) }
func encInt8(e *Encoder, v reflect.Value)   { e.buf = append(e.buf, uint8(v.Int())) }
func encInt(e *Encoder, v reflect.Value)    { e.encUint(intToUint(v.Int())) }
func encFloat(e *Encoder, v reflect.Value)  { e.encUint(floatToUint(v.Float())) }
func encComplex(e *Encoder, v reflect.Value) {
	c := v.Complex()
	e.encUint(floatToUint(real(c)))
	e.encUint(floatToUint(imag(c)))
}
func encString(e *Encoder, v reflect.Value) {
	s := v.String()
	e.encUint(uint64(len(s)))
	e.buf = append(e.buf, s...)
}

//func encTime(e *Encoder, v reflect.Value) {
//t := v.Interface().(time.Time)

//}
