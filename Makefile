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
start: compile
	time -f '%Uu %Ss %er %MkB %C' ./eod_scanner -config .config.json

.PHONY: cpuprofile
cpuprofile: compile
	./eod_scanner -config .config.json -cpuprofile eod_scanner.prof

.PHONY: memprofile
memprofile: compile
	./eod_scanner -config .config.json -memprofile eod_scanner.mprof
