.PHONY: all
all: compile tests

.PHONY: compile
compile:
	go build

.PHONY: fmt
fmt:
	go fmt

.PHONY: tests
tests:
	go test

.PHONY: start
start:
	go run main.go
