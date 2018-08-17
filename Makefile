
# can specify V=1 on the line with `make` to get verbose output
V ?= 0
Q = $(if $(filter 1,$V),,@)

.DEFAULT_GOAL := help

deps:  ## download dependencies
	$Q go get ./...

generate:  ## run a go generate
	$Q go generate ./...

setup: deps generate  ## attempt to get everything set up to do a build (deps and generate)

test:  ## run go test
	$Q go test ./...

vet:  ## run various linters and vetters
	$Q golint ./...
	$Q go vet ./...
	$Q gometalinter -D gas -D gocyclo -D goconst -e .pb.go -e _easyjson.go --warn-unmatched-nolint --enable-gc --deadline 180s ./...

help:  ## Show the help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' ./Makefile