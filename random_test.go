package gotiny

import (
	cr "crypto/rand"
	"math/rand"
	"reflect"
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

	t = [...]kind{
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
	vBool = [...]reflect.Value{
		reflect.ValueOf(true),
		reflect.ValueOf(false),
	}
)

func randValue(rt reflect.Type) reflect.Value {

	if rt == rts[kTime] {
		return reflect.ValueOf(time.Unix(0, time.Now().UnixNano()+rand.Int63n(1<<20)-1<<19))
	} else if rt == rts[kBytes] {
		bs := make([]byte, rand.Intn(maxLength))
		n, _ := cr.Read(bs)
		return reflect.ValueOf(bs[:n])
	} else if rt == rts[kStruckEmpty] {
		return reflect.ValueOf(struct{}{})
	}

	switch rt.Kind() {
	case reflect.Bool:
		return vBool[(rand.Intn(2))]
	case reflect.Int:
		if ptr1Size == 4 {
			return reflect.ValueOf(int(rand.Int63n(1<<32) - 1<<31))
		} else {
			return reflect.ValueOf(int(rand.Uint64() - 1<<63))
		}
	case reflect.Int8:
		return reflect.ValueOf(int8(rand.Intn(1<<8) - 1<<7))
	case reflect.Int16:
		return reflect.ValueOf(int16(rand.Intn(1<<16) - 1<<15))
	case reflect.Int32:
		return reflect.ValueOf(int32(rand.Int63n(1<<32) - 1<<31))
	case reflect.Int64:
		return reflect.ValueOf(int64(rand.Uint64() - 1<<63))
	case reflect.Uint:
		if ptr1Size == 4 {
			return reflect.ValueOf(uint(rand.Uint32()))
		} else {
			return reflect.ValueOf(uint(rand.Uint64()))
		}

	case reflect.Uint8:
		return reflect.ValueOf(uint8(rand.Intn(1 << 8)))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(rand.Intn(1 << 16)))
	case reflect.Uint32:
		return reflect.ValueOf(rand.Uint32())
	case reflect.Uint64:
		return reflect.ValueOf(rand.Uint64())
	case reflect.Uintptr:
		if ptr1Size == 4 {
			return reflect.ValueOf(uintptr(rand.Uint32()))
		} else {
			return reflect.ValueOf(uintptr(rand.Uint64()))
		}
	case reflect.Float32:
		return reflect.ValueOf(rand.Float32())
	case reflect.Float64:
		return reflect.ValueOf(rand.Float64())
	case reflect.Complex64:
		return reflect.ValueOf(complex(rand.Float32(), rand.Float32()))
	case reflect.Complex128:
		return reflect.ValueOf(complex(rand.Float64(), rand.Float64()))
	case reflect.String:
		return reflect.ValueOf(randString(rand.Intn(maxLength)))
	case reflect.Array:
		l := rt.Len()
		v := reflect.New(rt).Elem()
		for i := 0; i < l; i++ {
			v.Index(i).Set(randValue(rt.Elem()))
		}
	case reflect.Interface:
		return reflect.ValueOf(randValue(rt.Elem()).Interface())

	case reflect.Map:
		l := rand.Intn(maxLength)
		m := reflect.MakeMapWithSize(rt, l)
		for i := 0; i < l; i++ {
			m.SetMapIndex(randValue(rt.Key()), randValue(rt.Elem()))
		}
	case reflect.Pointer:
		v := reflect.New(rt)
		v.Elem().Set(randValue(rt.Elem()))
		return v
	case reflect.Slice:
		l := rand.Intn(maxLength)
		s := reflect.MakeSlice(rt, l, l)
		for i := 0; i < l; i++ {
			s.Index(i).Set(randValue(rt.Elem()))
		}
	case reflect.Struct:
		n := rt.NumField()
		v := reflect.New(rt).Elem()
		for i := 0; i < n; i++ {
			v.Field(i).Set(randValue(rt.Field(i).Type))
		}
	default:
		return reflect.New(rt).Elem()
	}
	return reflect.New(rt).Elem()
}

func randType() reflect.Type {
	t := randKind()
	if t < kArray {
		return rts[t]
	}
	switch t {
	case kArray:
		l := rand.Intn(maxLength)
		return reflect.ArrayOf(l, randType())
	case kInterface:
		return reflect.TypeOf(reflect.New(randType()).Elem().Interface())
	case kMap:
		return reflect.MapOf(randType(), randType())
	case kPointer:
		return reflect.PointerTo(randType())
	case kSlice:
		return reflect.SliceOf(randType())
	case kStruct:
		n := rand.Intn(maxLength) / 3
		fs := make([]reflect.StructField, n)
		for i := 0; i < n; i++ {
			fs[i] = reflect.StructField{
				Name: randString(10),
				Type: randType(),
			}
		}
		return reflect.StructOf(fs)
	default:
		panic("invalid kind")
	}
}

func randKind() kind {
	return kind(rand.Intn(int(kNumber)))
}
