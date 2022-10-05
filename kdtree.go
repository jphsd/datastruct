package datastruct

import (
	"fmt"
	"math"
	"sort"
)

type KDTree struct {
	Dims  int
	Root  *KDNode
	Nodes []*KDNode
	Dist  func(a, b []float64) float64
}

func NewKDTree(dims int, points ...[]float64) *KDTree {
	t := &KDTree{dims, nil, nil, defaultDist}
	t.Root = t.subtree(0, 0, points...)
	return t
}

func defaultDist(a, b []float64) float64 {
	da, db := len(a), len(b)
	if da > db {
		da = db
	}
	val := 0.0
	for i := 0; i < da; i++ {
		v := a[i] - b[i]
		val += v * v
	}
	return val
}

func (t *KDTree) subtree(axis, depth int, points ...[]float64) *KDNode {
	n := len(points)
	if n == 0 {
		return nil
	}
	if n == 1 {
		node := &KDNode{Point: points[0], Depth: depth}
		node.Id = len(t.Nodes)
		t.Nodes = append(t.Nodes, node)
		return node
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i][axis] < points[j][axis]
	})
	mid := n / 2
	pt := points[mid]
	next := (axis + 1) % t.Dims
	ndepth := depth + 1
	node := &KDNode{len(t.Nodes), pt, nil, nil, depth}
	t.Nodes = append(t.Nodes, node)
	node.Left = t.subtree(next, ndepth, points[:mid]...)
	node.Right = t.subtree(next, ndepth, points[mid+1:]...)
	return node
}

func (t *KDTree) String() string {
	return printHelper(t.Root)
}

func printHelper(n *KDNode) string {
	if n != nil && (n.Left != nil || n.Right != nil) {
		return fmt.Sprintf("{%s, %s, %s}", printHelper(n.Left), n, printHelper(n.Right))
	}
	if n == nil {
		return "nil"
	}
	return n.String()
}

func (t *KDTree) Insert(pt []float64) int {
	if len(pt) < t.Dims {
		return -1
	}
	var node *KDNode
	if t.Root == nil {
		t.Root = &KDNode{Point: pt}
		node = t.Root
	} else {
		node = t.Root.Insert(t.Dims, pt, 0)
	}
	node.Id = len(t.Nodes)
	t.Nodes = append(t.Nodes, node)
	return node.Id
}

// Poor man's expensive
func (t *KDTree) RemoveByPoint(pt []float64) bool {
	pts := t.Points()
	npts := [][]float64{}
	found := false
	for _, p := range pts {
		reject := true
		for i := 0; i < t.Dims; i++ {
			if !Within(p[i], pt[i], 0.000001) {
				reject = false
				break
			}
		}
		if !reject {
			npts = append(npts, pt)
		} else {
			found = true
			break
		}
	}
	if !found {
		return false
	}

	t.Nodes = nil
	t.Root = t.subtree(0, 0, npts...)
	return true
}

func (t *KDTree) RemoveById(id int) bool {
	n := len(t.Nodes)
	if id >= n {
		return false
	}
	n--
	pts := t.Points()
	copy(pts[id:], pts[id+1:])
	pts[len(pts)-1] = nil
	pts = pts[:len(pts)-1]

	t.Nodes = nil
	t.Root = t.subtree(0, 0, pts...)
	return true
}

func (t *KDTree) Balance() {
	pts := t.Points()
	t.Nodes = nil
	t.Root = t.subtree(0, 0, pts...)
}

// Insertion order
func (t *KDTree) Points() [][]float64 {
	n := len(t.Nodes)
	res := make([][]float64, n)
	for i, node := range t.Nodes {
		res[i] = node.Point
	}
	return res
}

func (t *KDTree) InOrderPoints() [][]float64 {
	if t.Root == nil {
		return [][]float64{}
	}
	return t.Root.Points()
}

// Find up to k nearest points to pt
// Returns points, distances, indices to kdnodes
func (t *KDTree) KNN(pt []float64, k int) ([][]float64, []float64, []int) {
	inds := []int{}
	res := [][]float64{}
	ds := []float64{}

	if t.Root == nil || k < 1 || len(pt) < t.Dims {
		return res, ds, inds
	}

	pq := NewPriorityList()
	t.knnHelper(pt, k, t.Root, 0, pq)

	for i := 0; i < k && i < len(*pq); i++ {
		id := (*pq)[i].Id
		inds = append(inds, id)
		res = append(res, t.Nodes[id].Point)
		ds = append(ds, (*pq)[i].Priority)
	}

	return res, ds, inds
}

func (t *KDTree) knnHelper(pt []float64, k int, node *KDNode, axis int, pq *PriorityList) {
	if node == nil {
		return
	}

	path := []*KDNode{}
	curr := node

	// Walk the tree to find insertion path of pt
	for curr != nil {
		path = append(path, curr)
		if pt[axis] < curr.Point[axis] {
			curr = curr.Left
		} else {
			curr = curr.Right
		}
		axis = (axis + 1) % t.Dims
	}

	// Retrace the path to create point list
	axis = (axis - 1 + t.Dims) % t.Dims
	for path, curr = last(path); curr != nil; path, curr = last(path) {
		dist := t.Dist(pt, curr.Point)
		checked := kthDistance(pq, k-1)
		if dist < checked {
			pq.Insert(PriorityItem{dist, curr.Id})
			checked = kthDistance(pq, k-1)
		}

		// check other side of plane
		if t.planeDist(pt, curr.Point, axis) < checked {
			var next *KDNode
			if pt[axis] < curr.Point[axis] {
				next = curr.Right
			} else {
				next = curr.Left
			}
			t.knnHelper(pt, k, next, (axis+1)%t.Dims, pq)
		}
		axis = (axis - 1 + t.Dims) % t.Dims
	}
}

// Find all points within d of pt. Note d must in in the same space as the Dist function
func (t *KDTree) DNN(pt []float64, d float64) ([][]float64, []float64, []int) {
	inds := []int{}
	res := [][]float64{}
	dists := []float64{}

	if t.Root == nil || d <= 0 || len(pt) < t.Dims {
		return res, dists, inds
	}

	pq := NewPriorityList()
	t.dnnHelper(pt, d, t.Root, 0, pq)

	for i := 0; i < len(*pq) && (*pq)[i].Priority <= d; i++ {
		id := (*pq)[i].Id
		inds = append(inds, id)
		res = append(res, t.Nodes[id].Point)
		dists = append(dists, (*pq)[i].Priority)
	}

	return res, dists, inds
}

func (t *KDTree) dnnHelper(pt []float64, d float64, node *KDNode, axis int, pq *PriorityList) {
	if node == nil {
		return
	}

	path := []*KDNode{}
	curr := node

	// Walk the tree to find insertion path of pt
	for curr != nil {
		path = append(path, curr)
		if pt[axis] < curr.Point[axis] {
			curr = curr.Left
		} else {
			curr = curr.Right
		}
		axis = (axis + 1) % t.Dims
	}

	// Retrace the path to create point list
	axis = (axis - 1 + t.Dims) % t.Dims
	for path, curr = last(path); curr != nil; path, curr = last(path) {
		dist := t.Dist(pt, curr.Point)
		if dist <= d {
			pq.Insert(PriorityItem{dist, curr.Id})
		}

		// check other side of plane
		if t.planeDist(pt, curr.Point, axis) <= d {
			var next *KDNode
			if pt[axis] < curr.Point[axis] {
				next = curr.Right
			} else {
				next = curr.Left
			}
			t.dnnHelper(pt, d, next, (axis+1)%t.Dims, pq)
		}
		axis = (axis - 1 + t.Dims) % t.Dims
	}
}

func last(p []*KDNode) ([]*KDNode, *KDNode) {
	l := len(p) - 1
	if l < 0 {
		return p, nil
	}
	return p[:l], p[l]
}

func (t *KDTree) planeDist(pt, cur []float64, axis int) float64 {
	a := make([]float64, len(pt))
	b := make([]float64, len(cur))
	a[axis] = pt[axis]
	b[axis] = cur[axis]
	return t.Dist(a, b)
}

func kthDistance(pq *PriorityList, k int) float64 {
	if len(*pq) <= k {
		return math.MaxFloat64
	}
	return (*pq)[k].Priority
}

type KDNode struct {
	Id    int
	Point []float64
	Left  *KDNode
	Right *KDNode
	Depth int
}

func (n *KDNode) String() string {
	res := "{"
	// res := fmt.Sprintf("{%d: ", n.Depth)
	for i, v := range n.Point {
		if i == 0 {
			res += fmt.Sprintf("%f", v)
		} else {
			res += fmt.Sprintf(", %f", v)
		}
	}
	return res + "}"
}

func (n *KDNode) Points() [][]float64 {
	res := [][]float64{}
	if n.Left != nil {
		res = n.Left.Points()
	}
	res = append(res, n.Point)
	if n.Right != nil {
		res = append(res, n.Right.Points()...)
	}
	return res
}

func (n *KDNode) Insert(dims int, p []float64, axis int) *KDNode {
	if len(p) < dims {
		return nil
	}
	var node *KDNode
	if p[axis] < n.Point[axis] {
		if n.Left == nil {
			node = &KDNode{Point: p, Depth: n.Depth + 1}
			n.Left = node
		} else {
			node = n.Left.Insert(dims, p, (axis+1)%dims)
		}
	} else {
		if n.Right == nil {
			node = &KDNode{Point: p, Depth: n.Depth + 1}
			n.Right = node
		} else {
			node = n.Right.Insert(dims, p, (axis+1)%dims)
		}
	}
	return node
}

// Within returns true if the two values are within e of each other.
func Within(d1, d2, e float64) bool {
	d := d1 - d2
	if d < 0.0 {
		d = -d
	}
	return d < e
}
