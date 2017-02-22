package gotiny

import (
	"reflect"
	"sync"
)

type encEng func(*Encoder, reflect.Value) //编码器

var (
	rt2Eng = map[reflect.Type]encEng{
		reflect.TypeOf((*string)(nil)).Elem():  encString,
		reflect.TypeOf((*bool)(nil)).Elem():    encBool,
		reflect.TypeOf((*uint8)(nil)).Elem():   encUint8,
		reflect.TypeOf((*int8)(nil)).Elem():    encInt8,
		reflect.TypeOf((*int)(nil)).Elem():     encInt,
		reflect.TypeOf((*uint)(nil)).Elem():    encUint,
		reflect.TypeOf((*int16)(nil)).Elem():   encInt,
		reflect.TypeOf((*int32)(nil)).Elem():   encInt,
		reflect.TypeOf((*int64)(nil)).Elem():   encInt,
		reflect.TypeOf((*uint16)(nil)).Elem():  encUint,
		reflect.TypeOf((*uint32)(nil)).Elem():  encUint,
		reflect.TypeOf((*uint64)(nil)).Elem():  encUint,
		reflect.TypeOf((*uintptr)(nil)).Elem(): encUint,
		reflect.TypeOf((*float64)(nil)).Elem(): encFloat,
		reflect.TypeOf((*float32)(nil)).Elem(): encFloat,
		//reflect.TypeOf((*complex64)(nil)).Elem():encComplex,
		//reflect.TypeOf((*complex128)(nil)).Elem():encComplex,
	}
	encLock sync.RWMutex
)

func GetEncEng(rt reflect.Type) (eng encEng) {
	encLock.RLock()
	eng = rt2Eng[rt]
	encLock.RUnlock()
	if eng != nil {
		return eng
	}
	encLock.Lock()
	defer encLock.Unlock()
	return buildEncEng(rt)
}

func buildEncEng(rt reflect.Type) (engine encEng) {
	//todo 循环类型和循环值处理
	// 循环类型 type x *x
	// 套嵌  type  a { b *a }
	//接口类型处理
	// 实现了 BinaryMarshaler
	// TextMarshaler
	//GobEncoder
	// 接口的处理
	engine = rt2Eng[rt]
	if engine != nil {
		return engine
	}
	if mfunc,_,yes:= implementsInterface(rt) ;yes{
		engine = func(e *Encoder,v reflect.Value){
			buf:= mfunc.Call([]reflect.Value{v})[0].Bytes()
			e.encUint(uint64(len(buf)))
			e.buf = append(e.buf, buf...)
		}
		goto end
	}

	switch rt.Kind() {
	case reflect.Complex64, reflect.Complex128:
		engine = encComplex
	case reflect.Ptr:
		eEng := buildEncEng(rt.Elem())
		engine = func(e *Encoder, v reflect.Value) {
			isNotNil := !v.IsNil()
			e.encBool(isNotNil)
			if isNotNil {
				eEng(e, v.Elem())
			}
		}
	case reflect.Array:
		eEng := buildEncEng(rt.Elem())
		l := rt.Len()
		engine = func(e *Encoder, v reflect.Value) {
			for i := 0; i < l; i++ {
				eEng(e, v.Index(i))
			}
		}
	case reflect.Slice:
		eEng := buildEncEng(rt.Elem())
		engine = func(e *Encoder, v reflect.Value) {
			isNotNil := !v.IsNil()
			e.encBool(isNotNil)
			if isNotNil {
				l := v.Len()
				e.encUint(uint64(l))
				for i := 0; i < l; i++ {
					eEng(e, v.Index(i))
				}
			}
		}
	case reflect.Map:
		kEng, eEng := buildEncEng(rt.Key()), buildEncEng(rt.Elem())
		//http://blog.csdn.net/hificamera/article/details/51701804
		engine = func(e *Encoder, v reflect.Value) {
			isNotNil := !v.IsNil()
			e.encBool(isNotNil)
			if isNotNil {
				keys := v.MapKeys()
				l := len(keys)
				e.encUint(uint64(l))
				for _, key := range keys {
					kEng(e, key)
					eEng(e, v.MapIndex(key))
				}
			}
		}
	case reflect.Struct:
		nf := rt.NumField()
		if nf > 0 {
			engs := make([]func(*Encoder, reflect.Value), nf)
			for i := 0; i < nf; i++ {
				engs[i] = buildEncEng(rt.Field(i).Type)
			}
			engine = func(e *Encoder, v reflect.Value) {
				for i := 0; i < nf; i++ {
					engs[i](e, v.Field(i))
				}
			}
		} else {
			engine = encignore
		}
	default:
		engine = encignore
	}
end:
	rt2Eng[rt] = engine
	return engine
}
