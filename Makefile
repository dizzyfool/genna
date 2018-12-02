.PHONY: vendor lint

vendor:
	@GO111MODULE=on go mod vendor

lint:
	@golangci-lint run ./...