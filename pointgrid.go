package datastruct

import (
	"fmt"
	"math"
)

// PointGrid provides a simple way of finding points that are close to each other.
type PointGrid struct {
	Rows, Columns int       // Grid cells down and across
	Min, Max      []float64 // Bounds of the space
	Wrap          bool      // Whether points should be wrapped in Min, Max
	w, h, dx, dy  float64
	grid          [][][]int
	n             int
}

// NewPointGrid creates a new PointGrid with the supplied attributes.
func NewPointGrid(rows, columns int, bounds [][]float64, wrap bool) *PointGrid {
	validateBounds(bounds)
	if rows < 1 {
		rows = 1
	}
	if columns < 1 {
		columns = 1
	}
	grid := make([][][]int, rows)
	for i := 0; i < rows; i++ {
		grid[i] = make([][]int, columns)
	}
	w, h := bounds[1][0]-bounds[0][0], bounds[1][1]-bounds[0][1]
	dx, dy := w/float64(columns), h/float64(rows)
	return &PointGrid{rows, columns, bounds[0], bounds[1], wrap, w, h, dx, dy, grid, 0}
}

// Location returns the grid location of a point (if it is within range).
func (g *PointGrid) Location(p []float64) (int, int, error) {
	x, y := p[0], p[1]
	if g.Wrap {
		for x < g.Min[0] {
			x += g.w
		}
		for x > g.Max[0] {
			x -= g.w
		}
		for y < g.Min[1] {
			y += g.h
		}
		for y > g.Max[1] {
			y -= g.h
		}
	} else if x < g.Min[0] || x > g.Max[0] || y < g.Min[1] || y > g.Max[1] {
		return 0, 0, fmt.Errorf("point out of range: %f,%f {%f,%f, %f,%f}", p[0], p[1], g.Min[0], g.Min[1], g.Max[0], g.Max[1])
	}
	c := int(math.Floor((x - g.Min[0]) / g.dx))
	r := int(math.Floor((y - g.Min[1]) / g.dy))
	return r, c, nil
}

// Add adds a point to the appropriate cell and returns its location and index.
func (g *PointGrid) Add(p []float64) (int, int, int, error) {
	r, c, err := g.Location(p)
	if err != nil {
		return 0, 0, -1, err
	}

	n := g.n
	g.grid[r][c] = append(g.grid[r][c], n)
	g.n++
	return r, c, n, nil
}

// Cell returns all the points in it.
func (g *PointGrid) Cell(row, column int) []int {
	if g.Wrap {
		if row < 0 {
			row += g.Rows
		} else if row >= g.Rows {
			row -= g.Rows
		}
		if column < 0 {
			column += g.Columns
		} else if column >= g.Columns {
			column -= g.Columns
		}
	} else if row < 0 || row >= g.Rows || column < 0 || column >= g.Columns {
		return []int{}
	}
	return g.grid[row][column]
}

// AdjacentCells returns the points of the cell and the cells adjacent to it.
func (g *PointGrid) AdjacentCells(row, column int) []int {
	res := make([]int, 0)
	for r := -1; r < 2; r++ {
		for c := -1; c < 2; c++ {
			res = append(res, g.Cell(row+r, column+c)...)
		}
	}
	return res
}

// Len returns the number of points stored in the grid.
func (g *PointGrid) Len() int {
	return g.n
}

func validateBounds(bounds [][]float64) {
	if bounds[0][0] > bounds[1][0] {
		bounds[0][0], bounds[1][0] = bounds[1][0], bounds[0][0]
	}
	if bounds[0][1] > bounds[1][1] {
		bounds[0][1], bounds[1][1] = bounds[1][1], bounds[0][1]
	}
}
