package database

import (
	"context"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// NewQueryLogger helper struct for query logging
type QueryLogger struct {
	logger *zap.Logger
}

// NewQueryLogger creates new helper struct for query logging
func NewQueryLogger(logger *zap.Logger) QueryLogger {
	return QueryLogger{logger: logger}
}

// BeforeQuery stores start time in custom data array
func (ql QueryLogger) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	if event.Stash != nil {
		event.Stash["startedAt"] = time.Now()
	}
	return ctx, nil
}

// AfterQuery calculates execution time and print it with formatted query
func (ql QueryLogger) AfterQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	query, err := event.FormattedQuery()
	if err != nil {
		ql.logger.Error("formatted query error", zap.Error(err))
	}

	var since time.Duration
	if event.Stash != nil {
		if v, ok := event.Stash["startedAt"]; ok {
			if startAt, ok := v.(time.Time); ok {
				since = time.Since(startAt)
			}
		}
	}

	ql.logger.Debug(query, zap.Duration("duration", since))
	return ctx, nil
}

// NewDatabase creates database connection
func NewDatabase(url string, logger *zap.Logger) (orm.DB, error) {
	options, err := pg.ParseURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "parsing connection url error")
	}

	client := pg.Connect(options)
	client.AddQueryHook(NewQueryLogger(logger))

	return client, nil
}
