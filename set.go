package datastruct

import "fmt"

// Set represents a set of integer elements.
type Set map[int]bool

// NewSet returns a new set with the provided elements in it.
func NewSet(elts ...int) Set {
	res := make(Set)
	for _, e := range elts {
		res[e] = true
	}
	return res
}

// Add adds the element e to the set and returns true if it wasn't already in the set. Adding
// an element when it already exists is a no-op signified by false.
func (s Set) Add(e int) bool {
	_, ok := s[e]
	if !ok {
		s[e] = true
		return true
	}
	return false
}

// Remove removes element e from the set and returns true if it was in the set.
// Removing a non-existent element is a no-op signified by false.
func (s Set) Remove(e int) bool {
	_, ok := s[e]
	if !ok {
		return false
	}
	delete(s, e)
	return true
}

// Element returns true if the set contains the element e.
func (s Set) Element(e int) bool {
	return s[e]
}

// Empty returns true if the set is the empty set.
func (s Set) Empty() bool {
	return len(s) == 0
}

// Union returns a new set containing the union of the set and b (OR).
func (s Set) Union(b Set) Set {
	return Union(s, b)
}

// Intersection returns a new set containing the intersection of the set and b (AND).
func (s Set) Intersection(b Set) Set {
	return Intersection(s, b)
}

// Difference returns a new set containing only the elements in either the set or b but not in both (XOR).
func (s Set) Difference(b Set) Set {
	return Difference(s, b)
}

// Sub returns a new set containing the elements in the set which are not in b (SUB).
func (s Set) Sub(b Set) Set {
	return Sub(s, b)
}

// Contains returns true if b is completely contained in the set.
func (s Set) Contains(b Set) bool {
	return Contains(s, b)
}

// Disjoint returns true if the set and b share no elements in common.
func (s Set) Disjoint(b Set) bool {
	return Disjoint(s, b)
}

// String returns a string representation of the set.
func (s Set) String() string {
	if len(s) == 0 {
		return "{}"
	}
	res := "{"
	first := true
	for k := range s {
		if first {
			res += fmt.Sprintf("%d", k)
			first = false
		} else {
			res += fmt.Sprintf(", %d", k)
		}
	}
	return res + "}"
}

// Slice returns a slice representation of the set.
func (s Set) Slice() []int {
	n := len(s)
	res := make([]int, n)
	i := 0
	for k := range s {
		res[i] = k
		i++
	}
	return res
}

// Union returns a new set containing the union of a and b (OR).
func Union(a, b Set) Set {
	res := make(Set)
	for e := range a {
		res[e] = true
	}
	for e := range b {
		res[e] = true
	}
	return res
}

// Intersection returns a new set containing the intersection of a and b (AND).
func Intersection(a, b Set) Set {
	res := make(Set)
	la, lb := len(a), len(b)
	if la < lb {
		for e := range a {
			if b[e] {
				res[e] = true
			}
		}
	} else {
		for e := range b {
			if a[e] {
				res[e] = true
			}
		}
	}
	return res
}

// Difference returns a new set containing only the elements in either a or b but not in both (XOR).
func Difference(a, b Set) Set {
	// return Sub(Union(a, b), Intersection(a, b))
	res := make(Set)
	for e := range a {
		if !b[e] {
			res[e] = true
		}
	}
	for e := range b {
		if !a[e] {
			res[e] = true
		}
	}
	return res
}

// Sub returns a new set containing the elements in a which are not in b (SUB).
func Sub(a, b Set) Set {
	res := make(Set)
	for e := range a {
		if !b[e] {
			res[e] = true
		}
	}
	return res
}

// Contains returns true if b is completely contained in a.
func Contains(a, b Set) bool {
	for e := range b {
		if !a[e] {
			return false
		}
	}
	return true
}

// Disjoint returns true if a and b share no elements in common.
func Disjoint(a, b Set) bool {
	// return Intersection(a, b).Empty()
	la, lb := len(a), len(b)
	if la < lb {
		for e := range a {
			if b[e] {
				return false
			}
		}
	} else {
		for e := range b {
			if a[e] {
				return false
			}
		}
	}
	return true
}
