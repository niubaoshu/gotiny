package gotiny

import (
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type decEng func(*Decoder, unsafe.Pointer) // 解码器

var (
	rt2decEng = map[reflect.Type]decEng{
		reflect.TypeFor[bool]():       decBool,
		reflect.TypeFor[int]():        decInt,
		reflect.TypeFor[int8]():       decInt8,
		reflect.TypeFor[int16]():      decInt16,
		reflect.TypeFor[int32]():      decInt32,
		reflect.TypeFor[int64]():      decInt64,
		reflect.TypeFor[uint]():       decUint,
		reflect.TypeFor[uint8]():      decUint8,
		reflect.TypeFor[uint16]():     decUint16,
		reflect.TypeFor[uint32]():     decUint32,
		reflect.TypeFor[uint64]():     decUint64,
		reflect.TypeFor[uintptr]():    decUintptr,
		reflect.TypeFor[float32]():    decFloat32,
		reflect.TypeFor[float64]():    decFloat64,
		reflect.TypeFor[complex64]():  decComplex64,
		reflect.TypeFor[complex128](): decComplex128,
		reflect.TypeFor[[]byte]():     decBytes,
		reflect.TypeFor[string]():     decString,
		reflect.TypeFor[time.Time]():  decTime,
		reflect.TypeFor[struct{}]():   decIgnore,
		reflect.TypeOf(nil):           decIgnore,
	}

	decEngines = []decEng{
		reflect.Bool:       decBool,
		reflect.Int:        decInt,
		reflect.Int8:       decInt8,
		reflect.Int16:      decInt16,
		reflect.Int32:      decInt32,
		reflect.Int64:      decInt64,
		reflect.Uint:       decUint,
		reflect.Uint8:      decUint8,
		reflect.Uint16:     decUint16,
		reflect.Uint32:     decUint32,
		reflect.Uint64:     decUint64,
		reflect.Uintptr:    decUintptr,
		reflect.Float32:    decFloat32,
		reflect.Float64:    decFloat64,
		reflect.Complex64:  decComplex64,
		reflect.Complex128: decComplex128,
		reflect.String:     decString,
	}
	decLock sync.RWMutex
)

func getDecEngine(rt reflect.Type) decEng {
	decLock.RLock()
	engine := rt2decEng[rt]
	decLock.RUnlock()
	if engine != nil {
		return engine
	}
	decLock.Lock()
	buildDecEngine(rt, &engine)
	decLock.Unlock()
	return engine
}

func buildDecEngine(rt reflect.Type, engPtr *decEng) {
	engine, has := rt2decEng[rt]
	if has {
		*engPtr = engine
		return
	}

	if _, engine = implementOtherSerializer(rt); engine != nil {
		rt2decEng[rt] = engine
		*engPtr = engine
		return
	}

	kind := rt.Kind()
	var eEng decEng
	switch kind {
	case reflect.Ptr:
		et := rt.Elem()
		defer buildDecEngine(et, &eEng)
		engine = func(d *Decoder, p unsafe.Pointer) {
			if d.decIsNotNil() {
				if isNil(p) {
					//*(*unsafe.Pointer)(p) = unsafe.Pointer(reflect.New(et).Elem().UnsafeAddr())
					*(*unsafe.Pointer)(p) = reflect.New(et).UnsafePointer()
				}
				eEng(d, *(*unsafe.Pointer)(p))
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
			}
		}
	case reflect.Array:
		l, et := rt.Len(), rt.Elem()
		size := et.Size()
		defer buildDecEngine(et, &eEng)
		engine = func(d *Decoder, p unsafe.Pointer) {
			for i := 0; i < l; i++ {
				eEng(d, unsafe.Add(p, i*int(size)))
			}
		}
	case reflect.Slice:
		et := rt.Elem()
		size := et.Size()
		defer buildDecEngine(et, &eEng)
		engine = func(d *Decoder, p unsafe.Pointer) {
			header := (*sliceHeader)(p)
			if d.decIsNotNil() {
				l := d.decLength()
				if isNil(p) || header.cap < l {
					*header = sliceHeader{data: reflect.MakeSlice(rt, l, l).UnsafePointer(), len: l, cap: l}
				} else {
					header.len = l
				}
				for i := 0; i < l; i++ {
					eEng(d, unsafe.Add(header.data, uintptr(i)*size))
				}
			} else if !isNil(p) {
				*header = sliceHeader{data: nil, len: 0, cap: 0}
			}
		}
	case reflect.Map:
		kt, vt := rt.Key(), rt.Elem()
		var kEng, vEng decEng
		defer buildDecEngine(kt, &kEng)
		defer buildDecEngine(vt, &vEng)
		engine = func(d *Decoder, p unsafe.Pointer) {
			if d.decIsNotNil() {
				l := d.decLength()
				v := reflect.NewAt(rt, p).Elem()
				if isNil(p) {
					v = reflect.MakeMapWithSize(rt, l)
					*(*unsafe.Pointer)(p) = v.UnsafePointer()
				}
				key, val := reflect.New(kt).Elem(), reflect.New(vt).Elem()
				for i := 0; i < l; i++ {
					kEng(d, unsafe.Pointer(key.UnsafeAddr()))
					vEng(d, unsafe.Pointer(val.UnsafeAddr()))
					v.SetMapIndex(key, val)
					key.SetZero()
					val.SetZero()
				}
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
			}
		}
	case reflect.Struct:
		fields, offs := getFieldType(rt, 0)
		nf := len(fields)
		fEngines := make([]decEng, nf)
		defer func() {
			for i := 0; i < nf; i++ {
				buildDecEngine(fields[i], &fEngines[i])
			}
		}()
		engine = func(d *Decoder, p unsafe.Pointer) {
			for i := 0; i < nf && i < len(offs); i++ {
				fEngines[i](d, unsafe.Add(p, offs[i]))
			}
		}
	case reflect.Interface:
		engine = func(d *Decoder, p unsafe.Pointer) {
			if d.decIsNotNil() {
				var name string
				decString(d, unsafe.Pointer(&name))
				et, has := name2type[name]
				if !has {
					panic("unknown typ:" + name)
				}
				v := reflect.NewAt(rt, p).Elem()
				if v.IsNil() || v.Elem().Type() != et {
					ev := reflect.New(et).Elem()
					getDecEngine(et)(d, getUnsafePointer(ev))
					v.Set(ev)
				} else {
					getDecEngine(et)(d, getUnsafePointer(v.Elem()))
				}
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
			}
		}
	case reflect.Chan, reflect.Func, reflect.Invalid, reflect.UnsafePointer:
		panic("not support " + rt.String() + " type")
	default:
		engine = decEngines[kind]
	}
	rt2decEng[rt] = engine
	*engPtr = engine
}
