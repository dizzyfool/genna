package genna

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

// queryLogger helper struct for query logging
type queryLogger struct {
	logger *zap.Logger
}

// newQueryLogger creates new helper struct for query logging
func newQueryLogger(logger *zap.Logger) queryLogger {
	return queryLogger{logger: logger}
}

// BeforeQuery stores start time in custom data array
func (ql queryLogger) BeforeQuery(event *pg.QueryEvent) {
	event.Data = make(map[interface{}]interface{})
	event.Data["startedAt"] = time.Now()
}

// AfterQuery calculates execution time and print it with formatted query
func (ql queryLogger) AfterQuery(event *pg.QueryEvent) {
	query, err := event.FormattedQuery()
	if err != nil {
		ql.logger.Error("formatted query error", zap.Error(err))
	}

	var since time.Duration
	if event.Data != nil {
		if v, ok := event.Data["startedAt"]; ok {
			if startAt, ok := v.(time.Time); ok {
				since = time.Since(startAt)
			}
		}
	}

	ql.logger.Debug(query, zap.Duration("duration", since))
}

// newDatabase creates database connection
func newDatabase(url string, logger *zap.Logger) (orm.DB, error) {
	options, err := pg.ParseURL(url)
	if err != nil {
		return nil, xerrors.Errorf("parsing connection url error: %w", err)
	}

	client := pg.Connect(options)

	if logger != nil {
		client.AddQueryHook(newQueryLogger(logger))
	}

	return client, nil
}
