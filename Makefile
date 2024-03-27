.PHONY: build

build:
	go build -o bonuse-service -v ./cmd/main.go

.PHONY: test
test:
	go clean -testcache
	go test -v -race -cover -timeout 30s ./internal/repository/mongodb ./internal/service

.PHONY: lint
lint:
	golangci-lint run
.DEFAULT_GOAL := build