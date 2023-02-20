package gotiny

import (
	"reflect"
	"unsafe"
)

type Decoder struct {
	buf     []byte   // buffer
	index   int      // index of the next byte to be used in the buffer
	boolPos byte     // index of the next bool to be read in the buffer, i.e. buf[boolPos]
	boolBit byte     // bit position of the next bool to be read in buf[boolPos]
	engines []decEng // collection of decoders
	length  int      // number of decoders

}

func Unmarshal(buf []byte, is ...any) int {
	return NewDecoderWithPtr(is...).Decode(buf, is...)
}

func NewDecoderWithPtr(is ...any) *Decoder {
	l := len(is)
	engines := make([]decEng, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(is[i])
		if rt.Kind() != reflect.Ptr {
			panic("must a pointer type!")
		}
		engines[i] = getDecEngine(rt.Elem())
	}
	return &Decoder{
		length:  l,
		engines: engines,
	}
}

func NewDecoder(is ...any) *Decoder {
	l := len(is)
	engines := make([]decEng, l)
	for i := 0; i < l; i++ {
		engines[i] = getDecEngine(reflect.TypeOf(is[i]))
	}
	return &Decoder{
		length:  l,
		engines: engines,
	}
}

func NewDecoderWithType(ts ...reflect.Type) *Decoder {
	l := len(ts)
	des := make([]decEng, l)
	for i := 0; i < l; i++ {
		des[i] = getDecEngine(ts[i])
	}
	return &Decoder{
		length:  l,
		engines: des,
	}
}

func (d *Decoder) reset() int {
	index := d.index
	d.index = 0
	d.boolPos = 0
	d.boolBit = 0
	return index
}

// is is pointer of variable
func (d *Decoder) Decode(buf []byte, is ...any) int {
	d.buf = buf
	engines := d.engines
	for i := 0; i < len(engines) && i < len(is); i++ {
		engines[i](d, (*[2]unsafe.Pointer)(unsafe.Pointer(&is[i]))[1])
	}
	return d.reset()
}

// ps is a unsafe.Pointer of the variable
func (d *Decoder) DecodePtr(buf []byte, ps ...unsafe.Pointer) int {
	d.buf = buf
	engines := d.engines
	for i := 0; i < len(engines) && i < len(ps); i++ {
		engines[i](d, ps[i])
	}
	return d.reset()
}

func (d *Decoder) DecodeValue(buf []byte, vs ...reflect.Value) int {
	d.buf = buf
	engines := d.engines
	for i := 0; i < len(engines) && i < len(vs); i++ {
		engines[i](d, unsafe.Pointer(vs[i].UnsafeAddr()))
	}
	return d.reset()
}
