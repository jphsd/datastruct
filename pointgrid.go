package datastruct

import (
	"fmt"
	"math"
)

// PointGrid provides a simple way of finding points that are close to eachother.
type PointGrid struct {
	Rows, Columns int       // Grid cells down and across
	Min, Max      []float64 // Bounds of the space
	Wrap          bool      // Whether points should be wrapped in Min, Max
	w, h, dx, dy  float64
	grid          [][][][]float64
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
	grid := make([][][][]float64, rows)
	for i := 0; i < rows; i++ {
		grid[i] = make([][][]float64, columns)
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

// Add adds a point to the appropriate cell and returns its location.
func (g *PointGrid) Add(p []float64) (int, int, error) {
	r, c, err := g.Location(p)
	if err != nil {
		return 0, 0, err
	}

	g.grid[r][c] = append(g.grid[r][c], p)
	g.n++
	return r, c, nil
}

// Cell returns all the points in it.
func (g *PointGrid) Cell(row, column int) [][]float64 {
	roffs, coffs := 0.0, 0.0
	if g.Wrap {
		if row < 0 {
			row += g.Rows
			roffs = -g.h
		} else if row >= g.Rows {
			row -= g.Rows
			roffs = g.h
		}
		if column < 0 {
			column += g.Columns
			coffs = -g.w
		} else if column >= g.Columns {
			column -= g.Columns
			coffs = g.w
		}
		pts := g.grid[row][column]
		n := len(pts)
		res := make([][]float64, n)
		for i := 0; i < n; i++ {
			// Transform point to wrapped point so distance calcs work
			res[i] = []float64{pts[i][0] + coffs, pts[i][1] + roffs}
		}
		return res
	}
	if row < 0 || row >= g.Rows || column < 0 || column >= g.Columns {
		return [][]float64{}
	}
	return g.grid[row][column]
}

// AdjacentCells returns the points of the cell and the cells adjacent to it.
func (g *PointGrid) AdjacentCells(row, column int) [][]float64 {
	res := make([][]float64, 0)
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

// NearestPoint looks in the cell containing the point and adjacent cells for the closest point.
// Returns the point (if any) and its distance or an error. A very poor implementation of nearest neighbor.
func (g *PointGrid) NearestPoint(p []float64) ([]float64, float64, error) {
	r, c, err := g.Location(p)
	if err != nil {
		return nil, -1, err
	}

	points := g.AdjacentCells(r, c)
	np := len(points)
	if np == 0 {
		return nil, -1, nil
	}

	nearest := points[0]
	px, py := p[0], p[1]
	dx, dy := points[0][0]-px, points[0][1]-py
	nd := dx*dx + dy*dy
	for i := 1; i < np; i++ {
		dx, dy = points[i][0]-px, points[i][1]-py
		d := dx*dx + dy*dy
		if d < nd {
			nearest = points[i]
			nd = d
		}
	}
	return nearest, math.Sqrt(nd), nil
}

func validateBounds(bounds [][]float64) {
	if bounds[0][0] > bounds[1][0] {
		bounds[0][0], bounds[1][0] = bounds[1][0], bounds[0][0]
	}
	if bounds[0][1] > bounds[1][1] {
		bounds[0][1], bounds[1][1] = bounds[1][1], bounds[0][1]
	}
}
