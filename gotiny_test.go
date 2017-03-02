package gotiny

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"
	"unsafe"
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
	cirTyp *cirTyp

	cirStruct struct {
		a   int
		cir *cirStruct
	}
	cirmap map[int]cirmap
	A      struct {
		Name     string
		BirthDay time.Time
		Phone    string
		Siblings int
		Spouse   bool
		Money    float64
	}
)

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
	vptr        = &vint
	vsliceptr   = &vbytes
	vptrslice   = []*int{&vint, &vint, &vint}
	vnilptr     *int
	vnilptrptr  = &vnilptr
	vtime       = time.Now()
	vsliceStr   = []baseTyp{
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

	vcirStruct  = cirStruct{a: 1, cir: nil}
	v2cirStruct = cirStruct{a: 1, cir: &vcirStruct}

	vcirmap  = cirmap{1: nil}
	v2cirmap = cirmap{2: vcirmap}

	vAstruct = genA()

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
		vptr,
		vsliceptr,
		vptrslice,
		vnilptr,
		vnilptrptr,
		vtime,
		vsliceStr,
		vslicestring,
		varray,
		vcir,
		v2cir,
		vcirStruct,
		v2cirStruct,
		vcirmap,
		v2cirmap,
		vAstruct,
	}

	ptrs = []unsafe.Pointer{
		getPtr(&vbool),
		getPtr(&vfbool),
		getPtr(&vint8),
		getPtr(&vint16),
		getPtr(&vint32),
		getPtr(&vint64),
		getPtr(&v2int64),
		getPtr(&v3int64),
		getPtr(&vint),
		getPtr(&vuint),
		getPtr(&vuint8),
		getPtr(&vuint16),
		getPtr(&vuint32),
		getPtr(&vuint64),
		getPtr(&v2uint64),
		getPtr(&v3uint64),
		getPtr(&vuintptr),
		getPtr(&vfloat32),
		getPtr(&vfloat64),
		getPtr(&vcomp64),
		getPtr(&vcomp128),
		getPtr(&vstring),
		getPtr(&base),
		getPtr(&vbytes),
		getPtr(&vsliecbytes),
		getPtr(&vmap),
		getPtr(&v2map),
		getPtr(&v3map),
		getPtr(&v4map),
		getPtr(&v5map),
		getPtr(&vptr),
		getPtr(&vsliceptr),
		getPtr(&vptrslice),
		getPtr(&vnilptr),
		getPtr(&vnilptrptr),
		getPtr(&vtime),
		getPtr(&vsliceStr),
		getPtr(&vslicestring),
		getPtr(&varray),
		getPtr(&vcir),
		getPtr(&v2cir),
		getPtr(&vcirStruct),
		getPtr(&v2cirStruct),
		getPtr(&vcirmap),
		getPtr(&v2cirmap),
		getPtr(&vAstruct),
	}

	e = NewEncoder(vs...)
	d = NewDecoder(vs...)

	vals    = make([]reflect.Value, len(vs))
	types   = make([]reflect.Type, len(vs))
	retVals = make([]reflect.Value, len(vs))
	retPtrs = make([]unsafe.Pointer, len(vs))

	// buf     = make([]byte, 0, 1024)
	// network = bytes.NewBuffer(buf) // Stand-in for a network connection
	// //network bytes.Buffer
	// enc = gob.NewEncoder(network) // Will write to network.
	// dec = gob.NewDecoder(network) // Will read from network.
)

func init() {
	fmt.Fprintln(ioutil.Discard, time.Now())
	//v2cir = vcir
	//vcir = &v2cir
	//vcirStruct.cir = &vcirStruct
	fmt.Println("total", len(vs), len(ptrs), "value")

	if len(vs) != len(ptrs) {
		log.Fatal(" vs ptrs 不相等")
	}
	for i := 0; i < len(vs); i++ {
		types[i] = reflect.TypeOf(vs[i])
		vals[i] = reflect.NewAt(types[i], ptrs[i]).Elem()

		//var vp reflect.Value
		//if i == len(vs)-3 {
		//	a := 2
		//	vp = reflect.ValueOf(&a)
		//} else
		// if i == len(vs)-2 {
		// 	a := make([]byte, 15)
		// 	vp = reflect.ValueOf(&a)
		// } else if i == len(vs)-1 {
		// 	//a := map[int]int{111: 233, 6: 7}
		// 	a := map[int]int{}
		// 	vp = reflect.ValueOf(&a)
		// } else {
		//}
		retVals[i] = reflect.New(types[i]).Elem()
		retPtrs[i] = unsafe.Pointer(retVals[i].UnsafeAddr())
	}

	//ee := NewEncoder(0)
	//ret := ee.Encodes(vs...)
	//fmt.Println("gotiny length:", len(ret))

	// buf := make([]byte, 0, 1024)
	// network := bytes.NewBuffer(buf) // Stand-in for a network connection
	// enc := gob.NewEncoder(network)  // Will write to network.
	// for i := 0; i < len(vs); i++ {
	// 	enc.Encode(vs[i])
	// }
	// fmt.Println("stdgob length:", len(network.Bytes()))

	e.SetBuf(make([]byte, 0, 2048))
}

// Test basic operations in a safe manner.
func TestBasicEncoderDecoder(t *testing.T) {
	//fmt.Println(vs...)
	e.Reset()
	b := e.EncodeByUPtr(ptrs...)
	//t.Logf("%v\n", b)
	fmt.Printf("length: %d \n", len(b))
	d.ResetWith(b)
	d.DecodeByUPtr(retPtrs...)
	for i, result := range retVals {
		r := result.Interface()
		//fmt.Printf("%T: expected %v got %v ,%T\n", vs[i], vs[i], r, r)
		if !reflect.DeepEqual(vs[i], r) {
			t.Log(i)
			t.Fatalf("%T: expected %#v got %#v ,%T\n", vs[i], vs[i], r, r)
		}
	}

	d.ResetWith(e.EncodeValues(vals...))
	d.DecodeValues(retVals...)
	for i, result := range retVals {
		r := result.Interface()
		//fmt.Printf("%T: expected %v got %v ,%T\n", vs[i], vs[i], r, r)
		if !reflect.DeepEqual(vs[i], r) {
			t.Fatalf("%T: expected %#v got %#v ,%T\n", vs[i], vs[i], r, r)
		}
	}
}

// func BenchmarkStdEncode(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < 1000; j++ {
// 			for i := 0; i < len(vs); i++ {
// 				enc.Encode(vs[i])
// 			}
// 		}
// 	}
// }

// func BenchmarkStdDecode(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < 1000; j++ {
// 			for i := 0; i < len(presults); i++ {
// 				dec.Decode(presults[i])
// 				//err := dec.Decode(presults[i])
// 				//if err != nil {
// 				//	b.Fatal(j, err.Error())
// 				//}
// 			}
// 		}
// 	}
//}

func BenchmarkEncodes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < 1000; i++ {
			e.EncodeByUPtr(ptrs...)
		}
	}
}

func BenchmarkDecodes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < 1000; i++ {
			d.Reset()
			d.DecodeByUPtr(retPtrs...)
		}
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
