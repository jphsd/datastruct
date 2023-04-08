//go:build ignore

package main

import (
	"fmt"
	"github.com/jphsd/datastruct"
	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"math/rand"
)

func main() {
	width, height := 1000, 1000

	n := 0

	for {
		img := image.NewRGBA(width, height, color.White)
		r, c := 10, 10
		gbb := [][]float64{{100, 100}, {900, 900}}

		// Create grid
		grid := datastruct.NewBBGrid(r, c, gbb)

		// Create region
		p1 := []float64{rand.Float64() * float64(width), rand.Float64() * float64(height)}
		p2 := []float64{rand.Float64() * float64(width), rand.Float64() * float64(height)}
		bb := [][]float64{p1, p2}

		// Add region to grid
		grid.Add(1, bb)

		// Draw grid in black
		w, h := gbb[1][0]-gbb[0][0], gbb[1][1]-gbb[0][1]
		dx := w / float64(c)
		dy := h / float64(r)
		y := gbb[0][1]
		gshape := &g2d.Shape{}
		for i := 0; i < r+1; i++ {
			x := gbb[0][0]
			p1 := []float64{x, y}
			p2 := []float64{x + w, y}
			gshape.AddPaths(g2d.Line(p1, p2))
			y += dy
		}
		x := gbb[0][0]
		for i := 0; i < c+1; i++ {
			y := gbb[0][1]
			p1 := []float64{x, y}
			p2 := []float64{x, y + h}
			gshape.AddPaths(g2d.Line(p1, p2))
			x += dx
		}
		g2d.DrawShape(img, gshape, g2d.BlackPen)

		// Fill cells with #ids > 0 in blue
		y = gbb[0][1]
		cshape := &g2d.Shape{}
		for i := 0; i < r; i++ {
			x = gbb[0][0]
			for j := 0; j < c; j++ {
				ids := grid.Cell(i, j)
				if len(ids) > 0 {
					x1, y1 := x+dx, y+dy
					cshape.AddPaths(g2d.Polygon([]float64{x, y}, []float64{x1, y}, []float64{x1, y1}, []float64{x, y1}))
				}
				x += dx
			}
			y += dy
		}
		g2d.RenderColoredShape(img, cshape, color.Blue)

		// Draw region in green
		rshape := g2d.NewShape(g2d.Polygon(boundsToPoints(bb)...))
		g2d.RenderColoredShape(img, rshape, color.Green)

		// Run 1000 points
		hshape, mshape := &g2d.Shape{}, &g2d.Shape{}
		for i := 0; i < 1000; i++ {
			p := []float64{rand.Float64()*w + gbb[0][0], rand.Float64()*h + gbb[0][1]}
			rr, cc, err := grid.Location(p)
			if err != nil {
				fmt.Printf("%v out of bounds\n", p)
				continue
			}
			ids := grid.Cell(rr, cc)
			if len(ids) == 0 {
				mshape.AddPaths(g2d.Circle(p, 2))
			} else {
				hshape.AddPaths(g2d.Circle(p, 2))
			}
		}
		g2d.RenderColoredShape(img, mshape, color.Cyan)
		g2d.RenderColoredShape(img, hshape, color.Red)

		image.SaveImage(img, fmt.Sprintf("bbgrid-%d", n))
		n++
	}
}

func boundsToPoints(bb [][]float64) [][]float64 {
	x1, x2 := bb[0][0], bb[1][0]
	y1, y2 := bb[0][1], bb[1][1]
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	return [][]float64{{x1, y1}, {x2, y1}, {x2, y2}, {x1, y2}}
}
