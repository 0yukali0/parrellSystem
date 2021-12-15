package object

import (
	"simulation/common"
	"simulation/reader"
	"strconv"
)

type Event struct {
	Status    string
	Index 	  int
	TimeStamp uint64
	JobMeta   *Job
}

func NewEvent(eventType string, j *reader.JobDes, basicSub string) *Event {
	var timeStamp, sub, running, processNum uint64
	var err error
	sub, err = strconv.ParseUint(j.Sub, 10, 64)
	common.Check(err)
	sub2, err := strconv.ParseUint(basicSub, 10, 64)
	common.Check(err)
	sub = sub - sub2
	timeStamp = sub
	running, err = strconv.ParseUint(j.Running, 10, 64)
	common.Check(err)
	processNum, err = strconv.ParseUint(j.Allocation, 10, 64)
	common.Check(err)
	id := j.Name
	e := &Event{
		Status:    eventType,
		TimeStamp: timeStamp,
		JobMeta:   NewJob(id, sub, running, processNum),
	}
	e.JobMeta.SetManager(e)
	return e
}

// no waiting job will do it
func (e *Event) ToNextStep() (bool, uint64){
	eventType := e.Status
	if eventType == common.EventStatus[0] {
		e.Status =  common.EventStatus[1]
		job := e.GetJob()
		e.TimeStamp = job.Submission + job.ExecutionTime // finish time
		job.ResourceGetTime = job.Submission
	} else {
		return true ,e.GetJob().GetWaitingTime()
	}
	return false, 0
}

func (e *Event) ToNextStepInWaiting() (bool, uint64){
	eventType := e.Status
	if eventType == common.EventStatus[0] {
		e.Status =  common.EventStatus[1]
		job := e.GetJob()
		e.TimeStamp = job.ResourceGetTime + job.ExecutionTime
	} else {
		return true ,e.GetJob().GetWaitingTime()
	}
	return false, 0
}

func (e *Event) GetJob() *Job{
	return e.JobMeta
}
