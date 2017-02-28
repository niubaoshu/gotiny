package gotiny

import (

	//"encoding"
	//"encoding/gob"

	"reflect"
	"sync"
	"unsafe"
)

type (
	decEng func(*Decoder, unsafe.Pointer)
)

var (
	rt2decEng = map[reflect.Type]decEng{
		reflect.TypeOf((*string)(nil)).Elem():  decString,
		reflect.TypeOf((*bool)(nil)).Elem():    decBool,
		reflect.TypeOf((*int)(nil)).Elem():     decInt,
		reflect.TypeOf((*int8)(nil)).Elem():    decUint8,
		reflect.TypeOf((*int16)(nil)).Elem():   decInt16,
		reflect.TypeOf((*int32)(nil)).Elem():   decInt32,
		reflect.TypeOf((*int64)(nil)).Elem():   decInt64,
		reflect.TypeOf((*uint)(nil)).Elem():    decUint,
		reflect.TypeOf((*uint8)(nil)).Elem():   decUint8,
		reflect.TypeOf((*uint16)(nil)).Elem():  decUint16,
		reflect.TypeOf((*uint32)(nil)).Elem():  decUint32,
		reflect.TypeOf((*uint64)(nil)).Elem():  decUint64,
		reflect.TypeOf((*uintptr)(nil)).Elem(): decUint,
		reflect.TypeOf((*float64)(nil)).Elem(): decFloat64,
		reflect.TypeOf((*float32)(nil)).Elem(): decFloat32,
		//reflect.TypeOf(nil):                    decignore,
	}
	englock sync.RWMutex
)

func GetDecEngine(rt reflect.Type) decEng {
	englock.RLock()
	engine, has := rt2decEng[rt]
	englock.RUnlock()
	if has {
		return engine
	}
	englock.Lock()
	engine = buildDecEngine(rt)
	englock.Unlock()

	return engine

}

func buildDecEngine(rt reflect.Type) decEng {
	engine, has := rt2decEng[rt]
	if has {
		return engine
	}
	//if _, fn, yes := implementsGob(rt); yes {
	//	engine = func(d *Decoder, p unsafe.Pointer) {
	//		length := int(d.decUint())
	//		start := d.index
	//		d.index += length
	//		fn(v.Addr().Interface().(gob.GobDecoder), d.buf[start:d.index])
	//	}
	//	goto end
	//}
	//
	//if _, fn, yes := implementsBin(rt); yes {
	//	engine = func(d *Decoder, p unsafe.Pointer) {
	//		length := int(d.decUint())
	//		start := d.index
	//		d.index += length
	//		fn(v.Addr().Interface().(encoding.BinaryUnmarshaler), d.buf[start:d.index])
	//	}
	//	goto end
	//}

	// if _, mfunc, yes := implementsInterface(rt); yes {
	// 	engine = func(d *Decoder, p unsafe.Pointer) {
	// 		length := d.decUint()
	// 		start := d.index
	// 		d.index += int(length)
	// 		mfunc.Call([]reflect.Value{v., reflect.ValueOf(d.buf[start:d.index])})
	// 	}
	// 	goto end
	// }

	switch rt.Kind() {
	case reflect.Complex64:
		engine = decComplex64
	case reflect.Complex128:
		engine = decComplex128
	case reflect.Ptr:
		et := rt.Elem()
		eEng := buildDecEngine(et)
		engine = func(d *Decoder, p unsafe.Pointer) {
			if d.decBool() {
				if isNil(p) {
					*(*uintptr)(p) = reflect.New(et).Elem().UnsafeAddr()
				}
				eEng(d, elem(p))
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
			}
		}
	case reflect.Array:
		l := rt.Len()
		et := rt.Elem()
		eEng := buildDecEngine(et)
		size := et.Size()
		engine = func(d *Decoder, p unsafe.Pointer) {
			for i := 0; i < l; i++ {
				eEng(d, unsafe.Pointer(uintptr(p)+uintptr(i)*size))
			}
		}
	case reflect.Slice:
		et := rt.Elem()
		eEng := buildDecEngine(et)
		size := et.Size()
		engine = func(d *Decoder, p unsafe.Pointer) {
			if d.decBool() {
				l := d.decLength()
				if isNil(p) || *(*int)(unsafe.Pointer(uintptr(p) + ptrSize + ptrSize)) < l {
					*(*uintptr)(p) = reflect.MakeSlice(rt, l, l).Pointer()
				}
				*(*int)(unsafe.Pointer(uintptr(p) + ptrSize)) = l
				*(*int)(unsafe.Pointer(uintptr(p) + ptrSize + ptrSize)) = l
				pp := *(*unsafe.Pointer)(p)
				for i := 0; i < l; i++ {
					eEng(d, unsafe.Pointer(uintptr(pp)+uintptr(i)*size))
				}
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
				*(*int)(unsafe.Pointer(uintptr(p) + ptrSize)) = 0
				*(*int)(unsafe.Pointer(uintptr(p) + ptrSize + ptrSize)) = 0
			}
		}
	case reflect.Map:
		kt, vt := rt.Key(), rt.Elem()
		kEng, vEng := buildDecEngine(kt), buildDecEngine(vt)
		// todo 原始key 值和value值重复使用
		engine = func(d *Decoder, p unsafe.Pointer) {
			if d.decBool() {
				l := d.decLength()
				if isNil(p) {
					*(*uintptr)(p) = reflect.MakeMap(rt).Pointer()
				}
				v := reflect.NewAt(rt, p).Elem()
				for i := 0; i < l; i++ {
					key, val := reflect.New(kt).Elem(), reflect.New(vt).Elem()
					kEng(d, unsafe.Pointer(key.UnsafeAddr()))
					vEng(d, unsafe.Pointer(val.UnsafeAddr()))
					v.SetMapIndex(key, val)
				}
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
			}
		}
	case reflect.Struct:
		nf := rt.NumField()
		if nf > 0 {
			engs, offs := make([]decEng, nf), make([]uintptr, nf)
			for i := 0; i < nf; i++ {
				field := rt.Field(i)
				engs[i] = buildDecEngine(field.Type)
				offs[i] = field.Offset
			}
			engine = func(d *Decoder, p unsafe.Pointer) {
				for i := 0; i < nf; i++ {
					engs[i](d, unsafe.Pointer(uintptr(p)+offs[i]))
				}
			}
		} else {
			engine = decignore
		}
	default:
		engine = decignore
	}
	//end:
	rt2decEng[rt] = engine
	return engine
}
