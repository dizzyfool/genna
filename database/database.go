package database

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type QueryLogger struct {
	logger *zap.Logger
}

func NewQueryLogger(logger *zap.Logger) QueryLogger {
	return QueryLogger{logger: logger}
}

func (ql QueryLogger) BeforeQuery(event *pg.QueryEvent) {
	event.Data["startedAt"] = time.Now()
}

func (ql QueryLogger) AfterQuery(event *pg.QueryEvent) {
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
