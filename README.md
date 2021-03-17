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
- [model](generators/model/README.md), that generates basic go-pg model
- [model-named](generators/named/README.md), same as basic but with named structs for columns and tables (author: [@Dionid](https://github.com/Dionid))
- [search](generators/search/README.md), that generates search structs for basic model
- [validation](generators/validate/README.md), that generates validate functions for basic model

Examples located in each generator
 