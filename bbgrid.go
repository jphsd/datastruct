package datastruct

import (
	"fmt"
	"math"
)

type BBGrid struct {
	Rows, Columns int       // Grid cells down and across
	Min, Max      []float64 // Bounds of the space
	dx, dy        float64
	grid          [][][]int //dead = -1
	n             int
}

func NewBBGrid(r, c int, bb [][]float64) *BBGrid {
	validateBounds(bb)
	dx, dy := bb[1][0]-bb[0][0], bb[1][1]-bb[0][1]
	dx /= float64(c)
	dy /= float64(r)
	grid := make([][][]int, r)
	for i := 0; i < r; i++ {
		grid[i] = make([][]int, c)
	}
	bbg := &BBGrid{r, c, bb[0], bb[1], dx, dy, grid, 0}
	return bbg
}

func (g *BBGrid) Add(id int, bb [][]float64) error {
	if id < 0 {
		return fmt.Errorf("id less than zero")
	}
	validateBounds(bb)
	if bb[1][0] < g.Min[0] || bb[0][0] > g.Max[0] || bb[1][1] < g.Min[1] || bb[0][1] > g.Max[1] {
		return fmt.Errorf("bounds don't overlap grid")
	}
	rmin, cmin := g.pointToRC(bb[0])
	rmax, cmax := g.pointToRC(bb[1])
	// Add id to cells it overlaps with
	for r := rmin; r < rmax+1; r++ {
		for c := cmin; c < cmax+1; c++ {
			g.grid[r][c] = append(g.grid[r][c], id)
		}
	}
	g.n++
	return nil
}

// This is expensive as every cell is checked
func (g *BBGrid) Remove(id int) {
	found := false
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Columns; c++ {
			for i, v := range g.grid[r][c] {
				if v == id {
					g.grid[r][c][i] = -1
					found = true
					break
				}
			}
		}
	}
	if found {
		g.n--
	}
	// Silently ignore nonexistent ids
}

func (g *BBGrid) Location(p []float64) (int, int, error) {
	if p[0] < g.Min[0] || p[0] > g.Max[0] || p[1] < g.Min[1] || p[1] > g.Max[1] {
		return -1, -1, fmt.Errorf("point outside of grid bounds")
	}
	r, c := g.pointToRC(p)
	return r, c, nil
}

func (g *BBGrid) Cell(r, c int) []int {
	if r < 0 || r >= g.Rows || c < 0 || c >= g.Columns {
		return nil
	}
	ids := g.grid[r][c]
	return ids
}

func (g *BBGrid) Range(bb [][]float64) []int {
	validateBounds(bb)
	// Clamp search to grid min/max
	if bb[0][0] < g.Min[0] {
		bb[0][0] = g.Min[0]
	}
	if bb[0][1] < g.Min[1] {
		bb[0][1] = g.Min[1]
	}
	if bb[1][0] < g.Max[0] {
		bb[1][0] = g.Max[0]
	}
	if bb[1][1] < g.Max[1] {
		bb[1][1] = g.Max[1]
	}
	r0, c0 := g.pointToRC(bb[0])
	r1, c1 := g.pointToRC(bb[1])
	idm := make(map[int]bool)
	for r := r0; r <= r1; r++ {
		for c := c0; c <= c1; c++ {
			for _, id := range g.grid[r][c] {
				idm[id] = true
			}
		}
	}

	res := make([]int, 0, len(idm))
	for k, _ := range idm {
		res = append(res, k)
	}
	return res
}

func (g *BBGrid) Len() int {
	return g.n
}

func (g *BBGrid) pointToRC(p []float64) (int, int) {
	// Bounds check already performed
	c := (int)(math.Floor((p[0] - g.Min[0]) / g.dx))
	r := (int)(math.Floor((p[1] - g.Min[1]) / g.dy))
	if c < 0 {
		c = 0
	}
	if c >= g.Columns {
		c = g.Columns - 1
	}
	if r < 0 {
		r = 0
	}
	if r >= g.Rows {
		r = g.Rows - 1
	}
	return r, c
}
