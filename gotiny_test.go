package gotiny_test

import (
	"io"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"fmt"

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
		inter       interface{}
		A
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

	tint int

	gotinytest string
)

func (tint) Read([]byte) (int, error)  { return 0, nil }
func (tint) Write([]byte) (int, error) { return 0, nil }
func (tint) Close() error              { return nil }

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
		inter:       interface{}(int(1)),
		A:           genA(),
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
	vint3       = tint(1234567)
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
	vslicebytes = [][]byte{[]byte("aaaaaaaaaaaaaaaaaaa"), []byte("bbbbbbbbbbbbbbb"), []byte("ccccccccccccc")}
	v2slice     = []int{1, 2, 3, 4, 5}
	v3slice     []byte
	varr        = [3]baseTyp{genBase(), genBase(), genBase()}
	vmap        = map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}
	v2map       = map[int]map[int]int{1: {2: 3, 3: 4}}
	v3map       = map[int][]byte{1: {2, 3, 3, 4}}
	v4map       = map[int]*int{1: &vint}
	v5map       = map[int]baseTyp{1: genBase(), 2: genBase()}
	v6map       = map[*int]baseTyp{&vint1: genBase(), &vint2: genBase()}
	v7map       = map[int][3]baseTyp{1: varr}
	vnilmap     map[int]int
	vptr        = &vint
	vsliceptr   = &vbytes
	vptrslice   = []*int{&vint, &vint, &vint}
	vnilptr     *int
	v2nilptr    []string
	vnilptrptr  = &vnilptr
	varrptr     = &varr
	vtime       = time.Now()

	vslicebase = []baseTyp{
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

	unsafePointer = unsafe.Pointer(&vtime)

	vcir        cirTyp
	v2cir       cirTyp = &vcir
	vcirStruct         = cirStruct{a: 1}
	v2cirStruct        = cirStruct{a: 1, cirStruct: &vcirStruct}

	vcirmap  = cirmap{1: nil}
	v2cirmap = cirmap{2: vcirmap}

	vAstruct = genA()

	vGotinyTest  = gotinytest("aaaaaaaaaaaaaaaaaaaaa")
	v2GotinyTest = &vGotinyTest

	vbinTest, _ = url.Parse("http://www.baidu.com/s?wd=234234")
	v2binTest   = &vbinTest

	v0interface interface{}
	vinterface  interface{}        = varray
	v1interface io.ReadWriteCloser = tint(2)
	v2interface io.ReadWriteCloser = os.Stdin
	v3interface interface{}        = &vinterface
	v4interface interface{}        = &v1interface
	v5interface interface{}        = &v2interface
	v6interface interface{}        = &v3interface
	v7interface interface{}        = &v0interface
	v8interface interface{}        = &vnilptr
	v9interface interface{}        = &v8interface

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
		vint3,
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
		vslicebytes,
		v2slice,
		v3slice,
		varr,
		vmap,
		v2map,
		v3map,
		v4map,
		v5map,
		v6map,
		v7map,
		vnilmap,
		vptr,
		vsliceptr,
		vptrslice,
		vnilptr,
		v2nilptr,
		vnilptrptr,
		varrptr,
		vtime,
		vslicebase,
		vslicestring,
		varray,
		vinterface,
		v1interface,
		v2interface,
		v3interface,
		v4interface,
		v5interface,
		v6interface,
		v7interface,
		v8interface,
		v9interface,
		unsafePointer,
		vcir,
		v2cir,
		vcirStruct,
		v2cirStruct,
		vcirmap,
		v2cirmap,
		vAstruct,
		vGotinyTest,
		v2GotinyTest,
		vbinTest,
		v2binTest,
		struct{}{},
	}

	length = len(vs)
	buf    = make([]byte, 0, 1<<13)
	e      = gotiny.NewEncoder(vs...)
	d      = gotiny.NewDecoder(vs...)
	c      = goutils.NewComparer()

	srci = make([]interface{}, length)
	reti = make([]interface{}, length)
	srcv = make([]reflect.Value, length)
	retv = make([]reflect.Value, length)
	srcp = make([]unsafe.Pointer, length)
	retp = make([]unsafe.Pointer, length)
	typs = make([]reflect.Type, length)
)

func init() {
	e.AppendTo(buf)

	for i := 0; i < length; i++ {
		typs[i] = reflect.TypeOf(vs[i])
		srcv[i] = reflect.ValueOf(vs[i])

		tempi := reflect.New(typs[i])
		tempi.Elem().Set(srcv[i])
		srci[i] = tempi.Interface()

		tempv := reflect.New(typs[i])
		retv[i] = tempv.Elem()
		reti[i] = tempv.Interface()

		srcp[i] = unsafe.Pointer(reflect.ValueOf(&srci[i]).Elem().InterfaceData()[1])
		retp[i] = unsafe.Pointer(reflect.ValueOf(&reti[i]).Elem().InterfaceData()[1])
	}
	fmt.Printf("total %d value. buf length: %d, encode length: %d \n", length, cap(buf), len(gotiny.Encodes(srci...)))
}

func TestEncodeDecode(t *testing.T) {
	gotiny.Decodes(gotiny.Encodes(srci...), reti...)
	for i, r := range reti {
		Assert(t, srci[i], r)
	}
}

func TestInterface(t *testing.T) {
	d.Decode(e.Encode(srci...), reti...)
	for i, r := range reti {
		Assert(t, srci[i], r)
	}
}

func TestPtr(t *testing.T) {
	d.DecodePtr(e.EncodePtr(srcp...), retp...)
	for i, r := range reti {
		Assert(t, srci[i], r)
	}
}

func TestValue(t *testing.T) {
	d.DecodeValue(e.EncodeValue(srcv...), retv...)
	for i, r := range reti {
		Assert(t, srci[i], r)
	}
}

//func BenchmarkEncodes(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		buf = gotiny.Encodes(srci...)
//	}
//}
//
//func BenchmarkDecodes(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		gotiny.Decodes(buf, reti...)
//	}
//}
//
//func BenchmarkEncodesValue(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		e.EncodeValue(srcv...)
//	}
//}
//
//func BenchmarkDecodesValue(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		d.DecodeValue(buf, retv...)
//	}
//}
//
//func BenchmarkEncodesPtr(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		e.EncodePtr(srcp...)
//	}
//}
//
//func BenchmarkDecodesPtr(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		d.DecodePtr(buf, retp...)
//	}
//}
//
//func BenchmarkEncodesInterface(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		e.Encode(srci...)
//	}
//}
//
//func BenchmarkDecodesInterface(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		d.Decode(buf, reti...)
//	}
//}

func Assert(t *testing.T, x, y interface{}) {
	if !c.DeepEqual(x, y) {
		t.Errorf("\n exp type = %T; value = %+v;\n got type = %T; value = %+v \n", x, x, y, y)
	}
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

func TestGetName(t *testing.T) {
	stdin := (interface {
		Read([]byte) (int, error)
		Write([]byte) (int, error)
	})(os.Stdin)
	items := []struct {
		ret string
		val interface{}
	}{
		{"int", int(1)},
		{"*int", (*int)(nil)},
		{"**int", (**int)(nil)},
		{"[]int", []int{}},
		{"[]time.Time", []time.Time{}},
		{"[]github.com/niubaoshu/gotiny.GoTinySerializer", []gotiny.GoTinySerializer{}},
		{"*interface {}", (*interface{})(nil)},
		{"map[int]string", map[int]string{}},
		{"struct { a struct { int; b int; dec []github.com/niubaoshu/gotiny.Decoder; abb interface {}; c io.ReadWriteCloser } }",
			struct {
				a struct {
					int
					b   int
					dec []gotiny.Decoder
					abb interface{}
					c   io.ReadWriteCloser
				}
			}{}},
		{"struct {}", struct{}{}},
		{"*interface { Read([]uint8) (int, error); Write([]uint8) (int, error) }", &stdin},
	}
	for _, item := range items {
		r := reflect.TypeOf(item.val)
		if string(gotiny.GetName(item.val)) != item.ret {
			t.Logf("real: %s , exp: %s", gotiny.GetName(item.val), item.ret)
			t.Fatalf("string:%s,name:%s,pkgpath:%s,fmt %T", r.String(), r.Name(), r.PkgPath(), item.val)
		}
	}
}
