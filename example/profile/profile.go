package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"runtime/pprof"
	"time"

	"github.com/niubaoshu/gotiny"
	"github.com/niubaoshu/goutils"
)

type str struct {
	A map[int]map[int]string
	B []bool
	c int
}

type ET0 struct {
	s str
	F map[int]map[int]string
}

var (
	//_   = rand.Intn(1)
	now = time.Now()
	a   = "234234"
	i   = map[int]map[int]string{
		1: map[int]string{
			1: a,
		},
	}
	strs = `抵制西方的司法独立，有什么错？有人说马克思主义还是西方的，有本事别用啊。这都是犯了形而上学的错误，任何理论、思想都必须和中国国情相结合，和当前实际相结合。全部照搬照抄的教条主义王明已经试过一次，结果怎么样？歪解周强讲话，不是蠢就是别有用心，蠢的可以教育，别有用心就该打倒`
	st   = str{A: i, B: []bool{true, false, false, false, false, true, true, false, true, false, true}, c: 234234}
	//st     = str{c: 234234}
	et0      = ET0{s: st, F: i}
	stp      = &st
	stpp     = &stp
	nilslice []byte
	slice    = []byte{1, 2, 3}
	mapt     = map[int]int{0: 1, 1: 2, 2: 3, 3: 4}
	nilmap   map[int][]byte
	nilptr   *map[int][]string
	inta          = 2
	ptrint   *int = &inta
	nilint   *int
	vs       = []interface{}{
		ptrint,
		strs,
		`习近平离京对瑞士联邦进行国事访问
		出席世界经济论坛2017年年会并访问在瑞士的国际组织
		新华社北京1月15日电1月15日上午，国家主席习近平乘专机离开北京，应以洛伊特哈德为主席的瑞士联邦委员会邀请，对瑞士进行国事访问；应世界经济论坛创始人兼执行主席施瓦布邀请，出席在达沃斯举行的世界经济论坛2017年年会；应联合国秘书长古特雷斯、世界卫生组织总干事陈冯富珍、国际奥林匹克委员会主席巴赫邀请，访问联合国日内瓦总部、世界卫生组织、国际奥林匹克委员会。
		陪同习近平出访的有：习近平主席夫人彭丽媛，中共中央政治局委员、中央政策研究室主任王沪宁，中共中央政治局委员、中央书记处书记、中央办公厅主任栗战书，国务委员杨洁篪等。返回腾讯网首页>>`,
		true,
		false,
		int(123456),
		int8(123),
		int16(-12345),
		int32(123456),
		int64(-1234567),
		int64(1<<63 - 1),
		int64(rand.Int63()),
		uint(123),
		uint8(123),
		uint16(12345),
		uint32(123456),
		uint64(1234567),
		uint64(1<<64 - 1),
		uint64(rand.Uint32() * rand.Uint32()),
		uintptr(12345678),
		float32(1.2345),
		float64(1.2345678),
		complex64(1.2345 + 2.3456i),
		complex128(1.2345678 + 2.3456789i),
		string("hello,日本国"),
		string("9b899bec35bc6bb8"),
		inta,
		[][][][3][][3]int{{{{{{2, 3}}}}}},
		map[int]map[int]map[int]map[int]map[int]map[int]map[int]map[int]int{1: {1: {1: {1: {1: {1: {1: {1: 2}}}}}}}},
		map[int]map[int]int{1: {2: 3}},
		[][]bool{},
		[]byte("hello，中国人"),
		[][]byte{[]byte("hello"), []byte("world")},
		[4]string{"2324", "23423", "捉鬼", "《：LSESERsef色粉色问问我二维牛"},
		map[int]string{1: "h", 2: "h", 3: "nihao"},
		map[string]map[int]string{"werwer": {1: "呼呼喊喊"}, "汉字": {2: "世界"}},
		a,
		i,
		&i,
		st,
		stp,
		stpp,
		struct{}{},
		[][][]struct{}{},
		struct {
			a, C int
		}{1, 2},
		et0,
		[100]int{},
		now,
		ptrint,
		nilmap,
		nilslice,
		nilptr,
		nilint,
		slice,
		mapt,
	}
	e = gotiny.NewEncoder(vs...)
	d = gotiny.NewDecoder(vs...)

	spvals = make([]interface{}, len(vs))
	rpvals = make([]interface{}, len(vs))
	c      = goutils.NewComparer()

	buf = make([]byte, 0, 2048)
)

func init() {

	for i := 0; i < len(vs); i++ {
		typ := reflect.TypeOf(vs[i])
		temp := reflect.New(typ)
		temp.Elem().Set(reflect.ValueOf(vs[i]))
		spvals[i] = temp.Interface()

		if i == len(vs)-2 {
			a := make([]byte, 15)
			rpvals[i] = &a
		} else if i == len(vs)-1 {
			//a := map[int]int{111: 233, 6: 7}
			a := map[int]int{}
			rpvals[i] = &a
		} else {
			rpvals[i] = reflect.New(typ).Interface()
		}
	}
	e.AppendTo(buf[:0])
}

func main() {
	f, err := os.Create("cpuprofile.pprof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for i := 0; i < 1000; i++ {
		for i := 0; i < 1000; i++ {
			e.AppendTo(buf[:0])
			d.Decodes(e.Encodes(spvals...), rpvals...)
			for i, result := range rpvals {
				r := reflect.ValueOf(result).Elem().Interface()
				if Assert(vs[i], r) != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func Assert(x, y interface{}) error {
	if !c.DeepEqual(x, y) {
		return fmt.Errorf("\n exp type =  %T; value = %#v;\n got type = %T; value = %#v ", x, x, y, y)
	}
	return nil
}
