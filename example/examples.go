package main

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/raszia/gotiny"
)

func main() {

	marshalUnmarshalExample()
	encodeDecodeExample()
	marshalUnmarshalCompressExample()
	marshalUnmarshalEncryptExample()
	marshalUnmarshalCompressEncryptExample()
}

// marshal src and unmarshel the returned data to dst
// no compression
func marshalUnmarshalExample() {
	src1, src2 := "marshalUnmarshal", []byte(" Example!")
	dst1, dst2 := "", []byte{3, 4, 5}

	data := gotiny.Marshal(&src1, &src2)
	gotiny.Unmarshal(data, &dst1, &dst2)

	fmt.Println(dst1 + string(dst2)) // print "hello world!"
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

	fmt.Println(dst1 + string(dst2)) // print "hello world!"
}

// marshal src and unmarshel the returned data to dst
// with compression
func marshalUnmarshalCompressExample() {
	src1, src2 := "marshalUnmarshalCompress", []byte(" Example!")
	dst1, dst2 := "", []byte{3, 4, 5}

	data := gotiny.MarshalCompress(&src1, &src2)
	gotiny.UnmarshalCompress(data, &dst1, &dst2)

	fmt.Println(dst1 + string(dst2)) // print "hello world!"
}

// marshal src and unmarshel the returned data to dst
// with compression and encryption
func marshalUnmarshalCompressEncryptExample() {
	src1, src2 := "marshalUnmarshalCompressEncrypt", []byte(" Example!")
	dst1, dst2 := "", []byte{3, 4, 5}

	var str = "0123456789abcdef0123456789abcdef" // 32-byte hex string
	var key [32]byte

	// Convert the string to a byte slice
	bSlice, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}

	// Copy the byte slice into the array
	copy(key[:], bSlice)
	aesConfig := gotiny.NewAES256config(key)
	data := gotiny.MarshalCompressEncrypt(aesConfig, &src1, &src2)
	gotiny.UnmarshalCompressEncrypt(aesConfig, data, &dst1, &dst2)

	fmt.Println(dst1 + string(dst2)) // print "hello world!"
}

// marshal src and unmarshel the returned data to dst
// with compression and encryption
func marshalUnmarshalEncryptExample() {
	src1, src2 := "marshalUnmarshalEncrypt", []byte(" Example!")
	dst1, dst2 := "", []byte{3, 4, 5}

	var str = "0123456789abcdef0123456789abcdef" // 32-byte hex string
	var key [32]byte

	// Convert the string to a byte slice
	bSlice, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}

	// Copy the byte slice into the array
	copy(key[:], bSlice)
	aesConfig := gotiny.NewAES256config(key)
	data := gotiny.MarshalEncrypt(aesConfig, &src1, &src2)
	gotiny.UnmarshalEncrypt(aesConfig, data, &dst1, &dst2)

	fmt.Println(dst1 + string(dst2)) // print "hello world!"
}
