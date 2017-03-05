package gotiny

import (
	"encoding"
	"encoding/gob"
	"reflect"
	"sync"
	"unsafe"
)

type (
	encEng    func(*Encoder, unsafe.Pointer)  //编码器
	encEngPtr *func(*Encoder, unsafe.Pointer) //编码器指针
)

var (
	rt2Eng = map[reflect.Type]encEngPtr{
		reflect.TypeOf((*string)(nil)).Elem():  &encString,
		reflect.TypeOf((*bool)(nil)).Elem():    &encBool,
		reflect.TypeOf((*uint8)(nil)).Elem():   &encUint8,
		reflect.TypeOf((*int8)(nil)).Elem():    &encUint8,
		reflect.TypeOf((*int)(nil)).Elem():     &encInt,
		reflect.TypeOf((*uint)(nil)).Elem():    &encUint,
		reflect.TypeOf((*int16)(nil)).Elem():   &encInt16,
		reflect.TypeOf((*int32)(nil)).Elem():   &encInt32,
		reflect.TypeOf((*int64)(nil)).Elem():   &encInt64,
		reflect.TypeOf((*uint16)(nil)).Elem():  &encUint16,
		reflect.TypeOf((*uint32)(nil)).Elem():  &encUint32,
		reflect.TypeOf((*uint64)(nil)).Elem():  &encUint64,
		reflect.TypeOf((*uintptr)(nil)).Elem(): &encUint,
		reflect.TypeOf((*float32)(nil)).Elem(): &encFloat32,
		reflect.TypeOf((*float64)(nil)).Elem(): &encFloat64,
		reflect.TypeOf((*[]byte)(nil)).Elem():  &encBytes,
		reflect.TypeOf(nil):                    &encignore,
	}

	baseEncEng = []encEng{
		reflect.Bool:          encBool,
		reflect.String:        encString,
		reflect.Uint8:         encUint8,
		reflect.Int8:          encUint8,
		reflect.Int:           encInt,
		reflect.Uint:          encUint,
		reflect.Int16:         encInt16,
		reflect.Int32:         encInt32,
		reflect.Int64:         encInt64,
		reflect.Uint16:        encUint16,
		reflect.Uint32:        encUint32,
		reflect.Uint64:        encUint64,
		reflect.Uintptr:       encUint,
		reflect.Float32:       encFloat32,
		reflect.Float64:       encFloat64,
		reflect.Complex64:     encComplex64,
		reflect.Complex128:    encComplex128,
		reflect.Invalid:       encignore,
		reflect.Chan:          encignore,
		reflect.Func:          encignore,
		reflect.Interface:     encignore,
		reflect.UnsafePointer: encignore,
	}

	encLock sync.RWMutex
)

func getEncEngine(rt reflect.Type) encEng {
	encLock.RLock()
	engptr := rt2Eng[rt]
	encLock.RUnlock()
	if engptr != nil && *engptr != nil {
		return *engptr
	}
	encLock.Lock()
	engptr = buildEncEngine(rt)
	encLock.Unlock()
	return *engptr
}

func buildEncEngine(rt reflect.Type) encEngPtr {
	//todo  接口类型处理
	engine, has := rt2Eng[rt]
	if has {
		return engine
	}
	engine = new(func(*Encoder, unsafe.Pointer))
	rt2Eng[rt] = engine

	if fn, _, yes := implementsGob(rt); yes {
		*engine = func(e *Encoder, p unsafe.Pointer) {
			buf, _ := fn(reflect.NewAt(rt, p).Elem().Interface().(gob.GobEncoder))
			e.encUint(uint64(len(buf)))
			e.buf = append(e.buf, buf...)
		}
		return engine
	}

	if fn, _, yes := implementsBin(rt); yes {
		*engine = func(e *Encoder, p unsafe.Pointer) {
			buf, _ := fn(reflect.NewAt(rt, p).Elem().Interface().(encoding.BinaryMarshaler))
			e.encUint(uint64(len(buf)))
			e.buf = append(e.buf, buf...)
		}
		return engine
	}

	if fn, _, yes := implementsGotiny(reflect.PtrTo(rt)); yes {
		*engine = func(e *Encoder, p unsafe.Pointer) {
			e.buf = fn(reflect.NewAt(rt, p).Interface().(GoTinySerializer), e.buf)
		}
		return engine
	}

	var eEng encEngPtr
	switch rt.Kind() {
	case reflect.Ptr:
		eEng = buildEncEngine(rt.Elem())
		*engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encBool(isNotNil)
			if isNotNil {
				(*eEng)(e, elem(p))
			}
		}
	case reflect.Array:
		et := rt.Elem()
		eEng = buildEncEngine(et)
		l := rt.Len()
		size := et.Size()
		*engine = func(e *Encoder, p unsafe.Pointer) {
			for i := 0; i < l; i++ {
				(*eEng)(e, unsafe.Pointer(uintptr(p)+uintptr(i)*size))
			}
		}
	case reflect.Slice:
		et := rt.Elem()
		eEng = buildEncEngine(et)
		size := et.Size()
		*engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encBool(isNotNil)
			if isNotNil {
				l := *(*int)(unsafe.Pointer(uintptr(p) + ptrSize))
				e.encLength(l)
				pp := *(*unsafe.Pointer)(p)
				for i := 0; i < l; i++ {
					(*eEng)(e, unsafe.Pointer(uintptr(pp)+uintptr(i)*size))
				}
			}
		}
	case reflect.Map:
		kEng := buildEncEngine(rt.Key())
		eEng = buildEncEngine(rt.Elem())
		//http://blog.csdn.net/hificamera/article/details/51701804
		var encKey, encVal func(e *Encoder, v reflect.Value)
		if rt.Elem().Kind() == reflect.Map {
			eEng = func(eng encEngPtr) encEngPtr {
				emEng := func(e *Encoder, p unsafe.Pointer) {
					(*eng)(e, unsafe.Pointer(&p))
				}
				return &emEng
			}(eEng)
		}
		if rt.Key().Kind() == reflect.Ptr {
			encKey = func(e *Encoder, v reflect.Value) {
				p := unsafe.Pointer(v.Pointer())
				(*kEng)(e, unsafe.Pointer(&p))
			}
		} else {
			encKey = func(e *Encoder, v reflect.Value) {
				(*kEng)(e, (*refVal)(unsafe.Pointer(&v)).ptr)
			}
		}
		if rt.Elem().Kind() == reflect.Ptr {
			encVal = func(e *Encoder, v reflect.Value) {
				p := unsafe.Pointer(v.Pointer())
				(*eEng)(e, unsafe.Pointer(&p))
			}
		} else {
			encVal = func(e *Encoder, v reflect.Value) {
				(*eEng)(e, (*refVal)(unsafe.Pointer(&v)).ptr)
			}
		}
		*engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encBool(isNotNil)
			if isNotNil {
				e.encLength(*(*int)(*(*unsafe.Pointer)(p)))
				v := reflect.NewAt(rt, p).Elem()
				for _, key := range v.MapKeys() {
					encKey(e, key)
					encVal(e, v.MapIndex(key))
				}
			}
		}
	case reflect.Struct:
		nf := rt.NumField()
		if nf > 0 {
			engs, offs := make([]encEngPtr, nf), make([]uintptr, nf)
			for i := 0; i < nf; i++ {
				ft := rt.Field(i)
				engs[i] = buildEncEngine(ft.Type)
				offs[i] = ft.Offset
			}
			*engine = func(e *Encoder, p unsafe.Pointer) {
				for i := 0; i < nf; i++ {
					(*engs[i])(e, unsafe.Pointer(uintptr(p)+offs[i]))
				}
			}
		} else {
			*engine = encignore
		}
	default:
		*engine = baseEncEng[rt.Kind()]
	}
	return engine
}
