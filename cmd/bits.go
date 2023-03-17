//go:build ignore

package main

import (
	"fmt"
	"github.com/jphsd/datastruct"
	"math/rand"
)

func main() {
	n := 126
	bs := make([]bool, n)
	for i := 0; i < n; i++ {
		if rand.Intn(2) == 0 {
			bs[i] = true
		}
	}

	bits := datastruct.BitsFromSlice(bs)

	bs1 := bits.Slice()
	fmt.Printf("bs %d, bs1 %d\n", len(bs), len(bs1))

	for i := 0; i < n; i++ {
		if bs[i] != bs1[i] {
			fmt.Printf("%d: %v vs %v\n", bs[i], bs1[i])
		}
	}

	n = len(bs1)
	for i := 0; i < 100; i++ {
		ind := rand.Intn(n)
		bits.Set(ind)
		if !bits.Get(ind) {
			fmt.Printf("Attempt to set bit %d failed\n", ind)
		}
	}
	for i := 0; i < 100; i++ {
		ind := rand.Intn(n)
		bits.Clear(ind)
		if bits.Get(ind) {
			fmt.Printf("Attempt to clear bit %d failed\n", ind)
		}
	}
}
