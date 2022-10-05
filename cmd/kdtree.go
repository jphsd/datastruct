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
	kdt := datastruct.NewKDTree(2)
	// Create 500 points
	n := 500
	pts := make([][]float64, n)
	for i := 0; i < n; i++ {
		x := rand.Float64()*20 - 10
		y := rand.Float64()*20 - 10
		pts[i] = []float64{x, y}
		kdt.Insert(pts[i])
	}

	// Run 10 KNN tests (vs total sort) with k = 20
	k := 20
	for i := 0; i < 10; i++ {
		x := rand.Float64()*20 - 10
		y := rand.Float64()*20 - 10
		pt := []float64{x, y}
		fmt.Printf("Round %d: pt %f, %f\n", i, x, y)
		nnpts, _, inds := kdt.KNN(pt, k)
		fnpts := nearest(pt, pts)
		for j := 0; j < k; j++ {
			if !equal(nnpts[j], fnpts[j]) {
				fmt.Print("ERROR ")
			}
			fmt.Printf("%d %d: %f, %f (%d) --- %f, %f\n", i, j, nnpts[j][0], nnpts[j][1], inds[j], fnpts[j][0], fnpts[j][1])
		}
	}
}

func nearest(pt []float64, pts [][]float64) [][]float64 {
	n := len(pts)
	pis := make([]datastruct.PriorityItem, n)
	for i, p := range pts {
		pis[i] = datastruct.PriorityItem{dist(pt, p), i}
	}
	pq := datastruct.NewPriorityList(pis...)
	res := make([][]float64, n)
	for i, pi := range *pq {
		res[i] = pts[pi.Id]
	}
	return res
}

func dist(a, b []float64) float64 {
	x := a[0] - b[0]
	y := a[1] - b[1]
	return x*x + y*y
}

func dumpNode(node *datastruct.KDNode) {
	for i := 0; i < node.Depth; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("%d: %f, %f\n", node.Id, node.Point[0], node.Point[1])
	if node.Left != nil {
		dumpNode(node.Left)
	} else {
		for i := 0; i < node.Depth+1; i++ {
			fmt.Print(" ")
		}
		fmt.Println("<nil>")
	}
	if node.Right != nil {
		dumpNode(node.Right)
	} else {
		for i := 0; i < node.Depth+1; i++ {
			fmt.Print(" ")
		}
		fmt.Println("<nil>")
	}
}

func equal(a, b []float64) bool {
	sa := fmt.Sprintf("%f, %f", a[0], a[1])
	sb := fmt.Sprintf("%f, %f", b[0], b[1])
	return sa == sb
}
