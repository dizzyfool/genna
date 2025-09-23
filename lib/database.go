package genna

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// queryLogger helper struct for query logging
type queryLogger struct {
	logger *log.Logger
}

// newQueryLogger creates new helper struct for query logging
func newQueryLogger(logger *log.Logger) queryLogger {
	return queryLogger{logger: logger}
}

// BeforeQuery stores start time in custom data array
func (ql queryLogger) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	event.Stash = make(map[interface{}]interface{})
	event.Stash["startedAt"] = time.Now()

	return ctx, nil
}

// AfterQuery calculates execution time and print it with formatted query
func (ql queryLogger) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	query, err := event.FormattedQuery()
	if err != nil {
		ql.logger.Printf("formatted query error: %s", err)
	}

	var since time.Duration
	if event.Stash != nil {
		if v, ok := event.Stash["startedAt"]; ok {
			if startAt, ok := v.(time.Time); ok {
				since = time.Since(startAt)
			}
		}
	}

	ql.logger.Printf("query: %s, duration: %d", query, since)
	return nil
}

// newDatabase creates database connection
func newDatabase(url string, logger *log.Logger) (orm.DB, error) {
	options, err := pg.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parsing connection url error: %w", err)
	}

	client := pg.Connect(options)

	if logger != nil {
		client.AddQueryHook(newQueryLogger(logger))
	}

	return client, nil
}
