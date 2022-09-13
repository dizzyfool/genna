package bungen

import (
	"fmt"
	"log"

	"github.com/LdDl/bungen/model"
	"github.com/LdDl/bungen/util"
	"github.com/uptrace/bun"
)

// Bungen is  struct should be embedded to custom generator when bungen used as library
type Bungen struct {
	url string

	DB    *bun.DB
	Store *store

	Logger *log.Logger
}

// New creates Bungen
func New(url string, logger *log.Logger) Bungen {
	return Bungen{
		url:    url,
		Logger: logger,
	}
}

func (g *Bungen) Connect() error {
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
func (g *Bungen) Read(selected []string, followFK, useSQLNulls bool, customTypes model.CustomTypeMapping) ([]model.Entity, error) {
	if err := g.Connect(); err != nil {
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
			entities[i].AddColumn(c.Column(useSQLNulls, customTypes))
		}
	}

	for _, r := range relations {
		rel := r.Relation()
		if i, ok := index[util.Join(r.TargetSchema, r.TargetTable)]; ok {
			rel.AddEntity(&entities[i])
		}
		if i, ok := index[util.Join(r.SourceSchema, r.SourceTable)]; ok {
			entities[i].AddRelation(rel)
		}
	}
	return entities, nil
}
