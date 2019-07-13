package util

import (
	"testing"
)

func TestAnnotation_AddTag(t *testing.T) {
	type fields struct {
		tags []tag
	}
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name:   "Should add new tag",
			fields: fields{[]tag{}},
			args:   args{"tag", "value"},
			want:   1,
		},
		{
			name:   "Should append to existing tag",
			fields: fields{[]tag{{"tag", []string{"value1"}}}},
			args:   args{"tag", "value2"},
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAnnotation()
			a.tags = tt.fields.tags

			a.AddTag(tt.args.name, tt.args.value)
			if ln := len(a.tags); ln != tt.want {
				t.Errorf("Tags len = %v, want %v", ln, tt.want)
			}
		})
	}
}

func TestAnnotation_String(t *testing.T) {
	type fields struct {
		tags []tag
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should print one tag",
			fields: fields{[]tag{{"tag1", []string{"valueA"}}}},
			want:   `tag1:"valueA"`,
		},
		{
			name: "Should print several tags",
			fields: fields{[]tag{
				{"tag1", []string{"valueA"}},
				{"tag2", []string{"valueB"}},
				{"tag3", []string{"valueC"}},
			}},
			want: `tag1:"valueA" tag2:"valueB" tag3:"valueC"`,
		},
		{
			name: "Should print several tags with several values",
			fields: fields{[]tag{
				{"tag1", []string{"valueA"}},
				{"tag2", []string{"valueB", "valueC"}},
			}},
			want: `tag1:"valueA" tag2:"valueB,valueC"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Annotation{
				tags: tt.fields.tags,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("Annotation.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
