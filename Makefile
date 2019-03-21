SHELL := /bin/sh

GO := go
GO_OPTS := "-v"

.POSIX:
.PHONY: clean all deps build test

all: deps build test

build:
	${GO} build ${GO_OPS}  ./...

deps:
	dep ensure

clean:
	rm -rf vendor
	${GO} clean ${GO_OPS}  ./...

test: deps build
	${GO} test ${GO_OPS}  ./...
