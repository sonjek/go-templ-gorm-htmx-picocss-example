
# -------------------------------------------------------------------------------------------------
# main
# -------------------------------------------------------------------------------------------------

all: help

## build: Compile templ files and build application
.PHONY: build
build: get-deps generate-web
	CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -trimpath -o 'bin/app' ./cmd/app

## start: Build and start application
.PHONY: start
start: get-deps generate-web
	go run ./cmd/app

## build-docker: Build Docker container image with this app
.PHONY: build-docker
build-docker:
	docker build -t $(shell basename $(PWD)):latest --no-cache -f Dockerfile .

## run-docker: Run Docker container image with this app
.PHONY: run-docker
run-docker:
	docker run --rm -it -p 8089:8089 $(shell basename $(PWD)):latest

# -------------------------------------------------------------------------------------------------
# testing
# -------------------------------------------------------------------------------------------------

## test: Run unit tests
.PHONY: test
test: check-go
	@go test -v -count=1 ./...

# -------------------------------------------------------------------------------------------------
# tools && shared
# -------------------------------------------------------------------------------------------------

## tidy: Removes unused dependencies and adds missing ones
.PHONY: tidy
tidy: check-go
	go mod tidy

## update-deps: Update go dependencies
.PHONY: update-deps
update-deps: check-go
	go get -u ./...
	-@$(MAKE) tidy

## get-deps: Download go dependencies
.PHONY: get-deps
get-deps: check-go
	go mod download

## generate-web: Compile templ files via github.com/a-h/templ/cmd/templ
.PHONY: generate-web
generate-web: check-go
	go tool templ generate

## air: Build and start application in live reload mode via air
.PHONY: air
air: get-deps generate-web
	go tool air

## format: Fix code format issues
.PHONY: format
format:
	go run mvdan.cc/gofumpt@latest -w -l .

## deadcode: Run deadcode tool for find unreachable functions
deadcode:
	go run golang.org/x/tools/cmd/deadcode@latest -test ./...

## audit: Quality checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

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
