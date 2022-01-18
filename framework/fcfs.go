package framework

import (
	"container/heap"
	"simulation/object"
	"simulation/queue"
	"simulation/common"
	"fmt"
)

func FCFS() {
	submitAndFinishQueue := queue.GetEventsQueue()
	waitingQueue := queue.GetWaitingQueue()
	WaitingTotalTime := uint64(0)
	totalRequestLen := uint64(submitAndFinishQueue.Len())

	for pendingEventsNum := submitAndFinishQueue.Len() + waitingQueue.Len(); pendingEventsNum > 0; {
		if submitAndFinishQueue.Len() > 0 {
			event := heap.Pop(submitAndFinishQueue).(*object.Event)
			job := event.GetJob()

			// In this moment, event action
			switch event.GetStatus() {
			case "Submit":
				common.SetSystemClock(event.GetJob().GetSubmitTime())
				// FCFS base condition
				if waitingQueue.Len() > 0 {
					event.Handle("SubmitFail")
					heap.Push(waitingQueue, event)
					continue
				}

				if common.Allocate(job.GetAllocation(), job.Allocated) {
					event.Handle("SubmitSucess")
					heap.Push(submitAndFinishQueue, event)
				 } else {
					event.Handle("SubmitFail")
					heap.Push(waitingQueue, event)
				}
			case "Finish":
				common.SetSystemClock(event.GetJob().GetFinishTime())
				event.Handle("ReleaseResource")
				WaitingTotalTime += job.GetWaitingTime()
			}
		} else {
			break
		}

		for waitingQueue.Len() > 0 {
			event := heap.Pop(waitingQueue).(*object.Event)
			job := event.GetJob()
			if common.Allocate(job.GetAllocation(), job.Allocated) {
				event.Handle("WaitAndAllocated")
				heap.Push(submitAndFinishQueue, event)
				continue
			} else {
				heap.Push(waitingQueue, event)
				break
			}
		}
	}
		fmt.Printf("Average time of %v jobs are %v seconds\n", totalRequestLen, (float64(WaitingTotalTime)/float64(totalRequestLen)))
}