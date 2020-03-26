package genna

import (
	"fmt"
	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"
	"log"

	"github.com/go-pg/pg/v9/orm"
)

// Genna is  struct should be embedded to custom generator when genna used as library
type Genna struct {
	url string

	DB    orm.DB
	Store *store

	Logger *log.Logger
}

// New creates Genna
func New(url string, logger *log.Logger) Genna {
	return Genna{
		url:    url,
		Logger: logger,
	}
}

func (g *Genna) connect() error {
	var err error

	if g.DB == nil {
		if g.DB, err = newDatabase(g.url, g.Logger); err != nil {
			return fmt.Errorf("unable to connect to DB: %w", err)
		}

		g.Store = newStore(g.DB)
	}

	return nil
}

// Read reads database and gets entities with columns and relations
func (g *Genna) Read(selected []string, followFK bool, useSQLNulls bool, goPGVer int) ([]model.Entity, error) {
	if err := g.connect(); err != nil {
		return nil, err
	}

	tables, err := g.Store.Tables(selected)
	if err != nil {
		return nil, err
	}

	if len(tables) == 0 {
		return nil, fmt.Errorf("no tables found")
	}

	relations, err := g.Store.Relations(tables)
	if err != nil {
		return nil, err
	}

	if followFK {
		set := util.NewSet()
		for _, t := range tables {
			set.Add(util.Join(t.Schema, t.Name))
		}

		for _, r := range relations {
			t := r.Target()
			if set.Add(util.Join(t.Schema, t.Name)) {
				tables = append(tables, t)
			}
		}
	}

	tables = Sort(tables)

	columns, err := g.Store.Columns(tables)
	if err != nil {
		return nil, err
	}

	entities := make([]model.Entity, len(tables))
	index := map[string]int{}
	for i, t := range tables {
		index[util.Join(t.Schema, t.Name)] = i
		entities[i] = t.Entity()
	}

	for _, c := range columns {
		if i, ok := index[util.Join(c.Schema, c.Table)]; ok {
			entities[i].AddColumn(c.Column(useSQLNulls, goPGVer))
		}
	}

	for _, r := range relations {
		rel := r.Relation()
		if i, ok := index[util.Join(r.SourceSchema, r.SourceTable)]; ok {
			entities[i].AddRelation(rel)
		}
		if i, ok := index[util.Join(r.TargetSchema, r.TargetTable)]; ok {
			rel.AddEntity(&entities[i])
		}
	}

	return entities, nil
}
