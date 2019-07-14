# Genna - cli tool for generating go-pg models

[![Go Report Card](https://goreportcard.com/badge/github.com/dizzyfool/genna)](https://goreportcard.com/report/github.com/dizzyfool/genna)

Requirements:
- [go-pg](https://github.com/go-pg/pg)
- your PostgreSQL database

#### Idea

In most of the cases go-pg models represent database's tables and relations. Genna's main goal is to prepare those models by reading detailed information about PostrgeSQL database. The result should be several files with ready to use structs.

#### Usage

1. Install `go get github.com/dizzyfool/genna`
1. Read though help `genna -h`
1. Run `genna -c postgres://user:password@localhost:5432/yourdb -o ~/output/model.go`

#### Example

Create your database and tables in it

```sql
create table "users" (
    "userId"    serial      not null,
    "email"     varchar(64) not null,
    "activated" bool        not null default false,
    "name"      varchar(128),
    "countryId" integer,

    primary key ("userId")
);

create schema "geo";
create table geo."countries" (
    "countryId" serial     not null,
    "code"      varchar(3) not null,
    "coords"    integer[],

    primary key ("countryId")
);

alter table "users"
    add constraint "fk_user_country"
        foreign key ("countryId")
            references geo."countries" ("countryId") on update restrict on delete restrict;
```

Run generator

`genna -c postgres://user:password@localhost:5432/yourdb -o ~/output/model.go -t public.*,geo.* -f -s`

You should get following models on model package:

```go
//lint:file-ignore U1000 ignore unused code, it's generated
package model

var Columns = struct {
	User struct {
		ID, Activated, CountryID, Email, Name string

		Country string
	}
	GeoCountry struct {
		ID, Code, Coords string
	}
}{
	User: struct {
		ID, Activated, CountryID, Email, Name string

		Country string
	}{
		ID:        "userId",
		Activated: "activated",
		CountryID: "countryId",
		Email:     "email",
		Name:      "name",

		Country: "Country",
	},
	GeoCountry: struct {
		ID, Code, Coords string
	}{
		ID:     "countryId",
		Code:   "code",
		Coords: "coords",
	},
}

var Tables = struct {
	User struct {
		Name, Alias string
	}
	GeoCountry struct {
		Name, Alias string
	}
}{
	User: struct {
		Name, Alias string
	}{
		Name:  "users",
		Alias: "t",
	},
	GeoCountry: struct {
		Name, Alias string
	}{
		Name:  "geo.countries",
		Alias: "t",
	},
}

type User struct {
	tableName struct{} `sql:"users,alias:t" pg:",discard_unknown_columns"`

	ID        int     `sql:"userId,pk"`
	Activated bool    `sql:"activated,notnull"`
	CountryID *int    `sql:"countryId"`
	Email     string  `sql:"email,notnull"`
	Name      *string `sql:"name"`

	Country *GeoCountry `pg:"fk:countryId"`
}

type GeoCountry struct {
	tableName struct{} `sql:"geo.countries,alias:t" pg:",discard_unknown_columns"`

	ID     int    `sql:"countryId,pk"`
	Code   string `sql:"code,notnull"`
	Coords []int  `sql:"coords,array"`
}
```

If you choose to generate search filters another file will appear

```go
//lint:file-ignore U1000 ignore unused code, it's generated
package model

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// base filters

type applier func(query *orm.Query) (*orm.Query, error)

type search struct {
	custom map[string][]interface{}
}

func (s *search) apply(table string, values map[string]interface{}, query *orm.Query) *orm.Query {
	for field, value := range values {
		if value != nil {
			query.Where("?.? = ?", pg.F(table), pg.F(field), value)
		}
	}

	if s.custom != nil {
		for condition, params := range s.custom {
			query.Where(condition, params...)
		}
	}

	return query
}

func (s *search) with(condition string, params ...interface{}) {
	if s.custom == nil {
		s.custom = map[string][]interface{}{}
	}
	s.custom[condition] = params
}

// Searcher is interface for every generated filter
type Searcher interface {
	Apply(query *orm.Query) *orm.Query
	Q() applier
}

type UserSearch struct {
	search

	ID        interface{}
	Activated interface{}
	CountryID interface{}
	Email     interface{}
	Name      interface{}
}

func (s *UserSearch) Apply(query *orm.Query) *orm.Query {
	return s.apply(Tables.User.Alias, map[string]interface{}{
		Columns.User.ID:        s.ID,
		Columns.User.Activated: s.Activated,
		Columns.User.CountryID: s.CountryID,
		Columns.User.Email:     s.Email,
		Columns.User.Name:      s.Name,
	}, query)
}

func (s *UserSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type GeoCountrySearch struct {
	search

	ID   interface{}
	Code interface{}
}

func (s *GeoCountrySearch) Apply(query *orm.Query) *orm.Query {
	return s.apply(Tables.GeoCountry.Alias, map[string]interface{}{
		Columns.GeoCountry.ID:   s.ID,
		Columns.GeoCountry.Code: s.Code,
	}, query)
}

func (s *GeoCountrySearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}
``` 

Try it

```go
package model

import (
	"fmt"
	"testing"

	"github.com/go-pg/pg"
)

const AllColumns = "t.*"

func TestModel(t *testing.T) {
	// connecting to db
	options, _ := pg.ParseURL("postgres://user:password@localhost:5432/yourdb")
	db := pg.Connect(options)

	if _, err := db.Exec(`truncate table users; truncate table geo.countries cascade;`); err != nil {
		panic(err)
	}

	// objects to insert
	toInsert := []GeoCountry{
		GeoCountry{
			Code:   "us",
			Coords: []int{1, 2},
		},
		GeoCountry{
			Code:   "uk",
			Coords: nil,
		},
	}

	// inserting
	if _, err := db.Model(&toInsert).Insert(); err != nil {
		panic(err)
	}

	// selecting
	var toSelect []GeoCountry
	
	if err := db.Model(&toSelect).Select(); err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", toSelect)

	// user with fk
	newUser := User{
		Email:     "test@gmail.com",
		Activated: true,
		CountryID: &toSelect[0].ID,
	}

	// inserting
	if _, err := db.Model(&newUser).Insert(); err != nil {
		panic(err)
	}

	// selecting inserted user
	user := User{}
	m := db.Model(&user).
		Column(AllColumns, Columns.User.Country).
		Where(`? = ?`, pg.F(Columns.User.Email), "test@gmail.com")

	if err := m.Select(); err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", user)
	fmt.Printf("%#v\n", user.Country)

	// selecting inserted user with generated search
	user = User{}
	search := UserSearch{
		ID: newUser.ID,
	}
	m = db.Model(&user).Apply(search.Q())

	if err := m.Select(); err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", user)

}

```

#### Validation

Genna can generate simple validation functions for models with `--validator` flag. This function should check if values in model can be stored in database, and not intended to implement application logic

```sql
create type "enumvals" as enum ('one', 'two', 'three');

create table "example"
(
    "foreignKey"    int
        constraint "validationTest_userId_fkey"
            references users
            on update restrict on delete restrict,

    "notNullJSON"   jsonb    not null,
    "notNullHStore" hstore   not null,
    "enum"          enumvals not null,
    "limitedString" varchar(12)
)
```

```go
//lint:file-ignore U1000 ignore unused code, it's generated
package model

import (
	"unicode/utf8"
)

const (
	ErrEmptyValue = "empty"
	ErrMaxLength  = "len"
	ErrWrongValue = "value"
)

// ... Columns & Tables //

type Example struct {
	tableName struct{} `sql:"example,alias:t" pg:",discard_unknown_columns"`

	Enum          string                 `sql:"enum,notnull"`
	ForeignKey    *int                   `sql:"foreignKey"`
	LimitedString *string                `sql:"limitedString"`
	NotNullHStore map[string]string      `sql:"notNullHStore,hstore,notnull"`
	NotNullJSON   map[string]interface{} `sql:"notNullJSON,notnull"`
}

func (m Example) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	switch m.Enum {
	case "one", "two", "three":
	default:
		errors[Columns.Example.Enum] = ErrWrongValue
	}

	if m.ForeignKey != nil && *m.ForeignKey == 0 {
		errors[Columns.Example.ForeignKey] = ErrEmptyValue
	}

	if m.LimitedString != nil && utf8.RuneCountInString(*m.LimitedString) > 12 {
		errors[Columns.Example.LimitedString] = ErrMaxLength
	}

	if m.NotNullHStore == nil {
		errors[Columns.Example.NotNullHStore] = ErrEmptyValue
	}

	if m.NotNullJSON == nil {
		errors[Columns.Example.NotNullJSON] = ErrEmptyValue
	}

	return errors, len(errors) == 0
}
```  
 