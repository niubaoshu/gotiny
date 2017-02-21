package gotiny

import (
	//"fmt"
	"reflect"
	"sync"
	"unsafe"
)

type (
	decEngine    func(*Decoder, reflect.Value)
	decEngByType func(*Decoder) reflect.Value
)

var (
	rt2decEng = map[reflect.Type]decEngine{}
	rt2decEBT = map[reflect.Type]decEngByType{}

	englock sync.RWMutex
	ebtlock sync.RWMutex
)

func GetDecEngine(rt reflect.Type) decEngine {
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

func buildDecEngine(rt reflect.Type) decEngine {
	engine, has := rt2decEng[rt]
	if has {
		return engine
	}
	if rt.Implements(gobEncIF) && rt.Implements(gobDecIF) {
		engine = decGob
		goto end
	}

	if rt.Implements(binEncIF) && rt.Implements(binDecIF) {
		engine = decBin
		goto end
	}
	//if rt.Implements(txtEncIF) && rt.Implements(txtDecIF) {
	//	engine = decTxt
	//	goto end
	//}

	switch rt.Kind() {
	case reflect.Bool:
		engine = decBool
	case reflect.Uint8:
		engine = decUint8
	case reflect.Int8:
		engine = decInt8
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		engine = decUint
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		engine = decInt
	case reflect.Float32, reflect.Float64:
		engine = decFloat
	case reflect.Complex64, reflect.Complex128:
		engine = decComplex
	case reflect.String:
		engine = decString
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
			fielsEng := make([]decEngine, nf)
			for i := 0; i < nf; i++ {
				ft := rt.Field(i).Type
				if rt.Field(i).PkgPath == "" {
					fielsEng[i] = buildDecEngine(ft)
				} else {
					fielsEng[i] = func(decEng decEngine) decEngine {
						return func(d *Decoder, fv reflect.Value) {
							decEng(d, reflect.NewAt(ft, unsafe.Pointer(fv.UnsafeAddr())).Elem())
						}
					}(buildDecEngine(ft))
					// fielsEng[i] = func(d *Decoder, fv reflect.Value) {
					// 	buildDecEngine(ft)(d, reflect.NewAt(ft, unsafe.Pointer(fv.UnsafeAddr())).Elem())
					// }
				}
				//https://golang.org/pkg/reflect/#StructField
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

// func GetDecEBT(rt reflect.Type) decEngByType {
// 	ebtlock.RLock()
// 	engine, has := rt2decEBT[rt]
// 	ebtlock.RUnlock()
// 	if has {
// 		return engine
// 	}
// 	ebtlock.Lock()
// 	engine = buildDecEBT(rt)
// 	ebtlock.Unlock()
// 	return engine
// }

// func buildDecEBT(rt reflect.Type) decEngByType {
// 	engine, has := rt2decEBT[rt]
// 	if has {
// 		return engine
// 	}
// 	switch rt.Kind() {
// 	case reflect.Bool:
// 		engine = func(d *Decoder) reflect.Value {
// 			return reflect.ValueOf(d.decBool())
// 		}
// 	case reflect.Uint8:
// 		engine = func(d *Decoder) reflect.Value {
// 			return reflect.ValueOf(d.decUint8())
// 		}
// 	case reflect.Int8:
// 		engine = func(d *Decoder) reflect.Value {
// 			return reflect.ValueOf(d.decInt8())
// 		}
// 	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
// 		engine = func(d *Decoder) reflect.Value {
// 			return reflect.ValueOf(d.decUint())
// 		}
// 	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
// 		engine = func(d *Decoder) reflect.Value {
// 			return reflect.ValueOf(d.decInt())
// 		}
// 	case reflect.Float32, reflect.Float64:
// 		engine = func(d *Decoder) reflect.Value {
// 			return reflect.ValueOf(d.decFloat())
// 		}
// 	case reflect.Complex64, reflect.Complex128:
// 		engine = func(d *Decoder) reflect.Value {
// 			return reflect.ValueOf(d.decComplex())
// 		}
// 	case reflect.String:
// 		engine = func(d *Decoder) reflect.Value {
// 			return reflect.ValueOf(d.decString())
// 		}
// 	case reflect.Ptr:
// 		et := rt.Elem()
// 		eEng := GetDecEngine(et)
// 		enil := reflect.Zero(rt)
// 		engine = func(d *Decoder) reflect.Value {
// 			if d.decBool() {
// 				v := reflect.New(et)
// 				eEng(d, v.Elem())
// 				return v
// 			} else {
// 				return enil
// 			}
// 		}
// 	case reflect.Array:
// 		l := rt.Len()
// 		eEng := GetDecEngine(rt.Elem())
// 		engine = func(d *Decoder) reflect.Value {
// 			v := reflect.New(rt).Elem()
// 			for i := 0; i < l; i++ {
// 				eEng(d, v.Index(i))
// 			}
// 			return v
// 		}
// 	case reflect.Slice:
// 		// et := rt.Elem()
// 		// eEng := buildDecEBT(et)
// 		// ek := et.Kind()
// 		// needCopy := ek == reflect.Ptr || ek == reflect.Map || ek == reflect.Slice
// 		eEng := buildDecEBT(rt.Elem())
// 		enil := reflect.Zero(rt)
// 		engine = func(d *Decoder, v reflect.Value) {
// 			if d.decBool() {
// 				l := int(d.decUint())
// 				if v.IsNil() || v.Cap() < l {
// 					//slice 内部元素的内存重用,只重用内部类型是slice,ptr,map的情况，
// 					//内部是array和struct时，他们的成员是前者的时候虽然也可以重用，但成本较高,暂不考虑
// 					//nv := reflect.MakeSlice(rt, l, l)
// 					//if needCopy {
// 					//	reflect.Copy(nv, v)
// 					//}
// 					v.Set(reflect.MakeSlice(rt, l, l))
// 				} else {
// 					v.Set(v.Slice(0, l))
// 				}
// 				for i := 0; i < l; i++ {
// 					eEng(d, v.Index(i))
// 				}
// 			} else if !v.IsNil() {
// 				//v.Set(reflect.NewAt(v.Type(), unsafe.Pointer(uintptr(0))).Elem())
// 				v.Set(enil)
// 			}
// 		}
// 	case reflect.Map:
// 		kt, vt := rt.Key(), rt.Elem()
// 		kEng, vEng := buildDecEBT(kt), buildDecEBT(vt)
// 		enil := reflect.Zero(rt)
// 		engine = func(d *Decoder, v reflect.Value) {
// 			if d.decBool() {
// 				l := int(d.decUint())
// 				if v.IsNil() {
// 					v.Set(reflect.MakeMap(rt))
// 				}
// 				for i := 0; i < l; i++ {
// 					key, val := reflect.New(kt).Elem(), reflect.New(vt).Elem()
// 					kEng(d, key)
// 					vEng(d, val)
// 					v.SetMapIndex(key, val)
// 				}
// 			} else if !v.IsNil() {
// 				v.Set(enil)
// 			}
// 		}
// 	case reflect.Struct:
// 		nf := rt.NumField()
// 		if nf > 0 {
// 			fielsEng, fielsType := make([]decEngine, nf), make([]reflect.Type, nf)
// 			for i := 0; i < nf; i++ {
// 				ft := rt.Field(i).Type
// 				if rt.Field(i).PkgPath == "" {
// 					fielsEng[i] = buildDecEBT(ft)
// 				} else {
// 					fielsEng[i] = func(d *Decoder, fv reflect.Value) {
// 						buildDecEBT(fielsType[i])(d, reflect.NewAt(ft, unsafe.Pointer(fv.UnsafeAddr())).Elem())
// 					}
// 				}
// 				//https://golang.org/pkg/reflect/#StructField
// 			}
// 			engine = func(d *Decoder, v reflect.Value) {
// 				for i := 0; i < nf; i++ {
// 					fielsEng[i](d, v.Field(i))
// 				}
// 			}
// 		}
// 	}
// }
