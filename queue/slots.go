package queue

import (
	"simulation/object"
)

type SlotPQ []*object.Slot

func (pq SlotPQ) Len() int { return len(pq) }

func (pq SlotPQ) Less(i, j int) bool {
	return pq[i].GetStart() < pq[j].GetStart()
}

func (pq SlotPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
} 

func (pq *SlotPQ) Push(x interface{}) {
	n := len(*pq)
	item := x.(*object.Slot)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *SlotPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}