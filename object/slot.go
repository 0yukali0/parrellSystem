package object

import (
	"simulation/common"
)

type Request struct {
	start uint64
	end uint64
	allocation uint64
}

func NewRequest(start, end, allocation uint64) *Request {
	return &Request{
		start:		start,
		end:		end,
		allocation:	allocation,
	}
}

func (r *Request) GetStartTime() uint64 {
	return r.start
}

func (r *Request) GetFinishTime() uint64 {
	return r.end
}

func (r *Request) GetAllocation() uint64 {
	return r.allocation
}

type Slot struct {
	Start				uint64
	End					uint64
	Capicity			uint64
	Allocated			uint64
	Unused				uint64

	Index				int

	IsTrySuccess		bool
	IsParent			bool
	Parent				*Slot
	Child				[]*Slot
}

func NewSlot(start, end, cap, allocate uint64, parent *Slot) *Slot {
	s := &Slot{
		Start:			start,
		End:			end,
		Capicity:		common.GetSystemCapcity(),
		Allocated:		allocate,
		IsTrySuccess:	false,
		IsParent:		false,
		Parent:			parent,
		Child:			make([]*Slot, 0),
	}

	s.ComputeUnused()
	return s
}

func (s *Slot) Copy() *Slot {
	s2 := &Slot{
		Start:			s.GetStartTime(),
		End:			s.GetEndTime(),
		Capicity:		common.GetSystemCapcity(),
		Allocated:		s.GetAllocated(),
		IsTrySuccess:	false,
		IsParent:		false,
		Parent:			s,
		Child:			make([]*Slot, 0),
	}

	s2.ComputeUnused()
	return s2
}


func (s *Slot) TryAllocate(r *Request) bool {
	s.SetIsTrySuccess(r.GetAllocation() >= s.GetUnused())
	return s.GetIsTrySuccess()
}

/*
*	____________________
*	|				   |
*	| Allocatable cpu  |
*	|__________________|
*	|				   |
*	|   Allocated cpu  |
*	|__________________|
*	t0			   	   t1
*
*	Split table
*				|	end > t1	|	end == t1	|	end < t1	|
*	start >  t0	|		2		|		2		|		3		|
*	start == t0	|		1		|		1		|		2		|
*	start <  t0	|		1		|		1		|		2		|
*/

func (s *Slot) Allocate(r *Request) bool {
	if ok := s.GetIsTrySuccess();!ok {
		return ok
	}

	children := make([]*Slot, 0)
	front, mid, last := 0, 1, 2
	children = append(children, s.Copy())
	var startResult, finishResult bool

	if startResult = r.GetStartTime() > s.GetStartTime(); startResult {
		children = append(children, s.Copy())
		children[front].SetEndTime(r.GetStartTime())
		children[mid].SetStartTime(r.GetStartTime())
	} else {
		mid -= 1
		last -= 1
	}

	if finishResult = r.GetFinishTime() < s.GetEndTime(); finishResult {
		children = append(children, s.Copy())
		children[mid].SetEndTime(r.GetFinishTime())
		children[last].SetStartTime(r.GetFinishTime())
	}

	length := len(children)
	switch length {
	case 3:
		children[mid].AddAllocated(r.GetAllocation())
		children[mid].ComputeUnused()
	case 2:
		if startResult {
			children[mid].AddAllocated(r.GetAllocation())
			children[mid].ComputeUnused()
		}
		if finishResult {
			children[last].AddAllocated(r.GetAllocation())
			children[last].ComputeUnused()
		}
	default:
		children[front].AddAllocated(r.GetAllocation())
		children[front].ComputeUnused()
	}

	s.AddChildren(children)
	return s.GetIsTrySuccess()
}

func (s *Slot) AddChildren(children []*Slot) {
	s.SetIsParent(true)
	for _, child := range children {
		s.Child = append(s.Child, child)
	}
}

func (s *Slot) SetIsParent(yes bool) {
	s.IsParent = yes
}

func (s *Slot) SetIsTrySuccess(ok bool) {
	s.IsTrySuccess = ok
}

func (s *Slot) AddAllocated(req uint64) {
	s.Allocated += req
}

func (s *Slot) SetStartTime(start uint64) {
	s.Start = start
}

func (s *Slot) SetEndTime(end uint64)  {
	s.End = end
}

func (s *Slot) GetIsTrySuccess() bool {
	return s.IsTrySuccess
}

func (s *Slot) GetIsParent() bool {
	return s.IsParent
}

func (s *Slot) GetStartTime() uint64 {
	return s.Start
}


func (s *Slot) GetEndTime() uint64 {
	return s.End
}

func (s *Slot) GetCapcity() uint64 {
	return s.Capicity
}

func (s *Slot) GetAllocated() uint64 {
	return s.Allocated
}

func (s *Slot) GetUnused() uint64 {
	return s.Unused
}

func (s *Slot) ComputeUnused() uint64 {
	s.Unused = s.GetCapcity() - s.GetAllocated()
	return s.GetUnused()
}