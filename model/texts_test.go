package model

import (
	"testing"
)

func TestSingular(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should get normal singular",
			args: args{"dogs"},
			want: "dog",
		},
		{
			name: "Should get irregular singular",
			args: args{"children"},
			want: "child",
		},
		{
			name: "Should get non-countable",
			args: args{"fish"},
			want: "fish",
		},
		{
			name: "Should ignore non plural",
			args: args{"test"},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Singular(tt.args.input); got != tt.want {
				t.Errorf("Singular() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelName(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should generate from simple word",
			args: args{"users"},
			want: "User",
		},
		{
			name: "Should generate from simple word end with es",
			args: args{"companies"},
			want: "Company",
		},
		{
			name: "Should generate from simple word end with es",
			args: args{"glasses"},
			want: "Glass",
		},
		{
			name: "Should generate from non-countable",
			args: args{"audio"},
			want: "Audio",
		},
		{
			name: "Should generate from underscored",
			args: args{"user_orders"},
			want: "UserOrder",
		},
		{
			name: "Should generate from camelCased",
			args: args{"userOrders"},
			want: "UserOrder",
		},
		{
			name: "Should generate from plural in first place",
			args: args{"usersWithOrders"},
			want: "UserWithOrders",
		},
		{
			name: "Should generate from plural in last place",
			args: args{"usersWithOrders"},
			want: "UserWithOrders",
		},
		{
			name: "Should generate from abracadabra",
			args: args{"abracadabra"},
			want: "Abracadabra",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ModelName(tt.args.input); got != tt.want {
				t.Errorf("ModelName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructFieldName(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should generate from simple word",
			args: args{"title"},
			want: "Title",
		},
		{
			name: "Should generate from underscored",
			args: args{"short_title"},
			want: "ShortTitle",
		},
		{
			name: "Should generate from camelCased",
			args: args{"shortTitle"},
			want: "ShortTitle",
		},
		{
			name: "Should generate with underscored id",
			args: args{"location_id"},
			want: "LocationID",
		},
		{
			name: "Should generate with camelCased id",
			args: args{"locationId"},
			want: "LocationID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StructFieldName(tt.args.input); got != tt.want {
				t.Errorf("ColumnName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasUpper(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should detect upper case",
			args: args{"upPer"},
			want: true,
		},
		{
			name: "Should not detect only lower case",
			args: args{"lower"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasUpper(tt.args.input); got != tt.want {
				t.Errorf("HasUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceSuffix(t *testing.T) {
	type args struct {
		input   string
		suffix  string
		replace string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should replace suffix",
			args: args{"locationId", "Id", "ID"},
			want: "locationID",
		},
		{
			name: "Should not replace if not found",
			args: args{"location", "Id", "ID"},
			want: "location",
		},
		{
			name: "Should not replace if not suffix",
			args: args{"locationIdHere", "Id", "ID"},
			want: "locationIdHere",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceSuffix(tt.args.input, tt.args.suffix, tt.args.replace); got != tt.want {
				t.Errorf("ReplaceSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageName(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should generate valid package name with lower",
			args: args{"tesT"},
			want: "test",
		},
		{
			name: "Should generate valid package name with only letters",
			args: args{"te_sT$"},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PackageName(tt.args.input); got != tt.want {
				t.Errorf("PackageName() = %v, want %v", got, tt.want)
			}
		})
	}
}
