package framework

import (
	"container/heap"
	"simulation/object"
	"simulation/queue"
	"simulation/common"
	"fmt"
)
var (
	bk = make([]*object.Event, 0)
)

func Backfill() {
	submitAndFinishQueue := queue.GetEventsQueue()
	profile := queue.GetSlotQueue()
	WaitingTotalTime := uint64(0)
	totalRequestLen := uint64(submitAndFinishQueue.Len())

	for pendingEventsNum := submitAndFinishQueue.Len(); pendingEventsNum > 0; {
		//submitAndFinishQueue.Show()
		if submitAndFinishQueue.Len() > 0 {
			//fmt.Printf("Current system Time:%v\n", common.GetSystemClock())
			event := heap.Pop(submitAndFinishQueue).(*object.Event)
			job := event.GetJob()
			switch event.GetStatus() {
			case "Submit":
				common.SetSystemClock(event.GetTimeStamp())
				profile.Predicate()
				if !event.BackfillActive {
					startTime := profile.Allocate(job.GetExecutionTime(), job.GetAllocation())
					event.SetBackillExceptTime(startTime)
					event.HandleBackfillSupport()
					//profile.Show()
					heap.Push(submitAndFinishQueue, event)
					continue
				} 

				event.Handle("Backfill")
				job.SetResourceGetTime(common.GetSystemClock())
				job.ComputeFinishTime()
				job.ComputeWaitingTime()
				event.SetTimeStamp(job.GetFinishTime())
				heap.Push(submitAndFinishQueue, event)
			case "Finish":
				common.SetSystemClock(event.GetTimeStamp())
				profile.Predicate()
				WaitingTotalTime += job.GetResourceGetTime() - job.GetSubmitTime()
				bk = append(bk, event)
			}
		} else {
			break
		}
	}
		//7387
		fmt.Printf("Average time of %v jobs are %v seconds\n", totalRequestLen, (float64(WaitingTotalTime)/float64(totalRequestLen)))
		//detail()
}

func detail() {
	submitAndFinishQueue := queue.GetEventsQueue()
	for _, event := range bk {
		event.SetTimeStamp(event.GetJob().GetSubmitTime() + common.BaseSubmitTime)
		heap.Push(submitAndFinishQueue, event)
	}

	fmt.Printf("subTime // waitime // runtime // np // Id\n")
	for submitAndFinishQueue.Len() > 0 {
		event := heap.Pop(submitAndFinishQueue).(*object.Event)
		job := event.GetJob()
		fmt.Printf("%v(%v) // %v // %v // %v // %v\n",
		 job.GetSubmitTime() + common.BaseSubmitTime, job.GetSubmitTime(), job.GetWaitingTime(), job.GetExecutionTime(), job.GetAllocation(), job.Id)
	}
}