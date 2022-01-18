package queue

import (
	"simulation/object"
	"container/heap"
	"fmt"
)

var (
	SubmitAndFinishQueue = make(EventPQ, 0)
	WaintingQueue = make(EventPQ, 0)
)

func init() {
	heap.Init(GetWaitingQueue())
	heap.Init(GetEventsQueue())
}

func GetEventsQueue() *EventPQ{
	return &SubmitAndFinishQueue
}

func GetWaitingQueue() *EventPQ{
	return &WaintingQueue
}

type EventPQ []*object.Event

func (pq *EventPQ) Show() {
	bk := make([]*object.Event, 0)
	for index := 5;pq.Len() > 0 && index > 0;index-- {
		event := heap.Pop(pq).(*object.Event)
		if event.GetStatus() == "Finish" {
			fmt.Printf("F")
		}
		fmt.Printf("%v:%v-%v ", event.GetTimeStamp(), event.GetTimeStamp()+event.GetJob().GetExecutionTime(),event.GetJob().Id)
		bk = append(bk, event)
	}
	fmt.Println()
	for _, item := range bk {
		heap.Push(pq, item)
	}
}

func (pq EventPQ) Len() int { return len(pq) }

func (pq EventPQ) Less(i, j int) bool {
	/*
	if pq[i].TimeStamp == pq[j].TimeStamp {
		return pq[i].GetJob().Id < pq[j].GetJob().Id
	}
	*/
	return pq[i].GetTimeStamp() < pq[j].GetTimeStamp()
}

func (pq EventPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
} 

func (pq *EventPQ) Push(x interface{}) {
	n := len(*pq)
	item := x.(*object.Event)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *EventPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}