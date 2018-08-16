package gotiny

import "unsafe"

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
	var buf [maxVarintBytes]byte
	var n int
	for v > 0x7F {
		buf[n] = uint8(v) | 0x80
		v >>= 7
		n++
	}
	buf[n] = uint8(v)
	n++
	e.buf = append(e.buf, buf[0:n]...)
}

func (e *Encoder) encLength(v int)    { e.encUint(uint64(v)) }
func (e *Encoder) encString(s string) { e.encLength(len(s)); e.buf = append(e.buf, s...) }
func (e *Encoder) encIsNotNil(v bool) { e.encBool(v) }

func encIgnore(*Encoder, unsafe.Pointer)        {}
func encBool(e *Encoder, p unsafe.Pointer)      { e.encBool(*(*bool)(p)) }
func encInt(e *Encoder, p unsafe.Pointer)       { e.encUint(intToUint(int64(*(*int)(p)))) }
func encInt8(e *Encoder, p unsafe.Pointer)      { e.buf = append(e.buf, *(*uint8)(p)) }
func encInt16(e *Encoder, p unsafe.Pointer)     { e.encUint(intToUint(int64(*(*int16)(p)))) }
func encInt32(e *Encoder, p unsafe.Pointer)     { e.encUint(intToUint(int64(*(*int32)(p)))) }
func encInt64(e *Encoder, p unsafe.Pointer)     { e.encUint(intToUint(int64(*(*int64)(p)))) }
func encUint8(e *Encoder, p unsafe.Pointer)     { e.buf = append(e.buf, *(*uint8)(p)) }
func encUint16(e *Encoder, p unsafe.Pointer)    { e.encUint(uint64(*(*uint16)(p))) }
func encUint32(e *Encoder, p unsafe.Pointer)    { e.encUint(uint64(*(*uint32)(p))) }
func encUint64(e *Encoder, p unsafe.Pointer)    { e.encUint(uint64(*(*uint64)(p))) }
func encUint(e *Encoder, p unsafe.Pointer)      { e.encUint(uint64(*(*uint)(p))) }
func encUintptr(e *Encoder, p unsafe.Pointer)   { e.encUint(uint64(*(*uintptr)(p))) }
func encPointer(e *Encoder, p unsafe.Pointer)   { e.encUint(uint64(*(*uintptr)(p))) }
func encFloat32(e *Encoder, p unsafe.Pointer)   { e.encUint(uint64(float32ToUint32(p))) }
func encFloat64(e *Encoder, p unsafe.Pointer)   { e.encUint(float64ToUint(p)) }
func encString(e *Encoder, p unsafe.Pointer)    { e.encString(*(*string)(p)) }
func encComplex64(e *Encoder, p unsafe.Pointer) { e.encUint(*(*uint64)(p)) }
func encComplex128(e *Encoder, p unsafe.Pointer) {
	e.encUint(*(*uint64)(p))
	e.encUint(*(*uint64)(unsafe.Pointer(uintptr(p) + ptr1Size)))
}

func encBytes(e *Encoder, p unsafe.Pointer) {
	isNotNil := !isNil(p)
	e.encIsNotNil(isNotNil)
	if isNotNil {
		buf := *(*[]byte)(p)
		e.encLength(len(buf))
		e.buf = append(e.buf, buf...)
	}
}
