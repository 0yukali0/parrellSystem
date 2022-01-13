package main

import (
	"container/heap"

	"simulation/object"
	"simulation/reader"
	"simulation/queue"
	"simulation/common"
	"fmt"
)

func main() {
	result := reader.ReadFile(common.FilePath)
	submitAndFinishQueue := queue.GetEventsQueue()
	waitingQueue := queue.GetWaitingQueue()

	// assign event to event queue
	var firstSubmitTime string
	for idx, meta := range result {
		if idx == 0 {
			firstSubmitTime = meta.Submit
		}
		event := object.NewEvent(&meta, firstSubmitTime)

		if idx == 0 {
			common.SetSystemClock(event.GetTimeStamp())
		}
		heap.Push(submitAndFinishQueue, event)
	}
	
	WaitingTotalTime := uint64(0)
	totalRequestLen := uint64(submitAndFinishQueue.Len())

	for pendingEventsNum := submitAndFinishQueue.Len() + waitingQueue.Len(); pendingEventsNum > 0; {
		if submitAndFinishQueue.Len() > 0 {
			event := heap.Pop(submitAndFinishQueue).(*object.Event)
			job := event.GetJob()

			// In this moment, event action
			switch event.GetStatus() {
			case "Submit":
				// FCFS base condition
				common.SetSystemClock(event.GetJob().GetSubmitTime())
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
	fmt.Printf("Average time of %v jobs are %v seconds\n", totalRequestLen, (WaitingTotalTime/totalRequestLen))
}