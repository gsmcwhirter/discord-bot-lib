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
	$Q go get ./...

generate:  ## run a go generate
	$Q go generate ./...

setup: deps generate  ## attempt to get everything set up to do a build (deps and generate)

test:  ## run go test
	$Q go test ./...

version:  ## Print the version string and git sha that would be recorded if a release was built now
	$Q echo $(VERSION) $(GIT_SHA)

vet:  ## run various linters and vetters
	$Q golint ./...
	$Q go vet ./...
	$Q gometalinter -D gas -D gocyclo -D goconst -e .pb.go -e _easyjson.go --warn-unmatched-nolint --enable-gc --deadline 180s ./...

help:  ## Show the help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' ./Makefile