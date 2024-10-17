export PORT=4321
export 

.PHONY: build run

build:
	go generate ./...
	go build

run:
	./map