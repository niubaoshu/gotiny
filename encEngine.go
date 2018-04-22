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
		reflect.TypeOf((*struct{})(nil)).Elem():       &encIgnore,
		reflect.TypeOf(nil):                           &encIgnore,
	}
	encEngines = [...]encEng{
		reflect.Invalid:       encIgnore,
		reflect.Bool:          encBool,
		reflect.Int:           encInt,
		reflect.Int8:          encInt8,
		reflect.Int16:         encInt16,
		reflect.Int32:         encInt32,
		reflect.Int64:         encInt64,
		reflect.Uint:          encUint,
		reflect.Uint8:         encUint8,
		reflect.Uint16:        encUint16,
		reflect.Uint32:        encUint32,
		reflect.Uint64:        encUint64,
		reflect.Uintptr:       encUintptr,
		reflect.UnsafePointer: encPointer,
		reflect.Float32:       encFloat32,
		reflect.Float64:       encFloat64,
		reflect.Complex64:     encComplex64,
		reflect.Complex128:    encComplex128,
		reflect.String:        encString,
	}

	encLock sync.RWMutex

	type2name = map[reflect.Type]string{}
	name2type = map[string]reflect.Type{}
)

func getEncEngine(rt reflect.Type) encEng {
	encLock.RLock()
	engPtr := rt2encEng[rt]
	encLock.RUnlock()
	if engPtr != nil {
		eng := *engPtr
		if eng != nil {
			return eng
		}
	}
	encLock.Lock()
	engPtr = buildEncEngine(rt)
	encLock.Unlock()
	return *engPtr
}

func buildEncEngine(rt reflect.Type) encEngPtr {
	engPtr, has := rt2encEng[rt]
	if has {
		return engPtr
	}
	engPtr = new(func(*Encoder, unsafe.Pointer))
	rt2encEng[rt] = engPtr

	rtPtr := reflect.PtrTo(rt)
	if rtPtr.Implements(gobType) {
		*engPtr = func(e *Encoder, p unsafe.Pointer) {
			buf, err := reflect.NewAt(rt, p).Interface().(gob.GobEncoder).GobEncode()
			if err != nil {
				panic(err)
			}
			e.encLength(len(buf))
			e.buf = append(e.buf, buf...)
		}
		return engPtr
	}

	if rtPtr.Implements(binType) {
		*engPtr = func(e *Encoder, p unsafe.Pointer) {
			buf, err := reflect.NewAt(rt, p).Interface().(encoding.BinaryMarshaler).MarshalBinary()
			if err != nil {
				panic(err)
			}
			e.encLength(len(buf))
			e.buf = append(e.buf, buf...)
		}
		return engPtr
	}

	if rtPtr.Implements(tinyType) {
		*engPtr = func(e *Encoder, p unsafe.Pointer) {
			e.buf = reflect.NewAt(rt, p).Interface().(GoTinySerializer).GotinyEncode(e.buf)
		}
		return engPtr
	}

	kind := rt.Kind()
	switch kind {
	case reflect.Ptr:
		eEng := buildEncEngine(rt.Elem())
		*engPtr = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encBool(isNotNil)
			if isNotNil {
				(*eEng)(e, *(*unsafe.Pointer)(p))
			}
		}
	case reflect.Array:
		et, l := rt.Elem(), rt.Len()
		eEng := buildEncEngine(et)
		size := et.Size()
		*engPtr = func(e *Encoder, p unsafe.Pointer) {
			eng := *eEng
			for i := 0; i < l; i++ {
				eng(e, unsafe.Pointer(uintptr(p)+uintptr(i)*size))
			}
		}
	case reflect.Slice:
		et := rt.Elem()
		eEng, size := buildEncEngine(et), et.Size()
		*engPtr = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encBool(isNotNil)
			if isNotNil {
				header := (*sliceHeader)(p)
				l := header.len
				e.encLength(l)
				eng := *eEng
				for i := uintptr(0); i < uintptr(l); i++ {
					eng(e, unsafe.Pointer(uintptr(header.data)+i*size))
				}
			}
		}
	case reflect.Map:
		kEng, vEng := buildEncEngine(rt.Key()), buildEncEngine(rt.Elem())
		*engPtr = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encBool(isNotNil)
			if isNotNil {
				e.encLength(*(*int)(*(*unsafe.Pointer)(p)))
				v := reflect.NewAt(rt, p).Elem()
				// TODO flag&flagIndir 在编译时确定
				engKey, engVal := *kEng, *vEng
				keys := v.MapKeys()
				for i := 0; i < len(keys); i++ {
					val := v.MapIndex(keys[i])
					kv, vv := (*refVal)(unsafe.Pointer(&keys[i])), (*refVal)(unsafe.Pointer(&val))
					kp, vp := kv.ptr, vv.ptr
					if kv.flag&flagIndir == 0 {
						kp = unsafe.Pointer(&kv.ptr)
					}
					if vv.flag&flagIndir == 0 {
						vp = unsafe.Pointer(&vv.ptr)
					}
					engKey(e, kp)
					engVal(e, vp)
				}
			}
		}
	case reflect.Struct:
		nf := rt.NumField()
		engs, offs := make([]encEng, nf), make([]uintptr, nf)
		for i := 0; i < nf; i++ {
			field := rt.Field(i)
			engs[i] = *buildEncEngine(field.Type)
			offs[i] = field.Offset
		}
		*engPtr = func(e *Encoder, p unsafe.Pointer) {
			for i := 0; i < nf; i++ {
				engs[i](e, unsafe.Pointer(uintptr(p)+offs[i]))
			}
		}
	case reflect.Interface:
		if rt.NumMethod() > 0 {
			*engPtr = func(e *Encoder, p unsafe.Pointer) {
				isNotNil := !isNil(p)
				e.encBool(isNotNil)
				if isNotNil {
					v := reflect.ValueOf(*(*interface {
						M()
					})(p))
					et := v.Type()
					e.encString(getNameOfType(et))
					vv := (*refVal)(unsafe.Pointer(&v))
					vp := vv.ptr
					if vv.flag&flagIndir == 0 {
						vp = unsafe.Pointer(&vv.ptr)
					}
					(*buildEncEngine(et))(e, vp)
				}
			}
		} else {
			*engPtr = func(e *Encoder, p unsafe.Pointer) {
				isNotNil := !isNil(p)
				e.encBool(isNotNil)
				if isNotNil {
					v := reflect.ValueOf(*(*interface{})(p))
					et := v.Type()
					e.encString(getNameOfType(et))
					vv := (*refVal)(unsafe.Pointer(&v))
					vp := vv.ptr
					if vv.flag&flagIndir == 0 {
						vp = unsafe.Pointer(&vv.ptr)
					}
					(*buildEncEngine(et))(e, vp)
				}
			}
		}
	case reflect.Chan, reflect.Func:
		panic("not support " + rt.String() + " type")
	default:
		*engPtr = encEngines[kind]
	}
	return engPtr
}
