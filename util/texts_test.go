package util

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
			name: "Should get added non-countable",
			args: args{"sms"},
			want: "sms",
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

func TestEntityName(t *testing.T) {
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
			name: "Should generate from plural in last place",
			args: args{"usersWithOrders"},
			want: "UsersWithOrder",
		},
		{
			name: "Should generate from abracadabra",
			args: args{"abracadabra"},
			want: "Abracadabra",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EntityName(tt.args.input); got != tt.want {
				t.Errorf("EntityName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumnName(t *testing.T) {
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
			if got := ColumnName(tt.args.input); got != tt.want {
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
			want: "te_st",
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

func TestSanitize(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "should sanitize string contains special chars",
			s:    "te$t-Str1ng0§",
			want: "tet_Str1ng0",
		},
		{
			name: "should keep letters and numbers and dash",
			s:    "abcdef_12345-67890",
			want: "abcdef_12345_67890",
		},
		{
			name: "should add prefix if starting with number",
			s:    "1234abcdef",
			want: "T1234abcdef",
		},
		{
			name: "should add prefix if starting with number after sanitize",
			s:    "#1234abcdef",
			want: "T1234abcdef",
		},
		{
			name: "should add prefix if starting with dash",
			s:    "#-1234abcdef",
			want: "T_1234abcdef",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sanitize(tt.s); got != tt.want {
				t.Errorf("Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUpper(t *testing.T) {
	tests := []struct {
		name string
		c    byte
		want bool
	}{
		{
			name: "Should detect upper A",
			c:    'A',
			want: true,
		},
		{
			name: "Should detect upper Z",
			c:    'Z',
			want: true,
		},
		{
			name: "Should not detect lower z",
			c:    'z',
			want: false,
		},
		{
			name: "Should not detect 1",
			c:    '1',
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUpper(tt.c); got != tt.want {
				t.Errorf("IsUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLower(t *testing.T) {
	tests := []struct {
		name string
		c    byte
		want bool
	}{
		{
			name: "Should detect lower a",
			c:    'a',
			want: true,
		},
		{
			name: "Should detect lower z",
			c:    'z',
			want: true,
		},
		{
			name: "Should not detect upper Z",
			c:    'Z',
			want: false,
		},
		{
			name: "Should not detect 1",
			c:    '1',
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLower(tt.c); got != tt.want {
				t.Errorf("IsLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		name string
		c    byte
		want byte
	}{
		{
			name: "Should convert lower a to A",
			c:    'a',
			want: 'A',
		},
		{
			name: "Should convert lower z to Z",
			c:    'z',
			want: 'Z',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToUpper(tt.c); got != tt.want {
				t.Errorf("ToUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		name string
		c    byte
		want byte
	}{
		{
			name: "Should convert upper A to a",
			c:    'A',
			want: 'a',
		},
		{
			name: "Should convert upper Z to z",
			c:    'Z',
			want: 'z',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToLower(tt.c); got != tt.want {
				t.Errorf("ToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCamelCased(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "Should convert word to Word",
			s:    "word",
			want: "Word",
		},
		{
			name: "Should convert word_word to WordWord",
			s:    "word_word",
			want: "WordWord",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CamelCased(tt.s); got != tt.want {
				t.Errorf("CamelCased() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnderscore(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "Should convert Word to word",
			s:    "Word",
			want: "word",
		},
		{
			name: "Should convert WordWord to word_word",
			s:    "WordWord",
			want: "word_word",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Underscore(tt.s); got != tt.want {
				t.Errorf("Underscore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLowerFirst(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "Should convert Word to word",
			s:    "Word",
			want: "word",
		},
		{
			name: "Should convert WordWord to wordWord",
			s:    "WordWord",
			want: "wordWord",
		},
		{
			name: "Should convert 1WordWord to 1WordWord",
			s:    "1WordWord",
			want: "1WordWord",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LowerFirst(tt.s); got != tt.want {
				t.Errorf("LowerFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}
