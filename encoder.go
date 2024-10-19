package gotiny

import (
	"reflect"
)

type Encoder struct {
	buf     []byte //编码目的数组
	off     int
	boolPos int  //下一次要设置的bool在buf中的下标,即buf[boolPos]
	boolBit byte //下一次要设置的bool的buf[boolPos]中的bit位

	engines []encEng
	length  int
}

/*
Marshal serializes the value pointed to by the incoming pointer.
The argument must be a pointer, similar to the form &value
which is serializing value. If value itself is a pointer,
then you can pass in value directly,
which serializes the value pointed to by value.
*/
func Marshal(ps ...any) []byte {
	return NewEncoderWithPtr(ps...).encode(ps...)
}

// 创建一个编码ps 指向类型的编码器
func NewEncoderWithPtr(ps ...any) *Encoder {
	l := len(ps)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(ps[i])
		if rt.Kind() != reflect.Ptr {
			panic("must a pointer type!")
		}
		engines[i] = getEncEngine(rt.Elem())
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}
}

// 创建一个编码is 类型的编码器
func NewEncoder(is ...any) *Encoder {
	l := len(is)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		engines[i] = getEncEngine(reflect.TypeOf(is[i]))
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}
}

func NewEncoderWithType(ts ...reflect.Type) *Encoder {
	l := len(ts)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		engines[i] = getEncEngine(ts[i])
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}
}

// 入参是要编码值的指针
func (e *Encoder) encode(is ...any) []byte {
	engines := e.engines
	for i := 0; i < len(engines) && i < len(is); i++ {
		engines[i](e, reflect.ValueOf(is[i]).UnsafePointer())
	}
	return e.reset()
}

// vs 是持有要编码的值
func (e *Encoder) encodeValue(vs ...reflect.Value) []byte {
	engines := e.engines
	for i := 0; i < len(engines) && i < len(vs); i++ {
		engines[i](e, getUnsafePointer(vs[i]))
	}
	return e.reset()
}

// AppendTo 编码产生的数据将append到buf上
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
