package util

import (
	"fmt"
	"strings"
)

// Annotation is a simple helper used to build tags for structs
type Annotation struct {
	tags []tag
}

type tag struct {
	name   string
	values []string
}

// NewAnnotation creates annotation
func NewAnnotation() *Annotation {
	return &Annotation{}
}

// AddTag ads a tag if not exists, appends a value otherwise
func (a *Annotation) AddTag(name string, value string) *Annotation {
	for i, tag := range a.tags {
		if tag.name == name {
			a.tags[i].values = append(a.tags[i].values, value)
			return a
		}
	}
	a.tags = append(a.tags, tag{name, []string{value}})
	return a
}

func (a *Annotation) Len() int {
	return len(a.tags)
}

// String prints valid tag
func (a *Annotation) String() string {
	result := make([]string, 0)
	for _, tag := range a.tags {
		result = append(result, fmt.Sprintf(`%s:"%s"`, tag.name, strings.Join(tag.values, ",")))
	}

	return strings.Join(result, " ")
}
