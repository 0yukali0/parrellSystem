package object

import(
	"simulation/common"
)

func NewRootSlot(start, end, cap, allocate uint64) *Slot {
	s := &Slot{
		Start:			start,
		End:			end,
		Resource:		common.GetSystemCapcity(),
		IsTrySuccess:	false,
		Parent:			nil,
		Child:			make([]*Slot, 0),
	}
	return s
}

func (s *Slot) Copy() *Slot {
	s2 := &Slot{
		Start:			s.GetStart(),
		End:			s.GetEnd(),
		Resource:		s.GetResource(),
		IsTrySuccess:	false,
		Parent:			s,
		Child:			make([]*Slot, 0),
	}
	return s2
}