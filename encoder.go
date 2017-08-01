package gotiny

import (
	"reflect"
	"unsafe"
)

type Encoder struct {
	buf     []byte //编码目的数组
	boolPos int    //下一次要设置的bool在buf中的下标,即buf[boolPos]
	boolBit byte   //下一次要设置的bool的buf[boolPos]中的bit位

	encEngs []encEng
	length  int
}

func Encodes(is ...interface{}) []byte {
	return NewEncoderWithPtr(is...).Encode(is...)
}

func NewEncoderWithPtr(ps ...interface{}) *Encoder {
	l := len(ps)
	if l < 1 {
		panic("must have argument!")
	}
	engs := make([]encEng, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(ps[i])
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
		rt := reflect.TypeOf(is[i])
		engs[i] = getEncEngine(rt)
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
		engs[i] = getEncEngine(ts[i])
	}
	return &Encoder{
		length:  l,
		encEngs: engs,
	}
}

// 入参是要编码值的指针
func (e *Encoder) Encode(is ...interface{}) []byte {
	engs := e.encEngs
	for i := 0; i < e.length; i++ {
		engs[i](e, (*[2]unsafe.Pointer)(unsafe.Pointer(&is[i]))[1])
	}
	return e.reset()
}

// 入参是要编码的值得unsafe.Pointer 指针
func (e *Encoder) EncodePtr(ps ...unsafe.Pointer) []byte {
	engs := e.encEngs
	for i := 0; i < e.length; i++ {
		engs[i](e, ps[i])
	}
	return e.reset()
}

// vs 是持有要编码的值
func (e *Encoder) EncodeValue(vs ...reflect.Value) []byte {
	engs := e.encEngs
	for i := 0; i < e.length; i++ {
		engs[i](e, getPtr(vs[i]))
	}
	return e.reset()
}

// 编码产生的数据将append到buf上
func (e *Encoder) AppendTo(buf []byte) {
	e.buf = buf
}

func (e *Encoder) reset() []byte {
	buf := e.buf
	e.buf = nil
	e.boolBit = 0
	e.boolPos = 0
	return buf
}
