.PHONY: lint test

lint:
	@golangci-lint run ./...

test:
	psql postgres://genna:genna@localhost:5432/genna?sslmode=disable -f test_db.sql
	go test ./...