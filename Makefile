BUILD_DATE := `date -u +%Y%m%d`
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo v0.0.1)
GIT_SHA := $(shell git rev-parse HEAD)

APP_NAME := test-stress
PROJECT := github.com/gsmcwhirter/discord-bot-lib

# can specify V=1 on the line with `make` to get verbose output
V ?= 0
Q = $(if $(filter 1,$V),,@)

.DEFAULT_GOAL := help

build-stress: version generate
	$Q go build -v -ldflags "-X main.AppName=$(APP_NAME) -X main.BuildVersion=$(VERSION) -X main.BuildSHA=$(GIT_SHA) -X main.BuildDate=$(BUILD_DATE)" -o bin/$(APP_NAME) -race $(PROJECT)/cmd/$(APP_NAME)

deps:  ## download dependencies
	$Q go mod download
	$Q go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1
	$Q go get github.com/mailru/easyjson/easyjson
	$Q go get github.com/valyala/quicktemplate/qtc
	$Q go get golang.org/x/tools/cmd/stringer
	$Q go get golang.org/x/tools/cmd/goimports

generate:  ## run a go generate
	$Q go generate ./...

setup: deps generate  ## attempt to get everything set up to do a build (deps and generate)

test:  ## run go test
	$Q go test ./...

version:  ## Print the version string and git sha that would be recorded if a release was built now
	$Q echo $(VERSION) $(GIT_SHA)

vet:  deps ## run various linters and vetters
	$Q golangci-lint run -E deadcode,depguard,errcheck,gocritic,gofmt,goimports,gosec,govet,ineffassign,nakedret,prealloc,structcheck,typecheck,unconvert,varcheck ./...

help:  ## Show the help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' ./Makefile