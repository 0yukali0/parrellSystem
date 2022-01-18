package framework

import (
	"container/heap"
	"simulation/object"
	"simulation/queue"
	"simulation/common"
	"fmt"
)

func Preempt() {
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
				if common.PreemptActive {
					event.Handle("SubmitFail")
					heap.Push(waitingQueue, event)
					break
				}
				// FCFS base condition
				if waitingQueue.Len() > 0 {
					if common.BackfillActive {
						event.Handle("Backfill")
						heap.Push(submitAndFinishQueue, event)
						continue
					}
					event.Handle("SubmitFail")
					heap.Push(waitingQueue, event)
					continue
				}

				if common.Allocate(job.GetAllocation(), job.Allocated) {
					event.Handle("SubmitSucess")
					heap.Push(submitAndFinishQueue, event)
				 } else {
					if common.BackfillActive {
						event.Handle("Backfill")
						heap.Push(submitAndFinishQueue, event)
						continue
					}
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

		preemptFail := make([]*object.Event,0)
		for waitingQueue.Len() > 0 {
			event := heap.Pop(waitingQueue).(*object.Event)
			job := event.GetJob()
			if common.Allocate(job.GetAllocation(), job.Allocated) {
				event.Handle("WaitAndAllocated")
				heap.Push(submitAndFinishQueue, event)
				continue
			} else if common.PreemptActive {
				preemptFail = append(preemptFail, event)
			} else {
				heap.Push(waitingQueue, event)
				break
			}
		}

		// Put fail preempt items back to queue 
		if common.PreemptActive {
			for _, event := range preemptFail {
				heap.Push(waitingQueue, event)
			}
		}
	}
		fmt.Printf("Average time of %v jobs are %v seconds\n", totalRequestLen, (float64(WaitingTotalTime)/float64(totalRequestLen)))
}