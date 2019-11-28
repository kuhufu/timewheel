package pq

import (
	"fmt"
	"testing"
)

var comparator = func(a, b interface{}) int {
	pa := a.(*Item).Key.(int)
	pb := b.(*Item).Key.(int)

	switch {
	case pa < pb:
		return 1
	case pa > pb:
		return -1
	}
	return 0
}

func TestMinPQ_Min(t *testing.T) {
	pq := NewMinPQ(16, comparator)

	pq.Insert(&Item{
		Key: 10,
		Val: 10,
	})

	pq.Insert(&Item{
		Key: 9,
		Val: 9,
	})

	pq.Insert(&Item{
		Key: 10,
		Val: 10,
	})

	pq.Insert(&Item{
		Key: 11,
		Val: 11,
	})

	pq.Insert(&Item{
		Key: 1,
		Val: 1,
	})

	for !pq.IsEmpty() {
		fmt.Println(pq.DelMin())
	}
}

