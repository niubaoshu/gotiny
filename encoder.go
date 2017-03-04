package gotiny

import (
	"reflect"
	"unsafe"
)

type Encoder struct {
	buf     []byte //编码目的数组
	offset  int    //开始从buf的offset 写入编码的结果
	boolPos int    //下一次要设置的bool在buf中的下标,即buf[boolPos]
	boolBit byte   //下一次要设置的bool的buf[boolPos]中的bit位

	encEngs []encEng
	length  int
}

func Encodes(is ...interface{}) []byte {
	return NewEncoderWithPtr(is...).Encodes(is...)
}

func NewEncoderWithPtr(is ...interface{}) *Encoder {
	l := len(is)
	if l < 1 {
		panic("must have argument!")
	}
	engs := make([]encEng, l)
	for i := 0; i < l; i++ {
		engs[i] = GetEncEng(reflect.TypeOf(is[i]).Elem())
	}
	return &Encoder{
		length:  l,
		encEngs: engs,
	}
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
func NewEncoderWithTypes(ts ...reflect.Type) *Encoder {
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

func (e *Encoder) SetBuf(buf []byte) {
	e.buf = buf
	e.offset = len(buf)
	e.boolBit = 0
	e.boolPos = 0
}

// is is point of value slice
func (e *Encoder) Encodes(is ...interface{}) (buf []byte) {
	l, engs := e.length, e.encEngs
	for i := 0; i < l; i++ {
		engs[i](e, unsafe.Pointer(reflect.ValueOf(is[i]).Pointer()))
	}
	buf = e.buf
	e.Reset()
	return
}

func (e *Encoder) EncodeByUPtr(ps ...unsafe.Pointer) (buf []byte) {
	l, engs := e.length, e.encEngs
	for i := 0; i < l; i++ {
		engs[i](e, ps[i])
	}
	buf = e.buf
	e.Reset()
	return
}

func (e *Encoder) EncodeValues(vs ...reflect.Value) (buf []byte) {
	l, engs := e.length, e.encEngs
	for i := 0; i < l; i++ {
		engs[i](e, unsafe.Pointer(vs[i].UnsafeAddr()))
	}
	buf = e.buf
	e.Reset()
	return
}

func (e *Encoder) Reset() {
	e.buf = e.buf[:e.offset]
	e.boolBit = 0
	e.boolPos = 0
}
