package datastruct

import "sort"

type si struct {
	a, b int
}

// InvertSort returns the list of switches that need to be made to render a list of integers into a sorted one.
func InvertSort(values []int) []int {
	lv := len(values)
	l := make([]si, lv)

	for i, v := range values {
		l[i] = si{v, i}
	}

	sort.SliceStable(l, func(i, j int) bool {
		return l[i].a < l[j].a
	})

	res := make([]int, lv)
	for i, v := range l {
		res[i] = v.b
	}

	return res
}
