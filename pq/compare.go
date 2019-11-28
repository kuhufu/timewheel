package pq

type Comparator func(a, b interface{}) int

type Comparable interface {
	CompareTo(a interface{}) int
}
