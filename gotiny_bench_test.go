package gotiny

import (
	"math/rand"
	"testing"
)

func BenchmarkDecodeUint64(b *testing.B) {
	b.StopTimer()
	var ints = make([][]byte, 10000)
	for i := 0; i < len(ints); i++ {
		a := rand.Uint64()
		ints[i] = Marshal(&a)
	}
	d := Decoder{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		d.buf = ints[rand.Intn(10000)]
		d.index = 0
		d.decUint64()
	}
}

func BenchmarkEncodeUint64(b *testing.B) {
	b.StopTimer()
	e := Encoder{buf: make([]byte, 0, 600000000)}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e.encUint64(rand.Uint64())
	}
}
