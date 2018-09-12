package gotiny_test

import (
	"bytes"
	"encoding"
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
		fBool       bool
		fInt8       int8
		fInt16      int16
		fInt32      int32
		fInt64      int64
		fInt        int
		fUint8      uint8
		fUint16     uint16
		fUint32     uint32
		fUint64     uint64
		fUint       uint
		fUintptr    uintptr
		fFloat32    float32
		fFloat64    float64
		fComplex64  complex64
		fComplex128 complex128
		fString     string
		array       [3]uint32
		inter       interface{}
		A
	}

	A struct {
		Name     string
		BirthDay time.Time
		Phone    string `gotiny:"-"`
		Siblings int
		Spouse   bool
		Money    float64
	}

	cirTyp    *cirTyp
	cirStruct struct {
		a int
		*cirStruct
	}
	cirMap   map[int]cirMap
	cirSlice []cirSlice

	tint int

	gotinyTest string
)

func (tint) Read([]byte) (int, error)  { return 0, nil }
func (tint) Write([]byte) (int, error) { return 0, nil }
func (tint) Close() error              { return nil }

func (v *gotinyTest) GotinyEncode(buf []byte) []byte {
	return append(buf, gotiny.Marshal((*string)(v))...)
}

func (v *gotinyTest) GotinyDecode(buf []byte) int {
	return gotiny.Unmarshal(buf, (*string)(v))
}

func genBase() baseTyp {
	return baseTyp{
		fBool:       rand.Int()%2 == 0,
		fInt8:       int8(rand.Int()),
		fInt16:      int16(rand.Int()),
		fInt32:      int32(rand.Int()),
		fInt64:      int64(rand.Int()),
		fInt:        int(rand.Int()),
		fUint8:      uint8(rand.Int()),
		fUint16:     uint16(rand.Int()),
		fUint64:     uint64(rand.Int()),
		fUintptr:    uintptr(rand.Int()),
		fFloat32:    rand.Float32(),
		fFloat64:    rand.Float64(),
		fComplex64:  complex(rand.Float32(), rand.Float32()),
		fComplex128: complex(rand.Float64(), rand.Float64()),
		fString:     getRandomString(20 + rand.Intn(256)),
		array:       [3]uint32{rand.Uint32(), rand.Uint32()},
		inter:       interface{}(int(1)),
		A:           genA(),
	}
}

func genA() A {
	return A{
		Name:     getRandomString(16),
		BirthDay: time.Now(),
		//Phone:    getRandomString(10),
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
	v4uint64    = v2uint64 - uint64(rand.Intn(200))
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
	v3cir       cirTyp = &v2cir
	vcirStruct         = cirStruct{a: 1}
	v2cirStruct        = cirStruct{a: 1, cirStruct: &vcirStruct}
	vcirmap            = cirMap{1: nil}
	v2cirmap           = cirMap{2: vcirmap}
	v1cirSlice         = make([]cirSlice, 10)
	v2cirSlice         = append(v1cirSlice, v1cirSlice)
	v3cirSlice         = append(v2cirSlice, v1cirSlice)
	v4cirSlice         = append(v1cirSlice, v1cirSlice, v2cirSlice, v3cirSlice)

	vAstruct = genA()

	vGotinyTest  = gotinyTest("aaaaaaaaaaaaaaaaaaaaa")
	v2GotinyTest = &vGotinyTest

	vbinTest, _ = url.Parse("http://www.baidu.com/s?wd=234234")
	v2binTest   interface {
		encoding.BinaryMarshaler
		encoding.BinaryUnmarshaler
	} = vbinTest

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
		v4uint64,
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
		v3cir,
		vcirStruct,
		v2cirStruct,
		vcirmap,
		v2cirmap,
		v1cirSlice,
		v2cirSlice,
		v3cirSlice,
		v4cirSlice,
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
	fmt.Printf("total %d value. buf length: %d, encode length: %d \n", length, cap(buf), len(gotiny.Marshal(srci...)))
}

func TestEncodeDecode(t *testing.T) {
	gotiny.Unmarshal(gotiny.Marshal(srci...), reti...)
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

func TestMap(t *testing.T) {
	var sm = map[string]int{
		"a": 1,
	}
	var rm = map[string]int{
		//"b": 2,
	}
	buf := gotiny.Marshal(&sm)
	gotiny.Unmarshal(buf, &rm)
	Assert(t, sm, rm)
}

func TestHelloWorld(t *testing.T) {
	hello, world := []byte("hello, "), "world"
	hello2, world2 := []byte("1"), ""

	gotiny.Unmarshal(gotiny.Marshal(&hello, &world), &hello2, &world2)
	if !bytes.Equal(hello2, hello) || world2 != world {
		t.Error(hello2, world2)
	}
}

func Assert(t *testing.T, x, y interface{}) {
	if !c.DeepEqual(x, y) {
		e, g := indirect(x), indirect(y)
		t.Errorf("\n exp type = %T; value = %+v;\n got type = %T; value = %+v; \n", e, e, g, g)
	}
}

func indirect(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v.Interface()
}

func getRandomString(l int) string {
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result[i] = bytes[r.Intn(62)]
	}
	return string(result)
}

func TestGetName(t *testing.T) {
	stdin := (interface {
		Read([]byte) (int, error)
		Write([]byte) (int, error)
	})(os.Stdin)
	nt := newType()
	items := []struct {
		ret string
		val interface{}
	}{
		{"int", int(1)},
		{"github.com/niubaoshu/gotiny.Encoder", gotiny.Encoder{}},
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
		{"func(int) (int, error)", func(i int) (int, error) { return 0, nil }},
		{"func(int)", func(i int) {}},
		{"func(int) error", func(i int) error { return nil }},
		{"struct { A int }", nt},
		{"<nil>", nil},
	}
	for _, item := range items {
		r := reflect.TypeOf(item.val)
		if string(gotiny.GetName(item.val)) != item.ret {
			t.Logf("real: %s , exp: %s", gotiny.GetName(item.val), item.ret)
			t.Fatalf("string:%s,name:%s,pkgpath:%s,fmt %T", r.String(), r.Name(), r.PkgPath(), item.val)
		}
	}
}

func newType() struct {
	A int
} {
	return struct{ A int }{A: 1}
}
