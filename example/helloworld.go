package main

import (
	"fmt"

	"reflect"

	"github.com/niubaoshu/gotiny"
)

func main() {
	src1, src2 := "hello", []byte(" world!")
	ret1, ret2 := "", []byte{}
	gotiny.Decodes(gotiny.Encodes(&src1, &src2), &ret1, &ret2)
	fmt.Println(ret1 + string(ret2)) // print "hello world!"

	enc := gotiny.NewEncoder(src1, src2)
	dec := gotiny.NewDecoder(ret1, ret2)

	enc.EncodeValues(reflect.ValueOf(src1), reflect.ValueOf(src2))
	dec.ResetWith(enc.Bytes())
	ret1, ret2 = "", []byte{}
	dec.DecodeValues(reflect.ValueOf(&ret1).Elem(), reflect.ValueOf(&ret2).Elem())
	fmt.Println(ret1 + string(ret2)) // print "hello world!"
}
