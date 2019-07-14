package util

// Set stores only unique strings
type Set struct {
	elements []string
	index    map[string]struct{}
}

// NewSet creates Set
func NewSet() Set {
	return Set{
		elements: []string{},
		index:    map[string]struct{}{},
	}
}

// Add adds element to set
// return false if element already exists
func (s *Set) Add(element string) bool {
	if s.Exists(element) {
		return false
	}

	s.elements = append(s.elements, element)
	s.index[element] = struct{}{}

	return true
}

// Exists checks if element exists
func (s *Set) Exists(element string) bool {
	_, ok := s.index[element]
	return ok
}

// Elements return all elements from set
func (s *Set) Elements() []string {
	return s.elements
}

// Len gets elements count
func (s *Set) Len() int {
	return len(s.elements)
}
