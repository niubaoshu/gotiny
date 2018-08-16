package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/niubaoshu/gotiny"
)

type A struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

func BenchmarkGotinyMarshal(b *testing.B) {
	b.StopTimer()
	data := generate()
	b.ReportAllocs()
	e := gotiny.NewEncoder(A{})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e.Encode(data[rand.Intn(len(data))])
	}
}

func generate() []*A {
	a := make([]*A, 0, 1000)
	for i := 0; i < 1000; i++ {
		a = append(a, &A{
			Name:     randString(16),
			BirthDay: time.Now(),
			Phone:    randString(10),
			Siblings: rand.Intn(5),
			Spouse:   rand.Intn(2) == 1,
			Money:    rand.Float64(),
		})
	}
	return a
}

func randString(l int) string {
	buf := make([]byte, l)
	for i := 0; i < (l+1)/2; i++ {
		buf[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", buf)[:l]
}
