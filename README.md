# Genna - cli tool for generating go-pg models

[![Go Report Card](https://goreportcard.com/badge/github.com/dizzyfool/genna)](https://goreportcard.com/report/github.com/dizzyfool/genna)


Requirements:
- [go-pg](https://github.com/go-pg/pg)
- your PostgreSQL database

### Idea

In most of the cases go-pg models represent database's tables and relations. Genna's main goal is to prepare those models by reading detailed information about PostrgeSQL database. The result should be several files with ready to use structs.

### Usage

1. Install `go get github.com/dizzyfool/genna`
1. Read though help `genna -h`

Currently genna support 3 generators:
- [model](#Model), that generates basic go-pg model
- [search](#Search), that generates search structs for basic model
- [validate](#Validate), that generates validate functions for basic model

### Example

First create your database and tables in it

```sql
create table "projects"
(
    "projectId" serial      not null,
    "name"      varchar(64) not null,

    primary key ("projectId")
);

create table "users"
(
    "userId"    serial      not null,
    "email"     varchar(64) not null,
    "activated" bool        not null default false,
    "name"      varchar(128),
    "countryId" integer,

    primary key ("userId")
);

create schema "geo";
create table geo."countries"
(
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

#### Model

Run generator

`genna model -c postgres://user:password@localhost:5432/yourdb -o ~/output/model.go -t public.* -f`

You should get following models on model package:

```go
//lint:file-ignore U1000 ignore unused code, it's generated
package model

var Columns = struct {
	Project struct {
		ID, Name string
	}
	User struct {
		ID, Activated, CountryID, Email, Name string

		Country string
	}
	GeoCountry struct {
		ID, Code, Coords string
	}
}{
	Project: struct {
		ID, Name string
	}{
		ID:   "projectId",
		Name: "name",
	},
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
	Project struct {
		Name, Alias string
	}
	User struct {
		Name, Alias string
	}
	GeoCountry struct {
		Name, Alias string
	}
}{
	Project: struct {
		Name, Alias string
	}{
		Name:  "projects",
		Alias: "t",
	},
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

type Project struct {
	tableName struct{} `sql:"projects,alias:t" pg:",discard_unknown_columns"`

	ID   int    `sql:"projectId,pk"`
	Name string `sql:"name,notnull"`
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

#### Search

Run generator

`genna validate -c postgres://user:password@localhost:5432/yourdb -o ~/output/search.go -t public.* -f`

You should get following search on model package:

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

type ProjectSearch struct {
	search

	ID   *int
	Name *string
}

func (s *ProjectSearch) Apply(query *orm.Query) *orm.Query {
	return s.apply(Tables.Project.Alias, map[string]interface{}{
		Columns.Project.ID:   s.ID,
		Columns.Project.Name: s.Name,
	}, query)
}

func (s *ProjectSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type UserSearch struct {
	search

	ID        *int
	Activated *bool
	CountryID *int
	Email     *string
	Name      *string
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

	ID   *int
	Code *string
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

#### Validate

Generated functions should check if values in model can be stored in database, and not intended to implement application logic

Run generator

`genna search -c postgres://user:password@localhost:5432/yourdb -o ~/output/validate.go -t public.* -f`

You should get following search on model package:

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

func (m Project) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	if utf8.RuneCountInString(m.Name) > 64 {
		errors[Columns.Project.Name] = ErrMaxLength
	}

	return errors, len(errors) == 0
}

func (m User) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	if m.CountryID != nil && *m.CountryID == 0 {
		errors[Columns.User.CountryID] = ErrEmptyValue
	}

	if utf8.RuneCountInString(m.Email) > 64 {
		errors[Columns.User.Email] = ErrMaxLength
	}

	if m.Name != nil && utf8.RuneCountInString(*m.Name) > 128 {
		errors[Columns.User.Name] = ErrMaxLength
	}

	return errors, len(errors) == 0
}

func (m GeoCountry) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	if utf8.RuneCountInString(m.Code) > 3 {
		errors[Columns.GeoCountry.Code] = ErrMaxLength
	}

	return errors, len(errors) == 0
}
```

### Try

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
 