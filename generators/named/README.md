## Basic model generator with named structures

Author: [@Dionid](https://github.com/Dionid)

Use `model-named` sub-command to execute generator:

`genna model-named -h`

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

`genna model-named -c postgres://user:password@localhost:5432/yourdb -o ~/output/model.go -t public.* -f`

You should get following models on model package:

```go
//lint:file-ignore U1000 ignore unused code, it's generated
package model

type ColumnsProject struct {
	ID, Name string
}

type ColumnsUser struct {
	ID, Email, Activated, Name, CountryID string
	Country                               string
}

type ColumnsGeoCountry struct {
	ID, Code, Coords string
}

type ColumnsSt struct {
	Project    ColumnsProject
	User       ColumnsUser
	GeoCountry ColumnsGeoCountry
}

var Columns = ColumnsSt{
	Project: ColumnsProject{
		ID:   "projectId",
		Name: "name",
	},
	User: ColumnsUser{
		ID:        "userId",
		Email:     "email",
		Activated: "activated",
		Name:      "name",
		CountryID: "countryId",

		Country: "Country",
	},
	GeoCountry: ColumnsGeoCountry{
		ID:     "countryId",
		Code:   "code",
		Coords: "coords",
	},
}

type TableProject struct {
	Name, Alias string
}

type TableUser struct {
	Name, Alias string
}

type TableGeoCountry struct {
	Name, Alias string
}

type TablesSt struct {
	Project    TableProject
	User       TableUser
	GeoCountry TableGeoCountry
}

var Tables = TablesSt{
	Project: TableProject{
		Name:  "projects",
		Alias: "t",
	},
	User: TableUser{
		Name:  "users",
		Alias: "t",
	},
	GeoCountry: TableGeoCountry{
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
	Email     string  `sql:"email,notnull"`
	Activated bool    `sql:"activated,notnull"`
	Name      *string `sql:"name"`
	CountryID *int    `sql:"countryId"`

	Country *GeoCountry `pg:"fk:countryId"`
}

type GeoCountry struct {
	tableName struct{} `sql:"geo.countries,alias:t" pg:",discard_unknown_columns"`

	ID     int    `sql:"countryId,pk"`
	Code   string `sql:"code,notnull"`
	Coords []int  `sql:"coords,array"`
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
}

```
