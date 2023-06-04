//go:build ignore

package main

import (
	"fmt"
	"github.com/jphsd/datastruct"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	n := 20
	input := make([]int, n)
	for i := 0; i < n; i++ {
		input[i] = i
	}

	rand.Shuffle(n, func(i, j int) {
		input[i], input[j] = input[j], input[i]
	})

	fmt.Println("Input:")
	print(input)

	inv := datastruct.InvertSort(input)

	fmt.Println("Inverted:")
	print(inv)

	// Invert input
	output := make([]int, n)
	for i := 0; i < n; i++ {
		output[i] = input[inv[i]]
	}

	fmt.Println("Output:")
	print(output)
}

func print(vals []int) {
	for i, v := range vals {
		if i == 0 {
			fmt.Printf("%d", v)
		} else {
			fmt.Printf(", %d", v)
		}
	}
	fmt.Println("")
}
