package gotiny

import (
	"reflect"
	"unsafe"
)

type Decoder struct {
	buf     []byte //buf
	index   int    //下一个要读取的字节
	offset  int    //开始解码的偏移量
	boolean byte   //下一次要读取的bool在buf中的下标,即buf[boolPos]
	boolBit byte   //下一次要读取的bool的buf[boolPos]中的bit位

	decEngs []decEng //解码器集合
	length  int      //解码器数量
}

func Decodes(buf []byte, is ...interface{}) int {
	d := NewDecoderWithPtrs(is...)
	d.buf = buf
	d.Decodes(is...)
	return d.index
}

func NewDecoderWithPtrs(is ...interface{}) *Decoder {
	l := len(is)
	if l < 1 {
		panic("must have argument!")
	}
	des := make([]decEng, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(is[i])
		if rt.Kind() != reflect.Ptr {
			panic("must a pointer type!")
		}
		des[i] = getDecEngine(rt.Elem())
	}
	return &Decoder{
		length:  l,
		decEngs: des,
	}
}

func NewDecoder(is ...interface{}) *Decoder {
	l := len(is)
	if l < 1 {
		panic("must have argument!")
	}
	des := make([]decEng, l)
	for i := 0; i < l; i++ {
		des[i] = getDecEngine(reflect.TypeOf(is[i]))
	}
	return &Decoder{
		length:  l,
		decEngs: des,
	}
}

func NewDecoderWithTypes(ts ...reflect.Type) *Decoder {
	l := len(ts)
	if l < 1 {
		panic("must have argument!")
	}
	des := make([]decEng, l)
	for i := 0; i < l; i++ {
		des[i] = getDecEngine(ts[i])
	}
	return &Decoder{
		length:  l,
		decEngs: des,
	}
}

func (d *Decoder) Reset() {
	d.index = d.offset
	d.boolean = 0
	d.boolBit = 0
}

func (d *Decoder) ResetWith(b []byte) {
	d.buf = b
	d.Reset()
}
func (d *Decoder) Decodes(is ...interface{}) {
	engs := d.decEngs
	for i := 0; i < d.length; i++ {
		engs[i](d, (*[2]unsafe.Pointer)(unsafe.Pointer(&is[i]))[1])
	}
}

// is is pointer of value
func (d *Decoder) DecodeByUPtr(ps ...unsafe.Pointer) {
	engs := d.decEngs
	for i := 0; i < d.length; i++ {
		engs[i](d, ps[i])
	}
}

func (d *Decoder) DecodeValues(vs ...reflect.Value) {
	engs := d.decEngs
	for i := 0; i < d.length; i++ {
		engs[i](d, unsafe.Pointer(vs[i].UnsafeAddr()))
	}
}

func (d *Decoder) SetOff(off int) {
	d.offset = off
	d.index = off
}
