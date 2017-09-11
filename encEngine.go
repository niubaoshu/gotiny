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
		reflect.TypeOf((*struct{})(nil)).Elem():       &encignore,
		reflect.TypeOf(nil):                           &encignore,
	}
	eengs = [...]encEng{
		reflect.Invalid:       encignore,
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

	interTypes []reflect.Type
	interRTMap map[reflect.Type]int = map[reflect.Type]int{}
	interNames []string             = []string{}
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

	if reflect.PtrTo(rt).Implements(gobType) {
		*engine = func(e *Encoder, p unsafe.Pointer) {
			buf, err := reflect.NewAt(rt, p).Interface().(gob.GobEncoder).GobEncode()
			if err != nil {
				panic(err)
			}
			e.encLength(len(buf))
			e.buf = append(e.buf, buf...)
		}
		return engine
	}

	if reflect.PtrTo(rt).Implements(binType) {
		*engine = func(e *Encoder, p unsafe.Pointer) {
			buf, err := reflect.NewAt(rt, p).Interface().(encoding.BinaryMarshaler).MarshalBinary()
			if err != nil {
				panic(err)
			}
			e.encLength(len(buf))
			e.buf = append(e.buf, buf...)
		}
		return engine
	}

	if reflect.PtrTo(rt).Implements(tinyType) {
		*engine = func(e *Encoder, p unsafe.Pointer) {
			e.buf = reflect.NewAt(rt, p).Interface().(GoTinySerializer).GotinyEncode(e.buf)
		}
		return engine
	}

	kind := rt.Kind()
	switch kind {
	case reflect.Ptr:
		eEng := buildEncEngine(rt.Elem())
		*engine = func(e *Encoder, p unsafe.Pointer) {
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
		*engine = func(e *Encoder, p unsafe.Pointer) {
			for i := 0; i < l; i++ {
				(*eEng)(e, unsafe.Pointer(uintptr(p)+uintptr(i)*size))
			}
		}
	case reflect.Slice:
		et := rt.Elem()
		eEng, size := buildEncEngine(et), et.Size()
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
		eKey, eEng := buildEncEngine(rt.Key()), buildEncEngine(rt.Elem())
		*engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encBool(isNotNil)
			if isNotNil {
				e.encLength(*(*int)(*(*unsafe.Pointer)(p)))
				v := reflect.NewAt(rt, p).Elem()
				// TODO flag&flagIndir 在编译时确定
				for _, key := range v.MapKeys() {
					val := v.MapIndex(key)
					kv, vv := (*refVal)(unsafe.Pointer(&key)), (*refVal)(unsafe.Pointer(&val))
					kp, vp := kv.ptr, vv.ptr
					if kv.flag&flagIndir == 0 {
						kp = unsafe.Pointer(&kv.ptr)
					}
					if vv.flag&flagIndir == 0 {
						vp = unsafe.Pointer(&vv.ptr)
					}
					(*eKey)(e, kp)
					(*eEng)(e, vp)
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
	case reflect.Interface:
		if rt.NumMethod() > 0 {
			*engine = func(e *Encoder, p unsafe.Pointer) {
				isNotNil := !isNil(p)
				e.encBool(isNotNil)
				if isNotNil {
					v := reflect.ValueOf(*(*interface {
						M()
					})(p))
					et := v.Type()
					e.encLength(getRTID(et))
					eEng := buildEncEngine(et)
					vv := (*refVal)(unsafe.Pointer(&v))
					vp := vv.ptr
					if vv.flag&flagIndir == 0 {
						vp = unsafe.Pointer(&vv.ptr)
					}
					(*eEng)(e, vp)
				}
			}
		} else {
			*engine = func(e *Encoder, p unsafe.Pointer) {
				isNotNil := !isNil(p)
				e.encBool(isNotNil)
				if isNotNil {
					v := reflect.ValueOf(*(*interface{})(p))
					et := v.Type()
					e.encLength(getRTID(et))
					eEng := buildEncEngine(et)
					vv := (*refVal)(unsafe.Pointer(&v))
					vp := vv.ptr
					if vv.flag&flagIndir == 0 {
						vp = unsafe.Pointer(&vv.ptr)
					}
					(*eEng)(e, vp)
				}
			}
		}
	case reflect.Chan, reflect.Func:
		panic("not support " + rt.String() + " type")
	default:
		*engine = eengs[kind]
	}
	return engine
}

func getRTID(rt reflect.Type) int {
	if id, has := interRTMap[rt]; has {
		return id
	} else {
		id = len(interTypes)
		interTypes = append(interTypes, rt)
		interRTMap[rt] = id
		return id
	}
}

func Register(i interface{}) {
	register(reflect.TypeOf(i))
}

func register(rt reflect.Type) int {
	name := getName(rt)
	i := len(interNames)
	interNames = append(interNames, "")
	for i > 0 {
		if interNames[i-1] > name {
			interTypes[i] = interTypes[i-1]
			interNames[i] = interNames[i-1]
		}
		i--
	}
	interNames[i] = name
	interTypes[i] = rt
	interRTMap[rt] = i
	return i
}

func getName(rt reflect.Type) string {

	name := rt.String()

	// But for named types (or pointers to them), qualify with import path (but see inner comment).
	// Dereference one pointer looking for a named type.
	star := ""
	if rt.Name() == "" {
		if rt.Kind() == reflect.Ptr {
			star = "*"
			rt = rt.Elem()
		} else {
			panic("not support no named type " + name)
		}
	}
	if rt.Name() != "" {
		if rt.PkgPath() == "" {
			name = star + rt.Name()
		} else {
			name = star + rt.PkgPath() + "." + rt.Name()
		}
	}
	return name
}
