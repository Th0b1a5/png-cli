BINARY_NAME=png
VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")
LDFLAGS=-s -w -X main.version=$(VERSION)

.PHONY: all build release clean

all: build

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) main.go

release: clean
	mkdir -p releases
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o releases/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o releases/$(BINARY_NAME)-linux-arm64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o releases/$(BINARY_NAME)-darwin-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o releases/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o releases/$(BINARY_NAME)-windows-amd64.exe main.go
	GOOS=windows GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o releases/$(BINARY_NAME)-windows-arm64.exe main.go

push:
	git push
	git push --tags

clean:
	rm -rf releases $(BINARY_NAME)
