package datastruct

// PriorityItem contains the priority value and an id.
type PriorityItem struct {
	Priority float64
	Id       int
}

// PriorityList holds the prioritized list of items. The lower the priority, the closer to the start
// of the list the item is.
type PriorityList []PriorityItem

// NewPriorityList creates a new PriorityList with the items inserted. Lower values are inserted
// before higher ones.
func NewPriorityList(items ...PriorityItem) *PriorityList {
	res := &PriorityList{}
	for _, v := range items {
		res.Insert(v)
	}
	return res
}

// Insert inserts the item into the list at the correct point and returns that insertion point.
// Insertion is performed using a binary search and copy() for speed.
func (pq *PriorityList) Insert(v PriorityItem) int {
	if len(*pq) == 0 {
		// First entry
		*pq = append(*pq, v)
		return 0
	}
	n := len(*pq)
	res := pq.helper(v.Priority, 0, n-1, false)
	if res == n {
		// Insert at end
		*pq = append(*pq, v)
		return res
	}
	// Insert at res
	*pq = append(*pq, PriorityItem{})
	copy((*pq)[res+1:], (*pq)[res:]) // copy up
	(*pq)[res] = v
	return res
}

// Delete removes the entry in the list with the item (if found) and returns true. If the item isn't
// then false is returned. The priority value in the item is used to find where the item occurs in the
// list.
func (pq *PriorityList) Delete(v PriorityItem) bool {
	n := len(*pq)
	res := pq.helper(v.Priority, 0, n-1, false)
	if res == n {
		return false
	}
	for i := res; i >= 0; i-- {
		if (*pq)[i].Id == v.Id {
			return pq.DeleteEntry(i)
		}
		if (*pq)[i].Priority < v.Priority {
			break
		}
	}
	return false
}

// DeleteId removes the entry in the list with the item id (if found) and returns true. If the id
// isn't found then false is returned. This function uses a linear scan to find the id.
func (pq *PriorityList) DeleteId(id int) bool {
	for i := 0; i < len(*pq); i++ {
		if (*pq)[i].Id == id {
			return pq.DeleteEntry(i)
		}
	}
	return false
}

// DeleteEntry removes an entry from the list, compacts it and returns true. If the entry is not in range
// then false is returned.
func (pq *PriorityList) DeleteEntry(e int) bool {
	n := len(*pq)
	if e > n-1 {
		return false
	}
	if e == 0 {
		*pq = (*pq)[1:]
		return true
	}
	if e < n-1 {
		copy((*pq)[e:], (*pq)[e+1:]) // copy down
	}
	*pq = (*pq)[:n-1] // shrink
	return true
}

// Recursive binary search function.
func (pq *PriorityList) helper(v float64, s, e int, left bool) int {
	if v < (*pq)[s].Priority {
		// Before s
		return s
	}
	if v > (*pq)[e].Priority {
		// After e
		return e + 1
	}

	if left {
		if e-s == 1 {
			if v > (*pq)[s].Priority {
				return e
			} else {
				return s
			}
		}
		d := (e - s + 1) / 2
		if d == 0 {
			return s
		}
		if v > (*pq)[e-d].Priority {
			return pq.helper(v, e-d, e, true)
		}
		return pq.helper(v, s, e-d, true)
	}

	if e-s == 1 {
		if v < (*pq)[e].Priority {
			return e
		} else {
			return e + 1
		}
	}
	d := (e - s + 1) / 2
	if d == 0 {
		return e + 1
	}
	if v < (*pq)[s+d].Priority {
		return pq.helper(v, s, s+d, false)
	}
	return pq.helper(v, e-d, e, false)
}

// Where returns where an item of priority pri would be inserted. The bool left indicates
// for priorities of the same value whether the new one should be inserted to the left or
// right of the current ones.
func (pq *PriorityList) Where(v float64, left bool) int {
	return pq.helper(v, 0, len(*pq)-1, left)
}
