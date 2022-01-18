package object
/*
import (
	"fmt"
)
*/
func (s *Slot) TryAllocate(r *Request) bool {
	if s.GetResource() >= r.GetAllocation() {
		s.SetIsTrySuccess(true)
	} else {
		s.SetIsTrySuccess(false)
	}
	return s.GetIsTrySuccess()
}

func (s *Slot) Allocate(r *Request) {
	if s.GetEnd() <= r.GetFinishTime() {
		s.AllocateResource(r.GetAllocation())
	} else  {
		children := make([]*Slot, 0)
		for index := 0; index < 2;index++ {
			children = append(children, s.Copy())
		}
		//fmt.Printf("Last %v match %v-%v %v\n",s.GetStart(), r.GetStartTime(),r.GetFinishTime(), r.GetAllocation())
		children[0].SetEnd(r.GetFinishTime())
		children[0].AllocateResource(r.GetAllocation())
		children[1].SetStart(r.GetFinishTime())
		s.AddChildren(children)
	}
}