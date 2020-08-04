package timewheel

type Comparator func(a *Item, b *Item) int

var comparator = func(a, b *Item) int {
	pa := a.Key
	pb := b.Key

	switch {
	case pa > pb:
		return 1
	case pa < pb:
		return -1
	}
	return 0
}
