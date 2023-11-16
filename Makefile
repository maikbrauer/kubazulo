default: help

build:
	go build -v -o . ./...

test:
	go test -v ./...

clean:
	rm -rf bin/*

help:
	@echo 'Usage: make (build | clean)'