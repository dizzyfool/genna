# Genna - cli tool for generating go-pg models

[![Go Report Card](https://goreportcard.com/badge/github.com/dizzyfool/genna)](https://goreportcard.com/report/github.com/dizzyfool/genna)

Requirements:
- [go-pg](https://github.com/go-pg/pg)
- your PostgreSQL database

#### Idea

In most of the cases go-pg models represent database's tables and relations. Genna's main goal is to prepare those models by reading detailed information about PostrgeSQL database. The result should be several files with ready to use structs.

#### Usage

1. Install `go get github.com/dizzyfool/genna`
2. Read though help `genna -h`
3. Run `genna -c postgres://user:password@localhost:5432/yourdb -o ~/output`


 