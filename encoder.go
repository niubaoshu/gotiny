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
	e := NewEncoderWithPtr(is...)
	e.Encodes(is...)
	return e.Bytes()
}

func NewEncoderWithPtr(is ...interface{}) *Encoder {
	l := len(is)
	if l < 1 {
		panic("must have argument!")
	}
	engs := make([]encEng, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(is[i])
		if rt.Kind() != reflect.Ptr {
			panic("must a pointer type!")
		}
		engs[i] = getEncEngine(rt.Elem())
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
		engs[i] = getEncEngine(reflect.TypeOf(is[i]))
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
		if ts[i].Kind() != reflect.Ptr {
			panic("must a pointer type!")
		}
		engs[i] = getEncEngine(ts[i])
	}
	return &Encoder{
		length:  l,
		encEngs: engs,
	}
}

func (e *Encoder) ResetWith(buf []byte) {
	e.buf = buf
	e.offset = len(buf)
	e.boolBit = 0
	e.boolPos = 0
}

// 入参是要编码的值得指针
func (e *Encoder) Encodes(is ...interface{}) {
	l, engs := e.length, e.encEngs
	for i := 0; i < l; i++ {
		engs[i](e, unsafe.Pointer(reflect.ValueOf(is[i]).Pointer()))
	}
}

func (e *Encoder) EncodeByUPtrs(ps ...unsafe.Pointer) {
	l, engs := e.length, e.encEngs
	for i := 0; i < l; i++ {
		engs[i](e, ps[i])
	}
}

func (e *Encoder) EncodeValues(vs ...reflect.Value) {
	l, engs := e.length, e.encEngs
	for i := 0; i < l; i++ {
		engs[i](e, unsafe.Pointer(vs[i].UnsafeAddr()))
	}
}

func (e *Encoder) Reset() {
	e.buf = e.buf[:e.offset]
	e.boolBit = 0
	e.boolPos = 0
}

func (e *Encoder) Bytes() []byte {
	return e.buf
}
