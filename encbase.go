package gotiny

import (
	"unsafe"
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
func (e *Encoder) encLength(v int) { e.encUint(uint64(v)) }

func encignore(e *Encoder, p unsafe.Pointer)    {}
func encBool(e *Encoder, p unsafe.Pointer)      { e.encBool(*(*bool)(p)) }
func encInt16(e *Encoder, p unsafe.Pointer)     { e.encUint(intToUint(int64(*(*int16)(p)))) }
func encInt32(e *Encoder, p unsafe.Pointer)     { e.encUint(intToUint(int64(*(*int32)(p)))) }
func encInt64(e *Encoder, p unsafe.Pointer)     { e.encUint(intToUint(int64(*(*int64)(p)))) }
func encInt(e *Encoder, p unsafe.Pointer)       { e.encUint(intToUint(int64(*(*int)(p)))) }
func encUint8(e *Encoder, p unsafe.Pointer)     { e.buf = append(e.buf, *(*uint8)(p)) }
func encUint16(e *Encoder, p unsafe.Pointer)    { e.encUint(uint64(*(*uint16)(p))) }
func encUint32(e *Encoder, p unsafe.Pointer)    { e.encUint(uint64(*(*uint32)(p))) }
func encUint64(e *Encoder, p unsafe.Pointer)    { e.encUint(uint64(*(*uint64)(p))) }
func encUint(e *Encoder, p unsafe.Pointer)      { e.encUint(uint64(*(*uint)(p))) }
func encFloat32(e *Encoder, p unsafe.Pointer)   { e.encUint(floatToUint(float64(*(*float32)(p)))) }
func encFloat64(e *Encoder, p unsafe.Pointer)   { e.encUint(floatToUint(float64(*(*float64)(p)))) }
func encComplex64(e *Encoder, p unsafe.Pointer) { e.encUint(*(*uint64)(p)) }
func encComplex128(e *Encoder, p unsafe.Pointer) {
	e.encUint(*(*uint64)(p))
	e.encUint(*(*uint64)(unsafe.Pointer(uintptr(p) + 8)))
}

func encString(e *Encoder, p unsafe.Pointer) {
	s := *(*string)(p)
	e.encLength(len(s))
	e.buf = append(e.buf, s...)
}

//func encTime(e *Encoder, p unsafe.Pointer) {
//t := v.Interface().(time.Time)

//}
