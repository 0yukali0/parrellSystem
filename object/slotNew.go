package object

func NewRootSlot(start, end, cap uint64) *Slot {
	s := &Slot{
		Start:			start,
		End:			end,
		Resource:		cap,
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