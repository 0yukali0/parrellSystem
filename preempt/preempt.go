package preempt

import (
	"simulation/queue"
	"simulation/common"
	"simulation/object"
	"container/heap"
	//"fmt"
)

func Preempt(suspendSubmitTime uint64) {
	jobs := queue.GetJobsQueue()
	events := queue.GetEventsQueue()
	waitingJobs := make([]*object.Job, 0)
	for len(*jobs) > 0 {
		job := heap.Pop(jobs).(*object.Job)
		if common.Allocate(job.Allocation, job.Allocated) {
			job.Allocated = true
			event := job.GetManagerAndSetItRunning(suspendSubmitTime)
			event.ToNextStepInWaiting()
			//fmt.Printf("%10v, Non-FCFS-Allocate in %10v:	%3v:%3v	Sub:%10v	waitForStart:%10v\n", job.Id, event.TimeStamp, common.ProcessNum, job.Allocation, job.Submission, job.GetWaitingTimeBeforeRunning())
			heap.Push(events, event)
			continue
		}
		waitingJobs = append(waitingJobs, job)
	}
	for _,job := range waitingJobs {
		heap.Push(jobs, job)
	} 

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

		if event.Status == common.EventStatus[0] {
			ok := common.Allocate(event.GetJob().Allocation, event.GetJob().Allocated)
			if ok {
				job := event.GetJob()
				job.Allocated = true
				event.ToNextStep()
				//fmt.Printf("%10v, Non-FCFS-EAllocate in %10v:	%3v:%3v	Sub:%10v	waitForStart:%10v\n", job.Id, event.TimeStamp, common.ProcessNum, job.Allocation, job.Submission, job.GetWaitingTimeBeforeRunning())
				heap.Push(events, event)
			}
			heap.Push(jobs, event.GetJob())
		}
	}

	for _, event := range eventsBeforeLastRelease {
		heap.Push(events, event)
	}
	return
}