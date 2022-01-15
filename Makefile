PROJECT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SHELL = /bin/bash

build:	build-agent build-master

build-agent:
	@cd ${PROJECT_DIR}
	mkdir -p dist/
	go build -o dist/agent cmd/agent/main.go

build-master:
	@cd ${PROJECT_DIR}
	mkdir -p dist/
	go build -o dist/master cmd/master/main.go

clean:
	@cd ${PROJECT_DIR}
	rm -rf dist/
