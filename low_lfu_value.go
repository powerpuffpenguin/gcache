package gcache

import (
	"container/heap"
	"time"
)

type lfuValue interface {
	cacheValue
	GetCount() int
	SetCount(count int)
	Increment()
	SetIndex(index int)
	GetIndex() int
}

func newLFUValue(key, val interface{}, expiry time.Duration) lfuValue {
	if expiry > 0 {
		return &deadlineLFUValue{
			baseLFUValue: baseLFUValue{
				baseValue: baseValue{
					key:   key,
					value: val,
				},
				count: 1,
			},
			deadline: time.Now().Add(expiry),
		}
	}
	return &baseLFUValue{
		baseValue: baseValue{
			key:   key,
			value: val,
		},
		count: 1,
	}
}

type baseLFUValue struct {
	baseValue
	count int
	index int
}

func (v *baseLFUValue) GetCount() int {
	return v.count
}
func (v *baseLFUValue) SetCount(count int) {
	v.count = count
}
func (v *baseLFUValue) SetIndex(index int) {
	v.index = index
}
func (v *baseLFUValue) GetIndex() int {
	return v.index
}
func (v *baseLFUValue) Increment() {
	v.count++
}

type deadlineLFUValue struct {
	baseLFUValue
	deadline time.Time
}

func (v *deadlineLFUValue) IsDeleted() bool {
	return !v.deadline.After(time.Now())
}
func (v *deadlineLFUValue) SetDeadline(deadline time.Time) {
	v.deadline = deadline
}

type lfuValueHeap []lfuValue

func (a lfuValueHeap) Len() int {
	return len(a)
}
func (a lfuValueHeap) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
	a[i].SetIndex(i)
	a[j].SetIndex(j)
}
func (a lfuValueHeap) Less(i, j int) bool {
	return a[i].GetCount() < a[j].GetCount()
}
func (h *lfuValueHeap) Push(x interface{}) {
	v := x.(lfuValue)
	v.SetIndex(len(*h))
	*h = append(*h, v)
}
func (h *lfuValueHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	old[n-1] = nil
	return x
}

type lfuHeap struct {
	heap lfuValueHeap
}

func newLFUHeap(capacity int) *lfuHeap {
	heap := make(lfuValueHeap, 0, capacity)
	return &lfuHeap{
		heap: heap,
	}
}
func (h *lfuHeap) Push(val lfuValue) int {
	heap.Push(&h.heap, val)
	return val.GetIndex()
}
func (h *lfuHeap) Len() int {
	return h.heap.Len()
}

func (h *lfuHeap) Remove(i int) lfuValue {
	v := heap.Remove(&h.heap, i)
	if v == nil {
		return nil
	}
	return v.(lfuValue)
}
func (h *lfuHeap) Fix(i int) {
	heap.Fix(&h.heap, i)
}
func (h *lfuHeap) Clear() {
	for i := 0; i < len(h.heap); i++ {
		h.heap[i] = nil
	}
	h.heap = h.heap[:0]
}
