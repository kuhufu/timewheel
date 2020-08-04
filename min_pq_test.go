package timewheel

import (
	"fmt"
	"testing"
)

func TestMinPQ_Min(t *testing.T) {
	pq := NewMinPQ(16)

	pq.Insert(&Item{
		Key: 10,
		Val: nil,
	})

	pq.Insert(&Item{
		Key: 9,
		Val: nil,
	})

	pq.Insert(&Item{
		Key: 10,
		Val: nil,
	})

	pq.Insert(&Item{
		Key: 11,
		Val: nil,
	})

	pq.Insert(&Item{
		Key: 1,
		Val: nil,
	})

	for !pq.IsEmpty() {
		fmt.Println(pq.DelMin())
	}
}
