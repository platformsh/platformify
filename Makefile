PROJECT_NAME = platformify

.PHONY: test
test: ## Run unit tests
	go clean -testcache
	go test -v -race -mod=readonly -cover ./...

.PHONY: lint
lint: ## Run linter
	command -v golangci-lint >/dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run --timeout=10m --verbose

.PHONY: clean
clean: ## Remove the binary file from the root folder
	rm -rf $(PROJECT_NAME)

.PHONY: build
build: clean ## Compile the app into the root folder
	go build -o $(PROJECT_NAME) ./cmd/platformify/

.PHONY: snapshot
snapshot: ## Build snapshot
	command -v goreleaser >/dev/null || go install github.com/goreleaser/goreleaser@latest
	goreleaser build --snapshot --clean

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| cut -d ':' -f 1,2 \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
