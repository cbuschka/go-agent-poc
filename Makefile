PROJECT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SHELL = /bin/bash
ifeq (${GOPATH},)
        GOPATH := ${HOME}/go
endif

build:	build-agent build-master

build-agent:	generate-protocol
	@cd ${PROJECT_DIR}
	mkdir -p dist/
	go build -o dist/agent cmd/agent/main.go

build-master:	init
	@cd ${PROJECT_DIR}
	mkdir -p dist/
	go build -o dist/master cmd/master/main.go

generate-protocol:	init
	@cd ${PROJECT_DIR}
	@mkdir -p internal/protocol/generated/
	PATH="$$PATH:$$(go env GOPATH)/bin" protoc \
		--go_out=internal/protocol/generated/ --go_opt=paths=source_relative \
		--go-grpc_out=internal/protocol/generated/ --go-grpc_opt=paths=source_relative protocol.proto

init:
	@if [ -z "$$(protoc --version | cut -d ' ' -f 2)" ]; then \
		echo "protoc missing. you can find install instructions at https://grpc.io/docs/protoc-installation/" && \
		exit 1; \
	fi
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

clean:
	@cd ${PROJECT_DIR}
	rm -rf dist/
