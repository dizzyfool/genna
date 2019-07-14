package util

import (
	"reflect"
	"testing"
)

func TestSchemas(t *testing.T) {
	type args struct {
		tables []string
	}
	tests := []struct {
		name        string
		args        args
		wantSchemas []string
	}{
		{
			name: "Should get public schema",
			args: args{
				[]string{"users", "locations"},
			},
			wantSchemas: []string{PublicSchema},
		},
		{
			name: "Should get public schema from full table names",
			args: args{
				[]string{"users", "public.locations"},
			},
			wantSchemas: []string{PublicSchema},
		},
		{
			name: "Should get different schemas from full table names",
			args: args{
				[]string{"users.users", "users.locations", "orders.orders"},
			},
			wantSchemas: []string{"users", "orders"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSchemas := Schemas(tt.args.tables); !reflect.DeepEqual(gotSchemas, tt.wantSchemas) {
				t.Errorf("Schemas() = %v, want %v", gotSchemas, tt.wantSchemas)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name:  "Should split full name",
			args:  args{"public.users"},
			want:  "public",
			want1: "users",
		},
		{
			name:  "Should split simple name",
			args:  args{"users"},
			want:  PublicSchema,
			want1: "users",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Split(tt.args.input)
			if got != tt.want {
				t.Errorf("Split() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Split() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	type args struct {
		schema string
		table  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should join",
			args: args{"public", "users"},
			want: "public.users",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.args.schema, tt.args.table); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}
