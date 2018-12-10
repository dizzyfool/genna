package database

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"go.uber.org/zap"
)

// NewDatabase creates database connection
func NewDatabase(url string, logger *zap.Logger) (orm.DB, error) {
	options, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := pg.Connect(options)
	client.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			logger.Error("formatted query error", zap.Error(err))
		}

		logger.Debug(query,
			zap.String("caller", caller(event)),
			zap.Duration("duration", time.Since(event.StartTime)),
		)
	})

	return client, nil
}

func caller(event *pg.QueryProcessedEvent) string {
	dir, file := filepath.Split(event.File)
	return fmt.Sprintf("%s:%d", filepath.Join(filepath.Base(dir), file), event.Line)
}
