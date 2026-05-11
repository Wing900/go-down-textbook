BINARY=BoooookDown
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GOFLAGS=-trimpath -buildvcs=false
LDFLAGS=-s -w -X main.version=$(VERSION)

build:
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY).exe ./cmd/boooookdown

package-windows:
	powershell -ExecutionPolicy Bypass -File scripts/build-windows.ps1 -Version $(VERSION)

package-all:
	powershell -ExecutionPolicy Bypass -File scripts/package-all.ps1 -Version $(VERSION)

build-all:
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-windows-amd64.exe ./cmd/boooookdown
	GOOS=darwin  GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-darwin-amd64 ./cmd/boooookdown
	GOOS=darwin  GOARCH=arm64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-darwin-arm64 ./cmd/boooookdown
	GOOS=linux   GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-linux-amd64 ./cmd/boooookdown

clean:
	rm -rf bin/

run: build
	./bin/$(BINARY).exe

test:
	go vet ./...
	go build ./...

.PHONY: build build-all clean run test package-windows package-all
