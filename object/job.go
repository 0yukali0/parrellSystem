package object

type Job struct {
	Id uint64
	Submission uint64
	ExecutionTime  uint64

	WaitingTime uint64
	ResourceGetTime uint64
	Allocated bool

	FinishTime uint64

	Allocation uint64
	Index int
	Manager *Event
}

func NewJob(id uint64, submitTime, executionTime, allocation uint64) *Job{
	return &Job{
		Id:				id,
		Submission: 	submitTime,
		ExecutionTime:	executionTime,
		Allocation:		allocation,
		ResourceGetTime: 0,
		Manager:		nil,
		Allocated: 		false, 
	}
}

func (j *Job) GetSubmitTime() uint64 {
	return j.Submission
}

func (j *Job) GetExecutionTime() uint64 {
	return j.ExecutionTime
}

func (j *Job) GetResourceGetTime() uint64 {
	return j.ResourceGetTime
}

func (j *Job) GetAllocation() uint64 {
	return j.Allocation
}

func (j *Job) GetWaitingTime() uint64 {
	return j.WaitingTime
}

func (j *Job) GetFinishTime() uint64 {
	return j.FinishTime
}

func (j *Job) SetResourceGetTime(releaseTimestamp uint64) {
	j.ResourceGetTime = releaseTimestamp
	j.Allocated = true
}

func (j *Job) SetManager(e *Event) {
	j.Manager = e
}

func (j *Job) ComputeWaitingTime() {
	time1 := j.Submission
	time2 := j.ResourceGetTime
	j.WaitingTime = time2 - time1
}

func (j *Job) ComputeFinishTime() {
	j.FinishTime = j.GetResourceGetTime() + j.GetExecutionTime()
}
