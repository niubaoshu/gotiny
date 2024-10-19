package gotiny

import (
	"reflect"
	"unsafe"
)

type Decoder struct {
	buf     []byte //buf
	index   int    //下一个要使用的字节在buf中的下标
	boolPos byte   //下一次要读取的bool在buf中的下标,即buf[boolPos]
	boolBit byte   //下一次要读取的bool的buf[boolPos]中的bit位

	engines []decEng //解码器集合
	length  int      //解码器数量
}

func Unmarshal(buf []byte, is ...any) int {
	return NewDecoderWithPtr(is...).decode(buf, is...)
}

func NewDecoderWithPtr(is ...any) *Decoder {
	l := len(is)
	engines := make([]decEng, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(is[i])
		if rt.Kind() != reflect.Ptr {
			panic("the argument must be a pointer type!")
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

// Decode takes a byte slice and a variable number of pointers to variables.
// It decodes the byte slice into the variables.
// the arguments  must be a pointer type
// The return value is the number of bytes that were decoded.
func (d *Decoder) decode(buf []byte, is ...any) int {
	d.buf = buf
	engines := d.engines
	for i := 0; i < len(engines) && i < len(is); i++ {
		engines[i](d, reflect.ValueOf(is[i]).UnsafePointer())
	}
	return d.reset()
}

// DecodeValue takes a byte slice and a variable number of reflect.Values.
// It decodes the byte slice into the reflect.Values.
// The return value is the number of bytes that were decoded.
func (d *Decoder) decodeValue(buf []byte, vs ...reflect.Value) int {
	d.buf = buf
	engines := d.engines
	for i := 0; i < len(engines) && i < len(vs); i++ {
		engines[i](d, unsafe.Pointer(vs[i].UnsafeAddr()))
	}
	return d.reset()
}
