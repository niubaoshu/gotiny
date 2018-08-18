package gotiny

import (
	"time"
	"unsafe"
)

func (d *Decoder) decBool() (b bool) {
	if d.boolBit == 0 {
		d.boolBit = 1
		d.boolPos = d.buf[d.index]
		d.index++
	}
	b = d.boolPos&d.boolBit != 0
	d.boolBit <<= 1
	return
}

func (d *Decoder) decUint64() uint64 {
	buf, i := d.buf, d.index
	switch {
	case buf[i] < 0x80:
		d.index++
		return uint64(buf[i])
	case buf[i+1] < 0x80:
		d.index += 2
		return uint64(buf[i]) + uint64(buf[i+1]-1)<<7
	case buf[i+2] < 0x80:
		d.index += 3
		return uint64(buf[i]) + uint64(buf[i+1])<<7 + uint64(buf[i+2])<<14 - (1<<7 + 1<<14)
	case buf[i+3] < 0x80:
		d.index += 4
		return uint64(buf[i]) + uint64(buf[i+1])<<7 + uint64(buf[i+2])<<14 + uint64(buf[i+3])<<21 - (1<<7 + 1<<14 + 1<<21)
	case buf[i+4] < 0x80:
		d.index += 5
		return uint64(buf[i]) + uint64(buf[i+1])<<7 + uint64(buf[i+2])<<14 + uint64(buf[i+3])<<21 + uint64(buf[i+4])<<28 - (1<<7 + 1<<14 + 1<<21 + 1<<28)
	case buf[i+5] < 0x80:
		d.index += 6
		return uint64(buf[i]) + uint64(buf[i+1])<<7 + uint64(buf[i+2])<<14 + uint64(buf[i+3])<<21 + uint64(buf[i+4])<<28 + uint64(buf[i+5])<<35 -
			(1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35)
	case buf[i+6] < 0x80:
		d.index += 7
		return uint64(buf[i]) + uint64(buf[i+1])<<7 + uint64(buf[i+2])<<14 + uint64(buf[i+3])<<21 + uint64(buf[i+4])<<28 + uint64(buf[i+5])<<35 +
			uint64(buf[i+6])<<42 - (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35 + 1<<42)
	case buf[i+7] < 0x80:
		d.index += 8
		return uint64(buf[i]) + uint64(buf[i+1])<<7 + uint64(buf[i+2])<<14 + uint64(buf[i+3])<<21 + uint64(buf[i+4])<<28 + uint64(buf[i+5])<<35 +
			uint64(buf[i+6])<<42 + uint64(buf[i+7])<<49 - (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35 + 1<<42 + 1<<49)
	default:
		d.index += 9
		return uint64(buf[i]) + uint64(buf[i+1])<<7 + uint64(buf[i+2])<<14 + uint64(buf[i+3])<<21 + uint64(buf[i+4])<<28 + uint64(buf[i+5])<<35 +
			uint64(buf[i+6])<<42 + uint64(buf[i+7])<<49 + uint64(buf[i+8])<<56 - (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35 + 1<<42 + 1<<49 + 1<<56)
	}
}

func (d *Decoder) decUint16() uint16 {
	buf, i := d.buf, d.index
	x0 := buf[i]
	if x0 < 0x80 {
		d.index++
		return uint16(x0)
	}
	x1 := buf[i+1]
	if x1 < 0x80 {
		d.index += 2
		return uint16(x0) + uint16(x1-1)<<7
	}
	d.index += 3
	return uint16(x0) + uint16(x1)<<7 + uint16(buf[i+2])<<14 - (1<<7 + 1<<14)
}

func (d *Decoder) decUint32() uint32 {
	buf, i := d.buf, d.index
	x0 := buf[i]
	if x0 < 0x80 {
		d.index++
		return uint32(x0)
	}
	x1 := buf[i+1]
	if x1 < 0x80 {
		d.index += 2
		return uint32(x0) + uint32(x1-1)<<7
	}
	x2 := buf[i+2]
	if x2 < 0x80 {
		d.index += 3
		return uint32(x0) + uint32(x1)<<7 + uint32(x2)<<14 - (1<<7 + 1<<14)
	}
	x3 := buf[i+3]
	if x3 < 0x80 {
		d.index += 4
		return uint32(x0) + uint32(x1)<<7 + uint32(x2)<<14 + uint32(x3)<<21 - (1<<7 + 1<<14 + 1<<21)
	}
	d.index += 5
	return uint32(x0) + uint32(x1)<<7 + uint32(x2)<<14 + uint32(x3)<<21 + uint32(buf[i+4])<<28 - (1<<7 + 1<<14 + 1<<21 + 1<<28)
}

func (d *Decoder) decLength() int    { return int(d.decUint32()) }
func (d *Decoder) decIsNotNil() bool { return d.decBool() }

func decIgnore(*Decoder, unsafe.Pointer)        {}
func decBool(d *Decoder, p unsafe.Pointer)      { *(*bool)(p) = d.decBool() }
func decInt(d *Decoder, p unsafe.Pointer)       { *(*int)(p) = int(uint64ToInt64(d.decUint64())) }
func decInt8(d *Decoder, p unsafe.Pointer)      { *(*int8)(p) = int8(d.buf[d.index]); d.index++ }
func decInt16(d *Decoder, p unsafe.Pointer)     { *(*int16)(p) = uint16ToInt16(d.decUint16()) }
func decInt32(d *Decoder, p unsafe.Pointer)     { *(*int32)(p) = uint32ToInt32(d.decUint32()) }
func decInt64(d *Decoder, p unsafe.Pointer)     { *(*int64)(p) = uint64ToInt64(d.decUint64()) }
func decUint(d *Decoder, p unsafe.Pointer)      { *(*uint)(p) = uint(d.decUint64()) }
func decUint8(d *Decoder, p unsafe.Pointer)     { *(*uint8)(p) = d.buf[d.index]; d.index++ }
func decUint16(d *Decoder, p unsafe.Pointer)    { *(*uint16)(p) = d.decUint16() }
func decUint32(d *Decoder, p unsafe.Pointer)    { *(*uint32)(p) = d.decUint32() }
func decUint64(d *Decoder, p unsafe.Pointer)    { *(*uint64)(p) = d.decUint64() }
func decUintptr(d *Decoder, p unsafe.Pointer)   { *(*uintptr)(p) = uintptr(d.decUint64()) }
func decPointer(d *Decoder, p unsafe.Pointer)   { *(*uintptr)(p) = uintptr(d.decUint64()) }
func decFloat32(d *Decoder, p unsafe.Pointer)   { *(*float32)(p) = uint32ToFloat32(d.decUint32()) }
func decFloat64(d *Decoder, p unsafe.Pointer)   { *(*float64)(p) = uint64ToFloat64(d.decUint64()) }
func decTime(d *Decoder, p unsafe.Pointer)      { *(*time.Time)(p) = time.Unix(0, int64(d.decUint64())) }
func decComplex64(d *Decoder, p unsafe.Pointer) { *(*uint64)(p) = d.decUint64() }
func decComplex128(d *Decoder, p unsafe.Pointer) {
	*(*uint64)(p) = d.decUint64()
	*(*uint64)(unsafe.Pointer(uintptr(p) + ptr1Size)) = d.decUint64()
}

func decString(d *Decoder, p unsafe.Pointer) {
	l, val := int(d.decUint32()), (*string)(p)
	*val = string(d.buf[d.index : d.index+l])
	d.index += l
}

func decBytes(d *Decoder, p unsafe.Pointer) {
	bytes := (*[]byte)(p)
	if d.decIsNotNil() {
		l := int(d.decUint32())
		*bytes = d.buf[d.index : d.index+l]
		d.index += l
	} else if !isNil(p) {
		*bytes = nil
	}
}
