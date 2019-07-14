package genna

import (
	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"

	"github.com/go-pg/pg/orm"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type Genna struct {
	url string

	DB    orm.DB
	Store *store

	Logger *zap.Logger
}

func New(url string, logger *zap.Logger) Genna {
	return Genna{
		url:    url,
		Logger: logger,
	}
}

func (g *Genna) connect() error {
	var err error

	if g.DB == nil {
		if g.DB, err = newDatabase(g.url, g.Logger); err != nil {
			return xerrors.Errorf("unable to connect to DB: %w", err)
		}

		g.Store = newStore(g.DB)
	}

	return nil
}

func (g *Genna) Read(selected []string, followFK bool, useSqlNulls bool) ([]model.Entity, error) {
	if err := g.connect(); err != nil {
		return nil, err
	}

	tables, err := g.Store.Tables(selected)
	if err != nil {
		return nil, err
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
			entities[i].AddColumn(c.Column(useSqlNulls))
		}
	}

	for _, r := range relations {
		if i, ok := index[util.Join(r.SourceSchema, r.SourceTable)]; ok {
			entities[i].AddRelation(r.Relation())
		}
	}

	return entities, nil
}
