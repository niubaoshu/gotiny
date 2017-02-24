package gotiny

import (
	//"fmt"
	"encoding"
	"encoding/gob"
	"reflect"
	"sync"
	"unsafe"
)

type (
	decEng func(*Decoder, reflect.Value)
)

var (
	rt2decEng = map[reflect.Type]decEng{
		reflect.TypeOf((*string)(nil)).Elem():  decString,
		reflect.TypeOf((*bool)(nil)).Elem():    decBool,
		reflect.TypeOf((*uint8)(nil)).Elem():   decUint8,
		reflect.TypeOf((*int8)(nil)).Elem():    decInt8,
		reflect.TypeOf((*int)(nil)).Elem():     decInt,
		reflect.TypeOf((*uint)(nil)).Elem():    decUint,
		reflect.TypeOf((*int16)(nil)).Elem():   decInt,
		reflect.TypeOf((*int32)(nil)).Elem():   decInt,
		reflect.TypeOf((*int64)(nil)).Elem():   decInt,
		reflect.TypeOf((*uint16)(nil)).Elem():  decUint,
		reflect.TypeOf((*uint32)(nil)).Elem():  decUint,
		reflect.TypeOf((*uint64)(nil)).Elem():  decUint,
		reflect.TypeOf((*uintptr)(nil)).Elem(): decUint,
		reflect.TypeOf((*float64)(nil)).Elem(): decFloat,
		reflect.TypeOf((*float32)(nil)).Elem(): decFloat,
		//reflect.TypeOf((*complex64)(nil)).Elem():decComplex,
		//reflect.TypeOf((*complex128)(nil)).Elem():decComplex,
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
	if _, fn, yes := implementsGob(rt); yes {
		engine = func(d *Decoder, v reflect.Value) {
			length := int(d.decUint())
			start := d.index
			d.index += length
			fn(v.Addr().Interface().(gob.GobDecoder), d.buf[start:d.index])
		}
		goto end
	}

	if _, fn, yes := implementsBin(rt); yes {
		engine = func(d *Decoder, v reflect.Value) {
			length := int(d.decUint())
			start := d.index
			d.index += length
			fn(v.Addr().Interface().(encoding.BinaryUnmarshaler), d.buf[start:d.index])
		}
		goto end
	}

	// if _, mfunc, yes := implementsInterface(rt); yes {
	// 	engine = func(d *Decoder, v reflect.Value) {
	// 		length := d.decUint()
	// 		start := d.index
	// 		d.index += int(length)
	// 		mfunc.Call([]reflect.Value{v., reflect.ValueOf(d.buf[start:d.index])})
	// 	}
	// 	goto end
	// }

	switch rt.Kind() {
	case reflect.Complex64, reflect.Complex128:
		engine = decComplex
	case reflect.Ptr:
		et := rt.Elem()
		eEng := buildDecEngine(et)
		cnil := reflect.Zero(rt)
		engine = func(d *Decoder, v reflect.Value) {
			if d.decBool() {
				if v.IsNil() {
					v.Set(reflect.New(et))
				}
				eEng(d, v.Elem())
			} else if !v.IsNil() {
				v.Set(cnil)
			}
		}
	case reflect.Array:
		l := rt.Len()
		eEng := buildDecEngine(rt.Elem())
		engine = func(d *Decoder, v reflect.Value) {
			for i := 0; i < l; i++ {
				eEng(d, v.Index(i))
			}
		}
	case reflect.Slice:
		eEng := buildDecEngine(rt.Elem())
		cnil := reflect.Zero(rt)
		engine = func(d *Decoder, v reflect.Value) {
			if d.decBool() {
				l := int(d.decUint())
				if v.IsNil() || v.Cap() < l {
					v.Set(reflect.MakeSlice(rt, l, l))
				} else {
					v.Set(v.Slice(0, l))
				}
				for i := 0; i < l; i++ {
					eEng(d, v.Index(i))
				}
			} else if !v.IsNil() {
				//v.Set(reflect.NewAt(v.Type(), unsafe.Pointer(uintptr(0))).Elem())
				v.Set(cnil)
			}
		}
	case reflect.Map:
		kt, vt := rt.Key(), rt.Elem()
		kEng, vEng := buildDecEngine(kt), buildDecEngine(vt)
		cnil := reflect.Zero(rt)
		// todo 原始key 值和value值重复使用
		engine = func(d *Decoder, v reflect.Value) {
			if d.decBool() {
				l := int(d.decUint())
				if v.IsNil() {
					v.Set(reflect.MakeMap(rt))
				}
				for i := 0; i < l; i++ {
					key, val := reflect.New(kt).Elem(), reflect.New(vt).Elem()
					kEng(d, key)
					vEng(d, val)
					v.SetMapIndex(key, val)
				}
			} else if !v.IsNil() {
				v.Set(cnil)
			}
		}
	case reflect.Struct:
		nf := rt.NumField()
		if nf > 0 {
			fielsEng := make([]decEng, nf)
			for i := 0; i < nf; i++ {
				ft := rt.Field(i).Type
				//https://golang.org/pkg/reflect/#StructField
				if rt.Field(i).PkgPath == "" {
					fielsEng[i] = buildDecEngine(ft)
				} else {
					fielsEng[i] = func(decEng decEng) decEng {
						return func(d *Decoder, fv reflect.Value) {
							decEng(d, reflect.NewAt(ft, unsafe.Pointer(fv.UnsafeAddr())).Elem())
						}
					}(buildDecEngine(ft))
				}
			}
			engine = func(d *Decoder, v reflect.Value) {
				for i := 0; i < nf; i++ {
					fielsEng[i](d, v.Field(i))
				}
			}
		} else {
			engine = decignore
		}
	default:
		engine = decignore
	}
end:
	rt2decEng[rt] = engine
	return engine
}
