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

	finish := false
	isEventInThisMomentsFinished := false
	for !finish {
		if waitingQueue.Len() > 0 {
			event := heap.Pop(waitingQueue).(*object.Event)
			job := event.GetJob()
			if common.Allocate(job.GetAllocation(), job.Allocated) {
				event.Handle("WaitAndAllocated")
				heap.Push(submitAndFinishQueue, event)
				continue
			} else {
				isEventInThisMomentsFinished = true
				heap.Push(waitingQueue, event)
			}
		} else if waitingQueue.Len() == 0 {
			isEventInThisMomentsFinished = true
		}

		//fmt.Printf("EventQueue len: %v, WaitingQueue len %v\n", submitAndFinishQueue.Len(), waitingQueue.Len())
		if pendingEventsNum := submitAndFinishQueue.Len() + waitingQueue.Len(); pendingEventsNum == 0 {
			finish = true
			break
		}
	
		for submitAndFinishQueue.Len() > 0 {
			event := heap.Pop(submitAndFinishQueue).(*object.Event)
			job := event.GetJob()
			// find next timestamp and need waiting queue unlock it
			if event.GetTimeStamp() != common.GetSystemClock() {
				if isEventInThisMomentsFinished {
					common.SetSystemClock(event.GetTimeStamp())
				}
				heap.Push(submitAndFinishQueue, event)
				break
			}

			switch event.GetStatus() {
			case "Submit":
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
				event.Handle("ReleaseResource")
				WaitingTotalTime +=  job.GetWaitingTime()
			}
		}
	}
	fmt.Printf("Average time of %v jobs are %v seconds\n", totalRequestLen, (WaitingTotalTime/totalRequestLen))
}