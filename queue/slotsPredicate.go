package queue

import (
	"container/heap"
	"simulation/object"
	"simulation/common"
)

func (pq *SlotPQ) Predicate() {
	pq.cleanQueueBeforeSystemClocktime()
	pq.cleanParent()
}

func (pq *SlotPQ) cleanQueueBeforeSystemClocktime() {
	bk := make([]*object.Slot, 0)
	for currentTime := common.GetSystemClock();pq.Len() > 0; {
		slot := heap.Pop(pq).(*object.Slot)
		if currentTime >= slot.GetEnd() {
			continue
		}

		if currentTime >= slot.GetStart() {
			slot.SetStart(currentTime)
		}
		bk = append(bk, slot)
	}

	for _, slot := range bk {
		heap.Push(pq, slot)
	}
}

func (pq *SlotPQ) cleanParent() {
	bk := make([]*object.Slot, 0)
	for pq.Len() > 0 {
		target := heap.Pop(pq).(*object.Slot)
		if target.IsParent() {
			continue
		}
		bk = append(bk, target)
	}

	for _, item := range bk {
		heap.Push(pq, item)
	}
}