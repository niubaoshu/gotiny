package gotiny

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
)

var testNum = 10000

type (
	bTyp struct {
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
		inter       any
		tA
	}

	tA struct {
		Name string
		// Phone    string `gotiny:"-"`
		Siblings int
		Spouse   bool
		Money    float64
	}
	tint struct {
		a int16
	}

	cirTyp    *cirTyp
	cirStruct struct {
		b tint
		a int
		*cirStruct
	}
	cirMap   map[int]cirMap
	cirSlice []cirSlice

	gotinyTest string
)

func (*tint) Read([]byte) (int, error)  { return 0, nil }
func (*tint) Write([]byte) (int, error) { return 0, nil }
func (*tint) Close() error              { return nil }

func (t *tint) GotinyEncode(buf []byte) []byte { return append(buf, byte(t.a), byte(t.a>>8)) }
func (t *tint) GotinyDecode(buf []byte) int    { t.a = int16(buf[0]) + int16(buf[1])<<8; return 2 }

func (v *gotinyTest) GotinyEncode(buf []byte) []byte {
	return append(buf, Marshal((*string)(v))...)
}

func (v *gotinyTest) GotinyDecode(buf []byte) int {
	return Unmarshal(buf, (*string)(v))
}

func gentBase() []bTyp {
	n := 10
	base := make([]bTyp, n)
	for i := 0; i < n; i++ {
		base[i] = bTyp{
			fBool:       rand.Int()%2 == 0,
			fInt8:       int8(rand.Int()),
			fInt16:      int16(rand.Int()),
			fInt32:      int32(rand.Int()),
			fInt64:      int64(rand.Int()),
			fInt:        rand.Int(),
			fUint8:      uint8(rand.Int()),
			fUint16:     uint16(rand.Int()),
			fUint32:     uint32(rand.Int()),
			fUint64:     uint64(rand.Int()),
			fUint:       uint(rand.Int()),
			fUintptr:    uintptr(rand.Int()),
			fFloat32:    rand.Float32(),
			fFloat64:    rand.Float64(),
			fComplex64:  complex(rand.Float32(), rand.Float32()),
			fComplex128: complex(rand.Float64(), rand.Float64()),
			fString:     randString(20 + rand.Intn(256)),
			array:       [3]uint32{rand.Uint32(), rand.Uint32()},
			inter:       any(1),
			tA:          gentA(),
		}
	}
	return base
}

func gentA() tA {
	return tA{
		Name: randString(16),
		// Phone:    randString(10),
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
	v3int64     = rand.Int63()
	vint        = 123456
	vint1       = 123456
	vint2       = 1234567
	vint3       = tint{a: 123}
	vuint       = uint(123)
	vuint8      = uint8(123)
	vuint16     = uint16(12345)
	vuint32     = uint32(123456)
	vuint64     = uint64(1234567)
	v2uint64    = uint64(1<<64 - 1)
	v3uint64    = uint64(rand.Uint32() * rand.Uint32())
	v4uint64    = v2uint64 - uint64(rand.Intn(200))
	v5uint64    = uint64(1<<7 - 1)
	vuintptr    = uintptr(12345678)
	vfloat32    = float32(1.2345)
	vfloat64    = 1.2345678
	vcomp64     = complex(1.2345, 2.3456)
	vcomp128    = complex(1.2345678, 2.3456789)
	vstring     = "hello,日本国"
	a           = gentA()
	vbytes      = []byte("aaaaaaaaaaaaaaaaaaa")
	vslicebytes = [][]byte{[]byte("aaaaaaaaaaaaaaaaaaa"), []byte("bbbbbbbbbbbbbbb"), []byte("ccccccccccccc")}
	v2slice     = []int{1, 2, 3, 4, 5}
	v3slice     []byte
	varr        = [3][]bTyp{gentBase(), gentBase(), gentBase()}
	vmap        = map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}
	v2map       = map[int]map[int]int{1: {2: 3, 3: 4}}
	v3map       = map[int][]byte{1: {2, 3, 3, 4}}
	v4map       = map[int]*int{1: &vint}
	v5map       = map[int][]bTyp{1: gentBase(), 2: gentBase()}
	v7map       = map[int][3][]bTyp{1: varr}
	vnilmap     map[int]int
	vptr        = &vint
	vsliceptr   = &vbytes
	vptrslice   = []*int{&vint, &vint, &vint}
	vnilptr     *int
	v2nilptr    []string
	vnilptrptr  = &vnilptr
	varrptr     = &varr

	vslicebase = [][]bTyp{
		gentBase(),
		gentBase(),
		gentBase(),
	}
	vslicestring = []string{
		"aaaaaaaaa",
		"bbbbbbbbb",
		"ccccccccc",
	}

	varray = [3][]bTyp{
		gentBase(),
		gentBase(),
		gentBase(),
	}

	vcir        cirTyp
	v2cir       cirTyp = &vcir
	v3cir       cirTyp = &v2cir
	vcirStruct         = cirStruct{b: tint{a: 45}, a: 1}
	v2cirStruct        = cirStruct{b: tint{a: 22}, a: 1, cirStruct: &vcirStruct}
	vcirmap            = cirMap{1: nil}
	v2cirmap           = cirMap{2: vcirmap}
	v1cirSlice         = make([]cirSlice, 10)
	v2cirSlice         = append(v1cirSlice, v1cirSlice)
	v3cirSlice         = append(v2cirSlice, v1cirSlice)
	v4cirSlice         = append(v1cirSlice, v1cirSlice, v2cirSlice, v3cirSlice)

	vAstruct = gentA()

	vGotinyTest  = gotinyTest("aaaaaaaaaaaaaaaaaaaaa")
	v2GotinyTest = &vGotinyTest

	vbinTest, _ = url.Parse("https://www.baidu.com/s?wd=234234")
	v2binTest   interface {
		encoding.BinaryMarshaler
		encoding.BinaryUnmarshaler
	} = vbinTest

	v0interface  any
	vInterface   any                = varray
	v1interface  io.ReadWriteCloser = &tint{a: 1}
	v3interface  any                = &vInterface
	v4interface  any                = &v1interface
	v6interface  any                = &v3interface
	v7interface  any                = &v0interface
	v8interface  any                = &vnilptr
	v9interface  any                = &v8interface
	v10interface Serializer         = &tint{a: int16(rand.Intn(1<<16) - 1<<15)}
	v11interface io.ReadWriter      = &tint{a: int16(rand.Intn(1<<16) - 1<<15)}

	vs = []any{
		vbool,
		vfbool,
		false,
		true,
		[10]bool{false, true, true, false, true, true},
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
		v5uint64,
		vuintptr,
		vfloat32,
		vfloat64,
		vcomp64,
		vcomp128,
		vstring,
		a,
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
		v7map,
		vnilmap,
		vptr,
		vsliceptr,
		vptrslice,
		vnilptr,
		v2nilptr,
		vnilptrptr,
		varrptr,
		vslicebase,
		vslicestring,
		varray,
		vInterface,
		v1interface,
		v3interface,
		v4interface,
		v6interface,
		v7interface,
		v8interface,
		v9interface,
		v10interface,
		v11interface,
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
	tBuf   = make([]byte, 0, 1<<14)
	te     = NewEncoder(vs...)
	td     = NewDecoder(vs...)

	typs = make([]reflect.Type, length)

	srci = make([]any, length)
	reti = make([]any, length)
	srcv = make([]reflect.Value, length)
	retv = make([]reflect.Value, length)
)

// init function, prepare the data for testing
func init() {
	// Make sure the encoder append to the tBuf
	te.AppendTo(tBuf)
	// Create the reflect.Value and reflect.Type for the testing data
	for i := 0; i < length; i++ {
		typs[i] = reflect.TypeOf(vs[i])
		srcv[i] = reflect.ValueOf(vs[i])

		// Create a new value for srci, so that the srcv and srci are not the same
		tempI := reflect.New(typs[i])
		tempI.Elem().Set(srcv[i])
		srci[i] = tempI.Interface()

		// Create a new value for retv, so that the retv and reti are not the same
		tempV := reflect.New(typs[i])
		retv[i] = tempV.Elem()
		reti[i] = tempV.Interface()
	}
	// Print the summary of the testing data
	fmt.Printf("total %d value. tBuf length: %d, encode length: %d \n", length, cap(tBuf), len(Marshal(srci...)))
}

func TestMarshalUnmarshal(t *testing.T) {
	buf := Marshal(srci...)
	Unmarshal(buf, reti...)
	for i, r := range reti {
		Assert(t, buf, srci[i], r)
	}
}

func TestEncodeDecode(t *testing.T) {
	buf := te.encode(srci...)
	td.decode(buf, reti...)
	for i, r := range reti {
		Assert(t, buf, srci[i], r)
	}
}

func TestValue(t *testing.T) {
	td.decodeValue(te.encodeValue(srcv...), retv...)
	for i, r := range reti {
		Assert(t, tBuf, srci[i], r)
	}
}

func TestMap(t *testing.T) {
	var sm = map[string]int{
		"a": 1,
	}
	var rm = map[string]int{
		//"b": 2,
	}
	buf := Marshal(&sm)
	Unmarshal(buf, &rm)
	Assert(t, buf, sm, rm)
}

func TestHelloWorld(t *testing.T) {
	hello, world := []byte("hello, "), "world"
	hello2, world2 := []byte("1"), ""

	Unmarshal(Marshal(&hello, &world), &hello2, &world2)
	if !bytes.Equal(hello2, hello) || world2 != world {
		t.Error(hello2, world2)
	}
}

func TestGetName(t *testing.T) {
	stdin := (interface {
		Read([]byte) (int, error)
		Write([]byte) (int, error)
	})(os.Stdin)
	nt := newType()
	items := []struct {
		ret string
		val any
	}{
		{"int", 1},
		{"github.com/niubaoshu/gotiny.Encoder", Encoder{}},
		{"*int", (*int)(nil)},
		{"**int", (**int)(nil)},
		{"[]int", []int{}},
		{"[]time.Time", []time.Time{}},
		{"[]github.com/niubaoshu/gotiny.Serializer", []Serializer{}},
		{"*interface {}", (*any)(nil)},
		{"map[int]string", map[int]string{}},
		{"struct { a struct { int; b int; dec []github.com/niubaoshu/gotiny.Decoder; abb interface {}; c io.ReadWriteCloser } }",
			struct {
				a struct {
					int
					b   int
					dec []Decoder
					abb any
					c   io.ReadWriteCloser
				}
			}{}},
		{"struct {}", struct{}{}},
		{"*interface { Read([]uint8) (int, error); Write([]uint8) (int, error) }", &stdin},
		{"func(int) (int, error)", func(i int) (int, error) { return 0, nil }},
		{"func(int)", func(i int) {}},
		{"func(int) error", func(i int) error { return nil }},
		{"struct { a int }", nt},
		{"<nil>", nil},
	}
	for _, item := range items {
		r := reflect.TypeOf(item.val)
		if GetName(item.val) != item.ret {
			t.Logf("real: %s , exp: %s", GetName(item.val), item.ret)
			t.Fatalf("string:%s,name:%s,pkgpath:%s,fmt %T", r.String(), r.Name(), r.PkgPath(), item.val)
		}
	}
}

func newType() struct {
	a int
} {
	return struct{ a int }{a: 1}
}

func TestUint64(t *testing.T) {
	v := make([]uint64, testNum)
	for i := 0; i < testNum; i++ {
		v[i] = rand.Uint64()
	}
	buf := Marshal(&v)
	var vd []uint64
	Unmarshal(buf, &vd)

	for i := 0; i < testNum; i++ {
		Assert(t, buf, v[i], vd[i])
	}
}

func TestInt64(t *testing.T) {
	v := make([]int64, testNum)
	for i := 0; i < testNum; i++ {
		v[i] = rand.Int63()
		if rand.Intn(2)%2 == 0 {
			v[i] = -v[i]
		}
	}
	buf := Marshal(&v)
	var vd []int64
	Unmarshal(buf, &vd)

	for i := 0; i < testNum; i++ {
		Assert(t, buf, v[i], vd[i])
	}
}

func TestUint16(t *testing.T) {
	v := make([]uint16, testNum)
	for i := 0; i < testNum; i++ {
		v[i] = uint16(rand.Uint32())
	}
	buf := Marshal(&v)
	var vd []uint16
	Unmarshal(buf, &vd)

	for i := 0; i < testNum; i++ {
		Assert(t, buf, v[i], vd[i])
	}
}

func TestInt16(t *testing.T) {
	v := make([]int16, testNum)
	for i := 0; i < testNum; i++ {
		v[i] = int16(rand.Int31())
		if rand.Intn(2)%2 == 0 {
			v[i] = -v[i]
		}
	}
	buf := Marshal(&v)
	var vd []int16
	Unmarshal(buf, &vd)

	for i := 0; i < testNum; i++ {
		Assert(t, buf, v[i], vd[i])
	}
}

func TestUint32(t *testing.T) {
	v := make([]uint32, testNum)
	for i := 0; i < testNum; i++ {
		v[i] = rand.Uint32()
	}
	buf := Marshal(&v)
	var vd []uint32
	Unmarshal(buf, &vd)

	for i := 0; i < testNum; i++ {
		Assert(t, buf, v[i], vd[i])
	}
}

func TestInt32(t *testing.T) {
	v := make([]int32, testNum)
	for i := 0; i < testNum; i++ {
		v[i] = rand.Int31()
		if rand.Intn(2)%2 == 0 {
			v[i] = -v[i]
		}
	}
	buf := Marshal(&v)
	var vd []int32
	Unmarshal(buf, &vd)

	for i := 0; i < testNum; i++ {
		Assert(t, buf, v[i], vd[i])
	}
}

func TestBool(t *testing.T) {
	v := make([]bool, testNum)
	for i := 0; i < testNum; i++ {
		v[i] = rand.Intn(2)%2 == 0
	}
	buf := Marshal(&v)
	var vd []bool
	Unmarshal(buf, &vd)

	for i := 0; i < testNum; i++ {
		Assert(t, buf, v[i], vd[i])
	}
}

func TestTime(t *testing.T) {
	v := make([]time.Time, testNum)
	for i := 0; i < testNum; i++ {
		v[i] = time.Now()
	}
	buf := Marshal(&v)
	var vd []time.Time
	Unmarshal(buf, &vd)

	for i := 0; i < testNum; i++ {
		v[i].Equal(vd[i])
	}
}

func TestPointerMap(t *testing.T) {
	n := testNum / 100
	for i := 0; i < n; i++ {
		vMap := map[*int][]bTyp{&vint1: gentBase(), &vint2: gentBase()}
		vdMap := map[*int][]bTyp{}
		buf := Marshal(&vMap)
		Unmarshal(buf, &vdMap)
		vd := map[int][]bTyp{}
		for k, v := range vMap {
			vd[*k] = v
		}
		for k, v := range vdMap {
			Assert(t, buf, vd[*k], v)
		}
	}
}

func TestInterface(t *testing.T) {
	n := testNum / 100
	for i := 0; i < n; i++ {
		var v1 io.ReadWriter = &tint{a: int16(rand.Intn(1<<16) - 1<<15)}
		var v2 io.ReadWriter = bytes.NewBufferString(randString(10))
		var v3 any = &v1
		var d1 io.ReadWriter
		var d2 io.ReadWriter
		var d3 any
		buf := Marshal(&v1, &v2, &v3)
		Unmarshal(buf, &d1, &d2, &d3)
		Assert(t, buf, v1, d1)
		Assert(t, buf, v2, d2)
		Assert(t, buf, v3, d3)
	}
}

func indirect(i any) any {
	v := reflect.ValueOf(i)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v.Interface()
}

func randString(l int) string {
	bss := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result[i] = bss[r.Intn(62)]
	}
	return string(result)
}

func Assert(t *testing.T, buf []byte, x, y any) {
	if !reflect.DeepEqual(x, y) {
		e, g := indirect(x), indirect(y)
		t.Errorf("\nlength:%d \nexp type = %T; value = %+v;\ngot type = %T; value = %+v; \n", len(buf), e, e, g, g)
	}
}
