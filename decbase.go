package gotiny

import (
	"reflect"
)

func (d *Decoder) decBool() (b bool) {
	if d.boolBit == 0 {
		d.boolBit = 1
		d.boolean = d.buf[d.index]
		d.index++
	}
	b = d.boolean&d.boolBit != 0
	d.boolBit <<= 1
	return
}

func (d *Decoder) decBytes() (b []byte) {
	l := int(d.decUint())
	b = d.buf[d.index : d.index+l]
	d.index += l
	return
}

func (d *Decoder) decString() (s string) {
	l := int(d.decUint())
	s = string(d.buf[d.index : d.index+l])
	d.index += l
	return
}
func (d *Decoder) decUint() (u uint64) {
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

func (d *Decoder) decUintFast() uint64 {
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

func (d *Decoder) decUint8() (x uint8) { x = d.buf[d.index]; d.index++; return }
func (d *Decoder) decInt8() int8       { return int8(d.decUint8()) }
func (d *Decoder) decInt() int64       { return uintToInt(d.decUint()) }
func (d *Decoder) decFloat() float64   { return uintToFloat(d.decUint()) }
func (d *Decoder) decComplex() complex128 {
	return complex(uintToFloat(d.decUint()), uintToFloat(d.decUint()))
}

func decignore(d *Decoder, v reflect.Value)  {}
func decBool(d *Decoder, v reflect.Value)    { v.SetBool(d.decBool()) }
func decUint8(d *Decoder, v reflect.Value)   { v.SetUint(uint64(d.decUint8())) }
func decInt8(d *Decoder, v reflect.Value)    { v.SetInt(int64(d.decInt8())) }
func decUint(d *Decoder, v reflect.Value)    { v.SetUint(d.decUint()) }
func decInt(d *Decoder, v reflect.Value)     { v.SetInt(d.decInt()) }
func decFloat(d *Decoder, v reflect.Value)   { v.SetFloat(d.decFloat()) }
func decComplex(d *Decoder, v reflect.Value) { v.SetComplex(d.decComplex()) }
func decString(d *Decoder, v reflect.Value)  { v.SetString(d.decString()) }
