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
	rt2encEng = map[reflect.Type]encEngPtr{
		reflect.TypeOf((*bool)(nil)).Elem():           &encBool,
		reflect.TypeOf((*int)(nil)).Elem():            &encInt,
		reflect.TypeOf((*int8)(nil)).Elem():           &encInt8,
		reflect.TypeOf((*int16)(nil)).Elem():          &encInt16,
		reflect.TypeOf((*int32)(nil)).Elem():          &encInt32,
		reflect.TypeOf((*int64)(nil)).Elem():          &encInt64,
		reflect.TypeOf((*uint)(nil)).Elem():           &encUint,
		reflect.TypeOf((*uint8)(nil)).Elem():          &encUint8,
		reflect.TypeOf((*uint16)(nil)).Elem():         &encUint16,
		reflect.TypeOf((*uint32)(nil)).Elem():         &encUint32,
		reflect.TypeOf((*uint64)(nil)).Elem():         &encUint64,
		reflect.TypeOf((*uintptr)(nil)).Elem():        &encUintptr,
		reflect.TypeOf((*unsafe.Pointer)(nil)).Elem(): &encPointer,
		reflect.TypeOf((*float32)(nil)).Elem():        &encFloat32,
		reflect.TypeOf((*float64)(nil)).Elem():        &encFloat64,
		reflect.TypeOf((*complex64)(nil)).Elem():      &encComplex64,
		reflect.TypeOf((*complex128)(nil)).Elem():     &encComplex128,
		reflect.TypeOf((*[]byte)(nil)).Elem():         &encBytes,
		reflect.TypeOf((*string)(nil)).Elem():         &encString,
		reflect.TypeOf((*struct{})(nil)):              &encignore,
		reflect.TypeOf(nil):                           &encignore,
	}

	encLock sync.RWMutex
)

func getEncEngine(rt reflect.Type) encEng {
	encLock.RLock()
	engptr := rt2encEng[rt]
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
	//TODO  接口类型处理
	engine, has := rt2encEng[rt]
	if has {
		return engine
	}
	engine = new(func(*Encoder, unsafe.Pointer))
	rt2encEng[rt] = engine

	if fn, _, yes := implementsGob(rt); yes {
		*engine = func(e *Encoder, p unsafe.Pointer) {
			buf, err := fn(reflect.NewAt(rt, p).Elem().Interface().(gob.GobEncoder))
			if err != nil {
				panic(err)
			}
			e.encLength(len(buf))
			e.buf = append(e.buf, buf...)
		}
		return engine
	}

	if fn, _, yes := implementsBin(rt); yes {
		*engine = func(e *Encoder, p unsafe.Pointer) {
			buf, err := fn(reflect.NewAt(rt, p).Elem().Interface().(encoding.BinaryMarshaler))
			if err != nil {
				panic(err)
			}
			e.encLength(len(buf))
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
				(*eEng)(e, *(*unsafe.Pointer)(p))
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
				header := (*sliceHeader)(p)
				l := header.len
				e.encLength(l)
				for i := uintptr(0); i < uintptr(l); i++ {
					(*eEng)(e, unsafe.Pointer(uintptr(header.data)+i*size))
				}
			}
		}
	case reflect.Map:
		eKey := buildEncEngine(rt.Key())
		eEng = buildEncEngine(rt.Elem())
		*engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encBool(isNotNil)
			if isNotNil {
				e.encLength(*(*int)(*(*unsafe.Pointer)(p)))
				v := reflect.NewAt(rt, p).Elem()
				// TODO flag&flagIndir 在编译时确定
				for _, key := range v.MapKeys() {
					val := v.MapIndex(key)
					kv := (*refVal)(unsafe.Pointer(&key))
					vv := (*refVal)(unsafe.Pointer(&val))
					if kv.flag&flagIndir == 0 {
						(*eKey)(e, unsafe.Pointer(&kv.ptr))
					} else {
						(*eKey)(e, kv.ptr)
					}
					if vv.flag&flagIndir == 0 {
						(*eEng)(e, unsafe.Pointer(&vv.ptr))
					} else {
						(*eEng)(e, vv.ptr)
					}
				}
			}
		}
	case reflect.Struct:
		nf := rt.NumField()
		engs, offs := make([]encEngPtr, nf), make([]uintptr, nf)
		for i := 0; i < nf; i++ {
			field := rt.Field(i)
			engs[i] = buildEncEngine(field.Type)
			offs[i] = field.Offset
		}
		*engine = func(e *Encoder, p unsafe.Pointer) {
			for i := 0; i < nf; i++ {
				(*engs[i])(e, unsafe.Pointer(uintptr(p)+offs[i]))
			}
		}
	case reflect.Chan, reflect.Func, reflect.Interface:
		panic("not suport " + rt.String() + " type")
	}
	return engine
}
