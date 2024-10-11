.PHONY: ci
ci: tools test lint generate fmt tidy check_repository_unchanged

.PHONY: check_repository_unchanged
check_repository_unchanged:
	_tools/check_repository_unchanged.sh

.PHONY: generate
generate:
	go generate ./...

.PHONY: fmt
fmt:
	gosimports -l -w ./

.PHONY: test
test:
	go test -race ./...

.PHONY: start-services
start-services:
	docker compose rm -f -s -v
	docker compose up -d

.PHONY: test-integration
test-integration:
	FIRESTORE_EMULATOR_HOST=localhost:8200 go test -tags=test_integration -v ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	go vet ./...
	golangci-lint run ./...

.PHONY: tools
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1
	go install github.com/rinchsan/gosimports/cmd/gosimports@v0.3.5 # https://github.com/golang/go/issues/20818
