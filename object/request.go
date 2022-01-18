package object

type Request struct {
	start uint64
	duration uint64
	allocation uint64
}

func NewRequest(start, duration, allocation uint64) *Request {
	return &Request{
		start:		start,
		duration:	duration,
		allocation:	allocation,
	}
}

func (r *Request) GetStartTime() uint64 {
	return r.start
}

func (r *Request) SetStartTime(start uint64) {
	r.start = start
}

func (r *Request) GetDuration() uint64 {
	return r.duration
}

func (r *Request) GetFinishTime() uint64 {
	return r.GetStartTime() + r.GetDuration()
}

func (r *Request) GetAllocation() uint64 {
	return r.allocation
}