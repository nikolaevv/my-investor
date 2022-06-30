.PHONY: build
build:
	go build -o main -v ./cmd

.PHONY: run
run:
	go build -o main -v ./cmd && ./main

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: tinkoff
tinkoff:
	protoc --go-grpc_out=pkg/gen/proto/tinkoff/investapi/ --go_out=pkg/tinkoff/investapi/ --proto_path=proto/ proto/*.proto

.DEFAULT_GOAL := build