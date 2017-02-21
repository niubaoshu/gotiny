package gotiny

import (
	"reflect"
	"sync"
)

type TypeInfo struct {
	Length   int
	Engine   func(*Encoder, reflect.Value) //编码
	IsStatic bool                          // 该类型底层是否有过append操作
	Head     int
}

const (
	bit    = 1       // 1bit    1  bit
	byte1  = 1 << 3  // 8bit    1  byte
	byte2  = 3 << 3  // 16bit   3  byte
	byte4  = 5 << 3  // 32bit   5  byte
	byte8  = 10 << 3 // 64bit   10 byte
	byte16 = 20 << 3 // 128bit  20 byte
	word1  = 10 << 3 // 64bit   10 byte
)

var (
	baseLength = [...]int{
		reflect.Bool:          bit,
		reflect.Int:           word1,
		reflect.Ptr:           bit, // bool,ptr 编码1bit
		reflect.Int8:          byte1,
		reflect.Uint8:         byte1, //1个字节
		reflect.Int16:         byte2,
		reflect.Uint16:        byte2, //2个字节
		reflect.Int32:         byte4,
		reflect.Uint32:        byte4,
		reflect.Float32:       byte4, // 4个字节
		reflect.Int64:         byte8,
		reflect.Uint64:        byte8,
		reflect.Float64:       byte8,
		reflect.Complex64:     byte8, //8个字节
		reflect.Map:           word1 + bit,
		reflect.Slice:         word1 + bit, //一个机器字+1bit,长度和是否为nil
		reflect.String:        word1,       //string 不能为nil,只编码一个长度
		reflect.Uint:          word1,
		reflect.Uintptr:       word1,  //一个机器字
		reflect.Complex128:    byte16, //16个字节
		reflect.Interface:     0,
		reflect.Invalid:       0,
		reflect.Array:         0,
		reflect.Struct:        0, //数组不需要为长度编码 struct待拆解了后计算长度
		reflect.Chan:          0,
		reflect.Func:          0,
		reflect.UnsafePointer: 0, //不编码
	}

	baseInfo = [...]*TypeInfo{ //bit数
		reflect.Bool:          &TypeInfo{Length: baseLength[reflect.Bool], Head: baseLength[reflect.Bool], Engine: encBool, IsStatic: true},
		reflect.Uint8:         &TypeInfo{Length: baseLength[reflect.Uint8], Head: baseLength[reflect.Uint8], Engine: encUint8, IsStatic: true},
		reflect.Uint16:        &TypeInfo{Length: baseLength[reflect.Uint16], Head: baseLength[reflect.Uint16], Engine: encUint16, IsStatic: true},
		reflect.Uint32:        &TypeInfo{Length: baseLength[reflect.Uint32], Head: baseLength[reflect.Uint32], Engine: encUint32, IsStatic: true},
		reflect.Uint64:        &TypeInfo{Length: baseLength[reflect.Uint64], Head: baseLength[reflect.Uint64], Engine: encUint64, IsStatic: true},
		reflect.Uint:          &TypeInfo{Length: baseLength[reflect.Uint], Head: baseLength[reflect.Uint], Engine: encUint, IsStatic: true},
		reflect.Uintptr:       &TypeInfo{Length: baseLength[reflect.Uintptr], Head: baseLength[reflect.Uintptr], Engine: encUintptr, IsStatic: true},
		reflect.Int8:          &TypeInfo{Length: baseLength[reflect.Int8], Head: baseLength[reflect.Int8], Engine: encInt8, IsStatic: true},
		reflect.Int16:         &TypeInfo{Length: baseLength[reflect.Int16], Head: baseLength[reflect.Int16], Engine: encInt16, IsStatic: true},
		reflect.Int32:         &TypeInfo{Length: baseLength[reflect.Int32], Head: baseLength[reflect.Int32], Engine: encInt32, IsStatic: true},
		reflect.Int64:         &TypeInfo{Length: baseLength[reflect.Int64], Head: baseLength[reflect.Int64], Engine: encInt64, IsStatic: true},
		reflect.Int:           &TypeInfo{Length: baseLength[reflect.Int], Head: baseLength[reflect.Int], Engine: encInt, IsStatic: true},
		reflect.Float32:       &TypeInfo{Length: baseLength[reflect.Float32], Head: baseLength[reflect.Float32], Engine: encFloat32, IsStatic: true},
		reflect.Float64:       &TypeInfo{Length: baseLength[reflect.Float64], Head: baseLength[reflect.Float64], Engine: encFloat64, IsStatic: true},
		reflect.Complex64:     &TypeInfo{Length: baseLength[reflect.Complex64], Head: baseLength[reflect.Complex64], Engine: encComplex64, IsStatic: true},
		reflect.Complex128:    &TypeInfo{Length: baseLength[reflect.Complex128], Head: baseLength[reflect.Complex128], Engine: encComplex128, IsStatic: true},
		reflect.String:        &TypeInfo{Length: baseLength[reflect.String], Head: baseLength[reflect.String], Engine: encString, IsStatic: false},
		reflect.Invalid:       &TypeInfo{Length: 0, Head: 0, Engine: encignore, IsStatic: true},
		reflect.Interface:     &TypeInfo{Length: 0, Head: 0, Engine: encignore, IsStatic: true},
		reflect.Chan:          &TypeInfo{Length: 0, Head: 0, Engine: encignore, IsStatic: true},
		reflect.Func:          &TypeInfo{Length: 0, Head: 0, Engine: encignore, IsStatic: true},
		reflect.UnsafePointer: &TypeInfo{Length: 0, Head: 0, Engine: encignore, IsStatic: true},
	}
	rt2Info = map[reflect.Type]*TypeInfo{}
	encLock sync.RWMutex
)

func GetTypeInfo(rt reflect.Type) *TypeInfo {
	encLock.RLock()
	info, has := rt2Info[rt]
	encLock.RUnlock()
	if has {
		return info
	}
	encLock.Lock()
	defer encLock.Unlock()
	return buildInfo(rt)
}

func buildInfo(rt reflect.Type) *TypeInfo {
	//todo 循环类型和循环值处理
	// 循环类型 type x *x
	// 套嵌  type  a { b *a }
	//接口类型处理
	// 实现了 BinaryMarshaler
	// TextMarshaler
	//GobEncoder
	// 接口的处理
	info, has := rt2Info[rt]
	if has {
		return info
	}

	kind := rt.Kind()

	info = new(TypeInfo)
	if rt.Implements(gobEncIF) && rt.Implements(gobDecIF) {
		info.Engine = encGob
		//fmt.Println("encGob")
		goto end
	}

	if rt.Implements(binEncIF) && rt.Implements(binDecIF) {
		info.Engine = encBin
		goto end
	}
	 //if rt.Implements(txtEncIF) && rt.Implements(txtDecIF) {
	 //	info.Engine = encTxt
	 //	goto end
	 //}
	info.Length = baseLength[kind]
	switch kind {
	case reflect.Ptr:
		eInfo := buildInfo(rt.Elem())
		info.Length += eInfo.Length
		info.IsStatic = eInfo.IsStatic
		info.Head = eInfo.Head + baseLength[kind]
		info.Engine = func(e *Encoder, v reflect.Value) {
			isNotNil := !v.IsNil()
			e.encBool(isNotNil)
			if isNotNil {
				eInfo.Engine(e, v.Elem())
			}
		}
	case reflect.Array:
		eInfo := buildInfo(rt.Elem())
		l := rt.Len()
		info.Length += l * eInfo.Length
		info.IsStatic = eInfo.IsStatic
		lessOne := l - 1
		if eInfo.IsStatic {
			info.Engine = func(e *Encoder, v reflect.Value) {
				for i := 0; i < l; i++ {
					eInfo.Engine(e, v.Index(i))
				}
			}
			info.Head = l * eInfo.Length
		} else {
			info.Engine = func(e *Encoder, v reflect.Value) {
				reserved := e.reserved
				e.reserved = eInfo.Head
				for i := 0; i < lessOne; i++ {
					eInfo.Engine(e, v.Index(i))
				}
				e.reserved = reserved
				eInfo.Engine(e, v.Index(lessOne))
			}
			info.Head = eInfo.Head
		}
	case reflect.Slice:
		eInfo := buildInfo(rt.Elem())
		info.IsStatic = false
		if eInfo.IsStatic {
			info.Engine = func(e *Encoder, v reflect.Value) {
				isNotNil := !v.IsNil()
				e.encBool(isNotNil)
				if isNotNil {
					l := v.Len()
					//fmt.Printf("%#v\n", e)
					e.encUint(uint64(l))
					//fmt.Printf("%#v\n", e)
					e.reqLen += (l * eInfo.Length) //编译期间未规划的，现在规划
					//fmt.Printf("%#v\n", e)
					e.append(l*eInfo.Length + e.reserved) //检查到下一个变长类型空间是否足够，不够分配内存
					//fmt.Printf("%#v\n", e)
					//fmt.Println(l, eInfo.Length)
					for i := 0; i < l; i++ {
						//	fmt.Printf("encode bool %#v\n", e)
						eInfo.Engine(e, v.Index(i))
					}
					//fmt.Printf("encode bool %#v\n", e)
				} else {
					e.reqLen -= word1
				}
			}
			info.Head = baseLength[kind]
		} else {
			info.Engine = func(e *Encoder, v reflect.Value) {
				isNotNil := !v.IsNil()
				e.encBool(isNotNil)
				if isNotNil {
					l := v.Len()
					e.encUint(uint64(l))
					if l > 0 {
						e.reqLen += (l * eInfo.Length) //编译期间未规划的，现在规划
						reserved := e.reserved
						e.reserved = eInfo.Head
						lessOne := l - 1
						for i := 0; i < lessOne; i++ {
							eInfo.Engine(e, v.Index(i))
						}
						e.reserved = reserved
						eInfo.Engine(e, v.Index(lessOne))
					} else {
						e.append(e.reserved)
					}
				} else {
					e.reqLen -= word1
				}
			}
			info.Head = baseLength[kind] + eInfo.Head
		}
	case reflect.Map:
		kInfo, eInfo := buildInfo(rt.Key()), buildInfo(rt.Elem())
		info.IsStatic = false // 长度不固定，总是变长的
		//http://blog.csdn.net/hificamera/article/details/51701804
		slength := kInfo.Length + eInfo.Length
		if eInfo.IsStatic {
			info.Engine = func(e *Encoder, v reflect.Value) {
				isNotNil := !v.IsNil()
				e.encBool(isNotNil)
				if isNotNil {
					keys := v.MapKeys()
					l := len(keys)
					e.encUint(uint64(l))
					e.reqLen += l * slength
					e.append(l*slength + e.reserved)
					for _, key := range keys {
						kInfo.Engine(e, key)
						eInfo.Engine(e, v.MapIndex(key))
					}
				} else {
					e.reqLen -= word1
				}
			}
			info.Head = baseLength[kind]
		} else {
			info.Engine = func(e *Encoder, v reflect.Value) {
				isNotNil := !v.IsNil()
				e.encBool(isNotNil)
				if isNotNil {
					keys := v.MapKeys()
					l := len(keys)
					e.encUint(uint64(l))
					if l > 0 {
						e.reqLen += l * (slength)
						reserved := e.reserved
						e.reserved = slength
						lessOne := l - 1
						for i := 0; i < lessOne; i++ {
							kInfo.Engine(e, keys[i])
							eInfo.Engine(e, v.MapIndex(keys[i]))
						}
						e.reserved = reserved
						kInfo.Engine(e, keys[lessOne])
						eInfo.Engine(e, v.MapIndex(keys[lessOne]))
					} else {
						e.append(e.reserved)
					}
				} else {
					e.reqLen -= word1
				}
			}
			info.Head = baseLength[kind] + kInfo.Length + eInfo.Head
		}
	case reflect.Struct:
		nf := rt.NumField()
		if nf > 0 {
			engines, j, interval := make([]func(*Encoder, reflect.Value), nf), 0, [][2]int{{0, 0}}
			info.IsStatic = true // 默认为真，下面有假则为假
			var finfo *TypeInfo
			for i := 0; i < nf; i++ {
				finfo = buildInfo(rt.Field(i).Type)
				engines[i] = finfo.Engine
				info.Length += finfo.Length
				interval[j][0] += finfo.Head //Head is all
				if !finfo.IsStatic {
					info.IsStatic = finfo.IsStatic
					//if i != lessOne {
					interval = append(interval, [2]int{0, i}) //append 1 Length
					j++
					//}
					//interval[j][1] = i
				}
			}
			info.Head = interval[0][0]
			if info.IsStatic {
				info.Engine = func(e *Encoder, v reflect.Value) {
					for i := 0; i < nf; i++ {
						engines[i](e, v.Field(i))
					}
				}
			} else {
				info.Engine = func(e *Encoder, v reflect.Value) {
					//fmt.Printf("%#v\n", interval)
					reserved := e.reserved
					k, i := 1, 0 // interval 数组的第一个不用,第一个是head,由上一个负责越界检查
					for ; k < j; k++ {
						for e.reserved = interval[k][0]; i <= interval[k][1]; i++ {
							//fmt.Printf("%#v\n", e)
							engines[i](e, v.Field(i))
						}
					}
					for e.reserved = interval[k][0] + reserved; i < nf; i++ {
						//fmt.Printf("%#v\n", e)
						engines[i](e, v.Field(i))
					}
				}
			}
		} else {
			info = &TypeInfo{Length: 0, Head: 0, Engine: encignore, IsStatic: true}
		}
	default:
		info = baseInfo[kind]
	}
end:
	rt2Info[rt] = info
	return info
}
