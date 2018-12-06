package database

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
			wantSchemas: []string{"public"},
		},
		{
			name: "Should get public schema from full table names",
			args: args{
				[]string{"users", "public.locations"},
			},
			wantSchemas: []string{"public"},
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
