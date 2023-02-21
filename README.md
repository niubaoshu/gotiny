
# gotiny

   [![Go Report Card](https://goreportcard.com/badge/github.com/raszia/gotiny)](https://goreportcard.com/report/github.com/raszia/gotiny) [![CodeCov](https://codecov.io/gh/raszia/gotiny/branch/master/graph/badge.svg)](https://codecov.io/gh/raszia/gotiny) [![GoDoc](https://godoc.org/github.com/raszia/gotiny?status.svg)](https://godoc.org/github.com/raszia/gotiny) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/raszia/gotiny/blob/master/LICENSE)

gotiny is an efficient Go serialization library. By pre-generating encoding machines and reducing the use of the reflect library, gotiny improves efficiency and is almost as fast as serialization libraries that generate code.

## examples

```go

 package main

 import (
 "fmt"
 "reflect"

 "github.com/raszia/gotiny"
 )

func main() {

    marshalUnmarshalExample()
    encodeDecodeExample()
    marshalUnmarshalCompressExample()
}

 // marshal src and unmarshel the returned data to dst
 // no compression
func marshalUnmarshalExample() {
    src1, src2 := "marshalUnmarshal", []byte(" Example!")
    dst1, dst2 := "", []byte{3, 4, 5}

    data := gotiny.Marshal(&src1, &src2)
    gotiny.Unmarshal(data, &dst1, &dst2)

    fmt.Println(dst1 + string(dst2)) // print "marshalUnmarshal Example!"
}

 // encode the data using encoder and decode the data using decoder
 // no compression
func encodeDecodeExample() {
    src1, src2 := "encodeDecode", []byte(" Example!")
    dst1, dst2 := "", []byte{3, 4, 5}

    enc := gotiny.NewEncoder(src1, src2)
    dec := gotiny.NewDecoder(dst1, dst2)

    dst1, dst2 = "", []byte{
    3, 3, 3, 3, 3,
    3, 3, 3, 3, 3,
    4, 5, 6, 7, 44,
    7, 5, 6, 4, 7}

    dec.DecodeValue(enc.EncodeValue(reflect.ValueOf(src1),
    reflect.ValueOf(src2)),
    reflect.ValueOf(&dst1).Elem(),
    reflect.ValueOf(&dst2).Elem())

    fmt.Println(dst1 + string(dst2)) // print "encodeDecode Example!"
}

 // marshal src and unmarshel the returned data to dst
 // with compression
func marshalUnmarshalCompressExample() {
    src1, src2 := "marshalUnmarshalCompress", []byte(" Example!")
    dst1, dst2 := "", []byte{3, 4, 5}

    data := gotiny.MarshalCompress(&src1, &src2)
    gotiny.UnmarshalCompress(data, &dst1, &dst2)

    fmt.Println(dst1 + string(dst2)) // print "marshalUnmarshalCompress Example!"
}
```

## Features

- High efficiency: gotiny is over three times as fast as gob, the serialization library that comes with Golang. It is on par with other serialization frameworks that generate code and is even faster than some of them.
- Zero memory allocation except for map types.
- Supports encoding all built-in types and custom types, except func and chan types.
- Encodes non-exported fields of struct types. Non-encoding fields can be set using Golang tags.
- Strict type conversion: only types that are exactly the same are correctly encoded and decoded.
- Encodes nil values with types.
- Can handle cyclic types but not cyclic values. It will stack overflow.
- Decodes all types that can be encoded, regardless of the original and target values.
- Encoded byte strings do not contain type information, resulting in very small byte arrays.
- Encoded and Decode with compression (optional).

## Cannot process cyclic values. Does not support circular references. TODO

 type a *a
 var b a
 b = &b

## install

```bash
go get -u github.com/raszia/gotiny
```

## Encoding Protocol

### Boolean type

bool type takes up one bit, with the true value encoded as 1 and the false value encoded as 0. When bool type is encountered for the first time, a byte is allocated to encode the value into the least significant bit. When encountered for the second time, it is encoded into the second least significant bit. The ninth time a bool value is encountered, another byte is allocated to encode the value into the least significant bit, and so on.

### Integer type

- uint8 and int8 types are encoded as the next byte of the string.
- uint16,uint32,uint64,uint,uintptr are encoded using [Varints](https://developers.google.com/protocol-buffers/docs/encoding#varints) Encoding method.
- int16,int32,int64,int are converted to unsigned numbers using ZigZag and then encoded using [Varints](https://developers.google.com/protocol-buffers/docs/encoding#varints) Encoding.

### Floating point type

float32 and float64 are encoded using the encoding method for floating point types in [gob](https://golang.org/pkg/encoding/gob/)Encoding method for floating-point types.

### Complex number

- The complex64 type is forced to be converted to a uint64 and encoded using uint64 encoding
- complex128 type encodes the real and imaginary parts as float64 types.

### String type

The string type first encodes the length of the string by casting it to uint64 type and then encoding it. After that, it encodes the byte array of the string as is.

### Pointer type

For the pointer type, it checks whether it is nil. If it is nil, it encodes a false value of bool type and then ends. If it is not nil, it encodes a true value of bool type, and then dereferences the pointer and encodes it based on the type of the dereferenced object.

### Array and Slice type

It first casts the length to uint64 and encodes it using uint64 encoding. After that, it encodes each element based on its type.

### Map type

Similar to the above, it first encodes the length and then encodes each key and its corresponding value. It does this for each key-value pair in the map.

### Struct type

It encodes all the members of the struct based on their types, whether they are exported or not. The struct is strictly reconstructed.

### Types that implement interfaces

- For types that implement the BinaryMarshaler/BinaryUnmarshaler interfaces in the encoding package or the GobEncoder/GobDecoder interfaces in the gob package, the encoding and decoding is done using their implementation methods.
- For types that implement the GoTinySerialize interface in the gotiny.GoTinySerialize package, the encoding and decoding is done using their implementation methods.

## benchmark

[benchmark](https://github.com/niubaoshu/go_serialization_benchmarks)

### License

MIT
