# Makefile - kwtsms-cli cross-compilation targets
# Usage:
#   make build       - build for current platform
#   make build-all   - build for all 6 supported platforms into dist/
#   make clean       - remove build artifacts
#
# Requires: Go installed. No CGO. No external tools needed.
# CGO_ENABLED=0 ensures pure Go builds for true cross-compilation.

BINARY  = kwtsms-cli
DIST    = dist
LDFLAGS = -ldflags="-s -w"

.PHONY: build build-all clean test vet fmt

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY) .

build-all: clean
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64         go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-x64 .
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64         go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-arm64 .
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm GOARM=7   go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-armv7 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64         go build $(LDFLAGS) -o $(DIST)/$(BINARY)-macos-x64 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64         go build $(LDFLAGS) -o $(DIST)/$(BINARY)-macos-arm64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64         go build $(LDFLAGS) -o $(DIST)/$(BINARY)-windows-x64.exe .
	@echo "Built all targets:"
	@ls -lh $(DIST)/

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

clean:
	rm -rf $(DIST) $(BINARY) $(BINARY).exe
