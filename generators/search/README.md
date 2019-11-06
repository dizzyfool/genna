## Search model generator

**Basic model required, which can be generated with one of the model generators**

Use `search` sub-command to execute generator:

`genna search -h`

First create your database and tables in it

```sql
create table "projects"
(
    "projectId" serial not null,
    "name"      text   not null,

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

### Run generator

`genna search -c postgres://user:password@localhost:5432/yourdb -o ~/output/model.go -t public.* -f`

You should get following search structs on model package:

```go
//lint:file-ignore U1000 ignore unused code, it's generated
package model

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
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
	Email     *string
	Activated *bool
	Name      *string
	CountryID *int
}

func (s *UserSearch) Apply(query *orm.Query) *orm.Query {
	return s.apply(Tables.User.Alias, map[string]interface{}{
		Columns.User.ID:        s.ID,
		Columns.User.Email:     s.Email,
		Columns.User.Activated: s.Activated,
		Columns.User.Name:      s.Name,
		Columns.User.CountryID: s.CountryID,
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

### Try it

```go
package model

import (
	"fmt"
	"testing"

	"github.com/go-pg/pg/v9"
)

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

    code := "us"
    country := GeoCountry{}
    search := GeoCountrySearch{
        Code: &code,
    }
    m = db.Model(&country).Apply(search.Q())

    if err := m.Select(); err != nil {
        panic(err)
    }

	fmt.Printf("%#v\n", country)
}

```
