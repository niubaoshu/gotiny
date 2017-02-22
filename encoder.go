package gotiny

import (
	"reflect"
)

type Encoder struct {
	buf     []byte //编码目的数组
	offset  int    //从buf[offset]开始写数据
	boolPos int    //下一次要设置的bool在buf中的下标,即buf[boolPos]
	boolBit byte   //下一次要设置的bool的buf[boolPos]中的bit位

	encEngs []encEng
	length  int
}

func NewEncoder(is ...interface{}) *Encoder {
	l := len(is)
	if l < 1 {
		panic("must have argument!")
	}
	engs := make([]encEng, l)
	for i := 0; i < l; i++ {
		engs[i] = GetEncEng(reflect.TypeOf(is[i]))
	}
	return &Encoder{
		length:  l,
		encEngs: engs,
	}
}
func NewEncoderWithType(ts ...reflect.Type) *Encoder {
	l := len(ts)
	if l < 1 {
		panic("must have argument!")
	}
	engs := make([]encEng, l)
	for i := 0; i < l; i++ {
		engs[i] = GetEncEng(ts[i])
	}
	return &Encoder{
		length:  l,
		encEngs: engs,
	}
}

func (e *Encoder) SetOff(off int) {
	e.buf = e.buf[:off]
	e.offset = off
}

func (e *Encoder) SetBuf(buf []byte) {
	e.buf = buf
	e.Reset()
}

func (e *Encoder) Encodes(is ...interface{}) []byte {
	l, engs := e.length, e.encEngs
	for i := 0; i < l; i++ {
		engs[i](e, reflect.ValueOf(is[i]))
	}
	return e.buf
}

func (e *Encoder) EncodeValues(vs ...reflect.Value) []byte {
	l, engs := e.length, e.encEngs
	for i := 0; i < l; i++ {
		engs[i](e, vs[i])
	}
	return e.buf
}

func (e *Encoder) Reset() {
	e.buf = e.buf[:e.offset]
	e.boolBit = 0
	e.boolPos = 0
}
