package datastruct

import (
	"container/heap"
	"errors"
)

var (
	// ErrEmpty is returned when an attempt is made to pop an empty PriorityQueue
	ErrEmpty = errors.New("Empty queue")
)

// PriorityQueue wraps a minQueue (see example in container/heap) to a straight integer id
// based priority list.
type PriorityQueue struct {
	id2itm map[int]*pqitem
	items  minQueue
}

// NewPriorityQueue creates a new queue instance
func NewPriorityQueue() *PriorityQueue {
	id2item := make(map[int]*pqitem)
	items := minQueue{}
	heap.Init(&items)
	return &PriorityQueue{id2item, items}
}

// Len returns the number of entries in the queue
func (pq *PriorityQueue) Len() int {
	return len(pq.items)
}

// Insert a new id with priority or change the priority of an existin id
func (pq *PriorityQueue) Insert(id int, pri float64) {
	itm, ok := pq.id2itm[id]
	if ok {
		// Change an existing item's priority
		itm.priority = pri
		heap.Fix(&pq.items, itm.index)
		return
	}
	// New item
	itm = &pqitem{id, pri, -1}
	pq.id2itm[id] = itm
	heap.Push(&pq.items, itm)
}

// Pop returns the lowest priority id and remoes it from the queue
func (pq *PriorityQueue) Pop() (int, error) {
	if len(pq.items) == 0 {
		return 0, ErrEmpty
	}
	itm := heap.Pop(&pq.items).(*pqitem)
	pq.id2itm[itm.id] = nil
	delete(pq.id2itm, itm.id)
	return itm.id, nil
}

// Use container.heap to implement MinQueue (see example).
// MinQueue must support heap.Interface and sort.Interface

type pqitem struct {
	id       int
	priority float64 // The priority of the item in the queue.
	index    int     // Location of this tiem for heap.Fix()
}

type minQueue []*pqitem

func (mq minQueue) Len() int { return len(mq) }

func (mq minQueue) Less(i, j int) bool {
	return mq[i].priority < mq[j].priority
}

func (mq minQueue) Swap(i, j int) {
	mq[i], mq[j] = mq[j], mq[i]
	mq[i].index = i
	mq[j].index = j
}

func (mq *minQueue) Push(x any) {
	n := len(*mq)
	itm := x.(*pqitem)
	itm.index = n
	*mq = append(*mq, itm)
}

func (mq *minQueue) Pop() any {
	old := *mq
	n := len(old)
	itm := old[n-1]
	old[n-1] = nil
	itm.index = -1
	*mq = old[0 : n-1]
	return itm
}
