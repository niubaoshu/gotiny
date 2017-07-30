package gotiny

import (
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

func (d *Decoder) decUintslow() (u uint64) {
	buf, i := d.buf[d.index:], 0
	s := uint(0)
	for buf[i] > 0x7f {
		u |= (uint64(buf[i]&0x7f) << s)
		s += 7
		i++
	}
	u |= uint64(buf[i]) << s
	i++
	d.index += i
	return u
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
	decInt16      = func(d *Decoder, p unsafe.Pointer) { *(*int16)(p) = int16(uintToInt(d.decUint())) }
	decInt32      = func(d *Decoder, p unsafe.Pointer) { *(*int32)(p) = int32(uintToInt(d.decUint())) }
	decInt64      = func(d *Decoder, p unsafe.Pointer) { *(*int64)(p) = int64(uintToInt(d.decUint())) }
	decInt        = func(d *Decoder, p unsafe.Pointer) { *(*int)(p) = int(uintToInt(d.decUint())) }
	decUint8      = func(d *Decoder, p unsafe.Pointer) { *(*uint8)(p) = d.buf[d.index]; d.index++ }
	decUint16     = func(d *Decoder, p unsafe.Pointer) { *(*uint16)(p) = uint16(d.decUint()) }
	decUint32     = func(d *Decoder, p unsafe.Pointer) { *(*uint32)(p) = uint32(d.decUint()) }
	decUint64     = func(d *Decoder, p unsafe.Pointer) { *(*uint64)(p) = d.decUint() }
	decUint       = func(d *Decoder, p unsafe.Pointer) { *(*uint)(p) = uint(d.decUint()) }
	decUintptr    = func(d *Decoder, p unsafe.Pointer) { *(*uintptr)(p) = uintptr(d.decUint()) }
	decFloat32    = func(d *Decoder, p unsafe.Pointer) { *(*float32)(p) = float32(uintToFloat(d.decUint())) }
	decFloat64    = func(d *Decoder, p unsafe.Pointer) { *(*float64)(p) = uintToFloat(d.decUint()) }
	decComplex64  = func(d *Decoder, p unsafe.Pointer) { *(*uint64)(p) = d.decUint() }
	decComplex128 = func(d *Decoder, p unsafe.Pointer) {
		*(*uint64)(p) = d.decUint()
		*(*uint64)(unsafe.Pointer(uintptr(p) + 8)) = d.decUint()
	}
	decString = func(d *Decoder, p unsafe.Pointer) {
		l := d.decLength()
		var bytes []byte
		if *(*int)(next1Ptr(p)) < l { // len(str) < l
			bytes = make([]byte, l)
			*(*unsafe.Pointer)(p) = *(*unsafe.Pointer)(unsafe.Pointer(&bytes))
		} else {
			*(*unsafe.Pointer)(unsafe.Pointer(&bytes)) = *(*unsafe.Pointer)(p)
			*(*int)(next1Ptr(unsafe.Pointer(&bytes))) = l
			*(*int)(next2Ptr(unsafe.Pointer(&bytes))) = l
		}
		d.index += copy(bytes, d.buf[d.index:])
		*(*int)(next1Ptr(p)) = l
	}

	decBytes = func(d *Decoder, p unsafe.Pointer) {
		vptr := (*[]byte)(p)
		if d.decBool() {
			l := d.decLength()
			if *(*int)(next2Ptr(p)) < l { // cap(bytes) < l
				*vptr = make([]byte, l)
			} else {
				lenptr := (*int)(next1Ptr(p))
				if *lenptr > l {
					*lenptr = l
				}
			}
			d.index += copy(*vptr, d.buf[d.index:])
		} else if !isNil(p) {
			*vptr = nil
		}
	}
)
