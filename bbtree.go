package datastruct

import (
	"fmt"
	"math"
)

func Bounds(points ...[]float64) [][]float64 {
	min := []float64{math.MaxFloat64, math.MaxFloat64}
	max := []float64{-math.MaxFloat64, -math.MaxFloat64}
	for _, pt := range points {
		for i, v := range pt {
			if v < min[i] {
				min[i] = v
			}
			if v > max[i] {
				max[i] = v
			}
		}
	}
	return [][]float64{min, max}
}

// A 2D bounding box tree

type BBTree struct {
	Root *BBNode
}

func NewBBTree() *BBTree {
	// The root bb is the entire plane
	return &BBTree{NewBBNode(-1, nil, [][]float64{{-math.MaxFloat64, -math.MaxFloat64}, {math.MaxFloat64, math.MaxFloat64}})}
}

func (bt *BBTree) Insert(id int, bb [][]float64) {
	bt.Root.Decompose(id, bb)
}

func (bt *BBTree) Search(pt []float64) []*BBNode {
	path := bt.Root.Search(pt, nil)

	// Remove top level infinite bb and duplicates
	res := []*BBNode{}
	cur := -1
	for _, bn := range path {
		if cur == bn.Id {
			continue
		}
		cur = bn.Id
		res = append(res, bn)
	}
	return res
}

func (bt *BBTree) Leaves() [][][]float64 {
	return bt.Root.Leaves()
}

func (bt *BBTree) Flatten() [][][]float64 {
	return bt.Root.Flatten()
}

type BBNode struct {
	Id       int
	Depth    int
	BBox     [][]float64
	Parent   *BBNode
	Children []*BBNode
}

func NewBBNode(id int, parent *BBNode, bb [][]float64) *BBNode {
	/*
		// Sanity check - parent should contain bb
		if parent != nil &&
			(bb[0][0] < parent.BBox[0][0] || bb[0][0] > parent.BBox[1][0] ||
				bb[0][1] < parent.BBox[0][1] || bb[0][1] > parent.BBox[1][1] ||
				bb[1][0] < parent.BBox[0][0] || bb[1][0] > parent.BBox[1][0] ||
				bb[1][1] < parent.BBox[0][1] || bb[1][1] > parent.BBox[1][1]) {
			panic(fmt.Errorf("bb not contained in parent"))
		}
	*/

	if parent != nil {
		return &BBNode{id, parent.Depth + 1, bb, parent, nil}
	}

	return &BBNode{id, 0, bb, nil, nil}
}

func (bn *BBNode) Decompose(id int, bb [][]float64) {
	ibb := bn.Intersection(bb)
	if ibb == nil {
		return
	}

	if len(bn.Children) > 0 {
		// Run bb over children
		for _, child := range bn.Children {
			child.Decompose(id, bb)
		}
		return
	}

	// Dice up bn.BBox by the intersection with successive slices, biggest first.
	res := []*BBNode{NewBBNode(id, bn, ibb)}
	rem := bn.BBox
	for {
		chunk, remainder := sliceMax(rem, ibb)
		if chunk == nil {
			break
		}
		res = append(res, NewBBNode(bn.Id, bn, chunk))
		rem = remainder
	}

	bn.Children = res
}

// Assume ibb contained in bb - slice off the biggest piece we can.
// return chunk and remainder
func sliceMax(bb, ibb [][]float64) ([][]float64, [][]float64) {
	// Deal with infinities first
	if bb[0][0] == -math.MaxFloat64 {
		left, right := sliceInX(ibb[0][0], bb)
		return left, right
	}
	if bb[1][0] == math.MaxFloat64 {
		left, right := sliceInX(ibb[1][0], bb)
		return right, left
	}
	if bb[0][1] == -math.MaxFloat64 {
		top, bottom := sliceInY(ibb[0][1], bb)
		return top, bottom
	}
	if bb[1][1] == math.MaxFloat64 {
		top, bottom := sliceInY(ibb[1][1], bb)
		return bottom, top
	}

	max := 0.0
	imax := -1
	r := (bb[1][0] - bb[0][0]) / (bb[1][1] - bb[0][1])

	d := ibb[0][0] - bb[0][0]
	if d > max {
		max = d
		imax = 0
	}
	d = bb[1][0] - ibb[1][0]
	if d > max {
		max = d
		imax = 1
	}
	d = (ibb[0][1] - bb[0][1]) * r
	if d > max {
		max = d
		imax = 2
	}
	d = (bb[1][1] - ibb[1][1]) * r
	if d > max {
		imax = 3
	}

	switch imax {
	case -1:
		return nil, nil
	case 0:
		left, right := sliceInX(ibb[0][0], bb)
		return left, right
	case 1:
		left, right := sliceInX(ibb[1][0], bb)
		return right, left
	case 2:
		top, bottom := sliceInY(ibb[0][1], bb)
		return top, bottom
	case 3:
		top, bottom := sliceInY(ibb[1][1], bb)
		return bottom, top
	}

	// Shouldn't get here
	return nil, nil
}

// Assume x between min and max
func sliceInX(x float64, bb [][]float64) ([][]float64, [][]float64) {
	// L, R
	return [][]float64{bb[0], {x, bb[1][1]}}, [][]float64{{x, bb[0][1]}, bb[1]}
}

// Assume y between min and max
func sliceInY(y float64, bb [][]float64) ([][]float64, [][]float64) {
	// T, B
	return [][]float64{bb[0], {bb[1][0], y}}, [][]float64{{bb[0][0], y}, bb[1]}
}

func bbEmpty(bb [][]float64) bool {
	return bb[0][0] > bb[1][0] || bb[0][1] > bb[1][1] ||
		Within(bb[0][0], bb[1][0], 0.000001) || Within(bb[0][1], bb[1][1], 0.000001)
}

// Find the bounding box of the intersection of bb with bn.BBox
func (bn *BBNode) Intersection(bb [][]float64) [][]float64 {
	// Test for overlap and emptyness
	if bb[0][0] > bn.BBox[1][0] || bb[0][1] > bn.BBox[1][1] ||
		bb[1][0] < bn.BBox[0][0] || bb[1][1] < bn.BBox[0][1] ||
		bbEmpty(bb) || bbEmpty(bn.BBox) {
		return nil
	}

	// Intersection is the max of the mins and the min if the maxes
	minx := bb[0][0]
	if bn.BBox[0][0] > minx {
		minx = bn.BBox[0][0]
	}
	miny := bb[0][1]
	if bn.BBox[0][1] > miny {
		miny = bn.BBox[0][1]
	}
	maxx := bb[1][0]
	if bn.BBox[1][0] < maxx {
		maxx = bn.BBox[1][0]
	}
	maxy := bb[1][1]
	if bn.BBox[1][1] < maxy {
		maxy = bn.BBox[1][1]
	}

	if Within(maxx, minx, 0.000001) || Within(maxy, miny, 0.000001) {
		return nil
	}

	return [][]float64{{minx, miny}, {maxx, maxy}}
}

func (bn *BBNode) Search(pt []float64, path []*BBNode) []*BBNode {
	path = append(path, bn)

	if len(bn.Children) == 0 {
		return path
	}

	for _, child := range bn.Children {
		if child.Contains(pt) {
			return child.Search(pt, path)
		}
	}

	panic(fmt.Errorf("no children contain pt in search"))
	return path
}

func (bn *BBNode) Contains(pt []float64) bool {
	if pt[0] > bn.BBox[1][0] || pt[0] < bn.BBox[0][0] ||
		pt[1] < bn.BBox[0][1] || pt[1] > bn.BBox[1][1] {
		return false
	}
	return true
}

func (bn *BBNode) Leaves() [][][]float64 {
	if len(bn.Children) == 0 {
		if bn.Id == -1 {
			// Exclude ones that extend to infinity
			return nil
		}
		return [][][]float64{bn.BBox}
	}

	res := [][][]float64{}
	for _, child := range bn.Children {
		res = append(res, child.Leaves()...)
	}
	return res
}

func (bn *BBNode) Flatten() [][][]float64 {
	if len(bn.Children) == 0 {
		if bn.Id == -1 {
			// Exclude ones that extend to infinity
			return nil
		}
		return [][][]float64{bn.BBox}
	}

	res := [][][]float64{}
	if bn.Id != -1 {
		res = append(res, bn.BBox)
	}
	for _, child := range bn.Children {
		res = append(res, child.Flatten()...)
	}
	return res
}
