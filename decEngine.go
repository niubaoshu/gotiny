package gotiny

import (
	"encoding"
	"encoding/gob"
	"reflect"
	"sync"
	"unsafe"
)

type (
	decEng    func(*Decoder, unsafe.Pointer)
	decEngPtr *func(*Decoder, unsafe.Pointer)
)

var (
	rt2decEng = map[reflect.Type]decEngPtr{
		reflect.TypeOf((*bool)(nil)).Elem():           &decBool,
		reflect.TypeOf((*int)(nil)).Elem():            &decInt,
		reflect.TypeOf((*int8)(nil)).Elem():           &decInt8,
		reflect.TypeOf((*int16)(nil)).Elem():          &decInt16,
		reflect.TypeOf((*int32)(nil)).Elem():          &decInt32,
		reflect.TypeOf((*int64)(nil)).Elem():          &decInt64,
		reflect.TypeOf((*uint)(nil)).Elem():           &decUint,
		reflect.TypeOf((*uint8)(nil)).Elem():          &decUint8,
		reflect.TypeOf((*uint16)(nil)).Elem():         &decUint16,
		reflect.TypeOf((*uint32)(nil)).Elem():         &decUint32,
		reflect.TypeOf((*uint64)(nil)).Elem():         &decUint64,
		reflect.TypeOf((*uintptr)(nil)).Elem():        &decUintptr,
		reflect.TypeOf((*unsafe.Pointer)(nil)).Elem(): &decPointer,
		reflect.TypeOf((*float32)(nil)).Elem():        &decFloat32,
		reflect.TypeOf((*float64)(nil)).Elem():        &decFloat64,
		reflect.TypeOf((*complex64)(nil)).Elem():      &decComplex64,
		reflect.TypeOf((*complex128)(nil)).Elem():     &decComplex128,
		reflect.TypeOf((*[]byte)(nil)).Elem():         &decBytes,
		reflect.TypeOf((*string)(nil)).Elem():         &decString,
		reflect.TypeOf((*struct{})(nil)).Elem():       &decIgnore,
		reflect.TypeOf(nil):                           &decIgnore,
	}

	baseDecEngines = []decEng{
		reflect.Invalid:       decIgnore,
		reflect.Bool:          decBool,
		reflect.Int:           decInt,
		reflect.Int8:          decInt8,
		reflect.Int16:         decInt16,
		reflect.Int32:         decInt32,
		reflect.Int64:         decInt64,
		reflect.Uint:          decUint,
		reflect.Uint8:         decUint8,
		reflect.Uint16:        decUint16,
		reflect.Uint32:        decUint32,
		reflect.Uint64:        decUint64,
		reflect.Uintptr:       decUintptr,
		reflect.UnsafePointer: decPointer,
		reflect.Float32:       decFloat32,
		reflect.Float64:       decFloat64,
		reflect.Complex64:     decComplex64,
		reflect.Complex128:    decComplex128,
		reflect.String:        decString,
	}
	decLock sync.RWMutex
)

func getDecEngine(rt reflect.Type) decEng {
	decLock.RLock()
	engPtr := rt2decEng[rt]
	decLock.RUnlock()
	if engPtr != nil {
		eng := *engPtr
		if eng != nil {
			return eng
		}
	}
	decLock.Lock()
	engPtr = buildDecEngine(rt)
	decLock.Unlock()
	return *engPtr
}

func buildDecEngine(rt reflect.Type) decEngPtr {
	engPtr, has := rt2decEng[rt]
	if has {
		return engPtr
	}
	engPtr = new(func(*Decoder, unsafe.Pointer))
	rt2decEng[rt] = engPtr

	rtPtr := reflect.PtrTo(rt)
	if rtPtr.Implements(gobType) {
		*engPtr = func(d *Decoder, p unsafe.Pointer) {
			length := d.decLength()
			start := d.index
			d.index += length
			if err := reflect.NewAt(rt, p).Interface().(gob.GobDecoder).GobDecode(d.buf[start:d.index]); err != nil {
				panic(err)
			}
		}
		return engPtr
	}
	if rtPtr.Implements(binType) {
		*engPtr = func(d *Decoder, p unsafe.Pointer) {
			length := d.decLength()
			start := d.index
			d.index += length
			if err := reflect.NewAt(rt, p).Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(d.buf[start:d.index]); err != nil {
				panic(err)
			}
		}
		return engPtr
	}

	if rtPtr.Implements(tinyType) {
		*engPtr = func(d *Decoder, p unsafe.Pointer) {
			d.index += reflect.NewAt(rt, p).Interface().(GoTinySerializer).GotinyDecode(d.buf[d.index:])
		}
		return engPtr
	}

	kind := rt.Kind()
	switch kind {
	case reflect.Ptr:
		et := rt.Elem()
		eEng := buildDecEngine(et) // TODO 可以考虑在生成编码机的时候解引用掉子编码机，下同
		*engPtr = func(d *Decoder, p unsafe.Pointer) {
			if d.decBool() {
				if isNil(p) {
					*(*unsafe.Pointer)(p) = unsafe.Pointer(reflect.New(et).Elem().UnsafeAddr())
				}
				(*eEng)(d, *(*unsafe.Pointer)(p))
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
			}
		}
	case reflect.Array:
		l, et := rt.Len(), rt.Elem()
		eEng, size := buildDecEngine(et), et.Size()
		*engPtr = func(d *Decoder, p unsafe.Pointer) {
			eng := *eEng
			for i := 0; i < l; i++ {
				eng(d, unsafe.Pointer(uintptr(p)+uintptr(i)*size))
			}
		}
	case reflect.Slice:
		et := rt.Elem()
		eEng, size := buildDecEngine(et), et.Size()
		*engPtr = func(d *Decoder, p unsafe.Pointer) {
			header := (*sliceHeader)(p)
			if d.decBool() {
				l := d.decLength()
				if isNil(p) || header.cap < l {
					*header = sliceHeader{unsafe.Pointer(reflect.MakeSlice(rt, l, l).Pointer()), l, l}
				} else {
					header.len = l
				}
				eng := *eEng
				for i := uintptr(0); i < uintptr(l); i++ {
					eng(d, unsafe.Pointer(uintptr(header.data)+i*size))
				}
			} else if !isNil(p) {
				*header = sliceHeader{}
			}
		}
	case reflect.Map:
		kt, vt := rt.Key(), rt.Elem()
		kEng, vEng := buildDecEngine(kt), buildDecEngine(vt)
		*engPtr = func(d *Decoder, p unsafe.Pointer) {
			if d.decBool() {
				l := d.decLength()
				if isNil(p) {
					*(*unsafe.Pointer)(p) = unsafe.Pointer(reflect.MakeMap(rt).Pointer())
				}
				v := reflect.NewAt(rt, p).Elem()
				// TODO 考虑重用v中的key和value，可以重用v.Len()个
				engKey, engVal := *kEng, *vEng
				for i := 0; i < l; i++ {
					key, val := reflect.New(kt).Elem(), reflect.New(vt).Elem()
					engKey(d, unsafe.Pointer(key.UnsafeAddr()))
					engVal(d, unsafe.Pointer(val.UnsafeAddr()))
					v.SetMapIndex(key, val)
				}
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
			}
		}
	case reflect.Struct:
		nf := rt.NumField()
		engs, offs := make([]decEngPtr, nf), make([]uintptr, nf)
		for i := 0; i < nf; i++ {
			field := rt.Field(i)
			engs[i] = buildDecEngine(field.Type)
			offs[i] = field.Offset
		}
		*engPtr = func(d *Decoder, p unsafe.Pointer) {
			for i := 0; i < nf; i++ {
				(*engs[i])(d, unsafe.Pointer(uintptr(p)+offs[i]))
			}
		}
	case reflect.Interface:
		*engPtr = func(d *Decoder, p unsafe.Pointer) {
			if d.decBool() {
				name := ""
				decString(d, unsafe.Pointer(&name))
				et, has := name2type[name]
				if !has {
					panic("unknown typ:" + name)
				}
				v := reflect.NewAt(rt, p).Elem()
				var ev reflect.Value
				if v.IsNil() || v.Elem().Type() != et {
					ev = reflect.New(et).Elem()
				} else {
					ev = v.Elem()
				}
				vv := (*refVal)(unsafe.Pointer(&ev))
				vp := vv.ptr
				if vv.flag&flagIndir == 0 {
					vp = unsafe.Pointer(&vv.ptr)
				}
				(*buildDecEngine(et))(d, vp)
				v.Set(ev)
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
			}
		}
	case reflect.Chan, reflect.Func:
		panic("not support " + rt.String() + " type")
	default:
		*engPtr = baseDecEngines[kind]
	}
	return engPtr
}
