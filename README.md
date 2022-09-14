# Bungen - cli tool for generating Bun [Postgres driver] models

[![Go Report Card](https://goreportcard.com/badge/github.com/LdDl/bungen)](https://goreportcard.com/report/github.com/LdDl/bungen)

About:
* This tool has been adapted from Genna - https://github.com/dizzyfool/genna which has been made for [go-pg](https://github.com/go-pg/pg#postgresql-client-and-orm-for-golang) package. I think it will be cool to have same CLI but for [Bun](https://github.com/uptrace/bun#sql-first-golang-orm-for-postgresql-mysql-mssql-and-sqlite) which is new evolution of `go-pg`

* Although this CLI is targeted for Bun, you should be aware that it's still compatible with Postgres only.

* I've done a lot of plain code replacements: current tests are fine, but I haven't managed cases with multiple FK's, composite FK's. I do know that this tool needs to manage DEFAULT values e.g. `default:'SOME DEFAULT FUNCTION'` also, but this needs more affort (pull requests are welcome).

Requirements:
- [bun](https://github.com/uptrace/bun)
- your PostgreSQL database

### Idea

In most of the cases Bun [Postgres driver] models represent database's tables and relations. Bungen's main goal is to prepare those models by reading detailed information about PostrgeSQL database. The result should be several files with ready to use structs with [Bun](https://github.com/uptrace/bun) ORM package.

### Usage

1. Install `go get github.com/LdDl/bungen && go install github.com/LdDl/bungen@latest`
1. Read though help `bungen -h`

Currently bungen support 3 generators:
- [model](generators/model/README.md), that generates basic Bun [Postgres driver] model
- [model-named](generators/named/README.md), same as basic but with named structs for columns and tables (author: [@Dionid](https://github.com/Dionid))
- [search](generators/search/README.md), that generates search structs for basic model
- [validation](generators/validate/README.md), that generates validate functions for basic model

Examples located in each generator
 
## Thanks
- I am thankful to [Genna](https://github.com/dizzyfool/genna#genna---cli-tool-for-generating-go-pg-models) and its creator [@dizzyfool](https://github.com/dizzyfool). Its [contributors](https://github.com/dizzyfool/genna/graphs/contributors) should me mentioned also. This CLI saved a lot of time for me in the past.
- Big shoutouts to [Bun](https://github.com/uptrace/bun#sql-first-golang-orm-for-postgresql-mysql-mssql-and-sqlite) creators for great ORM package for Golang
