package util

import (
	"testing"
)

func TestSet_Add(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
		want bool
	}{
		{
			name: "Should add new element",
			old:  "old",
			new:  "new",
			want: true,
		},
		{
			name: "Should not add old element",
			old:  "old",
			new:  "old",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet()
			s.Add(tt.old)
			if got := s.Add(tt.new); got != tt.want {
				t.Errorf("Set.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Elements(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
		want []string
	}{
		{
			name: "Should get 2 elements",
			old:  "old",
			new:  "new",
			want: []string{"old", "new"},
		},
		{
			name: "Should get 1 element",
			old:  "old",
			new:  "old",
			want: []string{"old"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet()
			s.Add(tt.old)
			s.Add(tt.new)
			if got := s.Elements(); len(got) != len(tt.want) {
				t.Errorf("Set.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
