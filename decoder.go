package gotiny

import (
	"reflect"
)

const (
	maxVarintLen64 = 10
)

type Decoder struct {
	buf     []byte //buf
	index   int    //下一个要读取的字节
	boolean byte   //下一次要读取的bool在buf中的下标,即buf[boolPos]
	boolBit byte   //下一次要读取的bool的buf[boolPos]中的bit位

	//decEngine  decEngine
	decEngines []decEngine //解码器集合
	length     int         //解码器数量
}

//n is index
//func NewDecoder(b []byte) *Decoder {
//	return &Decoder{buf: b}
//}

func NewDecoder(is ...interface{}) *Decoder {
	l := len(is)
	des := make([]decEngine, l)
	for i := 0; i < l; i++ {
		des[i] = GetDecEngine(reflect.TypeOf(is[i]))
	}
	return &Decoder{
		length:     l,
		decEngines: des,
	}
}

func NewDecoderWithTypes(ts ...reflect.Type) *Decoder {
	l := len(ts)
	des := make([]decEngine, l)
	for i := 0; i < l; i++ {
		des[i] = GetDecEngine(ts[i])
	}
	return &Decoder{
		length:     l,
		decEngines: des,
	}
}

func (d *Decoder) GetUnusedBytes() []byte {
	return d.buf[d.index:]
}
func (d *Decoder) Reset() {
	d.index = 0
	d.boolean = 0
	d.boolBit = 0
}

func (d *Decoder) ResetWith(b []byte) {
	d.buf = b
	d.Reset()
}

// func (d *Decoder) DecodeByType(t reflect.Type) (v reflect.Value) {
// 	v = reflect.New(t).Elem()
// 	GetDecEngine(t)(d, v)
// 	return
// }

func (d *Decoder) DecodeByTypes(ts ...reflect.Type) (vs []reflect.Value) {
	vs = make([]reflect.Value, len(ts))
	for i := 0; i < d.length; i++ {
		vs[i] = reflect.New(ts[i]).Elem()
		d.decEngines[i](d, vs[i])
	}
	return
}

// func (d *Decoder) Decode(i interface{}) {
// 	v := reflect.ValueOf(i)
// 	if v.Kind() != reflect.Ptr || v.IsNil() { // must be a ptr but nilptr
// 		panic("totiny: only decode to pointer type, and not nilpointer")
// 	}
// 	GetDecEngine(v.Elem().Type())(d, v.Elem())
// }

// is is pointer of value
func (d *Decoder) Decodes(is ...interface{}) {
	for i := 0; i < d.length; i++ {
		d.decEngines[i](d, reflect.ValueOf(is[i]).Elem())
	}
}

// func (d *Decoder) DecodeValue(v reflect.Value) {
// 	GetDecEngine(v.Type())(d, v)
// }

func (d *Decoder) DecodeValues(vs ...reflect.Value) {
	for i := 0; i < d.length; i++ {
		d.decEngines[i](d, vs[i])
	}
}
