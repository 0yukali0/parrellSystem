package main

import (
	"simulation/reader"
	"simulation/queue"
	"simulation/object"
	"simulation/common"
	"simulation/backfill"
	"simulation/preempt"
	"container/heap"
	"fmt"
)

var (
	backfillActive = false
	preemptActive = true
)
var LastResourceReleasTime uint64
var WaitingTotalTime uint64

func main() {
	WaitingTotalTime = 0
	result := reader.ReadFile(common.FilePath)
	events := queue.GetEventsQueue()
	jobs := queue.GetJobsQueue()
	heap.Init(jobs)
	heap.Init(events)
	var basic string
	// assign event to event queue
	for idx, meta := range result {
		if idx == 0 {
			basic = meta.Sub
		}
		event := object.NewEvent(common.EventStatus[0],&meta, basic)
		heap.Push(events, event)
	}

	totalLen := uint64(events.Len())
	finish := false
	NextBackfill := true
	for !finish {
		// 1.is job queue empty?
		if len := jobs.Len();len > 0 {
			job := heap.Pop(jobs).(*object.Job)
			event := job.GetManagerAndSetItRunning(LastResourceReleasTime)
			// 2. if it's not empty, check whether resource is enough.
			if common.Allocate(job.Allocation, job.Allocated) {
				job.Allocated = true
				event.ToNextStepInWaiting()
				//fmt.Printf("%10v, WAllocate in %10v:	%3v:%3v	Sub:%10v	waitForStart:%10v\n", job.Id, event.TimeStamp, common.ProcessNum, job.Allocation, job.Submission, job.GetWaitingTimeBeforeRunning())
				heap.Push(events, event)
				continue
			} else {
				heap.Push(jobs, job)
			}
		}

		// 4. Is all queue clean?
		if len := events.Len() + jobs.Len(); len == 0 {
			fmt.Printf("Out!\n")
			finish = true
			continue
		}

		// 3. when job queue check all or empty, handle event queue
		event := heap.Pop(events).(*object.Event)
		eventType := event.Status
		timeStamp := event.TimeStamp
		if eventType == common.EventStatus[0] {
			//Keep FCFS
			if !jobs.IsEmpty() {
				//fmt.Printf("%10v, Waiting for waiting queue and req %v\n", event.GetJob().Id, event.GetJob().Allocation)
				heap.Push(jobs, event.GetJob())
				if NextBackfill && backfillActive{
					backfill.Backfilling(event.TimeStamp)
					NextBackfill = false
				}
				if preemptActive {
					preempt.Preempt(event.TimeStamp)
				}
				continue
			}
	
			// resource isn't enough
			ok := common.Allocate(event.GetJob().Allocation, event.GetJob().Allocated)
			if !ok {
				//fmt.Printf("%10v, Waiting for cpu %v\n", event.GetJob().Id, event.GetJob().Allocation)
				heap.Push(jobs, event.GetJob())
				continue
			}
	
			// sufficient cpu
			job := event.GetJob()
			job.Allocated = true
			//fmt.Printf("%10v, EAllocate in %10v:	%3v:%3v	sub:%10v\n", job.Id, event.TimeStamp, common.ProcessNum, job.Allocation, job.Submission)
			event.ToNextStep()
			heap.Push(events, event)
			continue
		}

		NextBackfill = true
		timeStamp = event.TimeStamp
		for {
			job := event.GetJob()
			common.Release(job.Allocation,job.Allocated)
			job.Finish()
			WaitingTotalTime += job.GetWaitingTime()
			LastResourceReleasTime = job.GetWaitingTime() + job.Submission + job.ExecutionTime
			event.GetJob().Allocated = false
			//fmt.Printf("%10v, Release in %10v:	%3v:%3v	sub:%10v	waitForStart:%10v	GetTime:%10v	Finish:%10v\n", job.Id, event.TimeStamp , common.ProcessNum, job.Allocation, job.Submission, job.GetWaitingTimeBeforeRunning(), job.ResourceGetTime, job.Submission+job.SimulateWaitDuration+job.ExecutionTime)
			if events.Len() == 0 {
				break
			}
			event = heap.Pop(events).(*object.Event)
			eventType = event.Status
			if eventType == common.EventStatus[0] || timeStamp != event.TimeStamp{
				heap.Push(events, event)
				break
			}
		}
	}
	//58206 7387 7300
	fmt.Printf("Average time of %v jobs are %v seconds\n", totalLen,WaitingTotalTime/totalLen)
}
