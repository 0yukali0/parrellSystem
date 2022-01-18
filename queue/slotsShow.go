package queue

import (
	"simulation/object"

	"container/heap"
	"fmt"
)

func (pq *SlotPQ) Show() {
	bk := make([]*object.Slot, 0)
	for pq.Len() > 0 {
		slot := heap.Pop(pq).(*object.Slot)
		bk = append(bk, slot)
		fmt.Printf("(%v,%v-%v) ",slot.GetStart(), slot.GetEnd(), slot.GetResource())
	}
	fmt.Println()
	for _, item := range bk {
		heap.Push(pq, item)
	}
}