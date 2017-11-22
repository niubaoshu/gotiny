package gotiny

import (
	"encoding"
	"encoding/gob"
	"reflect"
	"strconv"
	"unsafe"
)

const ptr1Size = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

type refVal struct {
	typ  unsafe.Pointer
	ptr  unsafe.Pointer
	flag flag
}

type flag uintptr

//go:linkname flagIndir reflect.flagIndir
const flagIndir flag = 1 << 7

type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}

type stringHeader struct {
	data unsafe.Pointer
	len  int
}

func floatToUint(v float64) uint64 {
	return reverseByte(*(*uint64)(unsafe.Pointer(&v)))
}

func uintToFloat(u uint64) float64 {
	u = reverseByte(u)
	return *((*float64)(unsafe.Pointer(&u)))
}

func reverseByte(u uint64) uint64 {
	u = (u << 32) | (u >> 32)
	u = ((u << 16) & 0xFFFF0000FFFF0000) | ((u >> 16) & 0xFFFF0000FFFF)
	u = ((u << 8) & 0xFF00FF00FF00FF00) | ((u >> 8) & 0xFF00FF00FF00FF)
	return u
}

// int -5 -4 -3 -2 -1 0 1 2 3 4 5  6
// uint 9  7  5  3  1 0 2 4 6 8 10 12
func intToUint(v int64) uint64 {
	return uint64((v << 1) ^ (v >> 63))
}

// uint 9  7  5  3  1 0 2 4 6 8 10 12
// int -5 -4 -3 -2 -1 0 1 2 3 4 5  6
func uintToInt(u uint64) int64 {
	v := int64(u)
	return (-(v & 1)) ^ (v>>1)&0x7FFFFFFFFFFFFFFF
}

type gobInter interface {
	gob.GobEncoder
	gob.GobDecoder
}

type binInter interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

var (
	gobType  = reflect.TypeOf((*gobInter)(nil)).Elem()
	binType  = reflect.TypeOf((*binInter)(nil)).Elem()
	tinyType = reflect.TypeOf((*GoTinySerializer)(nil)).Elem()
)

func isNil(p unsafe.Pointer) bool {
	return *(*unsafe.Pointer)(p) == nil
}

// 只应该由指针来实现该接口
type GoTinySerializer interface {
	// 编码方法，将对象的序列化结果append到入参数并返回，方法不应该修改入参数值原有的值
	GotinyEncode([]byte) []byte
	// 解码方法，将入参解码到对象里并返回使用的长度。方法从入参的第0个字节开始使用，并且不应该修改入参中的任何数据
	GotinyDecode([]byte) int
}

func GetName(obj interface{}) string {
	return getNameByType(reflect.TypeOf(obj))
}
func getNameByType(rt reflect.Type) string {
	return string(getName([]byte(nil), rt))
}

func getName(prefix []byte, rt reflect.Type) []byte {
	if rt == nil || rt.Kind() == reflect.Invalid {
		return []byte("<nil>")
	}
	if rt.Name() == "" { //未命名的，组合类型
		switch rt.Kind() {
		case reflect.Ptr:
			return getName(append(prefix, '*'), rt.Elem())
		case reflect.Array:
			return getName(append(prefix, "["+strconv.Itoa(rt.Len())+"]"...), rt.Elem())
		case reflect.Slice:
			return getName(append(prefix, '[', ']'), rt.Elem())
		case reflect.Struct:
			prefix = append(prefix, "struct {"...)
			nf := rt.NumField()
			if nf > 0 {
				prefix = append(prefix, ' ')
			}
			for i := 0; i < nf; i++ {
				field := rt.Field(i)
				if field.Anonymous {
					prefix = getName(prefix, field.Type)
				} else {
					prefix = append(prefix, field.Name+" "...)
					prefix = getName(prefix, field.Type)
				}
				if i != nf-1 {
					prefix = append(prefix, ';', ' ')
				} else {
					prefix = append(prefix, ' ')
				}
			}
			return append(prefix, '}')
		case reflect.Map:
			prefix = append(prefix, "map["...)
			prefix = append(getName(prefix, rt.Key()), ']')
			return getName(prefix, rt.Elem())
		case reflect.Interface:
			prefix = append(prefix, "interface {"...)
			nm := rt.NumMethod()
			if nm > 0 {
				prefix = append(prefix, ' ')
			}
			for i := 0; i < nm; i++ {
				method := rt.Method(i)
				fn := getName([]byte(nil), method.Type)
				prefix = append(prefix, method.Name+string(fn[4:])...)
				if i != nm-1 {
					prefix = append(prefix, ';', ' ')
				} else {
					prefix = append(prefix, ' ')
				}
			}
			prefix = append(prefix, '}')
		case reflect.Func:
			prefix = append(prefix, "func("...)
			for i := 0; i < rt.NumIn(); i++ {
				prefix = getName(prefix, rt.In(i))
				if i != rt.NumIn()-1 {
					prefix = append(prefix, ',', ' ')
				}
			}
			prefix = append(prefix, ')')
			no := rt.NumOut()
			if no > 1 {
				prefix = append(prefix, ' ', '(')
			}
			for i := 0; i < no; i++ {
				prefix = getName(prefix, rt.Out(i))
				if i != no-1 {
					prefix = append(prefix, ',', ' ')
				}
			}
			if no > 1 {
				prefix = append(prefix, ')')
			}
		}
	}
	if rt.PkgPath() == "" {
		prefix = append(prefix, rt.Name()...)
	} else {
		prefix = append(prefix, rt.PkgPath()+"."+rt.Name()...)
	}
	return prefix
}

func getNameOfType(rt reflect.Type) string {
	if name, has := type2name[rt]; has {
		return name
	} else {
		return registerType(rt)
	}
}

func Register(i interface{}) string {
	return registerType(reflect.TypeOf(i))
}

func registerType(rt reflect.Type) string {
	name := getNameByType(rt)
	RegisterName(name, rt)
	return name
}

func RegisterName(name string, rt reflect.Type) {
	if name == "" {
		panic("attempt to register empty name")
	}

	if _, has := type2name[rt]; has {
		panic("gotiny: registering duplicate types for " + getNameByType(rt))
	}

	if _, has := name2type[name]; has {
		panic("gotiny: registering name" + name + " is exist")
	}
	name2type[name] = rt
	type2name[rt] = name
}
