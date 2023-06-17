package datastruct

import (
	"fmt"
	"math"
)

// Bounded region row, like BBGrid but for 1D regions
type BRRow struct {
	Columns  int     // Columns across
	Min, Max float64 // Bounds of the space
	dx       float64
	row      [][]int //dead = -1
	n        int
}

func NewBRRow(c int, b []float64) *BRRow {
	validateBounds1(b)
	dx := b[1] - b[0]
	dx /= float64(c)
	row := make([][]int, c)
	brr := &BRRow{c, b[0], b[1], dx, row, 0}
	return brr
}

func (g *BRRow) Add(id int, br []float64) error {
	if id < 0 {
		return fmt.Errorf("id less than zero")
	}
	validateBounds1(br)
	if br[0] < g.Min || br[1] > g.Max {
		return fmt.Errorf("bounds don't overlap row")
	}
	min := g.pointToC(br[0])
	max := g.pointToC(br[1])
	// Add id to columns it overlaps with
	for c := min; c < max+1; c++ {
		g.row[c] = append(g.row[c], id)
	}
	g.n++
	return nil
}

// This is expensive as every column is checked
func (g *BRRow) Remove(id int) {
	found := false
	for c := 0; c < g.Columns; c++ {
		for i, v := range g.row[c] {
			if v == id {
				g.row[c][i] = -1
				found = true
				break
			}
		}
	}
	if found {
		g.n--
	}
	// Silently ignore nonexistent ids
}

func (g *BRRow) Location(p float64) (int, error) {
	if p < g.Min || p > g.Max {
		return -1, fmt.Errorf("point outside of row bounds")
	}
	c := g.pointToC(p)
	return c, nil
}

func (g *BRRow) Cell(c int) []int {
	if c < 0 || c >= g.Columns {
		return nil
	}
	ids := g.row[c]
	return ids
}

func (g *BRRow) Len() int {
	return g.n
}

func (g *BRRow) pointToC(p float64) int {
	// Bounds check already performed
	c := (int)(math.Floor((p - g.Min) / g.dx))
	if c < 0 {
		c = 0
	}
	if c >= g.Columns {
		c = g.Columns - 1
	}
	return c
}

func validateBounds1(bounds []float64) {
	if bounds[1] < bounds[0] {
		bounds[0], bounds[1] = bounds[1], bounds[0]
	}
}
