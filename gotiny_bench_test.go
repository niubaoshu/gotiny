package gotiny_test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/niubaoshu/gotiny"
)

var (
	ba = 123456

	bptr = &ba
	bmap = map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}

	bvs = []interface{}{
		//true,
		//false,
		//int(-123456),
		//int8(-123),
		//int16(-12345),
		//int32(-123456),
		//int64(-123456),
		//uint(123456),
		//uint8(123),
		//uint16(12345),
		//uint32(123456),
		//uint64(123456),
		//uintptr(123456),
		//unsafe.Pointer(uintptr(123456)),
		//float32(1.23456),
		//float64(1.23456),
		//complex(1.23456, 1.23456),
		//complex(1.23456, 1.23456),
		//string("hello,世界!"),
		//[]byte("hello,golang"),
		//struct{}{},
		//time.Now(),
		//gotinytest("gotiny"),
		//bptr,
		//[6]int{1, 2, 3, 4, 5, 6},
		//[]int{1, 2, 3, 4, 5, 6},
		bmap,
		struct{ int }{123456},
	}
	blength = len(bvs)
	bsrcp   = make([]unsafe.Pointer, blength)
	bretp   = make([]unsafe.Pointer, blength)
	bbuf    []byte

	be = gotiny.NewEncoder(bvs...)
	bd = gotiny.NewDecoder(bvs...)
)

func init() {
	pvs := make([]interface{}, blength)
	for i := 0; i < blength; i++ {
		typ := reflect.TypeOf(bvs[i])
		ret := reflect.New(typ).Elem().Interface()
		bsrcp[i] = unsafe.Pointer(reflect.ValueOf(&bvs[i]).Elem().InterfaceData()[1])
		bretp[i] = unsafe.Pointer(reflect.ValueOf(&ret).Elem().InterfaceData()[1])
		val := reflect.New(typ)
		val.Elem().Set(reflect.ValueOf(bvs[i]))
		pvs[i] = val.Interface()
	}

	bbuf = gotiny.Encodes(pvs...)
}

func BenchmarkEncPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		be.EncodePtr(bsrcp...)
	}
}

func BenchmarkDecPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bd.DecodePtr(bbuf, bretp...)
	}
}
