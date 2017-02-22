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

	//decEng  decEng
	decEngs []decEng //解码器集合
	length  int      //解码器数量
}

func NewDecoder(is ...interface{}) *Decoder {
	l := len(is)
	if l < 1 {
		panic("must have argument!")
	}
	des := make([]decEng, l)
	for i := 0; i < l; i++ {
		des[i] = GetDecEngine(reflect.TypeOf(is[i]))
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
		des[i] = GetDecEngine(ts[i])
	}
	return &Decoder{
		length:  l,
		decEngs: des,
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

// is is pointer of value
func (d *Decoder) Decodes(is ...interface{}) {
	l, engs := d.length, d.decEngs
	for i := 0; i < l; i++ {
		engs[i](d, reflect.ValueOf(is[i]).Elem())
	}
}

func (d *Decoder) DecodeValues(vs ...reflect.Value) {
	l, engs := d.length, d.decEngs
	for i := 0; i < l; i++ {
		engs[i](d, vs[i])
	}
}
