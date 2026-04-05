GOLANGCI_LINT_PACKAGE ?= github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4

# -------------------------------------------------------------------------------------------------
# main
# -------------------------------------------------------------------------------------------------

all: help

## build: Compile templ files and build application
.PHONY: build
build: prepare-data
	CGO_ENABLED=0 go build -tags swagger -ldflags="-s -w -extldflags '-static'" -trimpath -o 'bin/app' ./cmd/app

## start: Build and start application
.PHONY: start
start: prepare-data
	go run -tags swagger ./cmd/app

## dev: Build and start application in live reload mode
.PHONY: dev
dev: air

## build-docker: Build Docker container image with this app
.PHONY: build-docker
build-docker:
	docker build -t $(shell basename $(PWD)):latest --no-cache -f Dockerfile .

## run-docker: Run Docker container image with this app
.PHONY: run-docker
run-docker:
	docker run --rm -it -p 3000:3000 $(shell basename $(PWD)):latest

## prepare-data: Prepare data for the application
.PHONY: prepare-data
prepare-data: .deps-stamp get-js-deps generate-web generate-swagger

.deps-stamp: go.mod go.sum
	go mod download
	@touch .deps-stamp

## get-js-deps: Install frontend dependencies using bun (locally if available and otherwise via Docker)
.PHONY: get-js-deps
get-js-deps:
	@command -v bun &> /dev/null && bun install || docker run --rm -v "$(PWD):/app" -w /app oven/bun:1.3 bun install
	@mkdir -p internal/web/static/js internal/web/static/css
	@cp node_modules/htmx.org/dist/htmx.min.js internal/web/static/js/
	@cp node_modules/htmx-ext-response-targets/dist/response-targets.min.js internal/web/static/js/
	@bun tailwindcss -i internal/web/static/css/input.css -o internal/web/static/css/style.css --minify
	@cp -r node_modules/ionicons/dist/ionicons internal/web/static/js/

# -------------------------------------------------------------------------------------------------
# testing
# -------------------------------------------------------------------------------------------------

## test: Run unit tests
.PHONY: test
test: check-go
	@go test -v -count=1 ./...

# -------------------------------------------------------------------------------------------------
# tools
# -------------------------------------------------------------------------------------------------

## generate-web: Compile templ files via github.com/a-h/templ/cmd/templ
.PHONY: generate-web
generate-web: check-go
	go tool templ generate

## generate-swagger: Generate swagger documentation via swaggo/swag
.PHONY: generate-swagger
generate-swagger: check-go
	go tool swag init --dir ./cmd/app,./internal/web/handlers -g main.go --parseDependency --parseInternal

## air: Build and start application in live reload mode via air
.PHONY: air
air: prepare-data
	go tool air

## lint: Run golangci-lint to lint Go files
.PHONY: lint
lint:
	go run $(GOLANGCI_LINT_PACKAGE) run

## lint-fix: Run golangci-lint to lint Go files and fix issues
.PHONY: lint-fix
lint-fix:
	go run $(GOLANGCI_LINT_PACKAGE) run --fix

## format: Run golangci-lint fmt to show code format issues
.PHONY: format
format:
	go run $(GOLANGCI_LINT_PACKAGE) fmt

## audit: Quality checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# -------------------------------------------------------------------------------------------------
# shared
# -------------------------------------------------------------------------------------------------

## tidy: Removes unused dependencies and adds missing ones
.PHONY: tidy
tidy: check-go
	go mod tidy

## update-deps: Update Go dependencies
.PHONY: update-deps
update-deps: check-go
	go get -u ./...
	-@$(MAKE) tidy

## get-deps: Download Go dependencies
.PHONY: get-deps
get-deps: check-go
	go mod download

## check-go: Check that Go is installed
.PHONY: check-go
check-go:
	@command -v go &> /dev/null || (echo "Please install GoLang" && false)

## help: Display help
.PHONY: help
help: Makefile
	@echo "Usage:  make COMMAND"
	@echo
	@echo "Commands:"
	@sed -n 's/^##//p' $< | column -ts ':' |  sed -e 's/^/ /'
