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

	encEngs    []encEng
	encEngVals []encEngVal
	length     int
}

func Encodes(is ...interface{}) []byte {
	e := NewEncoderWithPtr(is...)
	e.Encodes(is...)
	return e.Bytes()
}

func NewEncoderWithPtr(ps ...interface{}) *Encoder {
	l := len(ps)
	if l < 1 {
		panic("must have argument!")
	}
	engs := make([]encEng, l)
	engvals := make([]encEngVal, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(ps[i])
		if rt.Kind() != reflect.Ptr {
			panic("must a pointer type!")
		}
		engs[i] = getEncEngine(rt.Elem())
		engvals[i] = getValEncEng(rt.Elem())
	}
	return &Encoder{
		length:     l,
		encEngs:    engs,
		encEngVals: engvals,
	}
}

func NewEncoder(is ...interface{}) *Encoder {
	l := len(is)
	if l < 1 {
		panic("must have argument!")
	}
	engs := make([]encEng, l)
	engvals := make([]encEngVal, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(is[i])
		engs[i] = getEncEngine(rt)
		engvals[i] = getValEncEng(rt)
	}
	return &Encoder{
		length:     l,
		encEngs:    engs,
		encEngVals: engvals,
	}
}
func NewEncoderWithTypes(ts ...reflect.Type) *Encoder {
	l := len(ts)
	if l < 1 {
		panic("must have argument!")
	}
	engs := make([]encEng, l)
	engvals := make([]encEngVal, l)
	for i := 0; i < l; i++ {
		engs[i] = getEncEngine(ts[i])
		engvals[i] = getValEncEng(ts[i])
	}
	return &Encoder{
		length:     l,
		encEngs:    engs,
		encEngVals: engvals,
	}
}

// 入参是要编码值的指针
func (e *Encoder) Encodes(is ...interface{}) {
	engs := e.encEngs
	for i := 0; i < e.length; i++ {
		engs[i](e, (*[2]unsafe.Pointer)(unsafe.Pointer(&is[i]))[1])
	}
}

func (e *Encoder) EncodeByUPtrs(ps ...unsafe.Pointer) {
	engs := e.encEngs
	for i := 0; i < e.length; i++ {
		engs[i](e, ps[i])
	}
}

// vs 是持有要编码的值
func (e *Encoder) EncodeValues(vs ...reflect.Value) {
	engs := e.encEngVals
	for i := 0; i < e.length; i++ {
		engs[i](e, vs[i])
	}
}

func (e *Encoder) ResetWithBuf(buf []byte) {
	e.buf = buf
	e.offset = len(buf)
	e.boolBit = 0
	e.boolPos = 0
}

func (e *Encoder) Reset() {
	e.buf = e.buf[:e.offset]
	e.boolBit = 0
	e.boolPos = 0
}

func (e *Encoder) Bytes() []byte {
	return e.buf
}
