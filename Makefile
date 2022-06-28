.PHONY: build
build:
	go build -o main -v ./cmd/server

.PHONY: run
run:
	go build -o main -v ./cmd/server && ./main

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: tinkoff
tinkoff:
	protoc --go-grpc_out=pkg/tinkoff/investapi/ --go_out=pkg/tinkoff/investapi/ --proto_path=protos/ protos/*.proto

.DEFAULT_GOAL := build