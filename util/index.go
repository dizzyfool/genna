package util

import (
	"fmt"
)

// Index stores unique strings
type Index struct {
	index map[string]struct{}
}

// NewIndex creates Index
func NewIndex() Index {
	return Index{
		index: map[string]struct{}{},
	}
}

// Available checks if string already exists
func (i *Index) Available(s string) bool {
	_, ok := i.index[s]
	return !ok
}

// Add ads string to index
func (i *Index) Add(s string) {
	i.index[s] = struct{}{}
}

// GetNext get next available string if given already exists
func (i *Index) GetNext(s string) string {
	if i.Available(s) {
		return s
	}

	suffix := 1
	for {
		next := fmt.Sprintf("%s%d", s, suffix)
		if i.Available(next) {
			return next
		}
		suffix++
	}
}
