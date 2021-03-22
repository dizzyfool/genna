package model

import (
	"reflect"
	"testing"
)

func Test_parseCustomType(t *testing.T) {
	tests := []struct {
		name         string
		raw          string
		wantPgType   string
		wantGoType   string
		wantGoImport string
		wantErr      bool
	}{
		{
			name:         "should parse same package import",
			raw:          "uuid:Params",
			wantPgType:   "uuid",
			wantGoType:   "Params",
			wantGoImport: "",
			wantErr:      false,
		},
		{
			name:         "should parse std package import",
			raw:          "uuid:bytes.Buffer",
			wantPgType:   "uuid",
			wantGoType:   "bytes.Buffer",
			wantGoImport: "bytes",
		},
		{
			name:         "should parse std sub-package import",
			raw:          "uuid:encodings/json.RawMessage",
			wantPgType:   "uuid",
			wantGoType:   "json.RawMessage",
			wantGoImport: "encodings/json",
		},
		{
			name:         "should parse external import",
			raw:          "uuid:github.com/google/uuid.UUID",
			wantPgType:   "uuid",
			wantGoType:   "uuid.UUID",
			wantGoImport: "github.com/google/uuid",
		},
		{
			name:         "should parse external version import",
			raw:          "uuid:github.com/google/uuid/v4.UUID",
			wantPgType:   "uuid",
			wantGoType:   "uuid.UUID",
			wantGoImport: "github.com/google/uuid/v4",
		},
		{
			name:    "should error on incorrect format",
			raw:     "uuid;github.com/google/uuid.UUID",
			wantErr: true,
		},
		{
			name:    "should error on incorrect import format",
			raw:     "uuid:github.com/google/uuid/v4",
			wantErr: true,
		},
		{
			name:    "should error on incorrect import format v2",
			raw:     "uuid:github.com/google/uuid/v4.",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPgType, gotGoTyp, gotGoImport, err := parseCustomType(tt.raw)

			if err == nil == tt.wantErr {
				t.Errorf("parseCustomType() gotErr = %v", err)
			}
			if err != nil {
				return
			}

			if gotPgType != tt.wantPgType {
				t.Errorf("parseCustomType() gotPgType = %v, want %v", gotPgType, tt.wantPgType)
			}
			if gotGoTyp != tt.wantGoType {
				t.Errorf("parseCustomType() gotGoTyp = %v, want %v", gotGoTyp, tt.wantGoType)
			}
			if gotGoImport != tt.wantGoImport {
				t.Errorf("parseCustomType() gotGoImport = %v, want %v", gotGoImport, tt.wantGoImport)
			}
		})
	}
}

func TestParseCustomTypes(t *testing.T) {
	n := func(params ...string) CustomTypeMapping {
		ctm := CustomTypeMapping{}

		for i := 0; i < len(params); i += 3 {
			ctm.Add(params[i], params[i+1], params[i+2])
		}

		return ctm
	}

	tests := []struct {
		name    string
		args    []string
		want    CustomTypeMapping
		wantErr bool
	}{
		{
			name: "Should parse correct mapping",
			args: []string{"uuid:github.com/google/uuid.UUID", "point:src/db.Point"},
			want: n("uuid", "uuid.UUID", "github.com/google/uuid", "point", "db.Point", "src/db"),
		},
		{
			name:    "Should error on wrong format",
			args:    []string{"uuidgithub.com/google/uuid.UUID", "point:src/dbPoint."},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCustomTypes(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCustomTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCustomTypes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
