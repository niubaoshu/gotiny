package gotiny

import (
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type encEng func(*Encoder, unsafe.Pointer) //编码器

var (
	rt2encEng = map[reflect.Type]encEng{
		reflect.TypeFor[bool]():       encBool,
		reflect.TypeFor[int]():        encInt,
		reflect.TypeFor[int8]():       encInt8,
		reflect.TypeFor[int16]():      encInt16,
		reflect.TypeFor[int32]():      encInt32,
		reflect.TypeFor[int64]():      encInt64,
		reflect.TypeFor[uint]():       encUint,
		reflect.TypeFor[uint8]():      encUint8,
		reflect.TypeFor[uint16]():     encUint16,
		reflect.TypeFor[uint32]():     encUint32,
		reflect.TypeFor[uint64]():     encUint64,
		reflect.TypeFor[uintptr]():    encUintptr,
		reflect.TypeFor[float32]():    encFloat32,
		reflect.TypeFor[float64]():    encFloat64,
		reflect.TypeFor[complex64]():  encComplex64,
		reflect.TypeFor[complex128](): encComplex128,
		reflect.TypeFor[[]byte]():     encBytes,
		reflect.TypeFor[string]():     encString,
		reflect.TypeFor[time.Time]():  encTime,
		reflect.TypeFor[struct{}]():   encIgnore,
		reflect.TypeOf(nil):           encIgnore,
	}

	encEngines = [...]encEng{
		reflect.Bool:       encBool,
		reflect.Int:        encInt,
		reflect.Int8:       encInt8,
		reflect.Int16:      encInt16,
		reflect.Int32:      encInt32,
		reflect.Int64:      encInt64,
		reflect.Uint:       encUint,
		reflect.Uint8:      encUint8,
		reflect.Uint16:     encUint16,
		reflect.Uint32:     encUint32,
		reflect.Uint64:     encUint64,
		reflect.Uintptr:    encUintptr,
		reflect.Float32:    encFloat32,
		reflect.Float64:    encFloat64,
		reflect.Complex64:  encComplex64,
		reflect.Complex128: encComplex128,
		reflect.String:     encString,
	}

	encLock sync.RWMutex
)

func UnusedUnixNanoEncodeTimeType() {
	delete(rt2encEng, reflect.TypeOf((*time.Time)(nil)).Elem())
	delete(rt2decEng, reflect.TypeOf((*time.Time)(nil)).Elem())
}

func getEncEngine(rt reflect.Type) encEng {
	encLock.RLock()
	engine := rt2encEng[rt]
	encLock.RUnlock()
	if engine != nil {
		return engine
	}
	encLock.Lock()
	buildEncEngine(rt, &engine)
	encLock.Unlock()
	return engine
}

func buildEncEngine(rt reflect.Type, engPtr *encEng) {
	engine := rt2encEng[rt]
	if engine != nil {
		*engPtr = engine
		return
	}

	if engine, _ = implementOtherSerializer(rt); engine != nil {
		rt2encEng[rt] = engine
		*engPtr = engine
		return
	}

	kind := rt.Kind()
	var eEng encEng
	switch kind {
	case reflect.Ptr:
		defer buildEncEngine(rt.Elem(), &eEng)
		engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encIsNotNil(isNotNil)
			if isNotNil {
				eEng(e, *(*unsafe.Pointer)(p))
			}
		}
	case reflect.Array:
		et, l := rt.Elem(), rt.Len()
		size := et.Size()
		defer buildEncEngine(et, &eEng)
		engine = func(e *Encoder, p unsafe.Pointer) {
			for i := 0; i < l; i++ {
				eEng(e, unsafe.Add(p, i*int(size)))
			}
		}
	case reflect.Slice:
		et := rt.Elem()
		size := et.Size()
		defer buildEncEngine(et, &eEng)
		engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encIsNotNil(isNotNil)
			if isNotNil {
				header := (*sliceHeader)(p)
				l := header.len
				e.encLength(l)
				for i := 0; i < l; i++ {
					eEng(e, unsafe.Add(header.data, i*int(size)))
				}
			}
		}
	case reflect.Map:
		var kEng encEng
		defer buildEncEngine(rt.Key(), &kEng)
		defer buildEncEngine(rt.Elem(), &eEng)
		engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encIsNotNil(isNotNil)
			if isNotNil {
				v := reflect.NewAt(rt, p).Elem()
				e.encLength(v.Len())
				iter := v.MapRange()
				for iter.Next() {
					kEng(e, getUnsafePointer(iter.Key()))
					eEng(e, getUnsafePointer(iter.Value()))
				}
			}
		}
	case reflect.Struct:
		fields, offs := getFieldType(rt, 0)
		nf := len(fields)
		fEngines := make([]encEng, nf)
		defer func() {
			for i := 0; i < nf; i++ {
				buildEncEngine(fields[i], &fEngines[i])
			}
		}()
		engine = func(e *Encoder, p unsafe.Pointer) {
			for i := 0; i < len(fEngines) && i < len(offs); i++ {
				fEngines[i](e, unsafe.Add(p, offs[i]))
			}
		}
	case reflect.Interface:
		if rt.NumMethod() > 0 {
			engine = func(e *Encoder, p unsafe.Pointer) {
				isNotNil := !isNil(p)
				e.encIsNotNil(isNotNil)
				if isNotNil {
					v := reflect.ValueOf(*(*interface{ M() })(p))
					et := v.Type()
					e.encString(getNameOfType(et))
					getEncEngine(et)(e, getUnsafePointer(v))
				}
			}
		} else {
			engine = func(e *Encoder, p unsafe.Pointer) {
				isNotNil := !isNil(p)
				e.encIsNotNil(isNotNil)
				if isNotNil {
					v := reflect.ValueOf(*(*any)(p))
					et := v.Type()
					e.encString(getNameOfType(et))
					getEncEngine(et)(e, getUnsafePointer(v))
				}
			}
		}
	case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Invalid:
		panic("not support " + rt.String() + " type")
	default:
		engine = encEngines[kind]
	}
	rt2encEng[rt] = engine
	*engPtr = engine
}
