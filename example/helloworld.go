package main

import (
	"fmt"

	"github.com/niubaoshu/gotiny"
)

func main() {
	src1, src2 := "hello", []byte(" world!")
	ret1, ret2 := "", []byte{}
	gotiny.Decodes(gotiny.Encodes(&src1, &src2), &ret1, &ret2)
	fmt.Println(ret1 + string(ret2)) // print "hello world!"
}
