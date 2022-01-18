package queue

import (
	"container/heap"
	"simulation/object"
	"simulation/common"
	//"fmt"
)

var (
	SlotsQueue = make(SlotPQ, 0)
	Root = object.NewRootSlot(common.DefaultTimeStart, common.DefaultTimeLimit, common.GetSystemCapcity())
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
	for find := false;pq.Len() > 0 && !find; {
		slot := heap.Pop(pq).(*object.Slot)
		request := object.NewRequest(slot.GetStart(), duration, allocation)
		slot.TryAllocate(request)

		// for each slot, try sequence
		if slot.GetIsTrySuccess() {
			heap.Push(pq, slot)
			seq := true

			for durationSum := uint64(0);pq.Len() > 0 && durationSum < duration && !find && seq; {
				slot = heap.Pop(pq).(*object.Slot)
				durationSum += (slot.GetEnd() - slot.GetStart())
				slot.TryAllocate(request)
				if !slot.GetIsTrySuccess() {
					seq = false
				}

				// no seq, pop slot0, else add to condicate
				if !seq {
					for _, candicate := range candicates {
						bk = append(bk, candicate)
					}
					candicates = make([]*object.Slot, 0)
					bk = append(bk, slot)
					seq = true
					durationSum = 0
				} else {
					candicates = append(candicates, slot)
				}

				if seq && durationSum >= duration {
					find = true
					for index, cSlot := range candicates {
						if index == 0 {
							startTime = cSlot.GetStart()
							request.SetStartTime(startTime)
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
					break
				}
			}
		} else {
			bk = append(bk, slot)
			request.SetStartTime(slot.GetEnd())
		}
	}

	for _, slot := range bk {
		heap.Push(pq, slot)
	}
	return 
}