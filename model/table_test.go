package model

import (
	"reflect"
	"sort"
	"testing"
)

func TestTable_ModelName(t *testing.T) {
	type fields struct {
		Name   string
		Schema string
	}
	tests := []struct {
		name       string
		fields     fields
		withSchema bool
		want       string
	}{
		{
			name:   "Should generate from simple word",
			fields: fields{Name: "users"},
			want:   "User",
		},
		{
			name:   "Should generate from non-countable",
			fields: fields{Name: "audio"},
			want:   "Audio",
		},
		{
			name:   "Should generate from underscored",
			fields: fields{Name: "user_orders"},
			want:   "UserOrder",
		},
		{
			name:   "Should generate from camelCased",
			fields: fields{Name: "userOrders"},
			want:   "UserOrder",
		},
		{
			name:   "Should generate from plural in last place",
			fields: fields{Name: "usersWithOrders"},
			want:   "UsersWithOrder",
		},
		{
			name:   "Should generate from abracadabra",
			fields: fields{Name: "abracadabra"},
			want:   "Abracadabra",
		},
		{
			name:       "Should generate from simple word with public schema",
			fields:     fields{Name: "users", Schema: "public"},
			withSchema: true,
			want:       "User",
		},
		{
			name:       "Should generate from simple word with custom schema",
			fields:     fields{Name: "users", Schema: "users"},
			withSchema: true,
			want:       "UsersUser",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Name:   tt.fields.Name,
				Schema: tt.fields.Schema,
			}
			if got := tbl.ModelName(); got != tt.want {
				t.Errorf("Table.ModelName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_TableName(t *testing.T) {
	type fields struct {
		Schema string
		Name   string
	}
	tests := []struct {
		name   string
		fields fields
		quoted bool
		want   string
	}{
		{
			name:   "Should generate from public schema and simple table name, ignore quotes",
			fields: fields{PublicSchema, "users"},
			quoted: true,
			want:   "users",
		},
		{
			name:   "Should generate from non-public schema and simple table name, ignore quotes",
			fields: fields{"users", "users"},
			quoted: true,
			want:   "users.users",
		},
		{
			name:   "Should generate quoted and escaped from public schema and table name",
			fields: fields{PublicSchema, "userOrders"},
			quoted: true,
			want:   `\"userOrders\"`,
		},

		{
			name:   "Should generate quoted and escaped",
			fields: fields{"allUsers", "userOrders"},
			quoted: true,
			want:   `\"allUsers\".\"userOrders\"`,
		},
		{
			name:   "Should generate from public schema and simple table name",
			fields: fields{PublicSchema, "users"},
			want:   "users",
		},
		{
			name:   "Should generate from non-public schema and simple table name",
			fields: fields{"users", "users"},
			want:   "users.users",
		},
		{
			name:   "Should generate quoted and escaped from public schema and table name, no quotes",
			fields: fields{PublicSchema, "userOrders"},
			want:   `userOrders`,
		},

		{
			name:   "Should generate quoted and escaped, no quotes",
			fields: fields{"allUsers", "userOrders"},
			want:   `allUsers.userOrders`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema: tt.fields.Schema,
				Name:   tt.fields.Name,
			}
			if got := tbl.TableName(tt.quoted); got != tt.want {
				t.Errorf("Table.TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_ViewName(t *testing.T) {
	type fields struct {
		Schema string
		Name   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should generate from public schema and simple table name",
			fields: fields{PublicSchema, "users"},
			want:   `\"getUsers\"`,
		},
		{
			name:   "Should generate from non-public schema and simple table name",
			fields: fields{"users", "users"},
			want:   `users.\"getUsers\"`,
		},
		{
			name:   "Should generate quoted and escaped",
			fields: fields{"allUsers", "users"},
			want:   `\"allUsers\".\"getUsers\"`,
		},
		{
			name:   "Should generate from underscored",
			fields: fields{PublicSchema, "users_orders"},
			want:   `\"getUsersOrders\"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema: tt.fields.Schema,
				Name:   tt.fields.Name,
			}
			if got := tbl.ViewName(); got != tt.want {
				t.Errorf("Table.ViewName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_TableNameTag(t *testing.T) {
	type fields struct {
		Schema string
		Name   string
	}
	type args struct {
		withView  bool
		noDiscard bool
		noAlias   bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Should generate with default params",
			fields: fields{PublicSchema, "users"},
			args:   args{false, false, false},
			want:   `sql:"users,alias:t" pg:",discard_unknown_columns"`,
		},
		{
			name:   "Should generate without alias",
			fields: fields{PublicSchema, "users"},
			args:   args{false, false, true},
			want:   `sql:"users" pg:",discard_unknown_columns"`,
		},
		{
			name:   "Should generate with view",
			fields: fields{PublicSchema, "users"},
			args:   args{true, false, false},
			want:   `sql:"users,select:\"getUsers\",alias:t" pg:",discard_unknown_columns"`,
		},
		{
			name:   "Should generate with no discard and alias",
			fields: fields{PublicSchema, "users"},
			args:   args{false, true, true},
			want:   `sql:"users"`,
		},
		{
			name:   "Should generate with no discard and view and alias",
			fields: fields{PublicSchema, "users"},
			args:   args{true, true, true},
			want:   `sql:"users,select:\"getUsers\""`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema: tt.fields.Schema,
				Name:   tt.fields.Name,
			}
			if got := tbl.TableNameTag(tt.args.withView, tt.args.noDiscard, tt.args.noAlias); got != tt.want {
				t.Errorf("Table.TableNameTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_Validate(t *testing.T) {
	pkColumn := Column{
		Name: "userId",
		Type: TypeInt8,
		IsPK: true,
	}

	fkColumn := Column{
		Name: "locationId",
		Type: TypeInt8,
		IsFK: true,
	}

	invalidColumn := Column{
		Name: "locationId",
		Type: "unknown",
		IsFK: true,
	}

	validRelation := Relation{
		Type: HasOne,
		// other doesn't matter for now
	}

	type fields struct {
		Schema    string
		Name      string
		Columns   []Column
		Relations []Relation
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{

		{
			name: "Should not raise error on valid table",
			fields: fields{
				Schema:    "valid",
				Name:      "valid",
				Columns:   []Column{pkColumn, fkColumn},
				Relations: []Relation{validRelation},
			},
			wantErr: false,
		},
		{
			name: "Should raise error on empty name",
			fields: fields{
				Schema:  "valid",
				Name:    " ",
				Columns: []Column{pkColumn},
			},
			wantErr: true,
		},
		{
			name: "Should raise error on empty schema",
			fields: fields{
				Schema:  " ",
				Name:    "valid",
				Columns: []Column{pkColumn},
			},
			wantErr: true,
		},
		{
			name: "Should raise error on invalid name",
			fields: fields{
				Schema:  "valid",
				Name:    "#test",
				Columns: []Column{pkColumn},
			},
			wantErr: true,
		},
		{
			name: "Should raise error on invalid schema",
			fields: fields{
				Schema:  "#test",
				Name:    "valid",
				Columns: []Column{pkColumn},
			},
			wantErr: true,
		},
		{
			name: "Should raise error on empty columns",
			fields: fields{
				Schema: "valid",
				Name:   "valid",
			},
			wantErr: true,
		},
		{
			name: "Should raise error on invalid columns",
			fields: fields{
				Schema:  "valid",
				Name:    "valid",
				Columns: []Column{invalidColumn},
			},
			wantErr: true,
		},
		{
			name: "Should raise error on empty relations with fkey",
			fields: fields{
				Schema:  "valid",
				Name:    "valid",
				Columns: []Column{fkColumn},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema:    tt.fields.Schema,
				Name:      tt.fields.Name,
				Columns:   tt.fields.Columns,
				Relations: tt.fields.Relations,
			}
			if err := tbl.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Table.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTable_Imports(t *testing.T) {
	type fields struct {
		Schema    string
		Name      string
		Columns   []Column
		Relations []Relation
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Should not generate imports if only simple types",
			fields: fields{
				Columns: []Column{
					{
						Name: "userId",
						Type: TypeInt8,
					},
					{
						Name: "locationId",
						Type: TypeInt8,
					},
				},
			},
			want: []string{},
		},
		{
			name: "Should not generate imports without duplicates",
			fields: fields{
				Columns: []Column{
					{
						Name: "userId",
						Type: TypeInt8,
					},
					{
						Name: "createdAt",
						Type: TypeTimestamp,
					},
					{
						Name:       "deletedAt",
						Type:       TypeTimestamp,
						IsNullable: true,
					},
					{
						Name:       "updatedAt",
						Type:       TypeTimestamp,
						IsNullable: true,
					},
				},
			},
			want: []string{"time", "github.com/go-pg/pg"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema:    tt.fields.Schema,
				Name:      tt.fields.Name,
				Columns:   tt.fields.Columns,
				Relations: tt.fields.Relations,
			}
			got := tbl.Imports()

			sort.Strings(got)
			sort.Strings(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.Imports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_HasMultiplePKs(t *testing.T) {
	type fields struct {
		Schema    string
		Name      string
		Columns   []Column
		Relations []Relation
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Should detect multiple PKs",
			fields: fields{
				Columns: []Column{
					{
						Name: "userId",
						Type: TypeInt8,
						IsPK: true,
					},
					{
						Name: "locationId",
						Type: TypeInt8,
						IsPK: true,
					},
				},
			},
			want: true,
		},
		{
			name: "Should detect single PKs",
			fields: fields{
				Columns: []Column{
					{
						Name: "userId",
						Type: TypeInt8,
						IsPK: true,
					},
					{
						Name: "locationId",
						Type: TypeInt8,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema:    tt.fields.Schema,
				Name:      tt.fields.Name,
				Columns:   tt.fields.Columns,
				Relations: tt.fields.Relations,
			}
			if got := tbl.HasMultiplePKs(); got != tt.want {
				t.Errorf("Table.HasMultiplePKs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_JoinAlias(t *testing.T) {
	type fields struct {
		Schema    string
		Name      string
		Columns   []Column
		Relations []Relation
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should generate from simple word",
			fields: fields{Name: "users"},
			want:   "user",
		},
		{
			name:   "Should generate from non-countable",
			fields: fields{Name: "audio"},
			want:   "audio",
		},
		{
			name:   "Should generate from underscored",
			fields: fields{Name: "user_orders"},
			want:   "user_order",
		},
		{
			name:   "Should generate from camelCased",
			fields: fields{Name: "userOrders"},
			want:   "user_order",
		},
		{
			name:   "Should generate from plural in last place",
			fields: fields{Name: "usersWithOrders"},
			want:   "users_with_order",
		},
		{
			name:   "Should generate from abracadabra",
			fields: fields{Name: "abracadabra"},
			want:   "abracadabra",
		},
		{
			name:   "Should generate from simple word with public schema",
			fields: fields{Name: "users", Schema: "public"},
			want:   "user",
		},
		{
			name:   "Should generate from simple word with custom schema",
			fields: fields{Name: "users", Schema: "users"},
			want:   "users_user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema:    tt.fields.Schema,
				Name:      tt.fields.Name,
				Columns:   tt.fields.Columns,
				Relations: tt.fields.Relations,
			}
			if got := tbl.JoinAlias(); got != tt.want {
				t.Errorf("Table.JoinAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_SearchModelName(t *testing.T) {
	type fields struct {
		Schema    string
		Name      string
		Columns   []Column
		Relations []Relation
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should generate from simple word",
			fields: fields{Name: "users"},
			want:   "UserSearch",
		},
		{
			name:   "Should generate from non-countable",
			fields: fields{Name: "audio"},
			want:   "AudioSearch",
		},
		{
			name:   "Should generate from underscored",
			fields: fields{Name: "user_orders"},
			want:   "UserOrderSearch",
		},
		{
			name:   "Should generate from camelCased",
			fields: fields{Name: "userOrders"},
			want:   "UserOrderSearch",
		},
		{
			name:   "Should generate from plural in last place",
			fields: fields{Name: "usersWithOrders"},
			want:   "UsersWithOrderSearch",
		},
		{
			name:   "Should generate from abracadabra",
			fields: fields{Name: "abracadabra"},
			want:   "AbracadabraSearch",
		},
		{
			name:   "Should generate from simple word with public schema",
			fields: fields{Name: "users", Schema: "public"},
			want:   "UserSearch",
		},
		{
			name:   "Should generate from simple word with custom schema",
			fields: fields{Name: "users", Schema: "users"},
			want:   "UsersUserSearch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema:    tt.fields.Schema,
				Name:      tt.fields.Name,
				Columns:   tt.fields.Columns,
				Relations: tt.fields.Relations,
			}
			if got := tbl.SearchModelName(); got != tt.want {
				t.Errorf("Table.SearchModelName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_SearchImports(t *testing.T) {
	type fields struct {
		Schema    string
		Name      string
		Columns   []Column
		Relations []Relation
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Should not generate imports if only simple types",
			fields: fields{
				Columns: []Column{
					{
						Name: "userId",
						Type: TypeInt8,
					},
					{
						Name: "locationId",
						Type: TypeInt8,
					},
				},
			},
			want: []string{},
		},
		{
			name: "Should not generate imports without duplicates",
			fields: fields{
				Columns: []Column{
					{
						Name: "userId",
						Type: TypeInt8,
					},
					{
						Name: "createdAt",
						Type: TypeTimestamp,
					},
					{
						Name:       "deletedAt",
						Type:       TypeTimestamp,
						IsNullable: true,
					},
					{
						Name:       "updatedAt",
						Type:       TypeTimestamp,
						IsNullable: true,
					},
				},
			},
			want: []string{"time"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema:    tt.fields.Schema,
				Name:      tt.fields.Name,
				Columns:   tt.fields.Columns,
				Relations: tt.fields.Relations,
			}
			if got := tbl.SearchImports(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.SearchImports() = %v, want %v", got, tt.want)
			}
		})
	}
}
