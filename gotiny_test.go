package gotiny_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/niubaoshu/gotiny"
	"github.com/niubaoshu/goutils"
)

type (
	baseTyp struct {
		fbool       bool
		fint8       int8
		fint16      int16
		fint32      int32
		fint64      int64
		fint        int
		fuint8      uint8
		fuint16     uint16
		fuint32     uint32
		fuint64     uint64
		fuint       uint
		fuintptr    uintptr
		ffloat32    float32
		ffloat64    float64
		fcomplex64  complex64
		fcomplex128 complex128
		fstring     string
		array       [3]uint32
	}

	A struct {
		Name     string
		BirthDay time.Time
		Phone    string
		Siblings int
		Spouse   bool
		Money    float64
	}

	cirTyp    *cirTyp
	cirStruct struct {
		a int
		*cirStruct
	}
	cirmap map[int]cirmap

	gotinytest string
)

func (v *gotinytest) GotinyEncode(buf []byte) []byte {
	return append(buf, gotiny.Encodes((*string)(v))...)
}

func (v *gotinytest) GotinyDecode(buf []byte) int {
	return gotiny.Decodes(buf, (*string)(v))
}

func genBase() baseTyp {
	return baseTyp{
		fbool:       rand.Int()%2 == 0,
		fint8:       int8(rand.Int()),
		fint16:      int16(rand.Int()),
		fint32:      int32(rand.Int()),
		fint64:      int64(rand.Int()),
		fint:        int(rand.Int()),
		fuint8:      uint8(rand.Int()),
		fuint16:     uint16(rand.Int()),
		fuint64:     uint64(rand.Int()),
		fuintptr:    uintptr(rand.Int()),
		ffloat32:    rand.Float32(),
		ffloat64:    rand.Float64(),
		fcomplex64:  complex(rand.Float32(), rand.Float32()),
		fcomplex128: complex(rand.Float64(), rand.Float64()),
		fstring:     GetRandomString(20 + rand.Intn(256)),
		array:       [3]uint32{rand.Uint32(), rand.Uint32()},
	}
}

func genA() A {
	return A{
		Name:     GetRandomString(16),
		BirthDay: time.Now(),
		Phone:    GetRandomString(10),
		Siblings: rand.Intn(5),
		Spouse:   rand.Intn(2) == 1,
		Money:    rand.Float64(),
	}
}

var (
	vbool       = true
	vfbool      = false
	vint8       = int8(123)
	vint16      = int16(-12345)
	vint32      = int32(123456)
	vint64      = int64(-1234567)
	v2int64     = int64(1<<63 - 1)
	v3int64     = int64(rand.Int63())
	vint        = int(123456)
	vint1       = int(123456)
	vint2       = int(1234567)
	vuint       = uint(123)
	vuint8      = uint8(123)
	vuint16     = uint16(12345)
	vuint32     = uint32(123456)
	vuint64     = uint64(1234567)
	v2uint64    = uint64(1<<64 - 1)
	v3uint64    = uint64(rand.Uint32() * rand.Uint32())
	vuintptr    = uintptr(12345678)
	vfloat32    = float32(1.2345)
	vfloat64    = float64(1.2345678)
	vcomp64     = complex(1.2345, 2.3456)
	vcomp128    = complex(1.2345678, 2.3456789)
	vstring     = string("hello,日本国")
	base        = genBase()
	vbytes      = []byte("aaaaaaaaaaaaaaaaaaa")
	vsliecbytes = [][]byte{[]byte("aaaaaaaaaaaaaaaaaaa"), []byte("bbbbbbbbbbbbbbb"), []byte("ccccccccccccc")}
	vmap        = map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}
	v2map       = map[int]map[int]int{1: {2: 3, 3: 4}}
	v3map       = map[int][]byte{1: {2, 3, 3, 4}}
	temp        = 1
	v4map       = map[int]*int{1: &temp}
	v5map       = map[int]baseTyp{1: genBase(), 2: genBase()}
	v6map       = map[*int]baseTyp{&vint1: genBase(), &vint2: genBase()}
	vnilmap     map[int]int
	vptr        = &vint
	vsliceptr   = &vbytes
	vptrslice   = []*int{&vint, &vint, &vint}
	vnilptr     *int
	vnilptrptr  = &vnilptr
	vtime       = time.Now()
	vslicebase  = []baseTyp{
		genBase(),
		genBase(),
		genBase(),
	}
	vslicestring = []string{
		"aaaaaaaaa",
		"bbbbbbbbb",
		"ccccccccc",
	}

	varray = [3]baseTyp{
		genBase(),
		genBase(),
		genBase(),
	}
	vcir  cirTyp
	v2cir cirTyp = &vcir

	vcirStruct  = cirStruct{a: 1}
	v2cirStruct = cirStruct{a: 1, cirStruct: &vcirStruct}

	vcirmap  = cirmap{1: nil}
	v2cirmap = cirmap{2: vcirmap}

	vAstruct = genA()

	vGotinyTest  = gotinytest("aaaaaaaaaaaaaaaaaaaaa")
	v2GotinyTest = &vGotinyTest

	vs = []interface{}{
		vbool,
		vfbool,
		vint8,
		vint16,
		vint32,
		vint64,
		v2int64,
		v3int64,
		vint,
		vint1,
		vint2,
		vuint,
		vuint8,
		vuint16,
		vuint32,
		vuint64,
		v2uint64,
		v3uint64,
		vuintptr,
		vfloat32,
		vfloat64,
		vcomp64,
		vcomp128,
		vstring,
		base,
		vbytes,
		vsliecbytes,
		vmap,
		v2map,
		v3map,
		v4map,
		v5map,
		v6map,
		vnilmap,
		vptr,
		vsliceptr,
		vptrslice,
		vnilptr,
		vnilptrptr,
		vtime,
		vslicebase,
		vslicestring,
		varray,
		vcir,
		v2cir,
		vcirStruct,
		v2cirStruct,
		vcirmap,
		v2cirmap,
		vAstruct,
		vGotinyTest,
		v2GotinyTest,
		struct{}{},
	}

	e      = gotiny.NewEncoder(vs...)
	d      = gotiny.NewDecoder(vs...)
	length = len(vs)

	srci = make([]interface{}, length)
	reti = make([]interface{}, length)
	srcv = make([]reflect.Value, length)
	retv = make([]reflect.Value, length)
	srcp = make([]unsafe.Pointer, length)
	retp = make([]unsafe.Pointer, length)

	c = goutils.NewComparer()
)

func init() {
	fmt.Println("total", length, "value")
	for i := 0; i < length; i++ {
		typ := reflect.TypeOf(vs[i])
		srcv[i] = reflect.ValueOf(vs[i])
		tempv := reflect.New(typ)
		retv[i] = tempv.Elem()
		tempi := reflect.New(typ)
		tempi.Elem().Set(srcv[i])
		srci[i] = tempi.Interface()
		reti[i] = tempv.Interface()
		srcp[i] = unsafe.Pointer(reflect.ValueOf(&srci[i]).Elem().InterfaceData()[1])
		retp[i] = unsafe.Pointer(reflect.ValueOf(&reti[i]).Elem().InterfaceData()[1])
	}
	e.ResetWithBuf(make([]byte, 0, 2048))
}

func TestInterface(t *testing.T) {
	e.Reset()
	e.Encodes(srci...)
	d.ResetWith(e.Bytes())
	d.Decodes(reti...)
	for i, r := range reti {
		Assert(t, srci[i], r)
	}
}

func TestEncodeDecode(t *testing.T) {
	gotiny.Decodes(gotiny.Encodes(srci...), reti...)
	for i, r := range reti {
		Assert(t, srci[i], r)
	}
}

func TestPtr(t *testing.T) {
	e.EncodeByUPtrs(srcp...)
	b := e.Bytes()
	fmt.Printf("length: %d \n", len(b))
	d.ResetWith(b)
	d.DecodeByUPtr(retp...)
	for i, r := range reti {
		Assert(t, srci[i], r)
	}

}

func TestValue(t *testing.T) {
	e.Reset()
	e.EncodeValues(srcv...)
	d.ResetWith(e.Bytes())
	d.DecodeValues(retv...)
	for i, r := range reti {
		Assert(t, srci[i], r)
	}

}

//
//var enctest = NewEncoder(1)
//var enctest2 = NewEncoder(1, 2, 3, 4, 5, 6, 7)
//
//func BenchmarkFIEnc(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		encitest(1, 2, 3, 3, 3, 4, 5, 2, 23)
//	}
//}
//
//func BenchmarkJson(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		jsonf(1)
//	}
//}
//func jsonf(is interface{}) {
//	json.Marshal(is)
//}
//
//func encitest(is ...interface{}) {
//	enctest.Encodes(is...)
//}
//
//func BenchmarkFPEnc(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		encptest()
//	}
//}
//
//func encptest() {
//	a0, a1, a2, a3, a4, a5, a6 := 1, 2, 3, 4, 5, 6, 7
//	enctest2.EncodeByUPtrs(unsafe.Pointer(&a0), unsafe.Pointer(&a1), unsafe.Pointer(&a2),
//		unsafe.Pointer(&a3), unsafe.Pointer(&a4), unsafe.Pointer(&a5), unsafe.Pointer(&a6))
//}

var buf []byte

func BenchmarkEncodes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf = gotiny.Encodes(srci...)
	}
}

func BenchmarkDecodes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gotiny.Decodes(buf, reti...)
	}
}

func BenchmarkEncodesPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.Reset()
		e.EncodeByUPtrs(srcp...)
	}
}

func BenchmarkDecodesPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d.Reset()
		d.DecodeByUPtr(retp...)
	}
}
func BenchmarkEncodesValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.Reset()
		e.EncodeValues(srcv...)
	}
}

func BenchmarkDecodesValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d.Reset()
		d.DecodeValues(retv...)
	}
}
func BenchmarkEncodesInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.Reset()
		e.Encodes(srci...)
	}
}

func BenchmarkDecodesInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d.Reset()
		d.Decodes(reti...)
	}
}

// func BenchmarkFloatToUint(b *testing.B) {
// 	var f = 1.0
// 	for i := 0; i < b.N; i++ {
// 		floatToUint(f)
// 	}
// }
// func BenchmarkIntToUint(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		intToUint(1)
// 	}
// }

//func BenchmarkUintToInt(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		for j := 0; j < 100000; j++ {
//			uintToInt(uint64(i))
//		}
//	}
//}
//
// var (
// 	ee        = NewEncoder(0)
// 	maxuint64 = uint64(1<<64 - 1)
// )

// func BenchmarkEncUint(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		ee.encUint(maxuint64)
// 	}
// }

// func BenchmarkEncUint2(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		ee.encUint(maxuint64)
// 	}
// }

// func BenchmarkDecUint(b *testing.B) {
// 	b.StopTimer()
// 	dd := NewDecoder(ee.Bytes())
// 	dd.Reset()
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		dd.DecUint()
// 	}
// }

func Assert(t *testing.T, x, y interface{}) {
	if !c.DeepEqual(x, y) {
		t.Fatalf("\n exp type =  %T; value = %#v;\n got type = %T; value = %#v ", x, x, y, y)
	}
}

func getPtr(i interface{}) unsafe.Pointer {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		panic("不是指针")
	}
	return unsafe.Pointer(v.Elem().UnsafeAddr())
}
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
