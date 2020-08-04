package timewheel

import (
	"sync"
)

type Item struct {
	Key int64
	Val *Timer
}

//小堆
type MinPQ struct {
	//
	pq []*Item

	//队列item数量
	n int

	//比较器
	comparator Comparator

	sync.Mutex
}

func NewMinPQ(cap int) *MinPQ {
	q := &MinPQ{
		pq:         make([]*Item, cap+1),
		n:          0,
		comparator: comparator,
	}

	return q
}

func (m *MinPQ) IsEmpty() bool {
	return m.n == 0
}

func (m *MinPQ) Len() int {
	return m.n
}

func (m *MinPQ) Min() (item *Item, ok bool) {
	m.Lock()
	defer m.Unlock()

	return m.min()
}

func (m *MinPQ) min() (item *Item, ok bool) {
	if m.IsEmpty() {
		return nil, false
	}
	return m.pq[1], true
}

func (m *MinPQ) PullMinIf(match func(item *Item) bool) (item *Item, ok bool) {
	m.Lock()
	defer m.Unlock()

	item, ok = m.min()

	if ok && match(item) {
		return m.delMin()
	}

	return nil, false
}

func (m *MinPQ) resize(cap int) {
	//这边值得考虑一下，缩小的时候是否要重新申请内存
	tmp := make([]*Item, cap)
	copy(tmp, m.pq)
	m.pq = tmp
}

func (m *MinPQ) Insert(item *Item) {
	m.Lock()
	defer m.Unlock()

	if m.n == len(m.pq)-1 {
		m.resize(len(m.pq) * 2)
	}

	m.n++
	m.pq[m.n] = item
	m.swim(m.n)
}

func (m *MinPQ) DelMin() (item *Item, ok bool) {
	m.Lock()
	defer m.Unlock()

	return m.delMin()
}

func (m *MinPQ) delMin() (item *Item, ok bool) {
	if m.IsEmpty() {
		return nil, false
	}

	min := m.pq[1]
	m.swap(1, m.n)
	m.n--
	m.sink(1)
	m.pq[m.n+1] = nil
	if m.n > 0 && m.n == (len(m.pq)-1)/4 {
		m.resize(len(m.pq) / 2)
	}
	return min, true
}

func (m *MinPQ) swim(k int) {
	for k > 1 && m.greater(k/2, k) {
		m.swap(k/2, k)
		k = k / 2
	}
}

func (m *MinPQ) sink(k int) {
	for 2*k < m.n {
		j := 2 * k
		if j < m.n && m.greater(j, j+1) {
			j++
		}
		if !m.greater(k, j) {
			break
		}
		m.swap(k, j)
		k = j
	}
}

func (m *MinPQ) greater(i, j int) bool {
	return m.comparator(m.pq[i], m.pq[j]) == 1
}

func (m *MinPQ) swap(i, j int) {
	m.pq[i], m.pq[j] = m.pq[j], m.pq[i]
}
