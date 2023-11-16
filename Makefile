default: help

build:
	go build -v -o kubazulo

test:
	go test -v ./...

clean:
	rm -rf bin/*

help:
	@echo 'Usage: make (build | test | clean)'