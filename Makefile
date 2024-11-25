export PORT=4321

.PHONY: build run

build:
	go generate ./...
	go build

run:
	./map