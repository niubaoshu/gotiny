package gotiny

import (
	"reflect"
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

func (d *Decoder) decUint() uint64 {
	buf, i := d.buf, d.index
	if buf[i] < 0x80 {
		d.index++
		return uint64(buf[i])
	}
	// we already checked the first byte
	x := uint64(buf[i]) - 0x80
	i++

	b := uint64(buf[i])
	i++

	x += b << 7
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 7

	b = uint64(buf[i])
	i++
	x += b << 14
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 14

	b = uint64(buf[i])
	i++
	x += b << 21
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 21

	b = uint64(buf[i])
	i++
	x += b << 28
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 28

	b = uint64(buf[i])
	i++
	x += b << 35
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 35

	b = uint64(buf[i])
	i++
	x += b << 42
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 42

	b = uint64(buf[i])
	i++
	x += b << 49
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 49

	b = uint64(buf[i])
	i++
	x += b << 56
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 56

	b = uint64(buf[i])
	i++
	x += b << 63
done:
	d.index = i
	return x
}

func (d *Decoder) decLength() int    { return int(d.decUint()) }
func (d *Decoder) decIsNotNil() bool { return d.decBool() }

func decIgnore(*Decoder, unsafe.Pointer)      {}
func decBool(d *Decoder, p unsafe.Pointer)    { *(*bool)(p) = d.decBool() }
func decInt(d *Decoder, p unsafe.Pointer)     { *(*int)(p) = int(uintToInt(d.decUint())) }
func decInt8(d *Decoder, p unsafe.Pointer)    { *(*int8)(p) = int8(d.buf[d.index]); d.index++ }
func decInt16(d *Decoder, p unsafe.Pointer)   { *(*int16)(p) = int16(uintToInt(d.decUint())) }
func decInt32(d *Decoder, p unsafe.Pointer)   { *(*int32)(p) = int32(uintToInt(d.decUint())) }
func decInt64(d *Decoder, p unsafe.Pointer)   { *(*int64)(p) = int64(uintToInt(d.decUint())) }
func decUint(d *Decoder, p unsafe.Pointer)    { *(*uint)(p) = uint(d.decUint()) }
func decUint8(d *Decoder, p unsafe.Pointer)   { *(*uint8)(p) = d.buf[d.index]; d.index++ }
func decUint16(d *Decoder, p unsafe.Pointer)  { *(*uint16)(p) = uint16(d.decUint()) }
func decUint32(d *Decoder, p unsafe.Pointer)  { *(*uint32)(p) = uint32(d.decUint()) }
func decUint64(d *Decoder, p unsafe.Pointer)  { *(*uint64)(p) = d.decUint() }
func decUintptr(d *Decoder, p unsafe.Pointer) { *(*uintptr)(p) = uintptr(d.decUint()) }
func decPointer(d *Decoder, p unsafe.Pointer) { *(*uintptr)(p) = uintptr(d.decUint()) }
func decFloat32(d *Decoder, p unsafe.Pointer) {
	*(*float32)(p) = float32(uint32ToFloat32(uint32(d.decUint())))
}
func decFloat64(d *Decoder, p unsafe.Pointer)   { *(*float64)(p) = uintToFloat64(d.decUint()) }
func decComplex64(d *Decoder, p unsafe.Pointer) { *(*uint64)(p) = d.decUint() }
func decComplex128(d *Decoder, p unsafe.Pointer) {
	*(*uint64)(p) = d.decUint()
	*(*uint64)(unsafe.Pointer(uintptr(p) + ptr1Size)) = d.decUint()
}

func decString(d *Decoder, p unsafe.Pointer) {
	l, header := d.decLength(), (*reflect.StringHeader)(p)
	bytes := (*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: header.Data, Len: l, Cap: l}))
	if header.Len < l {
		*bytes = make([]byte, l)
		header.Data = (*reflect.SliceHeader)(unsafe.Pointer(bytes)).Data
	}
	header.Len = l
	d.index += copy(*bytes, d.buf[d.index:])
}

func decBytes(d *Decoder, p unsafe.Pointer) {
	bytes := (*[]byte)(p)
	if d.decIsNotNil() {
		l := d.decLength()
		header := (*reflect.SliceHeader)(p)
		if header.Cap < l {
			*bytes = make([]byte, l)
		} else {
			header.Len = l
		}
		d.index += copy(*bytes, d.buf[d.index:])
	} else if !isNil(p) {
		*bytes = nil
	}
}
