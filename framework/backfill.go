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
	profile := queue.GetSlotQueue()
	WaitingTotalTime := uint64(0)
	totalRequestLen := uint64(submitAndFinishQueue.Len())

	for pendingEventsNum := submitAndFinishQueue.Len(); pendingEventsNum > 0; {
		//submitAndFinishQueue.Show()
		if submitAndFinishQueue.Len() > 0 {
			event := heap.Pop(submitAndFinishQueue).(*object.Event)
			job := event.GetJob()

			switch event.GetStatus() {
			case "Submit":
				common.SetSystemClock(event.GetTimeStamp())
				if !event.BackfillActive {
					startTime := profile.Allocate(job.GetExecutionTime(), job.GetAllocation())
					event.SetBackillExceptTime(startTime)
					event.HandleBackfillSupport()
					heap.Push(submitAndFinishQueue, event)
					continue
				} 

				if common.Allocate(job.GetAllocation(), job.Allocated) {
					event.Handle("SubmitSucess")
				 }
				heap.Push(submitAndFinishQueue, event)
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