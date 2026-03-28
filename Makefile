.PHONY: build run clean test

BINARY=termnode

build:
	go build -o $(BINARY) .

build-mqtt:
	go build -tags mqtt -o $(BINARY) .

run: build
	./$(BINARY)

debug: build
	./$(BINARY) -debug

test:
	go test ./...

clean:
	rm -f $(BINARY) debug.log

deps:
	go mod tidy
