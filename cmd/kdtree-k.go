//go:build ignore

package main

import (
	"fmt"
	"github.com/jphsd/datastruct"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"math"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	kdt := datastruct.NewKDTree(2)
	// Create 500 points in [-10,10]^2
	n := 500
	pts := make([][]float64, n)
	for i := 0; i < n; i++ {
		x := rand.Float64()*20 - 10
		y := rand.Float64()*20 - 10
		pts[i] = []float64{x, y}
		kdt.Insert(pts[i])
	}

	for i := 0; i < 10; i++ {
		pt := []float64{0, 0}
		points, _, _ := kdt.KNN(pt, (i+1)*10)
		dump(points, pts, pt, i)
	}
}

func dump(points, pts [][]float64, pt []float64, n int) {
	width, height := 1000, 1000
	img := image.NewRGBA(width, height, color.White)
	xfm := g2d.ScaleAndInset(float64(width), float64(height), 10, 10, true, datastruct.Bounds(pts...))

	points = xfm.Apply(points...)
	pts = xfm.Apply(pts...)
	pt = xfm.Apply(pt)[0]

	oShape, dShape := &g2d.Shape{}, &g2d.Shape{}
	for _, p := range pts {
		oShape.AddPaths(g2d.Circle(p, 2))
	}
	for _, p := range points {
		dShape.AddPaths(g2d.Circle(p, 2))
	}
	fmt.Printf("Found %d points, expected %d\n", len(points), (n+1)*10)

	last := points[len(points)-1]
	dx, dy := pt[0]-last[0], pt[1]-last[1]
	r := math.Sqrt(dx*dx + dy*dy)
	dShape.AddPaths(g2d.Circle(pt, r))

	g2d.DrawShape(img, oShape, g2d.Black)
	g2d.DrawShape(img, dShape, g2d.Red)

	image.SaveImage(img, fmt.Sprintf("kdtreek-%d", n))
}
