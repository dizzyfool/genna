package model

import (
	"reflect"
	"sort"
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

func TestDiscloseSchemas(t *testing.T) {
	type args struct {
		tables []string
		models []Table
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Should disclose *",
			args: args{
				tables: []string{"*"},
				models: []Table{
					{Schema: "public", Name: "users"},
					{Schema: "public", Name: "locations"},
				},
			},
			want: []string{"public.users", "public.locations"},
		},
		{
			name: "Should disclose users.*",
			args: args{
				tables: []string{"users.*"},
				models: []Table{
					{Schema: "users", Name: "users"},
					{Schema: "users", Name: "locations"},
				},
			},
			want: []string{"users.users", "users.locations"},
		},
		{
			name: "Should disclose [users.*, user.users]",
			args: args{
				tables: []string{"users.*", "users.users"},
				models: []Table{
					{Schema: "users", Name: "users"},
					{Schema: "users", Name: "locations"},
				},
			},
			want: []string{"users.users", "users.locations"},
		},

		{
			name: "Should disclose [user.locations, public.*]",
			args: args{
				tables: []string{"users.locations", "public.*"},
				models: []Table{
					{Schema: "public", Name: "users"},
					{Schema: "public", Name: "locations"},
					{Schema: "users", Name: "locations"},
					{Schema: "users", Name: "tests"},
				},
			},
			want: []string{"public.users", "public.locations", "users.locations"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DiscloseSchemas(tt.args.tables, tt.args.models)

			sort.Strings(got)
			sort.Strings(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiscloseSchema() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestFollowFKs(t *testing.T) {
	locationRel := Relation{
		TargetSchema: "public",
		TargetTable:  "locations",
	}
	userRel := Relation{
		TargetSchema: "users",
		TargetTable:  "users",
	}

	type args struct {
		tables []string
		models []Table
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Should follow foreign key",
			args: args{
				tables: []string{"public.locations"},
				models: []Table{
					{
						Schema:    PublicSchema,
						Name:      "locations",
						Relations: []Relation{userRel},
					},
					{
						Schema: "users",
						Name:   "users",
					},
					{
						Schema: "users",
						Name:   "locations",
					},
				},
			},
			want: []string{"public.locations", "users.users"},
		},
		{
			name: "Should follow self recursive foreign key",
			args: args{
				tables: []string{"public.locations"},
				models: []Table{
					{
						Schema:    PublicSchema,
						Name:      "locations",
						Relations: []Relation{locationRel},
					},
					{
						Schema: "users",
						Name:   "users",
					},
					{
						Schema: "users",
						Name:   "locations",
					},
				},
			},
			want: []string{"public.locations"},
		},
		{
			name: "Should follow recursive foreign key",
			args: args{
				tables: []string{"public.locations"},
				models: []Table{
					{
						Schema:    PublicSchema,
						Name:      "locations",
						Relations: []Relation{userRel},
					},
					{
						Schema:    "users",
						Name:      "users",
						Relations: []Relation{locationRel},
					},
					{
						Schema: "users",
						Name:   "locations",
					},
				},
			},
			want: []string{"public.locations", "users.users"},
		},
		{
			name: "Should follow recursive foreign key",
			args: args{
				tables: []string{"public.locations"},
				models: []Table{
					{
						Schema:    PublicSchema,
						Name:      "locations",
						Relations: []Relation{userRel},
					},
					{
						Schema:    "users",
						Name:      "users",
						Relations: []Relation{locationRel},
					},
					{
						Schema: "users",
						Name:   "locations",
					},
				},
			},
			want: []string{"public.locations", "users.users"},
		},
		{
			name: "Should follow recursive foreign key without doubles",
			args: args{
				tables: []string{"public.locations", "users.users"},
				models: []Table{
					{
						Schema: PublicSchema,
						Name:   "locations",
					},
					{
						Schema:    "users",
						Name:      "users",
						Relations: []Relation{locationRel},
					},
					{
						Schema:    "users",
						Name:      "locations",
						Relations: []Relation{locationRel},
					},
				},
			},
			want: []string{"public.locations", "users.users"},
		},
		{
			name: "Should follow recursive deep foreign key",
			args: args{
				tables: []string{"users.locations"},
				models: []Table{
					{
						Schema:    PublicSchema,
						Name:      "locations",
						Relations: []Relation{userRel},
					},
					{
						Schema: PublicSchema,
						Name:   "ignored",
					},
					{
						Schema: "users",
						Name:   "users",
					},
					{
						Schema:    "users",
						Name:      "locations",
						Relations: []Relation{locationRel},
					},
				},
			},
			want: []string{"public.locations", "users.users", "users.locations"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FollowFKs(tt.args.tables, tt.args.models)

			sort.Strings(got)
			sort.Strings(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FollowFKs() = %v, want %v", got, tt.want)
			}
		})
	}
}
