default: help

build:
	go build -v -o bin/kubazulo .

clean:
	rm -rf bin/*

help:
	@echo 'Usage: make (build | clean)'