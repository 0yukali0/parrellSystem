package object

type Slot struct {
	Start				uint64
	End					uint64
	Resource			uint64

	Index				int 

	IsTrySuccess		bool
	Parent				*Slot
	Child				[]*Slot
}

func (s *Slot) SetStart(start uint64) {
	s.Start = start
}

func (s *Slot) GetStart() uint64 {
	return s.Start
}

func (s *Slot) SetEnd(end uint64) {
	s.End = end
}

func (s *Slot) GetEnd() uint64 {
	return s.End
}

func (s *Slot) AllocateResource(allocation uint64) {
	s.Resource -= allocation
}

func (s *Slot) ReleaseResource(allocation uint64) {
	s.Resource += allocation
}

func (s *Slot) GetResource() uint64 {
	return s.Resource
}

func (s *Slot) SetIsTrySuccess(ok bool) {
	s.IsTrySuccess = ok
}

func (s *Slot) GetIsTrySuccess() bool {
	return s.IsTrySuccess
}

func (s *Slot) IsParent() bool {
	return len(s.GetChildren()) > 0
}

func (s *Slot) AddChildren(children []*Slot) {
	for _, child := range children {
		child.Parent = s
		s.Child = append(s.Child, child)
	}
}

func (s *Slot) GetChildren() []*Slot {
	return s.Child
}
