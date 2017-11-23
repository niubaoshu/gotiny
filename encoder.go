package gotiny

import (
	"reflect"
	"unsafe"
)

type Encoder struct {
	buf     []byte //编码目的数组
	off     int
	boolPos int  //下一次要设置的bool在buf中的下标,即buf[boolPos]
	boolBit byte //下一次要设置的bool的buf[boolPos]中的bit位

	encEngs []encEng
	length  int
}

func Encodes(is ...interface{}) []byte {
	return NewEncoderWithPtr(is...).Encode(is...)
}

func NewEncoderWithPtr(ps ...interface{}) *Encoder {
	l := len(ps)
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
	engs, length := e.encEngs, e.length
	for i := 0; i < length; i++ {
		engs[i](e, (*eface)(unsafe.Pointer(&is[i])).data)
	}
	return e.reset()
}

// 入参是要编码的值得unsafe.Pointer 指针
func (e *Encoder) EncodePtr(ps ...unsafe.Pointer) []byte {
	engs, length := e.encEngs, e.length
	for i := 0; i < length; i++ {
		engs[i](e, ps[i])
	}
	return e.reset()
}

// vs 是持有要编码的值
func (e *Encoder) EncodeValue(vs ...reflect.Value) []byte {
	engs, length := e.encEngs, e.length
	for i := 0; i < length; i++ {
		v := (*refVal)(unsafe.Pointer(&vs[i]))
		p := v.ptr
		if v.flag&flagIndir == 0 {
			p = unsafe.Pointer(&v.ptr)
		}
		engs[i](e, p)
	}
	return e.reset()
}

// 编码产生的数据将append到buf上
func (e *Encoder) AppendTo(buf []byte) {
	e.off = len(buf)
	e.buf = buf
}

func (e *Encoder) reset() []byte {
	buf := e.buf
	e.buf = buf[:e.off]
	e.boolBit = 0
	e.boolPos = 0
	return buf
}
