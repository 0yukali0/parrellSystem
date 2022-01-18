package queue

import (
	"container/heap"
	"simulation/object"
	"simulation/common"
)

var (
	SlotsQueue = make(SlotPQ, 0)
	Root = object.NewRootSlot(common.DefaultTimeStart, common.DefaultTimeLimit, common.GetSystemCapcity(), 0)
)

func init() {
	heap.Init(GetSlotQueue())
	heap.Push(GetSlotQueue(), Root)
}

func GetSlotQueue() *SlotPQ {
	return &SlotsQueue
}

func (pq *SlotPQ) Allocate(duration, allocation uint64) (startTime uint64) {
	pq.Predicate()
	bk := make([]*object.Slot, 0)
	candicates := make([]*object.Slot, 0)
	for pq.Len() > 0 {
		slot := heap.Pop(pq).(*object.Slot)
		request := object.NewRequest(slot.GetStart(), slot.GetStart() + duration, allocation)
		slot.TryAllocate(request)

		// for each slot, try sequence
		if slot.GetIsTrySuccess() {
			heap.Push(pq, slot)
			seq := true
			for durationSum := uint64(0);pq.Len() > 0 && durationSum < duration; {
				slot = heap.Pop(pq).(*object.Slot)
				slot.TryAllocate(request)
				if !slot.GetIsTrySuccess() {
					seq = false
				}

				durationSum += (slot.GetEnd() - slot.GetStart())
				// no seq, pop slot0, else add to condicate
				if !seq {
					for _, candicate := range candicates {
						heap.Push(pq, candicate)
					}
					bk = append(bk, heap.Pop(pq).(*object.Slot))
					seq = true
					durationSum = 0
				} else {
					candicates = append(candicates, slot)
				}

				if seq && durationSum >= duration {
					for index, cSlot := range candicates {
						if index == 0 {
							startTime = cSlot.GetStart()
						}

						cSlot.Allocate(request)
						if cSlot.IsParent() {
							for _, child := range cSlot.GetChildren() {
								heap.Push(pq, child)
							}
						} else {
							heap.Push(pq, cSlot)
						}
					}
				}
			}
		} else {
			bk = append(bk, slot)
			request.SetStartTime(slot.GetStart())
		}
	}
	return 
}