package gotiny

import (
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type decEng func(*Decoder, unsafe.Pointer) // Decoder

var (
	rt2decEng = map[reflect.Type]decEng{
		reflect.TypeOf((*bool)(nil)).Elem():           decBool,
		reflect.TypeOf((*int)(nil)).Elem():            decInt,
		reflect.TypeOf((*int8)(nil)).Elem():           decInt8,
		reflect.TypeOf((*int16)(nil)).Elem():          decInt16,
		reflect.TypeOf((*int32)(nil)).Elem():          decInt32,
		reflect.TypeOf((*int64)(nil)).Elem():          decInt64,
		reflect.TypeOf((*uint)(nil)).Elem():           decUint,
		reflect.TypeOf((*uint8)(nil)).Elem():          decUint8,
		reflect.TypeOf((*uint16)(nil)).Elem():         decUint16,
		reflect.TypeOf((*uint32)(nil)).Elem():         decUint32,
		reflect.TypeOf((*uint64)(nil)).Elem():         decUint64,
		reflect.TypeOf((*uintptr)(nil)).Elem():        decUintptr,
		reflect.TypeOf((*unsafe.Pointer)(nil)).Elem(): decPointer,
		reflect.TypeOf((*float32)(nil)).Elem():        decFloat32,
		reflect.TypeOf((*float64)(nil)).Elem():        decFloat64,
		reflect.TypeOf((*complex64)(nil)).Elem():      decComplex64,
		reflect.TypeOf((*complex128)(nil)).Elem():     decComplex128,
		reflect.TypeOf((*[]byte)(nil)).Elem():         decBytes,
		reflect.TypeOf((*string)(nil)).Elem():         decString,
		reflect.TypeOf((*time.Time)(nil)).Elem():      decTime,
		reflect.TypeOf((*struct{})(nil)).Elem():       decIgnore,
		reflect.TypeOf(nil):                           decIgnore,
	}

	baseDecEngines = [...]decEng{
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

/*
getDecEngine is a function that retrieves a decoding engine for a given reflect.Type from a map.
The function takes a reflect.Type as its only argument and returns a decEng (decoding engine) struct.

The function first checks if the decoding engine already exists in the rt2decEng map.
If it does, it returns it. If it does not, it acquires a write lock using decLock,
and checks again if the engine exists in the map. If it does, it returns it.
If the engine still does not exist in the map, the function builds it using the buildDecEngine function,
adds it to the map, and then returns it.

This function is designed to be thread-safe, as it uses a lock to prevent multiple goroutines from
building the same decoding engine at the same time.
*/
func getDecEngine(rt reflect.Type) decEng {
	// Check if the decoding engine is already in the map. If it is, return it.
	engine := rt2decEng[rt]
	if engine != nil {
		return engine
	}

	// If the engine is not in the map, acquire the write lock and build it.
	// This is to prevent multiple goroutines from building the engine at the same time.
	decLock.Lock()
	defer decLock.Unlock()
	engine = rt2decEng[rt]
	if engine != nil {
		return engine
	}

	// If the engine is still not in the map, build it using the buildDecEngine function.
	buildDecEngine(rt, &engine)
	rt2decEng[rt] = engine
	return engine
}

// buildDecEngine takes a reflect type and a pointer to a decEng, and builds the decoding engine for that type.
func buildDecEngine(rt reflect.Type, engPtr *decEng) {
	// check if engine for the type already exists in the map
	engine, has := rt2decEng[rt]
	if has {
		*engPtr = engine
		return
	}
	// check if the type implements an unsupported serializer
	if _, engine = implementOtherSerializer(rt); engine != nil {
		rt2decEng[rt] = engine
		*engPtr = engine
		return
	}
	// determine the kind of the type
	kind := rt.Kind()
	var eEng decEng

	// build engine based on the kind of the type
	switch kind {
	// This case statement checks if the type is a pointer
	case reflect.Ptr:
		// et is set to the element type of the pointer type
		et := rt.Elem()
		// buildDecEngine is called with et and a pointer to eEng is passed to set it
		// to the new decoding engine
		defer buildDecEngine(et, &eEng)
		// a new function called "engine" is defined with parameters d *Decoder and p unsafe.Pointer
		engine = func(d *Decoder, p unsafe.Pointer) {
			// If the decoder is not nil
			if d.decIsNotNil() {
				// If the pointer is nil
				if isNil(p) {
					// A new element of et is created with reflect.New
					// It is then dereferenced with Elem() and converted to an unsafe.Pointer
					// It is then set to the value of p as an unsafe.Pointer
					*(*unsafe.Pointer)(p) = unsafe.Pointer(reflect.New(et).Elem().UnsafeAddr())
				}
				// The decoding engine eEng is called with d and the dereferenced value of p as an unsafe.Pointer
				eEng(d, *(*unsafe.Pointer)(p))
				// If the decoder is nil and the pointer is not nil
			} else if !isNil(p) {
				// Set the value of the pointer to nil
				*(*unsafe.Pointer)(p) = nil
			}
		}

	//the case where the type of the data being decoded is an array.
	case reflect.Array:
		// Get the length of the array and the element type.
		l, et := rt.Len(), rt.Elem()

		// Get the size of each element in the array.
		size := et.Size()

		// Create a deferred decoding engine for the element type.
		defer buildDecEngine(et, &eEng)

		// Define the decoding engine for the array.
		engine = func(d *Decoder, p unsafe.Pointer) {

			// Loop over each element in the array.
			for i := 0; i < l; i++ {
				// Pass the decoder and a pointer to the current element to the decoding engine for the element type.
				eEng(d, unsafe.Pointer(uintptr(p)+uintptr(i)*size))
			}
		}
	case reflect.Slice:
		et := rt.Elem()
		size := et.Size()
		defer buildDecEngine(et, &eEng)
		engine = func(d *Decoder, p unsafe.Pointer) {
			header := (*reflect.SliceHeader)(p)
			if d.decIsNotNil() {
				l := d.decLength()
				if isNil(p) || header.Cap < l {
					*header = reflect.SliceHeader{Data: reflect.MakeSlice(rt, l, l).Pointer(), Len: l, Cap: l}
				} else {
					header.Len = l
				}
				for i := 0; i < l; i++ {
					eEng(d, unsafe.Pointer(header.Data+uintptr(i)*size))
				}
			} else if !isNil(p) {
				*header = reflect.SliceHeader{}
			}
		}
	case reflect.Map:
		kt, vt := rt.Key(), rt.Elem()
		skt, svt := reflect.SliceOf(kt), reflect.SliceOf(vt)
		var kEng, vEng decEng
		defer buildDecEngine(kt, &kEng)
		defer buildDecEngine(vt, &vEng)
		engine = func(d *Decoder, p unsafe.Pointer) {
			if d.decIsNotNil() {
				l := d.decLength()
				var v reflect.Value
				if isNil(p) {
					v = reflect.MakeMapWithSize(rt, l)
					*(*unsafe.Pointer)(p) = unsafe.Pointer(v.Pointer())
				} else {
					v = reflect.NewAt(rt, p).Elem()
				}
				keys, vals := reflect.MakeSlice(skt, l, l), reflect.MakeSlice(svt, l, l)
				for i := 0; i < l; i++ {
					key, val := keys.Index(i), vals.Index(i)
					kEng(d, unsafe.Pointer(key.UnsafeAddr()))
					vEng(d, unsafe.Pointer(val.UnsafeAddr()))
					v.SetMapIndex(key, val)
				}
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
			}
		}
	case reflect.Struct:
		fields, offs := getFieldType(rt, 0)
		nf := len(fields)
		fEngines := make([]decEng, nf)
		defer func() {
			for i := 0; i < nf; i++ {
				buildDecEngine(fields[i], &fEngines[i])
			}
		}()
		engine = func(d *Decoder, p unsafe.Pointer) {
			for i := 0; i < len(fEngines) && i < len(offs); i++ {
				fEngines[i](d, unsafe.Pointer(uintptr(p)+offs[i]))
			}
		}
	case reflect.Interface:
		engine = func(d *Decoder, p unsafe.Pointer) {
			if d.decIsNotNil() {
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
				getDecEngine(et)(d, getUnsafePointer(&ev))
				v.Set(ev)
			} else if !isNil(p) {
				*(*unsafe.Pointer)(p) = nil
			}
		}
	case reflect.Chan, reflect.Func:
		panic("not support " + rt.String() + " type")
	default:
		engine = baseDecEngines[kind]
	}
	rt2decEng[rt] = engine
	*engPtr = engine
}
