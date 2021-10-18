## <font color="#FF4500" >gotiny 只是个玩具，不建议使用。</font>


# gotiny   [![Build status][travis-img]][travis-url] [![License][license-img]][license-url] [![GoDoc][doc-img]][doc-url] [![Go Report Card](https://goreportcard.com/badge/github.com/niubaoshu/gotiny)](https://goreportcard.com/report/github.com/niubaoshu/gotiny)
gotiny是一个注重效率的go语言序列化库。gotiny通过预先生成编码机和减少使用reflect库等方式来提高效率，几乎和生成代码的序列化库有同样高的速度。
## hello word 
    package main
    import (
   	    "fmt"
   	    "github.com/niubaoshu/gotiny"
    )
    
    func main() {
   	    src1, src2 := "hello", []byte(" world!")
   	    ret1, ret2 := "", []byte{}
   	    gotiny.Unmarshal(gotiny.Marshal(&src1, &src2), &ret1, &ret2)
   	    fmt.Println(ret1 + string(ret2)) // print "hello world!"
    }

## 特性
- 效率非常的高，是golang自带序列化库gob的3倍以上,和一般的生成代码序列化框架处于同一水平，甚至高于它们。
- 除map类型外0内存申请。
- 支持编码所有的除func,chan类型外的所有golang内置类型和自定义类型。
- struct 类型会编码非导出字段,可通过golang tag 的方式设置不编码。
- 严格的类型转换。gotiny中只有类型完全相同的才会正确编码和解码。
- 编码带类型的nil值。
- 可以处理循环类型，不能编码循环值，会栈溢出。
- 所有可以编码的类型都会完全的解码，不论原值是什么和目标值是什么。
- 编码生成的字节串不包含类型信息，生成的字节数组非常小。
## 无法处理循环值 不支持循环引用  TODO 
	type a *a
	var b a
	b = &b

## install
```bash
$ go get -u github.com/niubaoshu/gotiny
```

## 编码协议
### 布尔类型
bool类型占用一位，真值编码为1，假值编码为0。当第一次遇到bool类型时会申请一个字节，将值编入最低位，第二次遇到时编入次低位，第九次遇到bool值时再申请一个字节编入最低位，以此类推。
### 整数类型
- uint8和int8 类型作为一个字节编入字符串的下一个字节。
- uint16,uint32,uint64,uint,uintptr 采用[Varints](https://developers.google.com/protocol-buffers/docs/encoding#varints)编码方式。
- int16,int32,int64,int 采用ZigZag转换成一个无符号数后采用[Varints](https://developers.google.com/protocol-buffers/docs/encoding#varints)编码方式。

### 浮点数类型
float32和float64采用[gob](https://golang.org/pkg/encoding/gob/)中对浮点类型的编码方式。
### 复数类型
- complex64类型会强转为一个uint64后采用uint64的编码方式。
- complex128类型分别将虚实部分作为float64类型编码。

### 字符串类型
字符串类型先将字符串长度强转为uint64类型编码，然后将字符串字节数组自身原样编码。
### 指针类型
指针类型判断是否为nil，如果是nil，编入一个bool类型的false值后结束，如果不为nil，编入一个bool类型true值，之后将指针解引用，按解引用后的类型编码。
### array和slice类型
先将长度强转为一个uint64后采用uint64的编码方式编入，然后将每一个元素安装自身的类型编码。
### map类型
同上，先编入长度，然后编入一个健，后面跟健对应的值，在编入一个健，接着是值，以此类推。
### struct类型
将结构体的所有成员按其类型编码，无论是否导出，非导出的字段也会编码。结构体会严格还原。
### 实现接口的类型
- 对于实现encoding包BinaryMarshaler/BinaryUnmarshaler 或 实现 gob包GobEncoder/GobDecoder 接口的类型会用实现的方法编码。
- 对于实现了gotiny.GoTinySerialize包的类型将采用实现的方法编码和解码

## benchmark
[benchmark](https://github.com/niubaoshu/go_serialization_benchmarks)


### License
MIT

[travis-img]: https://travis-ci.org/niubaoshu/gotiny.svg?branch=master
[travis-url]: https://travis-ci.org/niubaoshu/gotiny
[license-img]: http://img.shields.io/badge/license-MIT-green.svg?style=flat-square
[license-url]: http://opensource.org/licenses/MIT
[doc-img]: http://img.shields.io/badge/GoDoc-reference-blue.svg?style=flat-square
[doc-url]: https://godoc.org/github.com/niubaoshu/gotiny
