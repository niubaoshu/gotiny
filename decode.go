package gotiny

// func (d *Decoder) Decodes(is ...interface{}) {
// 	for i, _ := range is {
// 		d.Decode(i)
// 	}
// }

// func (d *Decoder) Decode(i interface{}) {
// 	switch v := i.(type) {
// 	case *bool:
// 		d.decBool(v)
// 	case *uint8:
// 		d.decUint8(v)
// 	case *int8:
// 		d.decInt8(v)
// 	case *uint, *uint16, *uint32, *uint64, *uintptr:
// 		d.decUint(v)
// 	case *int, *int16, *int32, *int64:
// 		d.decInt(v)
// 	case *float32, *float64:
// 		d.decFloat(v)
// 	case *complex64, *complex128:
// 		d.decComplex(v)
// 	case *string:
// 		d.decString(v)
// 	case *[]byte:
// 		d.decSliceByte(v)
// 	default:
// 		d.decOther(v)
// 	}
// }

// func decBool(d *Decoder, v *bool) {
// 	if d.boolBit == 0 {
// 		d.boolBit = 1
// 		d.boolen = d.buff[d.offset]
// 		d.offset++
// 	}
// 	*v = d.boolen&d.boolBit != 0
// 	d.boolBit <<= 1
// }
// func decSliceByte(d *Decoder, v *[]byte) {
// 	var u uint64
// 	decUint(d, &u)
// 	*v = d.buff[d.offset : d.offset+int(u)]
// }
// func decUint8(d *Decoder, v *uint8) { *v = d.buff[d.offset]; d.offset++ }
// func decInt8(d *Decoder, v *int8)   { *v = int8(d.buff[d.offset]); d.offset++ }

// func decUint(d *Decoder, v *uint64) {
// 	var x uint64
// 	var s uint
// 	offset := d.offset
// 	buff := d.buff
// 	for buff[offset] >= 0x80 {
// 		x |= uint64(buff[offset]&0x7f) << s
// 		s += 7
// 		offset++
// 	}
// 	x |= uint64(buff[d.offset]&0x7f) << s
// 	d.offset = offset + 1
// 	*v = x
// }

// func decInt(d *Decoder, v *int64) {
// 	var ux uint64
// 	decUint(d, &ux)
// 	x := *v >> 1
// 	if *v&1 != 0 {
// 		x = ^x
// 	}
// 	*v = x
// }

// func decFloat(d *Decoder, v *float64) {
// 	decUint(d, v)
// 	var v uint64
// 	for i := 0; i < 8; i++ {
// 		v <<= 8
// 		v |= u & 0xFF
// 		u >>= 8
// 	}
// 	return *((*float64)(unsafe.Pointer(&v)))
// }

// func (d *Decoder) decComplex() complex128 {
// 	real := float64FromBits(d.decUint())
// 	imag := float64FromBits(d.decUint())
// 	return complex(real, imag)
// }

// func (d *Decoder) decString() string {
// 	l := int(d.decUint())
// 	d.offset += l
// 	return string(d.buff[d.offset : d.offset+l])
// }

// func float64FromBits(u uint64) float64 {
// 	var v uint64
// 	for i := 0; i < 8; i++ {
// 		v <<= 8
// 		v |= u & 0xFF
// 		u >>= 8
// 	}
// 	return *((*float64)(unsafe.Pointer(&v)))
// }

// func uvarint(buf []byte) (uint64, int) {
// 	var x uint64
// 	var s uint
// 	for i, b := range buf {
// 		if b < 0x80 {
// 			if i > 9 || i == 9 && b > 1 {
// 				return 0, -(i + 1) // overflow
// 			}
// 			return x | uint64(b)<<s, i + 1
// 		}
// 		x |= uint64(b&0x7f) << s
// 		s += 7
// 	}
// 	return 0, 0
// }

// func varint(buf []byte) (int64, int) {
// 	ux, n := uvarint(buf) // ok to continue in presence of error
// 	x := int64(ux >> 1)
// 	if ux&1 != 0 {
// 		x = ^x
// 	}
// 	return x, n
// }
// func (d *Decoder) decOther(v interface{}) {
// 	d.decodeValue(reflect.ValueOf(v))
// }
