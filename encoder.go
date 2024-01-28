package gotiny

import (
	"reflect"
	"unsafe"
)

type Encoder struct {
	buf     []byte   // the destination array for encoding
	off     int      // the current offset in buf
	boolPos int      // the next bool to be set in buf, i.e. buf[boolPos]
	boolBit byte     // the next bit to be set in buf[boolPos]
	engines []encEng // a collection of encoders
	length  int      // number of encoders
}

// Marshal marshal value
func Marshal(is ...any) []byte {
	return NewEncoderWithPtr(is...).Encode(is...)
}

// MarshalEncrypt marshal and encrypt value
func MarshalEncrypt(aesConfig *aesConfigStruct, is ...any) []byte {
	return NewEncoderWithPtr(is...).EncodeEncrypt(aesConfig, is...)
}

// MarshalCompress marshal and compress value
func MarshalCompress(is ...any) []byte {
	return NewEncoderWithPtr(is...).EncodeCompress(is...)
}

// MarshalCompressEncrypt marshal and compress and encrypt value
func MarshalCompressEncrypt(aesConfig *aesConfigStruct, is ...any) []byte {
	return NewEncoderWithPtr(is...).EncodeCompressEncrypt(aesConfig, is...)
}

// Create an encoder that points to the encoding of the given type.
func NewEncoderWithPtr(ps ...any) *Encoder {
	l := len(ps)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(ps[i])
		if rt.Kind() != reflect.Ptr {
			panic("must a pointer type!")
		}
		engines[i] = getEncEngine(rt.Elem())
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}
}

// Create an encoder of type Encoder.
func NewEncoder(is ...any) *Encoder {
	l := len(is)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		engines[i] = getEncEngine(reflect.TypeOf(is[i]))
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}
}

func NewEncoderWithType(ts ...reflect.Type) *Encoder {
	l := len(ts)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		engines[i] = getEncEngine(ts[i])
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}
}

// Encode encode value The input is a pointer to the value to be encoded.
func (e *Encoder) Encode(is ...any) []byte {
	engines := e.engines
	for i := 0; i < len(engines) && i < len(is); i++ {
		engines[i](e, (*[2]unsafe.Pointer)(unsafe.Pointer(&is[i]))[1])
	}
	return e.reset()
}

// EncodeEncrypt encrypt and encode value The input is a pointer to the value to be encoded.
func (e *Encoder) EncodeEncrypt(aesConfig *aesConfigStruct, is ...any) []byte {
	return aesConfig.Encrypt(e.Encode(is...))
}

// EncodeCompress compress and encode value The input is a pointer to the value to be encoded.
func (e *Encoder) EncodeCompress(is ...any) []byte {

	b := Gzip.Getbuffer()
	defer Gzip.Putbuffer(b)
	Gziper(b, e.Encode(is...))
	return b.Bytes()
}

// EncodeCompressEncrypt compress and ecrypt and encode value The input is a pointer to the value to be encoded.
func (e *Encoder) EncodeCompressEncrypt(aesConfig *aesConfigStruct, is ...any) []byte {
	return aesConfig.Encrypt(e.EncodeCompress(is...))
}

// an unsafe.Pointer to the value to be encoded.
func (e *Encoder) EncodePtr(ps ...unsafe.Pointer) []byte {
	engines := e.engines
	for i := 0; i < len(engines) && i < len(ps); i++ {
		engines[i](e, ps[i])
	}
	return e.reset()
}

// vs holds the value to be encoded.
func (e *Encoder) EncodeValue(vs ...reflect.Value) []byte {
	engines := e.engines
	for i := 0; i < len(engines) && i < len(vs); i++ {
		engines[i](e, getUnsafePointer(&vs[i]))
	}
	return e.reset()
}

// The encoded data will be appended to buf.
func (e *Encoder) AppendTo(buf []byte) {
	e.off = len(buf)
	e.buf = buf
}

func (e *Encoder) reset() []byte {
	buf := e.buf
	e.buf = buf[:e.off]
	e.boolBit = 0
	e.boolPos = 0
	return buf
}
