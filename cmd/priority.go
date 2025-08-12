//go:build ignore

package main

import (
	"fmt"
	"github.com/jphsd/datastruct"
)

func main() {
	l := datastruct.NewPriorityList()
	l.Insert(datastruct.NewPriorityItem(10, 1))
	l.Insert(datastruct.NewPriorityItem(9, 2))
	l.Insert(datastruct.NewPriorityItem(8, 3))
	l.Insert(datastruct.NewPriorityItem(7, 4))
	l.Insert(datastruct.NewPriorityItem(6, 5))
	for i, itm := range l.Slice() {
		fmt.Printf("%d: Pri %f Id %d\n", i, itm.Priority, itm.Id)
	}
	fmt.Printf("at %d\n", l.Insert(datastruct.NewPriorityItem(8, 6)))
	fmt.Printf("at %d\n", l.Insert(datastruct.NewPriorityItem(8, 7)))
	fmt.Printf("at %d\n", l.Insert(datastruct.NewPriorityItem(8, 8)))
	for i, itm := range l.Slice() {
		fmt.Printf("%d: Pri %f Id %d\n", i, itm.Priority, itm.Id)
	}

	for i := 0; i < 12; i++ {
		fmt.Printf("where %d right %d\n", i, l.Where(float64(i), false))
	}
	for i := 0; i < 12; i++ {
		fmt.Printf("where %d left %d\n", i, l.Where(float64(i), true))
	}

	l.Delete(datastruct.NewPriorityItem(7, 4))
	for i, itm := range l.Slice() {
		fmt.Printf("%d: Pri %f Id %d\n", i, itm.Priority, itm.Id)
	}
	l.Delete(datastruct.NewPriorityItem(8, 8))
	for i, itm := range l.Slice() {
		fmt.Printf("%d: Pri %f Id %d\n", i, itm.Priority, itm.Id)
	}

	pi := (*l)[4]
	pi.Priority -= 1
	l.ChangedPriority(pi)
	for i, itm := range l.Slice() {
		fmt.Printf("%d: Pri %f Id %d\n", i, itm.Priority, itm.Id)
	}
}
