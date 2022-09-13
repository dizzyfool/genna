.PHONY: lint test

lint:
	@golangci-lint run ./...

test:
	psql postgres://some_user:some_password@localhost:5432/some_db?sslmode=disable -f test_db.sql
	go test ./...