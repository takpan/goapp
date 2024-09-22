.DEFAULT_GOAL := goapp

.PHONY: all
all: clean goapp cli test

.PHONY: goapp
goapp:
	mkdir -p bin
	go build -o bin ./cmd/server

.PHONY: cli
cli:
	mkdir -p bin
	go build -o bin ./cmd/cli

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean
	rm -f bin/*
