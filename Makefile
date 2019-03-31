SHELL := /bin/sh

GO := go
GO_OPTS := "-v"

.POSIX:
.PHONY: clean all deps build test

all: deps build test

build:
	${GO} build ${GO_OPTS}  ./...

install:
	${GO} install ${GO_OPTS} ./...

deps:
	go mod download
	go mod vendor

clean:
	rm -rf vendor
	${GO} clean ${GO_OPTS}  ./...

test: deps build
	${GO} test ${GO_OPTS}  ./...
