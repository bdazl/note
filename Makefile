MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
RELEASE_VERSION := $(shell git describe --tags --always 2>/dev/null || git rev-parse --short HEAD)

# On Arch Linux, the package is extra/mingw-w64-gcc
WIN_CC := /usr/bin/x86_64-w64-mingw32-gcc
MODULE := github.com/bdazl/note

clean:
	rm -rf $(MAKEFILE_DIR)build

install-prereq:
	# https://github.com/mattn/go-sqlite3?tab=readme-ov-file#installation
	CGO_ENABLED=1 go install github.com/mattn/go-sqlite3

build-all: build-linux cross-build-windows

build-linux:
	mkdir -p build/amd64/linux
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X '$(MODULE)/cmd.Version=$(RELEASE_VERSION)'" -o build/amd64/linux/note

cross-build-windows:
	mkdir -p build/amd64/windows
	CC=$(WIN_CC) CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w -X '$(MODULE)/cmd.Version=$(RELEASE_VERSION)'" -o build/amd64/windows/note.exe
