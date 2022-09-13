package bungen

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// queryLogger helper struct for query logging
type queryLogger struct {
	logger log.Logger
}

// newQueryLogger creates new helper struct for query logging
func newQueryLogger(logger log.Logger) queryLogger {
	return queryLogger{logger: logger}
}

// BeforeQuery stores start time in custom data array
func (ql queryLogger) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {

	event.Stash = make(map[interface{}]interface{})
	event.Stash["startedAt"] = time.Now()

	return ctx
}

// AfterQuery calculates execution time and print it with formatted query
func (ql queryLogger) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	query := event.Operation()
	var since time.Duration
	if event.Stash != nil {
		if v, ok := event.Stash["startedAt"]; ok {
			if startAt, ok := v.(time.Time); ok {
				since = time.Since(startAt)
			}
		}
	}
	ql.logger.Printf("query: %s, duration: %d", query, since)
}

// newDatabase creates database connection
func newDatabase(dsn string, logger *log.Logger) (*bun.DB, error) {
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	// bun.WithDiscardUnknownColumns() instead of `discard_unknown_columns` tag
	client := bun.NewDB(pgdb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	if logger != nil {
		client.AddQueryHook(newQueryLogger(*logger))
	}

	return client, nil
}
