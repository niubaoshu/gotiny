package gotiny

import (
	cr "crypto/rand"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

const maxLength = 30

type kind uint

const (
	kBool kind = iota
	kInt
	kInt8
	kInt16
	kInt32
	kInt64
	kUint
	kUint8
	kUint16
	kUint32
	kUint64
	kUintptr
	kFloat32
	kFloat64
	kComplex64
	kComplex128
	kString
	kTime
	kBytes
	kStruckEmpty

	kArray
	kInterface
	kMap
	kPointer
	kSlice
	kStruct
	kNumber
)

var (
	rts = [...]reflect.Type{
		kBool:        reflect.TypeFor[bool](),
		kInt:         reflect.TypeFor[int](),
		kInt8:        reflect.TypeFor[int8](),
		kInt16:       reflect.TypeFor[int16](),
		kInt32:       reflect.TypeFor[int32](),
		kInt64:       reflect.TypeFor[int64](),
		kUint:        reflect.TypeFor[uint](),
		kUint8:       reflect.TypeFor[uint8](),
		kUint16:      reflect.TypeFor[uint16](),
		kUint32:      reflect.TypeFor[uint32](),
		kUint64:      reflect.TypeFor[uint64](),
		kUintptr:     reflect.TypeFor[uintptr](),
		kFloat32:     reflect.TypeFor[float32](),
		kFloat64:     reflect.TypeFor[float64](),
		kComplex64:   reflect.TypeFor[complex64](),
		kComplex128:  reflect.TypeFor[complex128](),
		kString:      reflect.TypeFor[string](),
		kTime:        reflect.TypeFor[time.Time](),
		kBytes:       reflect.TypeFor[[]byte](),
		kStruckEmpty: reflect.TypeFor[struct{}](),
	}

	t = []kind{
		kBool, kInt, kInt8, kInt16, kInt32, kInt64,
		kUint, kUint8, kUint16, kUint32, kUint64, kUintptr,
		kFloat32, kFloat64, kComplex64, kComplex128,
		kString, kTime, kBytes, kStruckEmpty,

		kArray, kArray, kArray, kArray, kArray,
		kInterface, kInterface, kInterface, kInterface, kInterface,
		kMap, kMap, kMap, kMap, kMap,
		kPointer, kPointer, kPointer, kPointer, kPointer,
		kSlice, kSlice, kSlice, kSlice, kSlice,
		kStruct, kStruct, kStruct, kStruct, kStruct,
	}
	it = []kind{
		kBool, kInt, kInt8, kInt16, kInt32, kInt64,
		kUint, kUint8, kUint16, kUint32, kUint64, kUintptr,
		kFloat32, kFloat64, kComplex64, kComplex128,
		kString, kTime, kBytes, kStruckEmpty,

		kArray, kArray, kArray, kArray, kArray,
		kMap, kMap, kMap, kMap, kMap,
		// kInterface, kInterface, kInterface, kInterface, kInterface,
		kPointer, kPointer, kPointer, kPointer, kPointer,
		kSlice, kSlice, kSlice, kSlice, kSlice,
		kStruct, kStruct, kStruct, kStruct, kStruct,
	}
	ct = []kind{
		kBool, kInt, kInt8, kInt16, kInt32, kInt64,
		kUint, kUint8, kUint16, kUint32, kUint64, kUintptr,
		kFloat32, kFloat64, kComplex64, kComplex128,
		kString, kTime,
		// kBytes, kStruckEmpty,
		kArray, kArray, kArray, kArray, kArray,
		// kMap, kMap, kMap, kMap, kMap,
		kInterface, kInterface, kInterface, kInterface, kInterface,
		kPointer, kPointer, kPointer, kPointer, kPointer,
		// kSlice, kSlice, kSlice, kSlice, kSlice,
		kStruct, kStruct, kStruct, kStruct, kStruct,
	}
	cit = []kind{
		kBool, kInt, kInt8, kInt16, kInt32, kInt64,
		kUint, kUint8, kUint16, kUint32, kUint64, kUintptr,
		kFloat32, kFloat64, kComplex64, kComplex128,
		kString, kTime,
		// kBytes, kStruckEmpty,
		kArray, kArray, kArray, kArray, kArray,
		// kMap, kMap, kMap, kMap, kMap,
		// kInterface, kInterface, kInterface, kInterface, kInterface,
		kPointer, kPointer, kPointer, kPointer, kPointer,
		// kSlice, kSlice, kSlice, kSlice, kSlice,
		kStruct, kStruct, kStruct, kStruct, kStruct,
	}
	vEmptyStruct = struct{}{}
)

func randValue(rt reflect.Type) reflect.Value {
	var v interface{}
	if rt == rts[kTime] {
		v = time.Unix(0, time.Now().UnixNano()+rand.Int63n(1<<20)-1<<19)
	} else if rt == rts[kBytes] {
		bs := make([]byte, rand.Intn(maxLength))
		n, _ := cr.Read(bs)
		v = bs[:n]
	} else if rt == rts[kStruckEmpty] {
		v = vEmptyStruct
	}

	switch rt.Kind() {
	case reflect.Bool:
		v = rand.Intn(2) == 0
	case reflect.Int:
		if ptr1Size == 4 {
			v = int(rand.Int63n(1<<32) - 1<<31)
		} else {
			v = int(rand.Uint64() - 1<<63)
		}
	case reflect.Int8:
		v = int8(rand.Intn(1<<8) - 1<<7)
	case reflect.Int16:
		v = int16(rand.Intn(1<<16) - 1<<15)
	case reflect.Int32:
		v = int32(rand.Intn(1<<32) - 1<<31)
	case reflect.Int64:
		v = int64(rand.Uint64() - 1<<63)
	case reflect.Uint:
		if ptr1Size == 4 {
			v = uint(rand.Uint32())
		} else {
			v = uint(rand.Uint64())
		}
	case reflect.Uint8:
		v = uint8(rand.Intn(1 << 8))
	case reflect.Uint16:
		v = uint16(rand.Intn(1 << 16))
	case reflect.Uint32:
		v = rand.Uint32()
	case reflect.Uint64:
		v = rand.Uint64()
	case reflect.Uintptr:
		if ptr1Size == 4 {
			v = uintptr(rand.Uint32())
		} else {
			v = uintptr(rand.Uint64())
		}
	case reflect.Float32:
		v = rand.Float32()
	case reflect.Float64:
		v = rand.Float64()
	case reflect.Complex64:
		v = complex(rand.Float32(), rand.Float32())
	case reflect.Complex128:
		v = complex(rand.Float64(), rand.Float64())
	case reflect.String:
		v = randString(rand.Intn(maxLength))
	case reflect.Array:
		l := rt.Len()
		v := reflect.New(rt).Elem()
		et := rt.Elem()
		for i := 0; i < l; i++ {
			v.Index(i).Set(randValue(et))
		}
		return v
	case reflect.Interface:
		return reflect.ValueOf(randValue(rt.Elem()).Interface())

	case reflect.Map:
		l := rand.Intn(maxLength)
		m := reflect.MakeMapWithSize(rt, l)
		for i := 0; i < l; i++ {
			m.SetMapIndex(randValue(rt.Key()), randValue(rt.Elem()))
		}
		return m
	case reflect.Pointer:
		v := reflect.New(rt).Elem()
		v.Set(randValue(rt.Elem()).Addr())
		return v
	case reflect.Slice:
		l := rand.Intn(maxLength)
		s := reflect.MakeSlice(rt, l, l)
		for i := 0; i < l; i++ {
			s.Index(i).Set(randValue(rt.Elem()))
		}
		return s
	case reflect.Struct:
		n := rt.NumField()
		v := reflect.New(rt).Elem()
		for i := 0; i < n; i++ {
			v.Field(i).Set(randValue(rt.Field(i).Type))
		}
		return v
	default:
		panic("unknown type")
	}
	fmt.Println(reflect.TypeOf(v))
	return reflect.ValueOf(v)
}

type state struct {
	r              *rand.Rand
	lt             kind
	needComparable bool
	depth          int
}

func newState() *state {
	return &state{rand.New(rand.NewSource(time.Now().UnixNano())), 30, false, 0}
}

func (s *state) randType() (rt reflect.Type) {
	s.depth++

	t := s.randKind()
	if t < kArray {
		s.depth--
		return rts[t]
	}
	switch t {
	case kArray:
		l := rand.Intn(maxLength)
		rt = reflect.ArrayOf(l, s.randType())
	case kInterface:
		rt = reflect.TypeOf(reflect.New(s.randType()).Elem().Interface())
	case kMap:
		s.needComparable = true
		key := s.randType()
		s.needComparable = false
		val := s.randType()
		rt = reflect.MapOf(key, val)
	case kPointer:
		rt = reflect.PointerTo(s.randType())
	case kSlice:
		rt = reflect.SliceOf(s.randType())
	case kStruct:
		n := rand.Intn(maxLength) / 3
		fs := make([]reflect.StructField, n)
		for i := 0; i < n; i++ {
			fs[i] = reflect.StructField{
				Name:    randFieldName(),
				Type:    s.randType(),
				PkgPath: "github.com/niubaoshu/gotiny",
			}
		}
		rt = reflect.StructOf(fs)
	default:
		panic("invalid kind")
	}
	s.depth--
	return rt
}

func (s *state) randKind() kind {
	ts := t
	if s.lt == kInterface {
		ts = it
	}
	if s.needComparable {
		ts = ct
	}
	if s.lt == kInterface && s.needComparable {
		ts = cit
	}
	if s.depth > 10 {
		ts = t[:kTime]
	}
	t := ts[s.r.Intn(len(ts))]
	s.lt = t
	return t
}

func TestRandType(t *testing.T) {
	r := newState()
	for i := 0; i < 1000000; i++ {
		r.randType()
	}
}

func TestRandValue(t *testing.T) {
	//r := newState()
	for i := 0; i < 10; i++ {
		//rt := r.randType()
		//t.Log(rt)
		//t.Log(randValue(rt))
	}
}

func randFieldName() string {
	if rand.Intn(2) == 0 {
		return "f" + randString(10)
	} else {
		return "F" + randString(10)
	}
}
