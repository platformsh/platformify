PROJECT_NAME = platformify
UPSUN_PROJECT_NAME = upsunify

.PHONY: test
test: generate ## Run unit tests
	go clean -testcache
	go test -v -race -mod=readonly -cover ./...

.PHONY: lint
lint: ## Run linter
	command -v golangci-lint >/dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59
	golangci-lint run --timeout=10m --verbose

.PHONY: lint-fix
lint-fix: ## Run linter fixes
	command -v golangci-lint >/dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52
	golangci-lint run --fix

.PHONY: generate
generate: ## Generate mock data
	command -v mockgen >/dev/null || go install github.com/golang/mock/mockgen@latest
	go generate ./...

.PHONY: clean
clean: ## Remove the binary file from the root folder
	rm -rf $(PROJECT_NAME) $(UPSUN_PROJECT_NAME)

.PHONY: build
build: clean build-upsun build-platform ## Compile the app into the root folder

build-upsun:
	go build -o $(UPSUN_PROJECT_NAME) ./cmd/upsunify/

build-platform:
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
