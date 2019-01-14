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
1. Run `genna -c postgres://user:password@localhost:5432/yourdb -o ~/output`

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

`genna -c postgres://user:password@localhost:5432/yourdb -o ~/output -t public.*,geo.* -f`

You should get following models on model package:

```go
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
		Activated: "activated,notnull",
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
		Name string
	}
	GeoCountry struct {
		Name string
	}
}{
	User: struct {
		Name string
	}{
		Name: "users",
	},
	GeoCountry: struct {
		Name string
	}{
		Name: "geo.countries",
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

Try it

```go
package model

import (
	"fmt"

	"github.com/go-pg/pg"
)

const AllColumns = "t.*"

func Test() {
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
	toSelect := []GeoCountry{}
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
	user := User{
		ID: newUser.ID,
	}
	m := db.Model(&user).
		Column(AllColumns, Columns.User.Country).
		Where(`? = ?`, pg.F(Columns.User.Email), "test@gmail.com")

	if err := m.Select(); err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", user)
	fmt.Printf("%#v\n", user.Country)

}

```

 