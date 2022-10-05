//go:build ignore

package main

import (
	"fmt"
	"github.com/jphsd/datastruct"
)

func main() {
	s1 := datastruct.NewSet(1, 2, 3, 4, 5)
	fmt.Printf("New %v\n", s1)

	ok := s1.Add(10)
	fmt.Printf("Add %v %v\n", ok, s1)

	ok = s1.Add(10)
	fmt.Printf("Add %v %v\n", ok, s1)

	ok = s1.Remove(5)
	fmt.Printf("Remove %v %v\n", ok, s1)

	ok = s1.Remove(5)
	fmt.Printf("Remove %v %v\n", ok, s1)

	ok = s1.Element(4)
	fmt.Printf("Element %v %v\n", ok, s1)

	ok = s1.Empty()
	fmt.Printf("Empty %v %v\n", ok, s1)

	s2 := datastruct.NewSet()
	ok = s2.Empty()
	fmt.Printf("Empty %v %v\n", ok, s2)

	s1 = datastruct.NewSet(1, 2, 3, 4, 5)
	s2 = datastruct.NewSet(4, 5, 6, 7, 8)

	s3 := s1.Union(s2)
	fmt.Printf("%v Union %v = %v\n", s1, s2, s3)

	s3 = s1.Intersection(s2)
	fmt.Printf("%v Intersection %v = %v\n", s1, s2, s3)

	s3 = s1.Difference(s2)
	fmt.Printf("%v Difference %v = %v\n", s1, s2, s3)

	// s3 = s1.Difference(s2)
	s3 = s1.Union(s2).Sub(s1.Intersection(s2))
	fmt.Printf("%v Difference %v = %v\n", s1, s2, s3)

	s3 = s1.Sub(s2)
	fmt.Printf("%v Sub %v = %v\n", s1, s2, s3)

	ok = s1.Contains(s3)
	fmt.Printf("%v %v Contains %v\n", ok, s1, s3)

	ok = s2.Contains(s3)
	fmt.Printf("%v %v Contains %v\n", ok, s2, s3)

	ok = s2.Disjoint(s3)
	fmt.Printf("%v %v Disjoint %v\n", ok, s2, s3)

	// ok = s2.Disjoint(s3)
	ok = s2.Intersection(s3).Empty()
	fmt.Printf("%v %v Disjoint %v\n", ok, s2, s3)

	ok = s1.Disjoint(s3)
	fmt.Printf("%v %v Disjoint %v\n", ok, s1, s3)

	// ok = s1.Disjoint(s3)
	ok = s1.Intersection(s3).Empty()
	fmt.Printf("%v %v Disjoint %v\n", ok, s1, s3)

	fmt.Printf("Slice %v %v\n", s1, s1.Slice())
}
