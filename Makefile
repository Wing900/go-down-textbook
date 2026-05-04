BINARY=go-down-textbook
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/$(BINARY).exe ./cmd/go-down-textbook

build-all:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/$(BINARY)-windows-amd64.exe ./cmd/go-down-textbook
	GOOS=darwin  GOARCH=amd64 go build -ldflags "-s -w" -o bin/$(BINARY)-darwin-amd64 ./cmd/go-down-textbook
	GOOS=darwin  GOARCH=arm64 go build -ldflags "-s -w" -o bin/$(BINARY)-darwin-arm64 ./cmd/go-down-textbook
	GOOS=linux   GOARCH=amd64 go build -ldflags "-s -w" -o bin/$(BINARY)-linux-amd64 ./cmd/go-down-textbook

clean:
	rm -rf bin/

run: build
	./bin/$(BINARY).exe

test:
	go vet ./...
	go build ./...

.PHONY: build build-all clean run test
