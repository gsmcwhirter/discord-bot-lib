BUILD_DATE := `date -u +%Y%m%d`
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo v0.0.1)
GIT_SHA := $(shell git rev-parse HEAD)

APP_NAME := test-stress
PROJECT := github.com/gsmcwhirter/discord-bot-lib

GOPROXY ?= https://proxy.golang.org

# can specify V=1 on the line with `make` to get verbose output
V ?= 0
Q = $(if $(filter 1,$V),,@)

.DEFAULT_GOAL := help

build-stress: version generate
	$Q GOPROXY=$(GOPROXY) go build -v -ldflags "-X main.AppName=$(APP_NAME) -X main.BuildVersion=$(VERSION) -X main.BuildSHA=$(GIT_SHA) -X main.BuildDate=$(BUILD_DATE)" -o bin/$(APP_NAME) -race $(PROJECT)/cmd/$(APP_NAME)

deps:  ## download dependencies
	$Q GOPROXY=$(GOPROXY) go mod tidy
	$Q GOPROXY=$(GOPROXY) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0

generate:  ## run a go generate
	$Q GOPROXY=$(GOPROXY) go generate ./...

setup: deps generate  ## attempt to get everything set up to do a build (deps and generate)

test:  ## run go test
	$Q GOPROXY=$(GOPROXY) go test ./...

test-coverage:
	$Q GOPROXY=$(GOPROXY) go test -coverprofile=coverage.out ./...
	$Q go tool cover -html=coverage.out

version:  ## Print the version string and git sha that would be recorded if a release was built now
	$Q echo $(VERSION) $(GIT_SHA)

vet: deps generate ## run various linters and vetters
	$Q bash -c 'for d in $$(go list -f {{.Dir}} ./...); do gofmt -s -w $$d/*.go; done'
	$Q bash -c 'for d in $$(go list -f {{.Dir}} ./...); do goimports -w -local $(PROJECT) $$d/*.go; done'
	$Q golangci-lint run -c .golangci.yml -E revive,gosimple,staticcheck ./...
	$Q golangci-lint run -c .golangci.yml -E asciicheck,contextcheck,depguard,durationcheck,errcheck,errname,gocritic,gofumpt,goimports,gosec,govet,ineffassign,nakedret,paralleltest,prealloc,predeclared,typecheck,unconvert,unused,whitespace ./...
	$Q golangci-lint run -c .golangci.yml -E godox ./... || true

help:  ## Show the help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' ./Makefile