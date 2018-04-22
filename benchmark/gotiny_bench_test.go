package benchmark

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"unsafe"

	"github.com/niubaoshu/gotiny"
)

var buf []byte

func init() {
	buf = gotiny.Encodes(genA())
	fmt.Println("buf length:", len(buf))
}

func BenchmarkEncode(b *testing.B) {
	a := genA()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gotiny.Encodes(a)
	}
}

func BenchmarkDecode(b *testing.B) {
	var a = genA()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gotiny.Decodes(buf, a)
	}
}

type (
	baseTyp struct {
		fBool         bool
		fInt8         int8
		fInt16        int16
		fInt32        int32
		fInt64        int64
		fInt          int
		fUint8        uint8
		fUint16       uint16
		fUint32       uint32
		fUint64       uint64
		fUint         uint
		fUintptr      uintptr
		fFloat32      float32
		fFloat64      float64
		fComplex64    complex64
		fComplex128   complex128
		fString       string
		unsafePointer unsafe.Pointer
	}

	A struct {
		Array    [10]baseTyp
		Slice    []baseTyp
		BirthDay time.Time
		inter    interface{}
		m        map[string]*baseTyp
	}
)

func genBase() baseTyp {
	return baseTyp{
		fBool:         rand.Int()%2 == 0,
		fInt8:         int8(rand.Int()),
		fInt16:        int16(rand.Int()),
		fInt32:        int32(rand.Int()),
		fInt64:        int64(rand.Int()),
		fInt:          int(rand.Int()),
		fUint8:        uint8(rand.Int()),
		fUint16:       uint16(rand.Int()),
		fUint64:       uint64(rand.Int()),
		fUintptr:      uintptr(rand.Int()),
		fFloat32:      rand.Float32(),
		fFloat64:      rand.Float64(),
		fComplex64:    complex(rand.Float32(), rand.Float32()),
		fComplex128:   complex(rand.Float64(), rand.Float64()),
		fString:       GetRandomString(20 + rand.Intn(256)),
		unsafePointer: unsafe.Pointer(nil),
	}
}

func genA() *A {
	a := &A{
		BirthDay: time.Now(),
		inter:    genBase(),
		m:        make(map[string]*baseTyp),
	}
	a.Slice = make([]baseTyp, len(a.Array))
	for i := 0; i < len(a.Array); i++ {
		a.Array[i] = genBase()
		a.Slice[i] = genBase()
		b := genBase()
		a.m[GetRandomString(len(a.Array))] = &b
	}
	return a
}

func GetRandomString(l int) string {
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result[i] = bytes[r.Intn(62)]
	}
	return string(result)
}
