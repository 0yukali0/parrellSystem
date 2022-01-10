package backfill

import (
	"simulation/queue"
	"simulation/common"
	"simulation/object"
	"container/heap"
	"fmt"
)

func Backfilling(suspendSubmitTime uint64) {
	tmp,find := findFirstWaitingJobGetRunning()
	if find == false {
		fmt.Println("Must wrong in backfilling when all task can run without backfill")
	}

	releaseTime := tmp[0]
	jobs := queue.GetJobsQueue()
	events := queue.GetEventsQueue()
	waitingJobs := make([]*object.Job, 0)
	//waiting queue
	for len(*jobs) > 0 {
		job := heap.Pop(jobs).(*object.Job)
		exceptTime := suspendSubmitTime + job.ExecutionTime
		if exceptTime <= releaseTime {
			if common.Allocate(job.Allocation, job.Allocated) {
				job.Allocated = true
				event := job.GetManagerAndSetItRunning(suspendSubmitTime)
				event.ToNextStepInWaiting()
				//fmt.Printf("%10v, BackAllocate in %10v:	%3v:%3v	Sub:%10v	waitForStart:%10v\n", job.Id, event.TimeStamp, common.ProcessNum, job.Allocation, job.Submission, job.GetWaitingTimeBeforeRunning())
				heap.Push(events, event)
				continue
			}
		}
		waitingJobs = append(waitingJobs, job)
	}

	for _,job := range waitingJobs {
		heap.Push(jobs, job)
	} 

	// event queue
	eventsBeforeLastRelease := make([]*object.Event, 0)
	for len(*events) > 0 {
		event := heap.Pop(events).(*object.Event)
		if event.TimeStamp > suspendSubmitTime {
			eventsBeforeLastRelease = append(eventsBeforeLastRelease, event)
			break
		}

		if event.Status == common.EventStatus[1]{
			eventsBeforeLastRelease = append(eventsBeforeLastRelease, event)
			continue
		}

		if event.Status == common.EventStatus[0] && event.TimeStamp <= suspendSubmitTime {
			exceptTime := suspendSubmitTime + event.GetJob().ExecutionTime
			if exceptTime <= releaseTime {
				ok := common.Allocate(event.GetJob().Allocation, event.GetJob().Allocated)
				if ok {
					job := event.GetJob()
					job.Allocated = true
					event.ToNextStep()
					heap.Push(events, event)
				}
			}
			heap.Push(jobs, event.GetJob())
		}
	}

	for _, event := range eventsBeforeLastRelease {
		heap.Push(events, event)
	}

	return
}

func findFirstWaitingJobGetRunning() (ReleaseTimes []uint64, find bool) {
	ReleaseTimes = make([]uint64,0)
	events := queue.GetEventsQueue()
	jobs := queue.GetJobsQueue()
	cpuNum := common.GetCurrentProcessNum()
	target := heap.Pop(jobs).(*object.Job)
	backUp := make([]*object.Event, 0)
	find = false

	for cpuNum < target.GetAllocation() {
		x := heap.Pop(events).(*object.Event)
		if x.Status == common.EventStatus[1] {
			allocation := x.GetJob().GetAllocation()
			cpuNum += allocation
			ReleaseTimes = append(ReleaseTimes, x.TimeStamp)
		}
		backUp = append(backUp, x)
		if len(*events) == 0 {
			break
		}
	}

	if cpuNum >= target.GetAllocation() {
		find = true
	}

	for _,x := range backUp {
		heap.Push(events,x)
	}

	heap.Push(jobs, target)
	return 
}