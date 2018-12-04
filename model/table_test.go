package model

import "testing"

func TestTable_ModelName(t *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should generate from simple word",
			fields: fields{"users"},
			want:   "User",
		},
		{
			name:   "Should generate from non-countable",
			fields: fields{"audio"},
			want:   "Audio",
		},
		{
			name:   "Should generate from underscored",
			fields: fields{"user_orders"},
			want:   "UserOrder",
		},
		{
			name:   "Should generate from camelCased",
			fields: fields{"userOrders"},
			want:   "UserOrder",
		},
		{
			name:   "Should generate from plural in first place",
			fields: fields{"usersWithOrders"},
			want:   "UserWithOrders",
		},
		{
			name:   "Should generate from plural in last place",
			fields: fields{"usersWithOrders"},
			want:   "UserWithOrders",
		},
		{
			name:   "Should generate from abracadabra",
			fields: fields{"abracadabra"},
			want:   "Abracadabra",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Name: tt.fields.Name,
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
		want   string
	}{
		{
			name:   "Should generate from public schema and simple table name",
			fields: fields{"public", "users"},
			want:   "users",
		},
		{
			name:   "Should generate from non-public schema and simple table name",
			fields: fields{"users", "users"},
			want:   "users.users",
		},
		{
			name:   "Should generate quoted and escaped from public schema and table name",
			fields: fields{"public", "userOrders"},
			want:   `\"userOrders\"`,
		},

		{
			name:   "Should generate quoted and escaped",
			fields: fields{"allUsers", "userOrders"},
			want:   `\"allUsers\".\"userOrders\"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema: tt.fields.Schema,
				Name:   tt.fields.Name,
			}
			if got := tbl.TableName(); got != tt.want {
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
			fields: fields{"public", "users"},
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
			fields: fields{"public", "users_orders"},
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
		noDiscard bool
		withView  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Should generate with default params",
			fields: fields{"public", "users"},
			args:   args{false, false},
			want:   `sql:"users",pg:",discard_unknown_columns"`,
		},
		{
			name:   "Should generate with view",
			fields: fields{"public", "users"},
			args:   args{false, true},
			want:   `sql:"users,select:\"getUsers\"",pg:",discard_unknown_columns"`,
		},
		{
			name:   "Should generate with no discard",
			fields: fields{"public", "users"},
			args:   args{true, false},
			want:   `sql:"users"`,
		},
		{
			name:   "Should generate with no discard and view",
			fields: fields{"public", "users"},
			args:   args{true, true},
			want:   `sql:"users,select:\"getUsers\""`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := Table{
				Schema: tt.fields.Schema,
				Name:   tt.fields.Name,
			}
			if got := tbl.TableNameTag(tt.args.noDiscard, tt.args.withView); got != tt.want {
				t.Errorf("Table.TableNameTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
