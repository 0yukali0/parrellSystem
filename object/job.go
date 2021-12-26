package object

import (
	"simulation/common"
	//"fmt"
)

type Job struct {
	Id string
	Submission uint64
	ExecutionTime  uint64
	ResourceGetTime uint64
	Allocation uint64
	Index int
	Manager *Event
	Allocated bool
	SimulateWaitDuration uint64
}

func NewJob(id string, sub, exe, alloc uint64) *Job{
	return &Job{
		Id:				id,
		Submission: 	sub,
		ExecutionTime:	exe,
		Allocation:		alloc,
		ResourceGetTime: 0,
		Manager:		nil,
		Allocated: 		false, 
	}
}

func (j *Job)SetManager(e *Event) {
	j.Manager = e
}

// start time in waiting queue 
func (j *Job)GetManagerAndSetItRunning(getResourceTime uint64) *Event{
	j.Manager.TimeStamp = getResourceTime + j.ExecutionTime
	j.Manager.Status = common.EventStatus[1]
	j.SetGetResourceTime(getResourceTime)
	return j.Manager
}

func (j *Job)SetGetResourceTime(releaseTimestamp uint64) {
	j.ResourceGetTime = releaseTimestamp
}

func (j *Job)Finish() {
	time1 := j.Submission
	time2 := j.ResourceGetTime
	j.SimulateWaitDuration = time2 - time1
}

func (j *Job)GetWaitingTime() uint64 {
	return j.SimulateWaitDuration
}

func (j *Job)GetWaitingTimeBeforeRunning() uint64{
	return j.ResourceGetTime - j.Submission
}
