package gotiny

import (
	"io"
	"math/rand"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

var (
	buf   []byte
	value = genA()
	e     *Encoder
	d     *Decoder
)

func init() {
	t := reflect.TypeOf(value).Elem()
	e = NewEncoderWithType(t)
	d = NewDecoderWithType(t)
	buf = e.Encode(value)
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Marshal(value)
	}
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Unmarshal(buf, value)
	}
}

func BenchmarkEncode2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.Encode(value)
	}
}

func BenchmarkDecode2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d.Decode(buf, value)
	}
}

type (
	baseTyp struct {
		FBool          bool
		FInt8          int8
		FInt16         int16
		FInt32         int32
		FInt64         int64
		FInt           int
		FUint8         uint8
		FUint16        uint16
		FUint32        uint32
		FUint64        uint64
		FUint          uint
		FUintptr       uintptr
		FFloat32       float32
		FFloat64       float64
		FComplex64     complex64
		FComplex128    complex128
		FString        string
		FUnsafePointer unsafe.Pointer
	}

	A struct {
		Array    [10]baseTyp
		Slice    []baseTyp
		BirthDay time.Time
		Inter    interface{}
		M        map[string]*baseTyp
	}
)

func genBase() baseTyp {
	return baseTyp{
		FBool:          rand.Int()%2 == 0,
		FInt8:          int8(rand.Int()),
		FInt16:         int16(rand.Int()),
		FInt32:         int32(rand.Int()),
		FInt64:         int64(rand.Int()),
		FInt:           int(rand.Int()),
		FUint8:         uint8(rand.Int()),
		FUint16:        uint16(rand.Int()),
		FUint64:        uint64(rand.Int()),
		FUintptr:       uintptr(rand.Int()),
		FFloat32:       rand.Float32(),
		FFloat64:       rand.Float64(),
		FComplex64:     complex(rand.Float32(), rand.Float32()),
		FComplex128:    complex(rand.Float64(), rand.Float64()),
		FString:        GetRandomString(20 + rand.Intn(256)),
		FUnsafePointer: unsafe.Pointer(nil),
	}
}

func genA() *A {
	ml := 10
	a := &A{
		BirthDay: time.Now(),
		Inter:    genBase(),
		M:        make(map[string]*baseTyp),
	}
	a.Slice = make([]baseTyp, len(a.Array))
	for i := 0; i < len(a.Array); i++ {
		a.Array[i] = genBase()
		a.Slice[i] = genBase()
	}

	for i := 0; i < ml; i++ {
		b := genBase()
		a.M[GetRandomString(10)] = &b
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

func BenchmarkDecodeUint64(b *testing.B) {
	var ints = make([][]byte, 10000)
	for i := 0; i < len(ints); i++ {
		a := rand.Uint64()
		ints[i] = Marshal(&a)
	}
	d := Decoder{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.buf = ints[rand.Intn(10000)]
		d.index = 0
		d.decUint64()
	}
}

func BenchmarkEncodeUint64(b *testing.B) {
	e := Encoder{buf: make([]byte, 0, 600000000)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.encUint64(rand.Uint64())
	}
}

func BenchmarkEncodeBool(b *testing.B) {
	l := 2000
	e := Encoder{buf: make([]byte, 0, 600000000)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < l*8; j++ {
			e.encBool(i%2 == 0)
		}
	}
}

func BenchmarkDecodeBool(b *testing.B) {
	l := 2000
	var ints = make([][]byte, 10000)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(ints); i++ {
		s := make([]byte, l)
		io.ReadFull(r, s)
		ints[i] = s
	}
	d := Decoder{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.buf = ints[rand.Intn(10000)]
		d.boolBit = 0
		d.boolPos = 0
		d.index = 0
		for j := 0; j < l*8; j++ {
			d.decBool()
		}
	}
}
