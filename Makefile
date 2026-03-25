.PHONY: build install test clean

build:
	go build -o bin/mobileops ./cmd/mobileops

install:
	go install ./cmd/mobileops

test:
	go test ./...

clean:
	rm -rf bin/
