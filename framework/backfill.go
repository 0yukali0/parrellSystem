package framework

import (
	"container/heap"
	"simulation/object"
	"simulation/queue"
	"simulation/common"
	"fmt"
)
func Backfill() {
	submitAndFinishQueue := queue.GetEventsQueue()
	WaitingTotalTime := uint64(0)
	totalRequestLen := uint64(submitAndFinishQueue.Len())

	for pendingEventsNum := submitAndFinishQueue.Len(); pendingEventsNum > 0; {
		if submitAndFinishQueue.Len() > 0 {
			event := heap.Pop(submitAndFinishQueue).(*object.Event)
			job := event.GetJob()

			// In this moment, event action
			switch event.GetStatus() {
			case "Submit":
				common.SetSystemClock(event.GetJob().GetSubmitTime())
				// FCFS base condition

				if common.Allocate(job.GetAllocation(), job.Allocated) {
					event.Handle("SubmitSucess")
					heap.Push(submitAndFinishQueue, event)
				 } else {
					if common.BackfillActive {
						event.Handle("Backfill")
						heap.Push(submitAndFinishQueue, event)
						continue
					}
				}
			case "Finish":
				common.SetSystemClock(event.GetJob().GetFinishTime())
				event.Handle("ReleaseResource")
				WaitingTotalTime += job.GetWaitingTime()
			}
		} else {
			break
		}
	}
		fmt.Printf("Average time of %v jobs are %v seconds\n", totalRequestLen, (float64(WaitingTotalTime)/float64(totalRequestLen)))
}