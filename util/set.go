package util

// Set stores only unique strings
type Set struct {
	elements []string
	index    map[string]struct{}
}

func NewSet() Set {
	return Set{
		elements: []string{},
		index:    map[string]struct{}{},
	}
}

func (s *Set) Add(element string) bool {
	if s.Exists(element) {
		return false
	}

	s.elements = append(s.elements, element)
	s.index[element] = struct{}{}

	return true
}

func (s *Set) Exists(element string) bool {
	_, ok := s.index[element]
	return ok
}

func (s *Set) Elements() []string {
	return s.elements
}

func (s *Set) Len() int {
	return len(s.elements)
}
