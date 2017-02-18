package gotiny

import (
	"reflect"
	"sync"
)

type (
	Pool interface {
		Get(min int) []byte //从Pool获得一个存储空间大于等于min的[]byte
		Put([]byte)         //将一个[]byte放回Pool
	}

	Encoder struct {
		buf    []byte //编码目的数组
		offset int    //从buf[offset]开始写数据
		index  int    //下一个要写的下标
		length int    //len(buf)，真实申请到的切片长度
		pool   Pool   //byte切片池

		reqLen   int //下一次申请buf的长度,单位是位,初始值是预算的长度，在此基础上进行增加
		reserved int //为之后的可预算长度的类型预留空间，单位是位

		boolPos int  //下一次要设置的bool在buf中的下标,即buf[boolPos]
		boolBit byte //下一次要设置的bool的buf[boolPos]中的bit位

		engine      func(...interface{})
		engineValue func(...reflect.Value)
	}
)

type bytePool struct {
	pool *sync.Pool
}

func newBytePool() *bytePool {
	return &bytePool{pool: &sync.Pool{}}
}

func (p *bytePool) Get(min int) []byte {
	for i := 0; i < 2; i++ {
		if buf := p.pool.Get(); buf == nil {
			return make([]byte, min)
		} else {
			b := buf.([]byte)
			if len(b) >= min {
				return b
			}
		}
	}
	return make([]byte, min)
}

func (p *bytePool) Put(buf []byte) {
	p.pool.Put(buf)
}

//off 是字节从缓存的off字节开始存储编码的内容
// func NewEncoderWithPool(p Pool, off int) *Encoder {
// 	return &Encoder{
// 		pool:   p,
// 		offset: off,
// 		index:  off,
// 	}
// }

func NewEncoder(is ...interface{}) (e *Encoder) {
	l := len(is)
	if l < 1 {
		panic("must have argument!")
	}
	ts := make([]reflect.Type, l)
	for i := 0; i < len(ts); i++ {
		ts[i] = reflect.TypeOf(is[i])
	}
	return NewEncoderWithType(ts...)
}
func NewEncoderWithType(ts ...reflect.Type) (e *Encoder) {
	l := len(ts)
	if l < 1 {
		panic("must have argument!")
	}
	engines, j, interval := make([]func(*Encoder, reflect.Value), l), 0, [][2]int{{0, 0}}
	isStatic := true // 默认为真，下面有假则为假
	length := 0
	var info *TypeInfo
	for i := 0; i < l; i++ {
		info = GetTypeInfo(ts[i])
		engines[i] = info.Engine
		length += info.Length
		interval[j][0] += info.Head //Head is all
		if !info.IsStatic {
			isStatic = info.IsStatic
			//if i != lessOne {
			interval = append(interval, [2]int{0, i}) //append 1 Length
			j++
			//}
			//interval[j][1] = i
		}
	}
	interval = append(interval, [2]int{0, l})
	j++
	if interval[0][0] > length { //interval[0][0] 就是head
		length = interval[0][0]
	}
	e = &Encoder{pool: newBytePool()}
	if isStatic {
		e.engine = func(is ...interface{}) {
			e.reset(length)
			for i := 0; i < l; i++ {
				engines[i](e, reflect.ValueOf(is[i]))
			}
		}
		e.engineValue = func(vs ...reflect.Value) {
			e.reset(length)
			for i := 0; i < l; i++ {
				engines[i](e, vs[i])
			}
		}
	} else {
		e.engine = func(is ...interface{}) {
			e.reset(length)
			for k, i := 1, 0; k < j; k++ {
				for e.reserved = interval[k][0]; i <= interval[k][1]; i++ {
					engines[i](e, reflect.ValueOf(is[i]))
				}
			}
		}
		e.engineValue = func(vs ...reflect.Value) {
			e.reset(length)
			for k, i := 1, 0; k < j; k++ {
				for e.reserved = interval[k][0]; i <= interval[k][1]; i++ {
					engines[i](e, vs[i])
				}
			}
		}
	}
	return
}

func (e *Encoder) SetOff(off int) {
	e.offset = off
}
func (e *Encoder) SetPool(p Pool) {
	e.pool = p
}

func (e *Encoder) Encodes(is ...interface{}) []byte {
	e.engine(is...)
	return e.buf[:e.index]
}

func (e *Encoder) EncodeValues(vs ...reflect.Value) []byte {
	e.engineValue(vs...)
	return e.buf[:e.index]
}

//检测是否有l的空间，没有就扩展,l的单位是字节
func (e *Encoder) append(l int) {
	if l&0x07 == 0 {
		l = l >> 3
	} else {
		l = l>>3 + 1
	}
	if e.index+l > e.length && e.reqLen > e.length<<3 {
		var buf []byte
		if e.reqLen&0x07 == 0 {
			buf = e.pool.Get(e.reqLen >> 3)
		} else {
			buf = e.pool.Get(e.reqLen>>3 + 1)
		}
		e.length = len(buf)
		copy(buf, e.buf)
		e.pool.Put(e.buf)
		e.buf = buf
	}
}

func (e *Encoder) reset(length int) {
	e.reqLen = length + e.offset<<3
	if e.reqLen > e.length<<3 {
		//fmt.Println(e.reqLen, e.length)
		if e.reqLen&0x07 == 0 {
			e.buf = e.pool.Get(e.reqLen >> 3)
		} else {
			e.buf = e.pool.Get(e.reqLen>>3 + 1)
		}
	}
	e.index = e.offset
	e.length = len(e.buf)
	e.reserved = 0
	e.boolBit = 0
	e.boolPos = 0
}
