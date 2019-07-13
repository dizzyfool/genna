package util

import (
	"testing"
)

func TestIndex_Available(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
		want bool
	}{
		{
			name: "Should check for available",
			old:  "old",
			new:  "new",
			want: true,
		},
		{
			name: "Should check for not available",
			old:  "old",
			new:  "old",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewIndex()
			i.Add(tt.old)
			if got := i.Available(tt.new); got != tt.want {
				t.Errorf("Index.Available() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_GetNext(t *testing.T) {
	tests := []struct {
		name  string
		index map[string]struct{}
		s     string
		want  string
	}{
		{
			name:  "Should get same if available",
			index: map[string]struct{}{"old": {}},
			s:     "new",
			want:  "new",
		},
		{
			name:  "Should get next",
			index: map[string]struct{}{"old": {}},
			s:     "old",
			want:  "old1",
		},
		{
			name:  "Should get second",
			index: map[string]struct{}{"old": {}, "old1": {}},
			s:     "old",
			want:  "old2",
		},
		{
			name:  "Should check for existing",
			index: map[string]struct{}{"old": {}, "old1": {}},
			s:     "old1",
			want:  "old11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewIndex()
			i.index = tt.index
			if got := i.GetNext(tt.s); got != tt.want {
				t.Errorf("Index.GetNext() = %v, want %v", got, tt.want)
			}
		})
	}
}
