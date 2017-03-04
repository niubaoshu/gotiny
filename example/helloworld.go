package main

import (
	"fmt"

	"github.com/niubaoshu/gotiny"
)

func main() {
	src := "hello world!"
	ret := ""
	gotiny.Decodes(gotiny.Encodes(&src), &ret)
	fmt.Println(ret) // print "hello world!"
}
