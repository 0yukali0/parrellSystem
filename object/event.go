package object

import (
	"strconv"
	"github.com/looplab/fsm"

	"simulation/reader"
	"simulation/common"
	"fmt"
)

type Event struct {
	Status    *fsm.FSM
	Index 	  int
	TimeStamp uint64
	JobMeta   *Job
}

func NewEvent(j *reader.JobDes, firstSubmitEventTime string) *Event {
	var submitTime, executionTime, processNum uint64
	var err error

	id := j.Name
	submitTime, err = strconv.ParseUint(j.Submit, 10, 64)
	common.Check(err)
	submitBase, err := strconv.ParseUint(firstSubmitEventTime, 10, 64)
	common.Check(err)
	submitTime -= submitBase
	executionTime, err = strconv.ParseUint(j.Running, 10, 64)
	common.Check(err)
	processNum, err = strconv.ParseUint(j.Allocation, 10, 64)
	common.Check(err)

	e := &Event{
		TimeStamp: submitTime,
		JobMeta:   NewJob(id, submitTime, executionTime, processNum),
	}

	e.Status = fsm.NewFSM(
		"Submit",
		fsm.Events{
			{Name:"SubmitSucess", Src: []string{"Submit"}, Dst:"Finish"},
			{Name:"SubmitFail", Src: []string{"Submit"}, Dst:"Waiting"},
			{Name:"WaitAndAllocated", Src: []string{"Waiting"}, Dst:"Finish"},
			{Name:"ReleaseResource", Src: []string{"Finish"}, Dst:"Release"},
		},
		fsm.Callbacks{
			"SubmitSucess": e.handleSubmitSucess,
			"SubmitFail": e.handleSubmitFail,
			"WaitAndAllocated": e.handleWaitAndAllocated,
			"ReleaseResource": e.handleReleaseResource,
		},
	)
	e.GetJob().SetManager(e)
	return e
}

func (e *Event) GetStatus() string {
	return e.Status.Current()
}

func (e *Event) GetJob() *Job{
	return e.JobMeta
}

func (e *Event) GetTimeStamp() uint64 {
	return e.TimeStamp
}

func (e *Event) SetTimeStamp(timeStamp uint64) {
	e.TimeStamp = timeStamp
}

func (e *Event) Handle(event string) {
	err := e.Status.Event(event)
	if err != nil {
		panic(e)
	}
}

func (e *Event) handleSubmitSucess(event *fsm.Event) {
	job := e.GetJob()
	job.SetResourceGetTime(e.GetTimeStamp())
	job.ComputeWaitingTime()
	job.ComputeFinishTime()
	fmt.Printf("%-6v EAllocate id:%-5v, %v in %v, cpu:%v,%v sub: %6v, getTime: %v, waiting: %v\n",
	 common.GetSystemClock(),
	 job.Id, e.Status.Current(), job.GetFinishTime(), common.GetCurrentProcessNum(), job.GetAllocation(), 
	 job.GetSubmitTime(), job.GetResourceGetTime(), job.GetWaitingTime())
	
	e.SetTimeStamp(job.GetFinishTime())
}

func (e *Event) handleSubmitFail (event *fsm.Event) {
	job := e.GetJob()
	fmt.Printf("%-6v StartWaiting id:%-5v, cpu:%v,%v sub: %6v\n",
	common.GetSystemClock(),
	job.Id, common.GetCurrentProcessNum(), job.GetAllocation(), 
	job.GetSubmitTime())
}

func (e *Event) handleWaitAndAllocated(event *fsm.Event) {
	job := e.GetJob()
	job.SetResourceGetTime(common.GetSystemClock())
	job.ComputeWaitingTime()
	job.ComputeFinishTime()
	fmt.Printf("%-6v WaitingEnd id:%-5v, %v in %v, cpu:%v,%v sub: %6v, getTime: %v, waiting: %v\n",
	common.GetSystemClock(),
	job.Id, e.Status.Current(), job.GetFinishTime(), common.GetCurrentProcessNum(), job.GetAllocation(), 
	job.GetSubmitTime(), job.GetResourceGetTime(), job.GetWaitingTime())
	e.SetTimeStamp(job.GetFinishTime())
}

func (e *Event) handleReleaseResource(event *fsm.Event) {
	job := e.GetJob()
	common.Release(job.GetAllocation(), job.Allocated)
	fmt.Printf("%-6v Release id:%-5v, %v, cpu:%v,%v sub: %6v, exe: %v, getTime: %v, waiting: %v, Finish: %v\n",
	common.GetSystemClock(),
	job.Id, e.Status.Current(), common.GetCurrentProcessNum(), job.GetAllocation(), 
	job.GetSubmitTime(), job.GetExecutionTime(), job.GetResourceGetTime(), job.GetWaitingTime(), job.GetFinishTime())
}
