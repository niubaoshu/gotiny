package gotiny

import "unsafe"

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

func (d *Decoder) decLength() int { return int(d.decUint()) }

var (
	decignore     = func(d *Decoder, p unsafe.Pointer) {}
	decBool       = func(d *Decoder, p unsafe.Pointer) { *(*bool)(p) = d.decBool() }
	decInt        = func(d *Decoder, p unsafe.Pointer) { *(*int)(p) = int(uintToInt(d.decUint())) }
	decInt8       = func(d *Decoder, p unsafe.Pointer) { *(*int8)(p) = int8(d.buf[d.index]); d.index++ }
	decInt16      = func(d *Decoder, p unsafe.Pointer) { *(*int16)(p) = int16(uintToInt(d.decUint())) }
	decInt32      = func(d *Decoder, p unsafe.Pointer) { *(*int32)(p) = int32(uintToInt(d.decUint())) }
	decInt64      = func(d *Decoder, p unsafe.Pointer) { *(*int64)(p) = int64(uintToInt(d.decUint())) }
	decUint       = func(d *Decoder, p unsafe.Pointer) { *(*uint)(p) = uint(d.decUint()) }
	decUint8      = func(d *Decoder, p unsafe.Pointer) { *(*uint8)(p) = d.buf[d.index]; d.index++ }
	decUint16     = func(d *Decoder, p unsafe.Pointer) { *(*uint16)(p) = uint16(d.decUint()) }
	decUint32     = func(d *Decoder, p unsafe.Pointer) { *(*uint32)(p) = uint32(d.decUint()) }
	decUint64     = func(d *Decoder, p unsafe.Pointer) { *(*uint64)(p) = d.decUint() }
	decUintptr    = func(d *Decoder, p unsafe.Pointer) { *(*uintptr)(p) = uintptr(d.decUint()) }
	decPointer    = func(d *Decoder, p unsafe.Pointer) { *(*uintptr)(p) = uintptr(d.decUint()) }
	decFloat32    = func(d *Decoder, p unsafe.Pointer) { *(*float32)(p) = float32(uintToFloat(d.decUint())) }
	decFloat64    = func(d *Decoder, p unsafe.Pointer) { *(*float64)(p) = uintToFloat(d.decUint()) }
	decComplex64  = func(d *Decoder, p unsafe.Pointer) { *(*uint64)(p) = d.decUint() }
	decComplex128 = func(d *Decoder, p unsafe.Pointer) {
		*(*uint64)(p) = d.decUint()
		*(*uint64)(unsafe.Pointer(uintptr(p) + ptr1Size)) = d.decUint()
	}

	decString = func(d *Decoder, p unsafe.Pointer) {
		l, header := d.decLength(), (*stringHeader)(p)
		bytes := (*[]byte)(unsafe.Pointer(&sliceHeader{header.data, l, l}))
		if header.len < l {
			*bytes = make([]byte, l)
			header.data = (*sliceHeader)(unsafe.Pointer(bytes)).data
		}
		header.len = l
		d.index += copy(*bytes, d.buf[d.index:])
	}

	decBytes = func(d *Decoder, p unsafe.Pointer) {
		bytes := (*[]byte)(p)
		if d.decBool() {
			l := d.decLength()
			header := (*sliceHeader)(p)
			if header.cap < l {
				*bytes = make([]byte, l)
			} else {
				header.len = l
			}
			d.index += copy(*bytes, d.buf[d.index:])
		} else if !isNil(p) {
			*bytes = nil
		}
	}
)
