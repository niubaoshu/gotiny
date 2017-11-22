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
	buf := e.buf
	for v >= 0x80 {
		buf = append(buf, uint8(v)|0x80)
		v >>= 7
	}
	e.buf = append(buf, uint8(v))
}

func (e *Encoder) encLength(v int)    { e.encUint(uint64(v)) }
func (e *Encoder) encString(s string) { e.encLength(len(s)); e.buf = append(e.buf, s...) }

var (
	encIgnore     = func(e *Encoder, p unsafe.Pointer) {}
	encBool       = func(e *Encoder, p unsafe.Pointer) { e.encBool(*(*bool)(p)) }
	encInt        = func(e *Encoder, p unsafe.Pointer) { e.encUint(intToUint(int64(*(*int)(p)))) }
	encInt8       = func(e *Encoder, p unsafe.Pointer) { e.buf = append(e.buf, *(*uint8)(p)) }
	encInt16      = func(e *Encoder, p unsafe.Pointer) { e.encUint(intToUint(int64(*(*int16)(p)))) }
	encInt32      = func(e *Encoder, p unsafe.Pointer) { e.encUint(intToUint(int64(*(*int32)(p)))) }
	encInt64      = func(e *Encoder, p unsafe.Pointer) { e.encUint(intToUint(int64(*(*int64)(p)))) }
	encUint8      = func(e *Encoder, p unsafe.Pointer) { e.buf = append(e.buf, *(*uint8)(p)) }
	encUint16     = func(e *Encoder, p unsafe.Pointer) { e.encUint(uint64(*(*uint16)(p))) }
	encUint32     = func(e *Encoder, p unsafe.Pointer) { e.encUint(uint64(*(*uint32)(p))) }
	encUint64     = func(e *Encoder, p unsafe.Pointer) { e.encUint(uint64(*(*uint64)(p))) }
	encUint       = func(e *Encoder, p unsafe.Pointer) { e.encUint(uint64(*(*uint)(p))) }
	encUintptr    = func(e *Encoder, p unsafe.Pointer) { e.encUint(uint64(*(*uintptr)(p))) }
	encPointer    = func(e *Encoder, p unsafe.Pointer) { e.encUint(uint64(*(*uintptr)(p))) }
	encFloat32    = func(e *Encoder, p unsafe.Pointer) { e.encUint(floatToUint(float64(*(*float32)(p)))) }
	encFloat64    = func(e *Encoder, p unsafe.Pointer) { e.encUint(floatToUint(float64(*(*float64)(p)))) }
	encComplex64  = func(e *Encoder, p unsafe.Pointer) { e.encUint(*(*uint64)(p)) }
	encString     = func(e *Encoder, p unsafe.Pointer) { e.encString(*(*string)(p)) }
	encComplex128 = func(e *Encoder, p unsafe.Pointer) {
		e.encUint(*(*uint64)(p))
		e.encUint(*(*uint64)(unsafe.Pointer(uintptr(p) + ptr1Size)))
	}

	encBytes = func(e *Encoder, p unsafe.Pointer) {
		isNotNil := !isNil(p)
		e.encBool(isNotNil)
		if isNotNil {
			buf := *(*[]byte)(p)
			e.encLength(len(buf))
			e.buf = append(e.buf, buf...)
		}
	}
)
